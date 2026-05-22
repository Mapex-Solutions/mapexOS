# Marketplace - Bounded Context

| Metadata | Value |
|---|---|
| Service | `plugin-marketplace-mock` |
| Owner | @thiagoanselmo |
| Last reviewed | 2026-04-21 |

## Purpose

Dev-only mock of the plugin marketplace CDN. The module exposes an Express server that serves static plugin manifests, registry catalogs, event schemas, and icons from the `public/` directory. It stands in for the real hosted marketplace during local development and CI, letting the workflow and triggers services fetch plugin metadata over HTTP without any external dependency. No business logic, no persistence, no auth; this module is intentionally thin and must not be deployed to production.

## Ubiquitous Language

- **Plugin manifest**: JSON document describing a plugin's node types, operations, and credentials (schema `mapex-plugin/v1`).
- **Registry**: Top-level catalog at `/plugins/registry.json` listing available plugins (schema `mapex-plugin-registry/v1`).
- **System plugin**: Built-in plugin shipped with the platform (`isSystem: true`), non-removable.
- **Public directory**: Filesystem root (`public/`) whose contents are served verbatim over HTTP.

## Published Events

None. This module is purely request/response over HTTP and does not publish domain or integration events.

## Consumed Events

None. The module does not subscribe to NATS or any messaging transport.

## Driving Ports

- `MarketplaceServicePort` — exposes configuration to the HTTP layer: `getPublicDirectory()` returns the absolute path to the static root, `getPort()` returns the TCP port the server binds to (default `3099`).
- HTTP interface (`interfaces/http/routes/marketplace_routes.ts`) registers a single static handler that serves `GET /plugins/registry.json`, `GET /plugins/{pluginId}/manifest.json`, `GET /plugins/{pluginId}/events.json`, and icon assets.

## Driven Ports

None. The service has no outbound dependencies; it reads files from the local filesystem via `express.static`, which is considered infrastructure provided by the framework rather than a domain-level driven port.

## Invariants

- The public directory path resolved at bootstrap must exist and be readable; otherwise `express.static` silently serves 404s.
- CORS is enabled globally so any origin can fetch manifests during dev.
- The port is fixed at `3099` to match client service configuration; changing it requires updating consumers.
- Manifest files on disk must conform to the `mapex-plugin/v1` DSL (see `docs/plugins/DSL_JSON.md`); the module does no schema validation.

## Cross-Context Interactions

- **workflow** service and **triggers** service fetch plugin manifests and the registry over HTTP at startup / on-demand to discover available node types and credential schemas.
- **Frontend (mapexOS app)** consumes the registry and icons when rendering the plugin picker in the workflow editor.
- In production these consumers will point at the real marketplace CDN; this mock is a drop-in replacement keyed by environment configuration.
