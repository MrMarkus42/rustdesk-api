[Deutsch](README.md) · **English**

# RustDesk API

API for self-hosted [RustDesk](https://rustdesk.com) servers, written in Go, with an admin interface. It implements the endpoints the RustDesk client uses for login, address books and device management.

This fork is based on [lejianwen/rustdesk-api](https://github.com/lejianwen/rustdesk-api) (v2.7), which is no longer maintained. It ships together with [MrMarkus42/rustdesk-server](https://github.com/MrMarkus42/rustdesk-server) as an s6 image.

## Features

- **Client API:** login (password, OIDC, GitHub, Google, LDAP), address books with tags and sharing, groups, online status, connection and file logs
- **Admin interface** (`/_admin/`) in **German** and English (also French, Spanish, Russian, Korean, Chinese): users, devices, address books, groups, OAuth, login log, server commands

### Changes compared to lejianwen

- **Security:** captcha and IP ban can no longer be bypassed with a forged `X-Forwarded-For`, random access tokens, no directory listing, updated dependencies (govulncheck and npm audit clean)
- **Login:** clients are no longer logged out after 7 days when `token-expire` is set higher
- **Admin interface:** German translation, English instead of Chinese as default, forms can clear fields again, OIDC with Amazon Cognito
- **Supply chain:** the admin interface is built from a reviewed, pinned commit of [MrMarkus42/rustdesk-api-web](https://github.com/MrMarkus42/rustdesk-api-web)

## Docker image

```
ghcr.io/mrmarkus42/rustdesk-api:<version>
```

To run it together with hbbs/hbbr, use the combined image `ghcr.io/mrmarkus42/rustdesk-server-s6` (see [rustdesk-server](https://github.com/MrMarkus42/rustdesk-server)).

## Important settings

Every setting in `conf/config.yaml` can be set via environment variable: prefix `RUSTDESK_API_`, section and key joined with `_`, dashes as `_` (example: `gin.trust-proxy` → `RUSTDESK_API_GIN_TRUST_PROXY`).

| Variable | Meaning |
|---|---|
| `RUSTDESK_API_RUSTDESK_ID_SERVER` | ID server handed to clients |
| `RUSTDESK_API_RUSTDESK_RELAY_SERVER` | Relay server |
| `RUSTDESK_API_RUSTDESK_API_SERVER` | Public address of the API (e.g. `https://rustdesk.example.org`) |
| `RUSTDESK_API_RUSTDESK_KEY_FILE` | Public server key (default `/data/id_ed25519.pub`) |
| `RUSTDESK_API_JWT_KEY` | Signs the login tokens, required for `MUST_LOGIN`. **Do not change it later** |
| `RUSTDESK_API_APP_TOKEN_EXPIRE` | Login validity (default `168h`) |
| `RUSTDESK_API_GIN_TRUST_PROXY` | Address/range of the reverse proxy (comma-separated, CIDR allowed) |
| `RUSTDESK_API_APP_BAN_THRESHOLD` | IP ban after n failed attempts (`0` = off). Behind a proxy only with `TRUST_PROXY` |
| `RUSTDESK_API_APP_REGISTER` | Self-registration (`false` recommended) |
| `RUSTDESK_API_LANG` | Language of API messages when the browser sends none (default `en`) |

Reset the admin password: `./apimain reset-admin-pwd <new-password>`

## Build

```bash
go build -o apimain ./cmd/apimain.go
```

The admin interface is built separately from [rustdesk-api-web](https://github.com/MrMarkus42/rustdesk-api-web) (`npm ci && npm run build`) and copied to `resources/admin/`. The workflows in `.github/workflows/` do this automatically.

## Origin and license

- Original project: [lejianwen/rustdesk-api](https://github.com/lejianwen/rustdesk-api)
- License: [MIT](LICENSE)
