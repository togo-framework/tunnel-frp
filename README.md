<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/tunnel-frp</h1>
  <p>frp (Fast Reverse Proxy) driver for togo tunnel — expose a local port via your own frps server.</p>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/tunnel-frp"><img src="https://pkg.go.dev/badge/github.com/togo-framework/tunnel-frp.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/tunnel-frp
```
<!-- /togo-header -->

**frp** ([Fast Reverse Proxy](https://github.com/fatedier/frp)) driver for togo's
[`tunnel`](https://github.com/togo-framework/tunnel) subsystem. Renders an `frpc`
TOML config for an HTTP proxy to your local port and runs `frpc` against the frps
server you control — self-hosted tunnels, no third party.

Requires the `frpc` binary and an `frps` server.

## Config

| Env | Meaning |
|-----|---------|
| `TUNNEL_DRIVER` | set to `frp` |
| `FRP_SERVER_ADDR` | frps server address (required) |
| `FRP_SERVER_PORT` | frps bind port (default `7000`) |
| `FRP_TOKEN` | frps auth token (optional) |
| `FRP_SUBDOMAIN` | HTTP vhost subdomain (one of subdomain/custom-domain) |
| `FRP_CUSTOM_DOMAIN` | full custom domain for the vhost |
| `FRP_SERVER_DOMAIN` | base domain for the subdomain URL (default: server addr) |
| `FRP_VHOST_PORT` | public HTTP vhost port (default `80`) |
| `FRPC_BIN` | path to `frpc` (default: `frpc` on PATH) |

```go
svc, _ := tunnel.FromKernel(k)
url, _ := svc.Start(ctx, "8080")   // → http://myapp.frp.example.com
defer svc.Stop(ctx)
```

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- togo-sponsors -->
