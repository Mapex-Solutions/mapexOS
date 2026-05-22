# Architecture

## Patterns
- Event‑driven execution via NATS JetStream
- Modular architecture with clean separation of concerns

## Project Structure
```
src/
├── modules/
│   ├── triggers/               # Trigger configuration management (CRUD)
│   └── events/                 # Trigger execution pipeline
│       └── infrastructure/
│           ├── registry/        # Executor factory registry
│           ├── technical/       # Technical executors (http, mqtt, rabbitmq, nats, websocket)
│           └── communications/  # Communication executors (email, teams, slack)
└── shared/
    └── configuration/           # Service configuration and defaults
```

## Domain Model (Triggers)
A Trigger is a **reusable integration action** stored in MongoDB and executed on demand.

Key fields:
- `triggerType`: execution adapter (`http`, `mqtt`, `rabbitmq`, `nats`, `websocket`, `email`, `teams`, `slack`).
- `category`: UX grouping (`technical` or `communication`).
- `config`: union type containing exactly one config block that matches `triggerType`.
- `enabled`: disabled triggers are ACKed and never executed.
- `isSystem`: global MAPEX templates (no `orgId`/`pathKey`).
- `isTemplate`: templates reusable across tenant hierarchies.
- `orgId`, `pathKey`: multi‑tenant scoping.

Template inheritance pattern:
- `isSystem=true`: global templates available to all orgs.
- `isTemplate=true`: vendor/customer templates inherited by descendants.
- `isSystem=false` and `isTemplate=false`: org‑local triggers.

## Execution Pipeline
1. NATS consumer receives `trigger.*.execute` (stream `TRIGGERS`).
2. Batch is processed with a **worker pool** (`TRIGGER_EXECUTOR_WORKERS`).
3. Trigger config loaded via cache‑aside (Redis → Mongo).
4. Placeholders in config are resolved from event payload.
5. Executor runs the action (HTTP, Email, MQTT, etc.).
6. Result is published to `events.trigger` for auditing.
7. Message is ACKed, NACKed (retry), or REJECTed (invalid payload).

Retry and DLQ behavior:
- Default retries: 5 attempts.
- Backoff: 1s, 5s, 30s, 2m, 10m.
- DLQ metadata includes `service=triggers` and `eventType=trigger.execute`.

## Placeholder Resolution
- Syntax: `{{path.to.field}}`.
- Resolution is **recursive** across nested objects and arrays.
- Data source is the `payload` of `TriggerExecuteEvent`.
- If you need `orgId`, `pathKey`, or other values, include them in the payload explicitly.

## Executor Registry
Executors are registered via a factory registry and selected by `triggerType`:
- Technical: `http`, `mqtt`, `rabbitmq`, `nats`, `websocket`
- Communication: `email`, `teams`, `slack`

Each executor implements a standard interface, making it easy to add new delivery mechanisms.

## Caching & Consistency
- Trigger configs are cached in Redis for fast execution (TTL: 60 minutes).
- Cache is automatically invalidated on update/delete.
- A counter cache is used for list/count endpoints (TTL: 6 hours).

## Related Contracts
See [API Specification](api-specification/index.md) for full DTO definitions.
