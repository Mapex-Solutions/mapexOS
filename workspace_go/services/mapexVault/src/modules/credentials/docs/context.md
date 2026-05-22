# Bounded Context: Credentials

**Service:** mapexVault
**Module path:** `src/modules/credentials/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

The Credentials bounded context owns the full lifecycle of third-party authentication material stored in the platform vault. It persists secrets under envelope encryption (Master Key wraps a per-record DEK which encrypts the data payload), exposes CRUD over HTTP, and keeps tokens alive by scheduling per-credential refresh timers on NATS JetStream. It also tracks Connections — the binding between an external account (e.g. Instagram, TikTok) and a stored Credential — and runs a reconcile safety-net loop that reseeds missing refresh timers after NATS drift or restarts.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Credential | Encrypted secret record (manual, oauth2, or userAndPass) with optional ProviderConfig driving token lifecycle | Workflow "credential reference" (just an ID pointer) |
| Connection | Link between an external provider account and a stored Credential | NATS connection / HTTP connection |
| ProviderConfig | Unencrypted HTTP request templates (LoginConfig, RefreshConfig) used by the vault to acquire/renew tokens | Plugin manifest `auth` block (source of truth on the frontend) |
| Envelope Encryption | Master Key encrypts DEK; DEK encrypts data — two-layer scheme keyed per record | Raw AES-GCM on the whole document |
| Refresh Schedule | Per-credential NATS scheduled message fired at `tokenExpiresAt - 15min` on `VAULT-SCHEDULE` | Reconcile timer (single global loop on `VAULT-RECONCILER`) |
| Reconciler | Hourly safety-net loop that reseeds refresh timers missing from `VAULT-SCHEDULE` | The individual refresh consumer |
| pathKey | Org-scoped path breadcrumb copied from RequestContext for multi-tenant filtering | MongoDB `_id` or credential `pluginId` |

## Published Events (driven — outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| Vault credential action (created/updated/deleted/error) | `mapexos.vault.events` (stream `MAPEX-VAULT`) | `{credentialId, action, timestamp}` (inferred — local map, no published contract) | (inferred) any service subscribing to `MAPEX-VAULT`; none known in-tree |
| Refresh schedule (self-delivery) | `mapexos.vault.schedule.{credentialId}` → target `mapexos.vault.schedule.fired` on stream `VAULT-SCHEDULE` | `{credentialId, credentialType}` | Own refresh consumer (`mapexVault-credential-refresh`) |
| Reconcile schedule (self-delivery) | `mapexos.vault.reconcile.schedule` → target `mapexos.vault.reconcile.fired` on stream `VAULT-RECONCILER` | `{trigger: "scheduled"}` with fixed MsgId `vault-reconcile` | Own reconcile consumer (queue group `{service}-VAULT-RECONCILE-GROUP`) |

## Consumed Events (driving — inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| Refresh fired | `mapexos.vault.schedule.fired` (stream `VAULT-SCHEDULE`) | `{credentialId, credentialType}` | Self — published by `publishRefreshSchedule` / `bootstrapSeed` / reconciler |
| Reconcile fired | `mapexos.vault.reconcile.fired` (stream `VAULT-RECONCILER`) | `{trigger}` | Self — published by `scheduleNextReconcile` |

## Driving Ports (what can call this module)
- HTTP external API (JWT auth + permission middleware) under `/api/v1/credentials`: `POST /`, `GET /`, `GET /:credentialId`, `PATCH /:credentialId`, `DELETE /:credentialId`, `POST /:credentialId/test`.
- HTTP internal API (API-key auth) under `/internal/credentials`: `GET /:credentialId/decrypt` — used by service-to-service callers (e.g. workflow, triggers) to obtain plaintext credential data plus `__pluginId` / `__credentialDefId` hints.
- NATS refresh consumer on `mapexos.vault.schedule.fired` (durable `mapexVault-credential-refresh`, BatchSize 1, DLQ eventType `credential-refresh`).
- NATS reconcile consumer on `mapexos.vault.reconcile.fired` (queue group, DuplicateWindow 10s, DLQ eventType `vault-reconcile`).
- Lifecycle hook `OnMount` (common.Mountable) — runs bootstrap seed + first reconcile timer when the DI container finishes wiring.

## Driven Ports (what this module requires)
- `repositories.CredentialRepository` — MongoDB persistence for `Credential` entities (Create, FindById, FindByIdAndUpdate, DeleteById, FindWithFilters, FindActiveWithTokenExpiry).
- `repositories.ConnectionRepository` — MongoDB persistence for `Connection` entities (Create, FindWithFilters, UpsertByAccount).
- `*envelope.EnvelopeService` — Master Key / DEK envelope encryption primitives (encryptData, decryptData helpers).
- `natsModel.Publisher` (name `core`) — publishes vault domain events on `MAPEX-VAULT`.
- `natsModel.ScheduleManager` (name `core`) — `PublishScheduled`, `PurgeStreamSubject`, `HasPendingMessages` on `VAULT-SCHEDULE` / `VAULT-RECONCILER`.
- Config values `ctx_timeout`, `internal_api_key`, `vault_reconcile_interval`, `service_name`.

## Invariants and Business Rules
- Encrypted fields (`encryptedDEK`, `dekNonce`, `encryptedData`, `dataNonce`) MUST never be serialised to API responses — all tagged `json:"-"` and responses go through `toCredentialResponse`.
- Token lifecycle fields (`tokenExpiresAt`, `lastRefreshedAt`, `refreshError`) are computed by the vault after a successful token exchange and MUST NOT be accepted from user input.
- `ProviderConfig` is stored unencrypted because the refresh consumer must read URL/method/paths without decrypting secrets; it MUST contain only templates and response-extraction paths, never secret values (secrets live only in the encrypted data map, resolved via `{{credential.*}}` placeholders at request time).
- Refresh schedules fire at `tokenExpiresAt - 15min` (`RefreshBufferMinutes`); schedules whose target time is already in the past are dropped, except during `bootstrapSeed` which reschedules expired credentials to `now + 30s`.
- Credential `status != active` short-circuits refresh (ack without action); token-request failures transition the credential to `status = error` with `refreshError` populated and publish an `error` vault event.
- Refresh is idempotent: `publishRefreshSchedule` purges the subject before publishing; `DeleteCredentialById` purges the credential's schedule subject before removing the record.
- Reconcile interval (default 3600s) MUST stay greater than the stream `DuplicateWindow` (10s) so the self-republishing timer is not dedup'd.
- Only credentials with a `LoginConfig` can be logged in via `HandleLogin`; only credentials with a `RefreshConfig` (falling back to `LoginConfig`) can be refreshed by the scheduled consumer.

## Known Cross-Context Interactions
- **workflow / triggers / js-executor** (inferred): consume `/internal/credentials/:id/decrypt` to obtain plaintext credential data keyed by pluginId when executing plugin nodes.
- **mapexos (core: users / orgs)**: request context (`OrgContext`, `OrgContextData.PathKey`) drives multi-tenant scoping on every CRUD operation — credentials and connections are filtered by `orgId` + `pathKey`.
- **permissions/vault**: external routes gate on `CredentialCreate | CredentialRead | CredentialUpdate | CredentialDelete` — enforced by `permissionMw.RequirePermission` before handlers run.
- **Plugin Marketplace (frontend + mock)** (inferred): supplies the plugin manifest `auth` block that the UI translates into the `ProviderConfig.LoginConfig` / `RefreshConfig` payloads posted here.
- **NATS infrastructure**: owns two dedicated streams (`VAULT-SCHEDULE`, `VAULT-RECONCILER`) with file storage + `AllowMsgSchedules` — migrations to these streams are part of the NATS JetStream migration ticket (TKT-2026-0014).
