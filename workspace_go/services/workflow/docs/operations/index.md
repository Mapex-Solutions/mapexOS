# Operations

## Start

```bash
./bin/workflow
```

## Smoke checks

- `GET /health` — should return `{"status": "healthy", ...}` with HTTP 200
- `GET /metrics` — should return Prometheus text format
- `POST /api/v1/workflows` — should create a workflow definition (requires auth)
- Publish a trigger message to `WORKFLOW-TRIGGER` stream — should create and execute an instance

## Dependencies

| Dependency | Required | Purpose |
|-----------|----------|---------|
| MongoDB | yes | Archive storage for completed instances + definition metadata |
| NATS | yes | JetStream streams, KV store, pub/sub |
| Redis | yes | Shared DB (auth middleware session/permission cache) |
| MinIO | yes | Source of truth for workflow definition JSON + code node scripts |

## Runtime Limits

| Limit | Value | Purpose |
|-------|-------|---------|
| MaxInlineSteps | 300 | Prevent infinite loops per execution cycle |
| MaxSubworkflowDepth | 10 | Prevent recursive subworkflow chains |
| MaxLoopIterations | 10,000 | Prevent unbounded loops |
| MaxFanoutBranches | 20 | Limit parallel branch count |
| InlineTimeoutSeconds | 30 | Timeout for a single execution cycle |

## Graceful Shutdown

Shutdown hooks are registered via `bootstrap/shutdown.go` with priority ordering:

| Priority | Component | Action |
|----------|-----------|--------|
| P0 | Fiber HTTP | Stop accepting requests, drain in-flight |
| P5 | MongoDB | Close connection pool |
| P5 | Redis | Close connection |
| P5 | NATS | Close connection (concurrent with MongoDB/Redis) |

## Benchmarks

See [Benchmarks](../benchmarks/index.md) for methodology and results.
