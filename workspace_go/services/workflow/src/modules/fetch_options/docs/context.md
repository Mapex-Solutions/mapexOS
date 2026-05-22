# Bounded Context: FetchOptions

**Service:** workflow
**Module path:** `src/modules/fetch_options/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose

Backend proxy that powers design-time dynamic dropdowns in the workflow editor. When a plugin's `NodePropertyDef` uses source type `fetchOptions` (e.g., "pick a Telegram chat"), the frontend calls this endpoint with `credentialId`, `pluginId`, `resourceKey`, and any `dependsOn` values. The service decrypts the credential via Vault, loads the plugin manifest's `FetchOptionsDef`, resolves templates, performs the outbound HTTP call to the provider, and extracts/transforms the response into `[{label, value}]` items. Pagination and search support are driven by the manifest.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| fetchOptions | Design-time options loader declared in a plugin manifest (used by dropdowns) | Runtime plugin operation (execution-time action) |
| resourceKey | Lookup key into `manifest.fetchOptions[key]` chosen by the frontend based on form state | Node operation id |
| dependsOn | Form-field values that parametrise the fetch (e.g., selected workspace id) | Node input edges |
| Action | Unified `ActionDef` contract (type + http/mqtt/nats/script block + output extractor) — reused at design-time | Runtime plugin action (same contract, different caller) |
| Credential | Vault-stored plaintext map decrypted on demand per request | Credential manifest definition (schema only) |

## Published Events (driven — outbound)

_None._

## Consumed Events (driving — inbound)

_None._

## Driving Ports (what can call this module)

- HTTP route under `/api/v1/load_options` (JWT auth, `CredentialRead` permission):
  - `POST /` — `FetchOptions(credentialId, pluginId, resourceKey, dependsOn)` returns `[]FetchOptionsItem`.

## Driven Ports (what this module requires)

- Vault client (credential decryption) (inferred — called via service dependency, same pattern as runtime's `VaultPort`).
- Plugin manifest loader (reads `manifest.fetchOptions[resourceKey]` from the `plugins` module / its TieredCache).
- Outbound HTTP client — executes the manifest-declared `ActionDef.http` request with template replacement.

## Invariants and Business Rules

- Never persists credentials; plaintext is resolved per request and discarded (inferred — no repositories in this module).
- Only `http` action type is supported at design-time; `script`/`mqtt`/`nats` in fetchOptions are out of scope here (inferred from frontend usage).
- Template placeholders (`{{credential.*}}`, `{{dependsOn.*}}`) are resolved before the outbound call.
- Response shape: post-extraction list items MUST be `{label, value}` — `valuePath`/`labelPath` in the manifest drive extraction.
- Pagination mode (`cursor` or `page`) and search (`search.param`, `search.minLength`) are honored only if declared in the manifest.

## Known Cross-Context Interactions

- Reads manifests from the **plugins** module (via its loader / TieredCache).
- Decrypts credentials via the **mapexVault** microservice (HTTP) — same dependency pattern used by runtime.
- Called exclusively by the frontend editor (design-time) — no service-to-service callers.
