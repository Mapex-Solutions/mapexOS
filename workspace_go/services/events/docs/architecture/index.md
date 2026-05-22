# Architecture

## Design
Modular architecture with clean separation of concerns.

## Project Structure
```
src/
├── modules/
│   ├── events/                 # Event storage + query APIs + NATS consumers
│   └── retention/              # Retention policy management
└── shared/
    └── configuration/           # Service configuration
```

## Module Initialization Order

1. **retention** — Manages retention policies per organization
2. **events** — Consumes events from NATS and stores in ClickHouse

> **IMPORTANT:** `retention` must initialize before `events` because the Events module requires retention policies to resolve TTL.

## Module Responsibilities
- `events`: ClickHouse persistence + query endpoints for event streams (7 consumers, 7 query endpoints)
- `retention`: per-org retention policies for event datasets, auto-creates defaults on new org creation

## Main Event Flow
```
NATS batch fetch (up to NATS_BATCH_SIZE messages)
  → Phase 1: parallel parse/validate/map (bounded worker pool)
  → Phase 2: ClickHouse bulk insert (single INSERT for valid entities)
  → Phase 3: resolve outcomes — ack (success), nack (insert failure), reject → DLQ (parse/validation failure)
```

## HOT Layer + ClickHouse
- ClickHouse is the hot analytics store for UI and debugging
- Streams are separated by event type (raw/router/trigger/jsexec/store)

## EVA (Entity-Value-Attribute)
- Processed events use EVA fields for flexible schema
- EVA enables fast filtering without ALTER TABLE
- EVA fields are resolved using tiered template cache (L0 RAM, L1 disk, L2 S3)
