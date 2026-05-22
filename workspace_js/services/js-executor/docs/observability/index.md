# Observability

## Health
`GET /health`

### Health status semantics
- `healthy` (HTTP 200): all critical checks OK
- `degraded` (HTTP 200): optional checks failed
- `unhealthy` (HTTP 503): at least one critical check failed

### Checks for this service
No external dependency checks are currently configured in the health handler. The endpoint reports service status only.

### Response shape
```json
{
  "status": "healthy",
  "service": "js-executor",
  "version": "1.0.0",
  "uptime": "48h32m15s",
  "timestamp": "2026-02-18T10:30:00Z",
  "lastCheckAt": "2026-02-18T10:29:55Z",
  "checks": {}
}
```

## Metrics
All service‑specific metrics are prefixed with `jsexec_`.

### Summary
| Metric | Type | Description |
|---|---|---|
| `jsexec_events_processed_total` | Counter | Events processed by consumer and status. |
| `jsexec_event_duration_seconds` | Histogram | End‑to‑end event processing duration by consumer. |
| `jsexec_payload_size_bytes` | Histogram | Incoming payload size distribution. |
| `jsexec_script_duration_seconds` | Histogram | V8 execution duration by script type. |
| `jsexec_script_errors_total` | Counter | Script errors by script type and error type. |
| `jsexec_compile_duration_seconds` | Histogram | Script compilation time by source. |
| `jsexec_asset_cache_total` | Counter | Asset cache hits by tier. |
| `jsexec_template_cache_total` | Counter | Template cache hits by tier. |
| `jsexec_bytecode_cache_total` | Counter | Bytecode cache hits by tier. |
| `jsexec_script_registry_total` | Counter | Script registry cache lookups by result. |
| `jsexec_piscina_completed_total` | Counter | Total Piscina tasks completed. |
| `jsexec_piscina_run_duration_seconds` | Histogram | Piscina worker execution time. |
| `jsexec_piscina_wait_duration_seconds` | Histogram | Queue wait time before worker pickup. |
| `jsexec_piscina_workers` | Gauge | Active worker thread count. |
| `jsexec_batch_size` | Histogram | Messages per NATS consumer batch. |
| `jsexec_nats_consumer_lag` | Gauge | Pending messages per consumer. |

### Labels
- `consumer`: NATS consumer name
- `status`: `success`/`error`
- `script_type`: `decode`, `validation`, `transform`
- `error_type`: script error classification
- `source`: `fresh` or `bytecode`
- `tier`: cache tier identifier
- `result`: `hit`/`miss`

### Default Node metrics
Node.js default metrics are enabled (heap, GC, event loop lag, active handles).

## Logs
Structured logs via shared logger. Use `LOG_LEVEL` to override environment defaults.
