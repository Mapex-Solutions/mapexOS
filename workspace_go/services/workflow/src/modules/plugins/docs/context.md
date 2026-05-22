# Bounded Context: Plugins

**Service:** workflow
**Module path:** `src/modules/plugins/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

Owns the `PluginManifest` root aggregate — declarative, JSON-only integration descriptors that expose `NodeTypeManifest`s, credential schemas, `fetchOptions` data loaders, plugin-level defaults, and lifecycle hooks (before/after/destroy). Provides CRUD + a TieredCache-backed loader (`PluginLoaderPort`) for fast manifest lookup by the runtime and by the fetch_options module. Coordinates cross-pod cache invalidation via a NATS FANOUT so all workflow pods drop stale L0/L1 entries when a manifest changes. Multi-tenant visibility mirrors the `RouteGroup` pattern (`isTemplate` + `orgId` + `pathKey`).

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| PluginManifest | Root aggregate: `pluginId` + nodeTypes + credentials + fetchOptions + hooks + defaults | WorkflowDefinition (user-authored DAG that references manifest node types) |
| NodeTypeManifest | One node-type declaration (e.g., `telegram/message`) — properties, handles, operations, hooks | WorkflowNode (an instance of a node type placed on a DAG) |
| ActionDef | Unified action contract (http/mqtt/nats/script + output extractor) reused by operations, fetchOptions, credential test, hooks | Runtime execution (the action is declaratively described, not a running task) |
| CredentialDef | Declarative auth-method schema (fields + test action) defined by a plugin | Stored credential (plaintext decrypted on demand elsewhere) |
| FANOUT invalidate | Cross-pod cache drop message on `mapexos.fanout.workflow.plugin.invalidate` | Stream-based retry (this is a fanout, not a durable queue) |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| PluginInvalidate | `mapexos.fanout.workflow.plugin.invalidate` (stream `FANOUT`) | `interfaces/message.PluginInvalidatePayload` `{pluginId, action}` | all workflow pods (self-subscribe) |

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| PluginInvalidate (self) | `mapexos.fanout.workflow.plugin.invalidate` | `PluginInvalidatePayload` | this module on any CRUD (cross-pod broadcast) |

## Driving Ports (what can call this module)

- HTTP routes under `/api/v1/plugins` (JWT): list (paginated + multi-tenant filter), get-by-id, get-enabled, create, update, delete.
- Cross-module Go API:
  - `PluginServicePort` — CRUD + `GetEnabledPlugins` for editor boot.
  - `PluginLoaderPort.GetManifest(pluginId)` / `GetAllEnabled()` / `Invalidate(...)` / `InvalidateAll()`.

## Driven Ports (what this module requires)

- `PluginManifestRepository` (MongoDB) — source of truth.
- `PluginLoaderPort` (TieredCache L0 RAM + L1 Disk + Mongo fallback, key `plugin:{pluginId}` and `plugins:all:enabled`).
- `natsModel.Fanout` (`core` bus) — publish + subscribe to `FANOUT`.

## Invariants and Business Rules

- On Create/Update/Delete: local cache invalidated first, then FANOUT published so peer pods drop their L0/L1 copies.
- `FANOUT` subscription failures MUST NOT panic the pod — only warn; other pods still self-invalidate on their own writes (inferred from `SubscribeFanout` warn-only path).
- Multi-tenant visibility: `isTemplate=true` shared down the org hierarchy; `isTemplate=false` scoped by `orgId`/`pathKey`.
- `NodeTypeManifest.Operations[*].Script` is never accepted via API PATCH/POST — scripts come only from audited manifests (enforced upstream of this module).
- `GetEnabledPlugins` uses the shared `plugins:all:enabled` cache key with a shorter TTL (inferred from loader comment).

## Known Cross-Context Interactions

- Consumed by **runtime** when resolving plugin nodes: manifest → action type → dispatch target (Triggers Service or JS Workflow Executor).
- Consumed by **fetch_options** (design-time) to locate the correct `FetchOptionsDef` by `resourceKey`.
- Consumed by **definitions** validation (`plugin_validation` domain service) to classify a definition as `valid` vs `plugin_missing`.
- Consumed by the frontend editor for the node-type palette and the dynamic form renderer.
