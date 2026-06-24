package frp

import (
	"strings"
	"testing"

	"github.com/togo-framework/tunnel"
)

func TestRenderSubdomain(t *testing.T) {
	c := config{ServerAddr: "frp.example.com", ServerPort: "7000", Token: "secret",
		ProxyName: "togo-8080", LocalPort: "8080", Subdomain: "myapp"}
	out := c.render()
	for _, want := range []string{
		`serverAddr = "frp.example.com"`,
		`serverPort = 7000`,
		`auth.token = "secret"`,
		`type = "http"`,
		`localPort = 8080`,
		`subdomain = "myapp"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("render() missing %q\n---\n%s", want, out)
		}
	}
}

func TestRenderCustomDomainNoToken(t *testing.T) {
	c := config{ServerAddr: "x", ServerPort: "7000", ProxyName: "p", LocalPort: "3000", CustomDomain: "app.example.com"}
	out := c.render()
	if strings.Contains(out, "auth.token") {
		t.Error("should not emit auth.token without a token")
	}
	if !strings.Contains(out, `customDomains = ["app.example.com"]`) {
		t.Errorf("missing customDomains\n%s", out)
	}
}

func TestPublicURL(t *testing.T) {
	cases := []struct {
		c    config
		want string
	}{
		{config{Subdomain: "myapp", ServerDomain: "frp.example.com", VhostPort: "80"}, "http://myapp.frp.example.com"},
		{config{Subdomain: "myapp", ServerDomain: "frp.example.com", VhostPort: "8080"}, "http://myapp.frp.example.com:8080"},
		{config{CustomDomain: "app.example.com", VhostPort: "80"}, "http://app.example.com"},
		{config{}, ""},
	}
	for _, tc := range cases {
		if got := tc.c.publicURL(); got != tc.want {
			t.Errorf("publicURL(%+v) = %q, want %q", tc.c, got, tc.want)
		}
	}
}

func TestDriverRegistered(t *testing.T) {
	found := false
	for _, n := range tunnel.Drivers() {
		if n == "frp" {
			found = true
		}
	}
	if !found {
		t.Fatal("frp driver not registered on tunnel base")
	}
}
