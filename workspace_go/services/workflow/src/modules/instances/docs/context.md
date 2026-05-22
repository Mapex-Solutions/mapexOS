# Bounded Context: Instances

**Service:** workflow
**Module path:** `src/modules/instances/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose

Owns the `WorkflowInstance` root aggregate — the parameterised configuration that binds a `WorkflowDefinition` to concrete external inputs, ownership (org/pathKey), retention policy, and execution options (`uniqueExecution`, `enabled`, `isTemplate`). One instance can spawn N executions (1:N). Provides CRUD + a TieredCache-backed loader used by the Runtime to resolve an `instanceId` before each run. Also exposes `POST /:instanceId/execute` as a synchronous HTTP trigger that delegates to the runtime's `ExecuteByInstanceID`.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| WorkflowInstance | Root aggregate: immutable config for a run (definitionId + externalInputs + ownership + retention) | WorkflowDefinition (template) or WorkflowExecution (run) |
| uniqueExecution | When true, only one execution may be in-flight per instance at a time (inferred — field on entity) | Singleton definition |
| retentionDays | Per-instance ClickHouse retention override (0 = fall back to org policy, propagated onto each execution) | Mongo TTL (fixed 3 days post-terminal) |
| externalInputs | Free-form map provided by the instance author; merged into the execution context on start | Event payload (runtime-bound payload) |
| InstanceLoader | `TieredCache(name:"instances")` wrapper (L0 RAM → L1 Disk → MongoDB fallback) | Direct repository |

## Published Events (driven — outbound)

_None directly._ The HTTP `execute` endpoint calls into the runtime service in-process; the runtime is what publishes downstream events.

## Consumed Events (driving — inbound)

_None._

## Driving Ports (what can call this module)

- HTTP routes under `/api/v1/workflow_instances` (JWT auth):
  - `GET /counter`, `GET /`, `GET /:instanceId`, `POST /`, `PUT /:instanceId`, `DELETE /:instanceId`.
  - `POST /:instanceId/execute` — synchronous trigger, returns `ExecuteResponseDTO` with `workflowUUID`/`status`/`errorInfo`.
- Cross-module Go API: `ports.InstanceLoaderPort.GetInstance(...)` / `Invalidate(...)` (consumed by runtime).

## Driven Ports (what this module requires)

- `InstanceRepository` (MongoDB) — CRUD source of truth.
- `InstanceLoaderPort` — TieredCache L0+L1 with Mongo fallback (key `instance:{id}`).
- `RuntimeServicePort.ExecuteByInstanceID` (cross-module) for the HTTP execute endpoint (inferred from runtime's port surface).

## Invariants and Business Rules

- `orgId` and `pathKey` are populated from `RequestContext` at create time (coverage middleware), not from the request body.
- `isSystem=false` for all user-created instances (system-authored instances bypass this path).
- Cache MUST be invalidated on every mutating CRUD (`UpdateInstanceById`, `DeleteInstanceById`) to prevent stale reads by the runtime (inferred from loader's `Invalidate` contract).
- `retentionDays` flows from instance → execution at creation and is carried through to the ClickHouse events record.
- `uniqueExecution` intent is enforced at runtime (signal/start dispatch), not at this module's boundary (inferred).

## Known Cross-Context Interactions

- Consumed by **runtime** via `InstanceLoaderPort` on every `handleNewInstance`/`signalOrStart` path.
- References **definitions** via `definitionId` + `definitionName` + `definitionVersion` (denormalized at instance creation).
- Feeds **archiver** indirectly: the instance's `retentionDays` becomes the execution's `retentionDays`, which the archiver propagates to ClickHouse.
