# tunnel-frp — docs

**frp.** Self-hosted tunnel via an `frpc` client against your frps server.

## Install

```bash
togo install togo-framework/tunnel-frp
```

Registers on the [`tunnel`](https://github.com/togo-framework/tunnel) base; select it with **tunnel.provider in togo.yaml (or TUNNEL_DRIVER)**, then use **`togo tunnel`**.

## Interface

`Tunnel` — `Start(ctx, addr) -> publicURL`, `Stop`, `Status`.

## Configuration

| Env var | Description |
|---|---|
| `FRP_SERVER_ADDR` | frps server address `host:port` (required). |
| `FRP_TOKEN` | frps auth token (required if the server sets one). |
| `FRP_SUBDOMAIN` | Subdomain to expose on the frps server. Optional. |
| `FRP_CUSTOM_DOMAIN` | Custom domain to bind on the frps server. Optional. |

## Usage & notes

Renders an `frpc` config and runs it against your `FRP_SERVER_ADDR`, exposing the local port via a subdomain or custom-domain vhost.

## Example

```bash
togo tunnel:start --provider frp
```

## Links

- [frp](https://github.com/fatedier/frp)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/tunnel-frp)
