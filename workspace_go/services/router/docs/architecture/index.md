# Architecture

## Patterns
- Event-driven execution via NATS JetStream
- Batch processing for high throughput

## Project Structure
```
src/
├── modules/
│   ├── routegroups/            # RouteGroup configuration management (CRUD + cache)
│   └── events/                 # Routing pipeline + match evaluation
└── shared/
    ├── configuration/          # Service configuration
    └── constants/              # Constants and defaults
```

## Module Responsibilities

### routegroups
Manages RouteGroup configuration via REST APIs. Provides external API with JWT authentication for users and internal API with API Key authentication for service-to-service communication. Configurations are cached in Redis for performance.

### events
Processes incoming route execution events and publishes results to downstream services based on configured routing rules. Handles asset cache invalidation across service instances.

## Domain Model (RouteGroups)
A RouteGroup defines how a stream of events is routed. It contains multiple routers, each with an optional match rule.

Key fields:
- `routers[]`: list of routing actions.
- `router.kind`: `rule_engine`, `lake_house`, `notification`, `trigger`, `save_event`.
- `router.match`: optional conditions using policy `all`/`any` and rules (`eq`/`neq`/`gt`/`gte`/`lt`/`lte`/`in`/`nin`).
- `isSystem`: global MAPEX templates (no orgId/pathKey).
- `isTemplate`: vendor/customer templates inherited by descendants.
- `orgId`, `pathKey`: multi-tenant scoping.

Template inheritance pattern:
- `isSystem=true`: global templates available to all orgs.
- `isTemplate=true`: vendor/customer templates inherited by descendants.
- `isSystem=false` and `isTemplate=false`: org-local route groups.

## Match Evaluation
Match rules are evaluated against the `event` payload using dot-path fields:
- Field path example: `payload.temperature`, `metadata.deviceType`.
- Operators: `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `in`, `nin`.
- Policy: `all` (AND), `any` (OR).
- If no match config is provided, the router **always matches**.

## Execution Pipeline
1. NATS consumer receives `route.execute` (stream `ROUTE-GROUPS`).
2. Messages are processed in batches for high throughput:
   - **Phase 1 (Parallel)**: Messages are processed concurrently by a bounded worker pool.
   - **Phase 2 (Flush)**: All buffered publishes are sent efficiently.
   - **Phase 3 (ACK)**: Each message is acknowledged, retried, or rejected.
3. Asset is resolved from TieredCache (L0 RAM -> L1 Disk -> L2 MinIO -> Fallback HTTP).
4. Each RouteGroup is loaded (Redis cache -> MongoDB fallback).
5. Routers are evaluated and published to downstream subjects.
6. Routing history is emitted to `events.router`.

## Subjects and Payloads

| Subject Pattern | Payload | Description |
|---|---|---|
| `ruleengine.{businessRuleId}.execute` | Enriched event with `businessRuleId` and metadata | Rule engine evaluation |
| `trigger.{triggerId}.execute` | `TriggerExecuteEvent` (original event as payload) | Trigger execution |
| `events.save` | `EventStoreDTO` with asset context and optional metadata | Event persistence |
| `events.lake_house` | Enriched event with asset identifiers and optional metadata | Lakehouse analytics |
| `events.notification` | Enriched event with `notificationId` and metadata | Notification dispatch |
| `events.router` | `RouterHistoryEvent` with match results per router | Routing audit/UI |

## Caching and Consistency
- Asset read model is stored in L2 (MinIO/S3). Router builds L1/L0 locally.
- Cache invalidation uses FANOUT subject `fanout.asset.invalidate`.
- Each service instance receives the invalidation message and clears L0+L1 for that asset.
- RouteGroup configs are cached in Redis (TTL: 60 minutes).
- Counter cache for list/count endpoints (TTL: 6 hours).

## Retry and DLQ
- Default retries: 5 attempts.
- Backoff: 1s, 5s, 30s, 2m, 10m.
- DLQ metadata includes `service=router` and `eventType=route.execute`.
- Invalid payloads (missing `orgId`, `assetUUID`, `event`, or invalid JSON) are REJECTed immediately (no retry).
- Processing errors are NACKed and retried.
