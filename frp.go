// Package frp is a togo tunnel driver for frp (Fast Reverse Proxy). It wraps the
// `frpc` client: it renders an frpc TOML config for an HTTP proxy to your local
// port and runs frpc against your frps server.
//
// Install: `togo install togo-framework/tunnel-frp`, set TUNNEL_DRIVER=frp and the
// FRP_* env. Requires the `frpc` binary (https://github.com/fatedier/frp) and an
// frps server you control.
package frp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/togo-framework/togo"
	"github.com/togo-framework/tunnel"
)

func init() {
	tunnel.RegisterDriver("frp", func(k *togo.Kernel) (tunnel.Tunnel, error) {
		if os.Getenv("FRP_SERVER_ADDR") == "" {
			return nil, fmt.Errorf("tunnel-frp: FRP_SERVER_ADDR not set")
		}
		return &driver{bin: envOr("FRPC_BIN", "frpc")}, nil
	})
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

// config holds the frpc settings resolved from env + the local port.
type config struct {
	ServerAddr   string
	ServerPort   string
	Token        string
	ProxyName    string
	LocalPort    string
	Subdomain    string // HTTP vhost subdomain (one of subdomain/customDomain)
	CustomDomain string
	ServerDomain string // base domain for the subdomain URL
	VhostPort    string // public HTTP vhost port (default 80)
}

func configFromEnv(port string) config {
	return config{
		ServerAddr:   os.Getenv("FRP_SERVER_ADDR"),
		ServerPort:   envOr("FRP_SERVER_PORT", "7000"),
		Token:        os.Getenv("FRP_TOKEN"),
		ProxyName:    envOr("FRP_PROXY_NAME", "togo-"+port),
		LocalPort:    port,
		Subdomain:    os.Getenv("FRP_SUBDOMAIN"),
		CustomDomain: os.Getenv("FRP_CUSTOM_DOMAIN"),
		ServerDomain: envOr("FRP_SERVER_DOMAIN", os.Getenv("FRP_SERVER_ADDR")),
		VhostPort:    envOr("FRP_VHOST_PORT", "80"),
	}
}

// render produces an frpc TOML config (frp v0.52+).
func (c config) render() string {
	var b strings.Builder
	fmt.Fprintf(&b, "serverAddr = %q\n", c.ServerAddr)
	fmt.Fprintf(&b, "serverPort = %s\n", c.ServerPort)
	if c.Token != "" {
		fmt.Fprintf(&b, "auth.method = \"token\"\nauth.token = %q\n", c.Token)
	}
	b.WriteString("\n[[proxies]]\n")
	fmt.Fprintf(&b, "name = %q\n", c.ProxyName)
	b.WriteString("type = \"http\"\n")
	b.WriteString("localIP = \"127.0.0.1\"\n")
	fmt.Fprintf(&b, "localPort = %s\n", c.LocalPort)
	switch {
	case c.CustomDomain != "":
		fmt.Fprintf(&b, "customDomains = [%q]\n", c.CustomDomain)
	case c.Subdomain != "":
		fmt.Fprintf(&b, "subdomain = %q\n", c.Subdomain)
	}
	return b.String()
}

// publicURL is the externally reachable URL for the proxy.
func (c config) publicURL() string {
	host := ""
	switch {
	case c.CustomDomain != "":
		host = c.CustomDomain
	case c.Subdomain != "" && c.ServerDomain != "":
		host = c.Subdomain + "." + c.ServerDomain
	default:
		return ""
	}
	if c.VhostPort != "" && c.VhostPort != "80" {
		return "http://" + host + ":" + c.VhostPort
	}
	return "http://" + host
}

type driver struct {
	bin string

	mu      sync.Mutex
	cmd     *exec.Cmd
	cfgPath string
	url     string
}

func (d *driver) Start(ctx context.Context, addr string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.cmd != nil {
		return d.url, nil
	}
	cfg := configFromEnv(tunnel.PortOf(addr))
	if cfg.Subdomain == "" && cfg.CustomDomain == "" {
		return "", fmt.Errorf("tunnel-frp: set FRP_SUBDOMAIN or FRP_CUSTOM_DOMAIN for the public hostname")
	}

	f, err := os.CreateTemp("", "frpc-*.toml")
	if err != nil {
		return "", err
	}
	if _, err := f.WriteString(cfg.render()); err != nil {
		f.Close()
		return "", err
	}
	f.Close()
	d.cfgPath = f.Name()

	cmd := exec.Command(d.bin, "-c", d.cfgPath)
	if err := cmd.Start(); err != nil {
		os.Remove(d.cfgPath)
		return "", fmt.Errorf("tunnel-frp: start %s: %w (is frpc installed?)", d.bin, err)
	}
	d.cmd = cmd
	d.url = cfg.publicURL()
	return d.url, nil
}

func (d *driver) Stop(context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.cfgPath != "" {
		os.Remove(d.cfgPath)
		d.cfgPath = ""
	}
	if d.cmd == nil || d.cmd.Process == nil {
		return nil
	}
	err := d.cmd.Process.Kill()
	_ = d.cmd.Wait()
	d.cmd = nil
	d.url = ""
	return err
}

func (d *driver) Status(context.Context) (tunnel.Status, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return tunnel.Status{Running: d.cmd != nil, URL: d.url}, nil
}
