# Plugin Marketplace Mock Service

Lightweight Express server that serves static plugin manifests, event schemas, and icons. Simulates a plugin marketplace CDN for development and testing.

## Service Info

| Field | Value |
|---|---|
| Port | `3099` |
| Entry Point | `src/main.ts` |
| Build | N/A (dev-only service) |
| Runtime | Node.js, Express, tsx |

## Structure

| Path | Purpose |
|---|---|
| `src/main.ts` | Express server with CORS, serves `public/` as static files |
| `public/plugins/registry.json` | Plugin catalog (schema: `mapex-plugin-registry/v1`) |
| `public/plugins/telegram/manifest.json` | Telegram plugin manifest with node types and operations |
| `public/plugins/telegram/events.json` | 22 trigger events in 6 categories |
| `public/plugins/telegram/icon.svg` | Telegram brand icon |

## Plugins

| Plugin | ID | Category | Nodes | Events | Credentials |
|---|---|---|---|---|---|
| Telegram | `telegram` | Messaging | 2 (message, chat) | 22 | Bot Token |

## Key Decisions

- **Static serving only**: No business logic, pure CDN mock via `express.static`.
- **Plugin DSL**: Manifests follow `mapex-plugin/v1` schema. See `docs/plugins/DSL_JSON.md` for spec.
- **System plugin**: Telegram is `isSystem: true` (built-in, non-removable).
- **Dev-only**: Not deployed to production. Real marketplace will be a CDN.
