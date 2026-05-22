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
  "service": "triggers",
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
All service‑specific metrics are prefixed with `triggers_`.

### Execution Pipeline Metrics
| Metric | Type | Meaning |
|---|---|---|
| `triggers_trigger_processed_total` | Counter | Total triggers processed by status (`success`, `error`, `disabled`, `no_executor`). |
| `triggers_trigger_processing_duration_seconds` | Histogram | End‑to‑end time for a trigger execution (from fetch to publish). |
| `triggers_trigger_batch_size` | Histogram | Number of messages per NATS batch. |
| `triggers_message_total` | Counter | Message lifecycle outcomes (`ack`, `nack`, `reject`, `dlq`). |
| `triggers_trigger_cache_total` | Counter | Cache results for trigger config lookup (`hit`, `miss`). |
| `triggers_placeholder_resolutions_total` | Counter | Placeholder resolution outcomes (`success`, `error`). |

### Executor Metrics
| Metric | Type | Labels | Meaning |
|---|---|---|---|
| `triggers_executor_duration_seconds` | Histogram | `type` | Execution time by executor type (http, email, mqtt, rabbitmq, nats, websocket, teams, slack). |
| `triggers_executor_total` | Counter | `type`, `status` | Executor attempts by type and status (`success`, `error`). |

### Publish Metrics
| Metric | Type | Meaning |
|---|---|---|
| `triggers_event_published_total` | Counter | Publish outcomes to `events.trigger` (`ok`, `error`). |
| `triggers_publish_duration_seconds` | Histogram | Publish latency to `events.trigger`. |

### Labels
- `status`: `success`, `error`, `disabled`, `no_executor`
- `result`: `ack`, `nack`, `reject`, `dlq`, `hit`, `miss`
- `type`: executor type (`http`, `email`, `mqtt`, `rabbitmq`, `nats`, `websocket`, `teams`, `slack`)

### Go runtime + process collectors
Controlled by:
- `METRICS_GO_COLLECTOR` (default: `true`)
- `METRICS_PROCESS_COLLECTOR` (default: `true`)

### Full instrumentation map
See `documentations/architecture/microservices/services/triggers/metrics.md` for placement details and cardinality.

## Logs
Structured logs via the shared microservices logger. Use `LOG_LEVEL` to override environment defaults.
