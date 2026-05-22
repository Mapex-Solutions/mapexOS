# Observability

## Health
`GET /health`

### Health status semantics
- `healthy` (HTTP 200): all critical checks OK
- `degraded` (HTTP 200): optional checks failed
- `unhealthy` (HTTP 503): at least one critical check failed

### Checks for this service
| Dependency | Criticality |
|---|---|
| MongoDB | critical |
| Redis (app) | critical |
| Redis (shared) | critical |
| NATS (core) | critical |

### Response shape
```json
{
  "status": "healthy",
  "service": "http_gateway",
  "version": "1.0.0",
  "uptime": "48h32m15s",
  "timestamp": "2026-02-18T10:30:00Z",
  "lastCheckAt": "2026-02-18T10:29:55Z",
  "checks": {
    "mongodb": { "connected": true, "latencyMs": 2 },
    "redis:app": { "connected": true, "latencyMs": 1 },
    "redis:shared": { "connected": true, "latencyMs": 1 },
    "nats:core": { "connected": true, "latencyMs": 3 }
  }
}
```

## Metrics
All service‑specific metrics are prefixed with `httpgw_`.

### Summary
| Metric | Type | Labels | Description |
|---|---|---|---|
| `httpgw_event_auth_total` | Counter | `auth_type`, `result` | Auth attempts by strategy and result (success/failure). |
| `httpgw_event_auth_duration_seconds` | Histogram | `auth_type` | Auth latency by strategy (captures JWKS latency). |
| `httpgw_event_auth_failures_total` | Counter | `auth_type` | Auth failures that triggered `events.raw` security events. |
| `httpgw_event_processed_total` | Counter | `status` | Events that passed auth and completed processing. |
| `httpgw_event_published_total` | Counter | `subject`, `status` | NATS publish attempts per subject (detects publish failures). |
| `httpgw_event_processing_duration_seconds` | Histogram | — | Processing duration excluding auth (isolates NATS publish latency). |
| `httpgw_event_payload_size_bytes` | Histogram | — | Incoming webhook body size distribution. |
| `httpgw_ds_operations_total` | Counter | `operation`, `status` | Data Source CRUD operations by operation and status. |
| `httpgw_ds_operation_duration_seconds` | Histogram | `operation` | Data Source CRUD latency (Mongo query performance). |
| `httpgw_ds_list_results_count` | Histogram | — | Items returned per Data Source list query. |
| `httpgw_ds_cache_total` | Counter | `result` | Data source cache lookups by result (hit/miss). |

### Labels
- `auth_type`: `oauth2`, `jwt`, `apiKey`, `ip_whitelist`, `none`
- `result`: `success`, `failure` (auth metrics) / `hit`, `miss` (cache metric)
- `status`: `success`, `error`
- `operation`: `list`, `create`, `read`, `update`, `delete`
- `subject`: `processor.js.execute`

### Go runtime + process collectors
Controlled by:
- `METRICS_GO_COLLECTOR` (default: `true`)
- `METRICS_PROCESS_COLLECTOR` (default: `true`)

### Full instrumentation map
See `documentations/architecture/microservices/services/http_gateway/metrics.md` for placement details, rationale, and cardinality.

## Logs
Structured logs via the shared microservices logger. Use `LOG_LEVEL` to override environment defaults.
