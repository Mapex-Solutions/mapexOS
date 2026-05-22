# Observability

## Health

`GET /health`

### Health status semantics

- `healthy` (HTTP 200): all critical checks OK
- `degraded` (HTTP 200): optional checks failed (e.g. MinIO)
- `unhealthy` (HTTP 503): at least one critical check failed

### Checks for this service

| Dependency | Criticality |
|-----------|-------------|
| MongoDB | critical |
| Redis Shared | critical |
| NATS Core | critical |
| MinIO Definitions | non-critical |

### Response shape

```json
{
  "status": "healthy",
  "service": "workflow",
  "version": "1.0.0",
  "uptime": "48h32m15s",
  "timestamp": "2026-03-09T10:30:00Z",
  "checks": {
    "mongodb": { "connected": true, "latencyMs": 2 },
    "redis_shared": { "connected": true, "latencyMs": 1 },
    "nats": { "connected": true, "latencyMs": 1 },
    "minio_definitions": { "connected": true, "latencyMs": 3 }
  }
}
```

## Metrics

All service-specific metrics are prefixed with `workflow_`.

### Summary

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `workflow_definition_operations_total` | Counter | `operation`, `status` | Definition CRUD operations |
| `workflow_definition_operation_duration_seconds` | Histogram | `operation` | Duration of definition operations |
| `workflow_definition_list_results_count` | Histogram | — | Items returned per list query |
| `workflow_definition_cache_total` | Counter | `result` | TieredCache hits/misses |
| `workflow_cache_invalidations_total` | Counter | `status` | Cache invalidation events |

### Labels

- `operation`: `create`, `read`, `update`, `delete`, `list`
- `status`: `success`, `error`, `not_found`
- `result`: `hit`, `miss`

## Logs

Structured logs via the shared microservices logger. Prefixes by layer:

### Consumer prefixes (NATS message handlers)

- `[CONSUMER:WorkflowTrigger]` — Trigger consumer lifecycle
- `[CONSUMER:WorkflowResume]` — Resume consumer lifecycle
- `[CONSUMER:WorkflowState]` — Archiver state consumer lifecycle
- `[CONSUMER:WorkflowSignal]` — Signal consumer lifecycle (instances module)

### Service prefixes (application layer)

- `[SERVICE:Runtime]` — DAG execution events
- `[SERVICE:Definition]` — Definition CRUD and MinIO script management
- `[SERVICE:Instances]` — Instance query and cancellation
- `[SERVICE:Reconciler]` — Timer sweep events

### Infrastructure prefixes (adapters)

- `[INFRA:Archiver]` — MongoDB batch persistence events
- `[INFRA:Archive]` — MongoDB repository operations (BulkInsert, BulkUpsert)
- `[INFRA:RuntimePublisher]` — NATS publish events (state, resume, code, subworkflow)
- `[INFRA:InstanceStateRepo]` — NATS KV operations (create, get, checkpoint)
- `[INFRA:DefinitionLoader]` — TieredCache lookups and population
- `[INFRA:DefinitionStorage]` — MinIO script/bytecode operations

### Module prefixes (initialization)

- `[MODULE:App]` — Module orchestrator
- `[MODULE:Definitions]` — Definition module init
- `[MODULE:Engine]` — Engine module init
- `[MODULE:Instances]` — Instances module init
- `[MODULE:Runtime]` — Runtime module init
- `[MODULE:Archiver]` — Archiver module init
- `[MODULE:Reconciler]` — Reconciler module init

### Bootstrap prefix

- `[APP:BOOTSTRAP]` — Infrastructure initialization (metrics, middleware, NATS)

Use `LOG_LEVEL` env var to override environment defaults.
