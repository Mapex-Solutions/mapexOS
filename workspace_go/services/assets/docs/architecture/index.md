# Architecture

## Design
DDD + Hexagonal modular architecture. Each module owns one bounded context and follows the canonical 4-layer tree (`application/`, `domain/`, `infrastructure/`, `interfaces/`) plus `module.go` and `docs/context.md`.

## Project Structure
```
src/
├── modules/
│   ├── assets/             # Asset CRUD + L2 read model + fanout invalidation
│   ├── assettemplates/     # Template CRUD + EVA + scripts + list name sync
│   ├── healthmonitor/      # Heartbeat ingest + presence consumer + scan-based offline detection
│   ├── mqttcerts/          # MQTT device certificate lifecycle (issue / revoke / list)
│   └── app/                # App bootstrap + cross-module wiring
└── shared/
    └── configuration/      # Service configuration
```

## Module Responsibilities

| Module | Responsibility |
|---|---|
| `assets` | Manage asset entities; write a single `AssetReadModel` to L2 (MinIO) carrying everything every consumer needs (including `Protocol.Mqtt.PasswordHash` and `CurrentCert.Serial` for the broker plugin); fanout invalidation on CUD; expose `GET /internal/assets/:uuid` as the L3 read-model fallback. |
| `assettemplates` | Manage templates (EVA mapping, scripts); L2 read-models; fanout invalidation; subscribe to `mapexos.lists.name_updated` for denormalized name sync |
| `healthmonitor` | Ingest heartbeats (NATS + HTTP), consume MQTT presence advisories from the broker plugin, schedule periodic offline scans, publish alerts |
| `mqttcerts` | Issue device certs locally using the RAM-cached intermediate CA (loaded from mapexVault via `OnMount`); revoke certs; HARD delete on asset removal; 30-day TTL audit |

There is NO MQTT auth-callout module / service / route. The mapex-mqtt-broker plugin reads `AssetReadModel` from its TieredCache (L1 Pebble → L2 MinIO → L3 HTTP via `/internal/assets/:uuid`) and decides every CONNECT locally — bcrypt-compare for password mode, serial-equality for cert mode. The platform's only HTTP path between broker and this service is that L3 read-model fallback, gated by the standard `apikeymw.ApiKeyAuthMiddleware` (`X-API-Key` header).

## Main Event Flows (NATS)

### Fanout invalidation (publish)
- `${env}.mapexos.fanout.asset.invalidate` — published by `assets` on asset CUD and by `mqttcerts` on issue/revoke (cert state change rides this subject). Consumed by Router, JS-Executor, Events, and the mapex-mqtt-broker plugin.
- `${env}.mapexos.fanout.template.invalidate` — published by `assettemplates` on template CUD. Consumed by Router, Events, JS-Executor.

### Heartbeat (subscribe — `healthmonitor`)
- Stream: `MAPEXOS-ASSETS-HEARTBEAT` (env-scoped)
- Subject: `${env}.mapexos.asset.heartbeat.>` (trailing token is `orgId`)
- Producers: `js-executor` (implicit heartbeat on every data event when `heartbeatMode='implicit'`) and `http_gateway` (explicit `POST /api/v1/heartbeat?ds={dsId}`)
- MQTT-protocol assets do NOT publish here — they use the presence path below.

### MQTT presence (subscribe — `healthmonitor`)
- Stream: `MAPEXOS-ASSETS-MQTT-PRESENCE` (env-scoped, CORE)
- Subject: `${env}.mapexos.mqtt.presence.advisory` (single subject; `Event` field = `connect` | `disconnect`)
- Producer: mapex-mqtt-broker plugin on every device CONNECT/DISCONNECT
- Two consumers with independent durables (`assets-mqtt-presence`, `assets-mqtt-presence-connect`) gate by `Event` in their handlers.
- Payload contract: `packages/contracts/services/assets/healthmonitor/presence.go::PresenceAdvisory`

### Scan scheduling (subscribe — `healthmonitor`)
- Stream: `MAPEXOS-ASSETS-HEALTH-MONITOR` (WorkQueue, AllowMsgSchedules)
- Schedule subject: `${env}.mapexos.healthmonitor.scan.schedule`
- Fire subject: `${env}.mapexos.healthmonitor.scan` — fires the next offline scan; service republishes the next schedule with `MsgId=hm-scan` so only one pending schedule exists across pods.

### List name sync (subscribe — `assettemplates`)
- `${env}.mapexos.lists.name_updated` — updates denormalized manufacturer/model/category names inside `assettemplates` documents.

## Read Model Cache (MinIO/S3)
- MongoDB stores the authoritative asset and template data.
- One `AssetReadModel` per asset is cached in MinIO/S3 (bucket `mapex-assets`, key `{orgId}/{assetUUID}.json`). It carries every field every consumer needs — public asset state PLUS `Protocol.Mqtt.PasswordHash` and `CurrentCert.Serial` so the broker plugin decides MQTT CONNECTs locally without an auth callout.
- Downstream services (Router, JS-Executor, Events, mapex-mqtt-broker plugin) fetch this single payload via their own TieredCache (L1/L2/L3) and project out the fields they need.

## MQTT PKI integration
- `mqttcerts` registers a `common.Mountable.OnMount` that pulls the intermediate CA from mapexVault `pki/` at service startup. Failure triggers an exponential-backoff retry goroutine; routes guarded by `RequireCAReady` return `503` until the CA lands in RAM.
- Device certs are signed locally with `crypto/x509` (ECDSA P-256, 128-bit random serial). The private key is generated on the device side or returned once at issuance — never persisted on the server.
- Issued cert metadata (serial, fingerprint, expiry, subject CN) lands on `Asset.CurrentCert` in Mongo and flows into the `AssetReadModel.CurrentCert` field that the broker plugin reads. The CONNECT decision for cert mode is `entry.CurrentCert.Serial == device.cert.Serial` — done LOCALLY by the broker plugin.
- Revocation writes a 30-day TTL row to `mqttRevokedCertificates` (Mongo TTL on `revokedAt`) and fans out invalidation so the broker plugin drops the cached entry and re-fetches the (now empty / replaced) `CurrentCert` on the next CONNECT.
- HARD delete on asset removal drops the cert state along with the asset entity (LGPD-friendly); no cascade revoke.

## Fallback & Invalidation
- **Fallback endpoints** repopulate L2 on cache miss (`/internal/assets/:uuid`, `/internal/templates/:id`).
- **Fanout invalidation** notifies consumers to invalidate L0/L1 when assets/templates/certs change.
