---
name: tunnel-frp
description: Expose a local togo app via a self-hosted frp server — set TUNNEL_DRIVER=frp + FRP_SERVER_ADDR and call tunnel.Start
---

# togo tunnel-frp

frp ([fast reverse proxy](https://github.com/fatedier/frp)) driver for the togo
`tunnel` subsystem — a self-hosted alternative to ngrok/Cloudflare.

## Setup

```bash
togo install togo-framework/tunnel
togo install togo-framework/tunnel-frp
```

1. Install the `frpc` binary and run an `frps` server you control.
2. `.env`:
   ```bash
   TUNNEL_DRIVER=frp
   FRP_SERVER_ADDR=frp.example.com:7000
   FRP_TOKEN=shared-secret
   FRP_TYPE=http
   FRP_SUBDOMAIN=myapp          # → http://myapp.frp.example.com
   ```
   For TCP: `FRP_TYPE=tcp` + `FRP_REMOTE_PORT=6000`.

## Use

```go
import (
	_ "github.com/togo-framework/tunnel"
	_ "github.com/togo-framework/tunnel-frp"
	"github.com/togo-framework/tunnel"
)

if tn, ok := tunnel.FromKernel(k); ok {
	url, _ := tn.Start(ctx, "8080")
	defer tn.Stop(ctx)
}
```

## Notes
- The driver writes a temporary `frpc` TOML config and runs `frpc -c`. Stop
  removes the config + kills frpc.
- HTTP vhost needs `FRP_SUBDOMAIN` or `FRP_CUSTOM_DOMAIN`; TCP needs
  `FRP_REMOTE_PORT`. Or set `FRP_PUBLIC_URL` to report a URL directly.
