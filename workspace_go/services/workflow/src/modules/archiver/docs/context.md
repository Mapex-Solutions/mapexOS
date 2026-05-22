# Bounded Context: Archiver

**Service:** workflow
**Module path:** `src/modules/archiver/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose

Archives workflow execution state into cold storage. The Archiver is the ONLY module that writes workflow execution documents to MongoDB; the Runtime never touches Mongo, working strictly against NATS KV + streams. It batch-consumes the `WORKFLOW-STATE` lifecycle stream, classifies events by type, writes the appropriate Mongo document shape, cleans up hot state in NATS KV on terminal events, and fans out to ClickHouse (via `EVENTS-WORKFLOW`) for long-term analytics. It also serves the HTTP query API for execution listing and detail.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Execution | A single run of a workflow instance (root aggregate persisted here) | Instance (the config) or Definition (the DAG template) |
| Lightweight stub | ~200B Mongo doc inserted on `created` for listing visibility | Full execution (~5–25KB) upserted on terminal |
| Hot state | Full execution JSON in NATS KV `exec.{uuid}` (authoritative while running) | MongoDB document (becomes authoritative on terminal) |
| Terminal event | `completed`, `failed`, or `cancelled` — triggers KV Get + BulkUpsertFull + KV Delete | `created`, `waiting`, `resumed` (lifecycle-only writes) |
| Backpressure mode | MongoManager signal (`Normal`/`Throttled`/`Backoff`) gating batch ingestion | NATS retry/backoff policy |

## Published Events (driven — outbound)

| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| WorkflowArchived (terminal history) | `mapexos.events.workflow` | workflow event payload (see `archiver.constant.EventsWorkflowSubject`) | events service (ClickHouse ingestion) |

## Consumed Events (driving — inbound)

| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| StateEvent | `mapexos.workflow.state.>` (stream `WORKFLOW-STATE`) | `shared/types.StateEvent` (aliased as `archiverMsg.StateEvent`) | runtime (RuntimePublisher) |

## Driving Ports (what can call this module)

- NATS batch consumer on `WORKFLOW-STATE` (`BatchSize=500`, per-message ACK/NACK via `BatchMessageHandlerV2`).
- HTTP routes under `/api/v1/workflow_executions` (JWT-authenticated):
  - `GET /` — paginated list with status/instance/definition filters.
  - `GET /:executionId` — detail; for non-terminal rows enriches with NATS KV hot state.

## Driven Ports (what this module requires)

- `ArchiveRepository` (MongoDB — bulk insert/upsert/update + find).
- `KVStore` (NATS KV `WORKFLOW-INSTANCES`) — read hot state on terminal, delete after archive.
- `Publisher` (NATS) — publish ClickHouse history events.
- `MongoManager` — backpressure signal + write-latency recording.

## Invariants and Business Rules

- Runtime MUST NEVER write to Mongo for executions — all Mongo writes go through this module.
- Terminal events MUST fetch full state from KV, upsert full document, then delete the KV key (in that order).
- Lightweight stubs on `created` MUST carry the pre-generated executionId (hex ObjectId from runtime) so terminal upsert hits the same document.
- Terminal docs MUST set `expireAt = now + 3d` for Mongo TTL auto-delete of hot storage.
- `retentionDays` mapping: `0` (uint16) becomes `-1` in the contract (org default); other values pass through.
- Backpressure: sleep 2s before batch when `MongoManager` reports `Backoff`.

## Known Cross-Context Interactions

- Consumes state events published by the Runtime module (same service).
- Produces `mapexos.events.workflow` consumed by the events service for ClickHouse cold storage and retention.
- Reads `runtime/application/ports.WorkflowExecution` entity type via the runtime ports (cross-module reference, no direct `domain/entities` import).
