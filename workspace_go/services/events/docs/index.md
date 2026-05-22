# Events Service Documentation

## Overview
The Events service is the **HOT event layer** of MapexOS. It consumes multiple NATS streams and persists event data in ClickHouse for low‑latency analytics and UI queries. It also applies **retention per organization** and supports **EVA (Entity‑Value‑Attribute)** to keep schema flexible without sacrificing query performance.

## Responsibilities
- Consume NATS event streams and persist to ClickHouse
- Serve query APIs for raw/router/trigger/businessrule/jsexec/store events
- Manage retention policies (per org and per event type)
- Resolve EVA fields for processed event queries

## Non‑Responsibilities
- Event ingestion from external systems (HTTP Gateway)
- Business rule evaluation (RuleEngine)
- Script execution (JS‑Executor)

## Primary Data Flow
1. NATS consumer receives a batch for a specific stream (up to `NATS_BATCH_SIZE` messages)
2. Phase 1: parallel parse/validate/map (bounded worker pool)
3. Phase 2: single ClickHouse bulk insert for the entire batch
4. Phase 3: per-message outcomes — ack (success), nack (insert failure), reject → DLQ (parse/validation failure)
5. Query APIs fetch stored events with filters and pagination

## EVA (Entity-Value-Attribute)
- EVA stores dynamic fields by numeric IDs (fast and compact)
- Avoids schema migrations while keeping queries fast
- Enables flexible analytics and ad-hoc querying without table changes
- EVA fields are resolved using tiered template cache (L0 RAM, L1 disk, L2 S3)

## Docs Map
- [Architecture](architecture/index.md)
- [Endpoints](endpoints/index.md)
- [Configuration](configuration/index.md)
- [Operations](operations/index.md)
- [Observability](observability/index.md)
- [Tests](tests/index.md)
- [Benchmarks](benchmarks/index.md)
