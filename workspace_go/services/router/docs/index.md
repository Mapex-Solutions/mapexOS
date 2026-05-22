# Router Service Documentation

## Overview
The Router service is the event routing layer of MapexOS. It consumes asset events from NATS JetStream, resolves asset context through the distributed TieredCache (L0 RAM, L1 Disk, L2 MinIO), evaluates RouteGroup match rules, and publishes routed events to downstream subjects including RuleEngine, Triggers, Event Store, Lakehouse, and Notifications. A single event can fan-out to multiple destinations based on per-router match conditions. Routing history is emitted for audit and UI visualization.

## Responsibilities
- Evaluate RouteGroup matches for incoming events using conditional match rules.
- Publish routing results to downstream NATS subjects.
- Manage RouteGroup configuration via REST API (external + internal MS-to-MS).
- Resolve asset context through TieredCache with automatic fallback.
- Handle cache invalidation via NATS FANOUT consumer.
- Emit routing execution history to `events.router` for audit.

## Non-Responsibilities
- Event ingestion (HTTP Gateway).
- Rule evaluation (RuleEngine).
- Trigger execution (Triggers service).
- Event persistence and analytics (Events service).
- Script execution (JS-Executor).

## Primary Data Flow
1. Consume `route.execute` from NATS JetStream (stream `ROUTE-GROUPS`).
2. Parse and validate required fields: `orgId`, `assetUUID`, `event`.
3. Resolve asset via TieredCache (L0 RAM -> L1 Disk -> L2 MinIO -> Fallback HTTP).
4. Load RouteGroups referenced by the asset (Redis cache -> MongoDB fallback).
5. For each router in the RouteGroup, evaluate match rules (policy `all`/`any`).
6. Publish matched events to downstream subjects (`events.save`, `events.lake_house`, `events.notification`, `ruleengine.{id}.execute`, `trigger.{id}.execute`).
7. Emit routing history to `events.router`.

## Docs Map
- [Architecture](architecture/index.md)
- [Endpoints](endpoints/index.md)
- [Configuration](configuration/index.md)
- [Operations](operations/index.md)
- [Observability](observability/index.md)
- [Tests](tests/index.md)
- [Benchmarks](benchmarks/index.md)
