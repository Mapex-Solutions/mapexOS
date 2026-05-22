# Observability

## Health
`GET /health`

### Health Status Semantics
- `healthy` (HTTP 200): all critical checks OK.
- `degraded` (HTTP 200): at least one optional check failed, all critical checks OK.
- `unhealthy` (HTTP 503): at least one critical check failed.

### Checks
| Dependency | Check Name | Criticality |
|---|---|---|
| MongoDB | `mongodb` | critical |
| Redis (app) | `redis:app` | critical |
| Redis (shared) | `redis:shared` | critical |
| NATS (core) | `nats:core` | critical |
| MinIO (assets) | `minio:assets` | optional |

Health check configuration: cache TTL = 10s, timeout = 5s.

### Response Schema
```json
{
  "status": "healthy",
  "service": "router",
  "version": "1.0.0",
  "uptime": "48h32m15s",
  "timestamp": "2026-02-18T10:30:00Z",
  "lastCheckAt": "2026-02-18T10:29:55Z",
  "checks": {
    "mongodb": { "connected": true, "latencyMs": 2 },
    "redis:app": { "connected": true, "latencyMs": 1 },
    "redis:shared": { "connected": true, "latencyMs": 1 },
    "nats:core": { "connected": true, "latencyMs": 3 },
    "minio:assets": { "connected": true, "latencyMs": 4, "critical": false }
  }
}
```

## Metrics
`GET /metrics`

All service-specific metrics are prefixed with `router_`.

### Summary
| Metric | Type | Description |
|---|---|---|
| `router_event_processed_total` | Counter | Events that completed processing, by status. |
| `router_event_processing_duration_seconds` | Histogram | End-to-end processing time per event. |
| `router_event_batch_size` | Histogram | Number of messages per batch fetch. |
| `router_message_total` | Counter | Message lifecycle outcomes (ack/nack/reject). |
| `router_asset_cache_total` | Counter | Asset cache lookups by tier. |
| `router_asset_cache_duration_seconds` | Histogram | Cache lookup latency by tier. |
| `router_cache_invalidations_total` | Counter | Cache invalidation operations by status. |
| `router_match_evaluations_total` | Counter | Match evaluations by result. |
| `router_match_rules_evaluated_total` | Counter | Total individual match rules evaluated. |
| `router_event_published_total` | Counter | NATS publish attempts by kind and status. |
| `router_publish_duration_seconds` | Histogram | NATS publish latency by kind. |
| `router_routegroup_operations_total` | Counter | RouteGroup CRUD operations by type and status. |
| `router_routegroup_operation_duration_seconds` | Histogram | RouteGroup CRUD operation latency. |
| `router_routegroup_list_results_count` | Histogram | Items returned per list query. |
| `router_routegroup_cache_total` | Counter | RouteGroup cache hit/miss ratio. |

### Labels
| Label | Metric(s) | Values |
|---|---|---|
| `status` | `router_event_processed_total`, `router_event_published_total`, `router_cache_invalidations_total`, `router_routegroup_operations_total` | `success`, `error` |
| `result` | `router_message_total` | `ack`, `nack`, `reject` |
| `result` | `router_match_evaluations_total` | `matched`, `unmatched`, `no_config` |
| `result` | `router_routegroup_cache_total` | `hit`, `miss` |
| `tier` | `router_asset_cache_total`, `router_asset_cache_duration_seconds` | `L0_RAM`, `L1_Disk`, `L2_MinIO`, `Fallback_HTTP`, `MISS` |
| `kind` | `router_event_published_total`, `router_publish_duration_seconds` | `save_event`, `rule_engine`, `lake_house`, `notification`, `trigger` |
| `operation` | `router_routegroup_operations_total`, `router_routegroup_operation_duration_seconds` | `list`, `create`, `read`, `update`, `delete` |

### Go Runtime and Process Collectors
Controlled by environment variables:
- `METRICS_GO_COLLECTOR` (default: `true`) -- enables Go runtime metrics.
- `METRICS_PROCESS_COLLECTOR` (default: `true`) -- enables process metrics.

## Logs
Structured logs via the shared microservices logger. Use `LOG_LEVEL` to override environment defaults (`debug` for dev, `info` for production).
