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
  "service": "mapexos",
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

---

## Metrics

`GET /metrics`

Prometheus metrics are exposed via a custom isolated registry with namespace `mapexos`. The `/metrics` endpoint is registered before global middlewares (CORS, Helmet) and requires no authentication.

### Layer 3: Service-Specific Metrics

All metrics below are declared in `src/bootstrap/metrics.go` and provided to the DIG container as `*MapexosMetrics`.

#### Auth Subsystem

| Metric | Type | Labels | Description |
|---|---|---|---|
| `mapexos_auth_attempts_total` | CounterVec | `method`, `status` | Authentication attempts. `method`: login, refresh, logout. `status`: success, failure. |
| `mapexos_auth_duration_seconds` | HistogramVec | `method` | Authentication operation latency by method. |
| `mapexos_session_operations_total` | CounterVec | `operation`, `status` | Session store operations. `operation`: store, get, invalidate. `status`: success, failure. |

#### User CRUD Subsystem

| Metric | Type | Labels | Description |
|---|---|---|---|
| `mapexos_user_operations_total` | CounterVec | `operation`, `status` | User CRUD operations. `operation`: create, get, get_by_id, update, delete, list, count. `status`: success, failure. |
| `mapexos_user_operation_duration_seconds` | HistogramVec | `operation` | User CRUD latency by operation. |
| `mapexos_user_list_results_count` | Histogram | (none) | Items returned per user list query. Buckets: 0, 1, 5, 10, 25, 50, 100, 250. |

#### Group CRUD Subsystem

| Metric | Type | Labels | Description |
|---|---|---|---|
| `mapexos_group_operations_total` | CounterVec | `operation`, `status` | Group CRUD operations. `operation`: create, get, get_by_id, update, delete, list, add_member, remove_member. `status`: success, failure. |
| `mapexos_group_operation_duration_seconds` | HistogramVec | `operation` | Group CRUD latency by operation. |
| `mapexos_group_list_results_count` | Histogram | (none) | Items returned per group list query. Buckets: 0, 1, 5, 10, 25, 50, 100, 250. |

#### Role CRUD Subsystem

| Metric | Type | Labels | Description |
|---|---|---|---|
| `mapexos_role_operations_total` | CounterVec | `operation`, `status` | Role CRUD operations. `operation`: create, get, get_by_id, update, delete, list. `status`: success, failure. |
| `mapexos_role_operation_duration_seconds` | HistogramVec | `operation` | Role CRUD latency by operation. |
| `mapexos_role_list_results_count` | Histogram | (none) | Items returned per role list query. Buckets: 0, 1, 5, 10, 25, 50, 100, 250. |

#### Membership CRUD Subsystem

| Metric | Type | Labels | Description |
|---|---|---|---|
| `mapexos_membership_operations_total` | CounterVec | `operation`, `status` | Membership CRUD operations. `operation`: create, get, get_by_id, update, delete, list. `status`: success, failure. |
| `mapexos_membership_operation_duration_seconds` | HistogramVec | `operation` | Membership CRUD latency by operation. |
| `mapexos_membership_list_results_count` | Histogram | (none) | Items returned per membership list query. Buckets: 0, 1, 5, 10, 25, 50, 100, 250. |

#### Organization CRUD Subsystem

| Metric | Type | Labels | Description |
|---|---|---|---|
| `mapexos_organization_operations_total` | CounterVec | `operation`, `status` | Organization CRUD operations. `operation`: create, get, get_by_id, update, delete, list, tree. `status`: success, failure. |
| `mapexos_organization_operation_duration_seconds` | HistogramVec | `operation` | Organization CRUD latency by operation. |
| `mapexos_organization_list_results_count` | Histogram | (none) | Items returned per organization list query. Buckets: 0, 1, 5, 10, 25, 50, 100, 250. |

#### Cache Subsystem

| Metric | Type | Labels | Description |
|---|---|---|---|
| `mapexos_cache_total` | CounterVec | `type`, `result` | Cache lookups. `type`: authorization, coverage, counter. `result`: hit, miss. |

### Go Runtime and Process Metrics (Optional)

Enabled via environment variables `METRICS_GO_COLLECTOR=true` and `METRICS_PROCESS_COLLECTOR=true`.

| Metric | Description |
|---|---|
| `go_goroutines` | Number of goroutines |
| `go_gc_duration_seconds` | GC pause duration |
| `go_memstats_alloc_bytes` | Allocated heap bytes |
| `go_memstats_heap_alloc_bytes` | Heap allocation bytes |
| `process_cpu_seconds_total` | Total CPU time |
| `process_resident_memory_bytes` | Resident memory (RSS) |
| `process_open_fds` | Open file descriptors |

### Metrics Summary

| Subsystem | Counters | Histograms | Total Series (est.) |
|---|---|---|---|
| Auth | 2 (attempts + session) | 1 | ~30 |
| User | 1 | 2 (duration + list) | ~50 |
| Group | 1 | 2 | ~60 |
| Role | 1 | 2 | ~40 |
| Membership | 1 | 2 | ~40 |
| Organization | 1 | 2 | ~50 |
| Cache | 1 | 0 | ~10 |
| **Total** | **8** | **11** | **~280** |

---

## Logs

Structured logs via the shared microservices logger. Use `LOG_LEVEL` to override environment defaults.

| Level | Description |
|---|---|
| `debug` | Detailed development logs |
| `info` | Operational events (startup, connections) |
| `warn` | Recoverable issues |
| `error` | Failures requiring attention |
