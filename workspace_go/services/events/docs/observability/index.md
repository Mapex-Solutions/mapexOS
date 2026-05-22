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
| ClickHouse | optional |
| MinIO | optional |

### Response shape
```json
{
  "status": "healthy",
  "service": "events",
  "version": "1.0.0",
  "uptime": "48h32m15s",
  "timestamp": "2026-02-18T10:30:00Z",
  "lastCheckAt": "2026-02-18T10:29:55Z",
  "checks": {
    "mongodb": { "connected": true, "latencyMs": 2 },
    "redis:app": { "connected": true, "latencyMs": 1 },
    "redis:shared": { "connected": true, "latencyMs": 1 },
    "nats:core": { "connected": true, "latencyMs": 3 },
    "clickhouse": { "connected": true, "latencyMs": 5, "critical": false },
    "minio": { "connected": true, "latencyMs": 4, "critical": false }
  }
}
```

## Metrics
All service‑specific metrics are prefixed with `events_`.

### Summary
| Metric | Type | Description |
|---|---|---|
| `events_event_processed_total` | Counter | Events processed by consumer and status. |
| `events_event_processing_duration_seconds` | Histogram | End‑to‑end batch processing time per consumer. |
| `events_event_batch_size` | Histogram | Batch size per consumer. |
| `events_message_total` | Counter | Message outcomes (ack/nack/reject/dlq) by consumer. |
| `events_clickhouse_insert_duration_seconds` | Histogram | ClickHouse bulk insert latency by table. |
| `events_clickhouse_insert_total` | Counter | ClickHouse insert attempts by table and status. |
| `events_clickhouse_insert_batch_size` | Histogram | Bulk insert batch size by table. |
| `events_template_cache_total` | Counter | Template cache lookups by result. |
| `events_eva_fields_mapped_total` | Counter | EVA fields resolved from templates. |
| `events_retention_cache_total` | Counter | Retention cache lookups by result. |

### Labels
- `consumer`: `raw`, `jsexec`, `dlq`, `router`, `businessrule`, `trigger`, `store`
- `status`: `success`, `error`
- `result`: `ack`, `nack`, `reject`, `dlq`
- `table`: `eventsRaw`, `eventsJsExecutor`, `eventsDLQ`, `eventsRouter`, `eventsBusinessRule`, `eventsTrigger`, `events`

### Go runtime + process collectors
Controlled by:
- `METRICS_GO_COLLECTOR` (default: `true`)
- `METRICS_PROCESS_COLLECTOR` (default: `true`)

### Full instrumentation map
See `documentations/architecture/microservices/services/events/metrics.md` for placement details and cardinality.

## Logs
Structured logs via the shared microservices logger. Use `LOG_LEVEL` to override environment defaults.
