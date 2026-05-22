# Observability

## Health

`GET /health`

### Health status semantics

- `healthy` (HTTP 200): all critical checks OK
- `degraded` (HTTP 200): optional checks failed
- `unhealthy` (HTTP 503): at least one critical check failed

### Checks for this service

| Dependency / Flag | Criticality | Notes |
|---|---|---|
| MongoDB | critical | |
| Redis (app) | critical | |
| Redis (shared) | critical | |
| NATS (core + JetStream) | critical | |
| MinIO | optional | L2 read-model storage |
| mapexVault | critical for cert issuance | Drives `caReady` — see below |
| `caReady` | flag | `true` once `mqttcerts.OnMount` has loaded the intermediate CA from mapexVault into RAM. `mqttcerts` HTTP routes return `503 ca_not_ready` while `caReady=false`. Does NOT mark the service unhealthy — the rest of the API (assets, templates, heartbeat) keeps serving. |

### Response shape

```json
{
  "status": "healthy",
  "service": "assets",
  "version": "1.0.0",
  "uptime": "48h32m15s",
  "timestamp": "2026-05-11T10:30:00Z",
  "lastCheckAt": "2026-05-11T10:29:55Z",
  "checks": {
    "mongodb": { "connected": true, "latencyMs": 2 },
    "redis:app": { "connected": true, "latencyMs": 1 },
    "redis:shared": { "connected": true, "latencyMs": 1 },
    "nats:core": { "connected": true, "latencyMs": 3 },
    "minio": { "connected": true, "latencyMs": 4, "critical": false },
    "mapexVault": { "connected": true, "latencyMs": 6 },
    "caReady": true
  }
}
```

## Metrics

`GET /metrics`

All service-specific metrics are prefixed with `assets_`. The endpoint is registered before global middlewares (no authentication required for Prometheus scraping).

### Summary

| Metric | Type | Labels | Description |
|---|---|---|---|
| `assets_asset_operations_total` | Counter | `operation`, `status` | Asset CRUD operations by type and result. |
| `assets_asset_operation_duration_seconds` | Histogram | `operation` | Asset CRUD operation latency. |
| `assets_asset_list_results_count` | Histogram | — | Items returned per asset list query. |
| `assets_asset_cache_total` | Counter | `result` | Asset cache lookups by result (hit/miss). |
| `assets_template_operations_total` | Counter | `operation`, `status` | Template CRUD operations by type and result. |
| `assets_template_operation_duration_seconds` | Histogram | `operation` | Template CRUD operation latency. |
| `assets_template_list_results_count` | Histogram | — | Items returned per template list query. |
| `assets_template_cache_total` | Counter | `result` | Template cache lookups by result (hit/miss). |
| `assets_mqttcerts_operations_total` | Counter | `operation`, `status` | MQTT cert issue/revoke/list outcomes. |
| `assets_mqttcerts_operation_duration_seconds` | Histogram | `operation` | Cert operation latency (signing + Mongo + fanout). |
| `assets_mqttcerts_ca_ready` | Gauge | — | `1` once the intermediate CA is mounted from mapexVault, `0` otherwise (also reflected in `/health`). |
| `assets_mqttcerts_ca_mount_attempts_total` | Counter | `result` | `OnMount` attempts by result (`success` / `error`); useful to alert on retry storms. |
| `assets_healthmonitor_heartbeats_total` | Counter | `origin`, `result` | Heartbeats received (`origin=implicit\|explicit`, `result=accepted\|dropped`). |
| `assets_healthmonitor_presence_total` | Counter | `event`, `result` | Broker presence advisories processed (`event=connect\|disconnect`). |
| `assets_healthmonitor_state_transitions_total` | Counter | `from`, `to` | Asset state flips (e.g. `online`→`offline`). |

### Labels

- `operation`: `list`, `create`, `read`, `update`, `delete` (CRUD) or `issue`, `revoke`, `list` (mqttcerts)
- `status`: `success`, `error`
- `result` (cache): `hit`, `miss`
- `result` (mqttcerts mount): `success`, `error`
- `origin` (heartbeats): `implicit`, `explicit`
- `event` (presence): `connect`, `disconnect`

### Metric groups

#### Asset CRUD (subsystem: `asset`)

Instruments the `assets` module — CRUD operations against the `assets` MongoDB collection and the TieredCache lookup that precedes reads.

| Metric | Buckets |
|---|---|
| `assets_asset_operations_total` | — |
| `assets_asset_operation_duration_seconds` | Default Prometheus buckets (0.005–10s) |
| `assets_asset_list_results_count` | 0, 1, 5, 10, 25, 50, 100, 250 |
| `assets_asset_cache_total` | — |

#### Asset Template CRUD (subsystem: `template`)

Instruments the `assettemplates` module — CRUD operations on the `assettemplates` collection and the TieredCache lookup for template read models.

| Metric | Buckets |
|---|---|
| `assets_template_operations_total` | — |
| `assets_template_operation_duration_seconds` | Default Prometheus buckets (0.005–10s) |
| `assets_template_list_results_count` | 0, 1, 5, 10, 25, 50, 100, 250 |
| `assets_template_cache_total` | — |

#### MQTT Certificates (subsystem: `mqttcerts`)

Instruments device cert issuance and revocation. Issuance latency is dominated by local ECDSA P-256 signing + Mongo write + fanout publish; the intermediate CA is read from RAM (no per-request mapexVault call).

| Metric | Buckets |
|---|---|
| `assets_mqttcerts_operations_total` | — |
| `assets_mqttcerts_operation_duration_seconds` | Default Prometheus buckets (0.005–10s) |
| `assets_mqttcerts_ca_ready` | — |
| `assets_mqttcerts_ca_mount_attempts_total` | — |

#### Health Monitor (subsystem: `healthmonitor`)

Instruments heartbeat ingestion (NATS + HTTP), broker presence consumers, and scheduled offline scans.

| Metric | Buckets |
|---|---|
| `assets_healthmonitor_heartbeats_total` | — |
| `assets_healthmonitor_presence_total` | — |
| `assets_healthmonitor_state_transitions_total` | — |

### Go runtime + process collectors

Controlled by environment variables:

- `METRICS_GO_COLLECTOR` (default: `true`) — enables `go_goroutines`, `go_gc_duration_seconds`, `go_memstats_*`, etc.
- `METRICS_PROCESS_COLLECTOR` (default: `true`) — enables `process_cpu_seconds_total`, `process_resident_memory_bytes`, `process_open_fds`, etc.

### Full instrumentation map

See `documentations/architecture/microservices/services/assets/metrics.md` for placement details and cardinality analysis.

## Logs

Structured logs via the shared microservices logger. Use `LOG_LEVEL` to override environment defaults.

| Level | When |
|---|---|
| `debug` | Used in `development` environment (default) |
| `info` | Default in non-development environments |
| `warn` | Degraded states (cache miss, optional dep unavailable, `mqttcerts.OnMount` retry) |
| `error` | Hard failures (MongoDB write error, NATS publish failure, mapexVault unreachable on first mount) |

Log lines follow the `[LAYER:Component]` convention — examples: `[SERVICE:Assets]`, `[SERVICE:MqttCerts]`, `[CONSUMER:MqttPresenceConnect]`, `[CONSUMER:MqttPresenceDisconnect]`, `[CONSUMER:AssetHeartbeat]`.
