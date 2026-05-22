# Operations

## Start
```bash
./bin/router
```

## Smoke Checks
- `GET /health` -- returns service health status
- `GET /metrics` -- returns Prometheus metrics
- `GET /api/v1/route_groups` with valid JWT -- returns paginated route groups

## Dependencies
- **MongoDB** -- RouteGroup persistence (database: `router`)
- **Redis** -- RouteGroup cache (DB 0) + shared cache for coverage/permissions (DB 5)
- **NATS JetStream** -- Route execution stream (`ROUTE-GROUPS`) + FANOUT invalidation stream (`FANOUT`)
- **MinIO/S3** -- Asset read model (L2 cache, bucket: `mapex-assets`)
- **MapexOS API** -- Permission middleware calls (`MAPEXOS_URL`)
- **Assets service** -- TieredCache fallback when L2 misses (`ASSETS_URL`)

## NATS Streams
- Route execution: stream `ROUTE-GROUPS`, subject `route.execute` (WorkQueue, durable)
- Cache invalidation: stream `FANOUT`, subject `fanout.asset.invalidate` (FANOUT, ephemeral)

## Scaling Guidance
- Scale horizontally by adding more instances when message processing lag increases.
- Configure `NATS_BATCH_SIZE` and `NATS_FETCH_TIMEOUT` to balance throughput and latency for your workload.

## Failure Handling
- Invalid payloads (bad JSON, missing `orgId`/`assetUUID`/`event`) are REJECTed (no retry).
- Processing errors are NACKed and retried up to 5 times with backoff: 1s, 5s, 30s, 2m, 10m.
- After max retries, messages are sent to DLQ with metadata: `service=router`, `eventType=route.execute`.

## Benchmarks
- [Benchmarks](../benchmarks/index.md)
