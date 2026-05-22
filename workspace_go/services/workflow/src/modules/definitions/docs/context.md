# Bounded Context: Definitions

**Service:** workflow
**Module path:** `src/modules/definitions/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose

Owns the authoring model for workflows: the `WorkflowDefinition` root aggregate — nodes, edges, variables, external inputs/signals, condition groups, retry policy, plugin references. Provides CRUD + validation + storage for definitions (MongoDB primary; MinIO as L2 for code-node scripts surfaced via TieredCache). Runs domain-level validators: DAG cycle detection, node schema validation, and plugin reference validation. Acts as the source-of-truth queried by Runtime (via TieredCache-backed loader) and consumed by the frontend editor.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| WorkflowDefinition | Root aggregate: DAG template (nodes+edges) + schema (states, inputs, signals) + metadata | WorkflowInstance (a parameterised config) or WorkflowExecution (a run) |
| FieldValue | Polymorphic value source (`event`, `state`, `variable`, `input`, `literal`, `nodeOutput`) | Plain JSON value |
| ConditionGroup | Recursive boolean tree with AND/OR/NAND/NOR logic over ConditionItems | Condition node at runtime (uses the group structure) |
| DefinitionStatus | `valid` / `plugin_missing` / `invalid` — set after validation | Execution status |
| Node script | JS source stored out-of-band in MinIO (L2), referenced by nodeId within the definition | Inline node config |

## Published Events (driven — outbound)

_None._ (CRUD-only module; no NATS publishers. Cache invalidation for plugins lives in the `plugins` module.)

## Consumed Events (driving — inbound)

_None._ (No NATS consumers.)

## Driving Ports (what can call this module)

- External HTTP routes under `/api/v1/workflow_definitions` (JWT): list, count, create, get-by-id, update, delete.
- Internal HTTP routes under `/internal/workflow-scripts` (API-key): `GetNodeScript` — used by js-workflow-executor when its TieredCache L2 (MinIO) misses.
- Cross-module Go API: `ports.DefinitionLoaderPort` (in runtime) aliases this module's `DefinitionRepository` via a TieredCache wrapper — runtime never imports this module's `domain/entities` directly.

## Driven Ports (what this module requires)

- `DefinitionRepository` (MongoDB) — source of truth for definitions.
- `DefinitionStoragePort` (MinIO provider) — object storage for node scripts (L2 of the script TieredCache).
- Domain services: `cycle_detector`, `node_validator`, `plugin_validation` (pure functions, no I/O).

## Invariants and Business Rules

- DAG must be acyclic (enforced by `cycle_detector`) — save rejects graphs with cycles.
- Every `nodeId` referenced by an edge must exist in `nodes`; validator enforces referential integrity.
- Plugin-backed nodes must reference an installed plugin; missing plugins set `status = plugin_missing` and populate `missingPlugins[]`.
- Multi-tenant visibility: `isTemplate=true` definitions are visible to child orgs; `isTemplate=false` are org-local (filtered via `pathKey`/`orgId`).
- Node scripts are write-through to MinIO; GET falls back to Mongo and re-populates MinIO on miss (inferred from `GetNodeScript` semantics).

## Known Cross-Context Interactions

- Consumed by the **runtime** module via `DefinitionLoaderPort` (TieredCache L0→L1→Mongo) before each DAG walk.
- Consumed by **js-workflow-executor** (internal `/internal/workflow-scripts` endpoint) for L2 fallback when its in-process script cache misses.
- Consumed by the frontend editor for CRUD and validation feedback.
- Plugin references validated against the **plugins** module's manifests at save time.
