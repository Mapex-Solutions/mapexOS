# Bounded Context: AssetTemplates

**Service:** assets
**Module path:** `src/modules/assettemplates/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-07-02

## Purpose
Owns the definition of an `AssetTemplate` — the reusable classification (manufacturer/model/category), schema (AvailableFields, DynamicFields with EVA `fieldId`) and scripts (validator, conversion, test, processor) that an Asset is bound to. Provides CRUD over HTTP, publishes scripts to MinIO (L2) so JS-Executor can consume them via TieredCache, and keeps denormalized classification names in sync by listening to list name-change events.

The module also owns **Template Update Plans** — scheduled, batched bulk migration of assets from one template (FROM) to another (TO). A plan snapshots the target asset ids at creation, self-schedules a NATS `@at` start timer, and on fire runs the migration in batches, switching each asset's template through the assets module (continue-on-failure per asset), then self-heals its counters from the per-asset execution records.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| AssetTemplate | Blueprint (classification + scripts + field schema) assigned to many Assets | `Asset` (the actual device instance) |
| System template | `IsSystem=true`, OrgID nil, visible to everyone (MAPEX global) | `IsTemplate=true` shared/vendor template scoped to an org |
| DynamicField | Field mapping with immutable `fieldId` (uint16) for EVA storage in ClickHouse | `AvailableFields` (string list used for Rule autocomplete only) |
| FieldId | Immutable numeric key per field; never reused, soft-deleted via `Status=0` | `NextFieldId` (auto-increment counter) |
| AvailableFields | Flat list of field names for rule/UI autocomplete, cached 24h in Redis | `DynamicFields` which carries typing and EVA info |
| FieldVocabulary | Curated, multi-tenant catalog of canonical dynamic-field names (English `value` + localized hint + type/unit/category) suggested in the authoring UI; read-only in this module | `DynamicField` (a concrete field mapping on a template) |
| Script | One of `ScriptValidator`, `ScriptConversion`, `ScriptTest`, `ScriptProcessor` — JS source executed by JS-Executor | — |
| List classification | External "list" entities (manufacturer/model/category) whose names are denormalized here | MongoDB `_id` references in `ManufacturerId`/`ModelId`/`CategoryId` |
| MigrationPlan | Aggregate for one bulk migration FROM one template TO another; holds only aggregate counters, the snapshotted asset ids, timing and status | `MigrationExecution` (the per-asset record) |
| MigrationExecution | One document per asset in a plan, carrying its `ExecStatus` and failure reason | `MigrationPlan` (holds counters only) |
| PlanStatus | Migration lifecycle `pending → scheduled → running → complete` (or `completed_with_errors`); `cancelled` from `pending`/`scheduled` | `AssetTemplate` status / DynamicField `Status` |
| ExecStatus | Per-asset outcome `pending → migrated` or `pending → failed` | `PlanStatus` (the plan-level lifecycle) |
| Stale-timer guard | 5s tolerance in `RunMigrationPlan` that ignores a duplicate/late `@at` fire so a plan runs at most once | `AllowMsgSchedules` (the JetStream native-schedule flag) |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| Template scripts write | MinIO bucket key `{orgId|mapexos_public}/{templateId}.json` (object storage, not NATS) | Scripts bundle (inferred from `TemplateStoragePort.WriteScripts`) | JS-Executor (via TieredCache L2) |
| Template cache invalidate | `mapexos.fanout.template.invalidate` (stream `FANOUT`) | `contracts/services/assets/assettemplates/types.go::TemplateInvalidatePayload` | Router, Events, JS-Executor |
| Migration start timer | `assettemplates.migration.schedule` (STATIC publish subject; delivered to target `assettemplates.migration.timer.start`) | scheduled self-publish (native NATS `@at`); plan id in payload + MsgId, never in the subject | this module's migration_timers consumer |

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| List name updated | `mapexos.lists.name_updated` (stream `MAPEXOS-LISTS`) | `Event{ListId, ListType, NewName, OrgId}` (`list_name_updated/types.go`) | Core `mapexos` service (lists module) |
| Migration start timer | `assettemplates.migration.timer.start` (stream `ASSETTEMPLATES`, WorkQueue, `AllowMsgSchedules`) | plan id in payload + MsgId | this module (self-scheduled) |

Handled `ListType` values: `asset_manufacturer`, `asset_model`, `asset_category` — any other type is acked and ignored. Consumer uses NATS Core connection with queue group `{service}-LIST-NAME-GROUP` and DLQ policy.

## Driving Ports (inbound — who calls this module)
- HTTP `/api/v1/asset_templates` (JWT auth): CRUD + `GET /counter` + `GET /:id/available_fields`
- HTTP `GET /api/v1/asset_templates/field-vocabulary` (JWT auth, `AssetTemplateList` permission): read-only field-name vocabulary for the caller's org (system standard unioned with the org's own entries), grouped by category and localized via `?lang` (en-US fallback)
- HTTP `/internal/templates` (API-Key auth): TieredCache fallback endpoints for JS-Executor
- HTTP `/api/v1/asset_templates/migrations` (JWT auth, gated by `templatemigrations.*` permissions): migration plan CRUD + per-plan executions listing
- NATS consumer on `mapexos.lists.name_updated` for denormalized-name sync
- NATS migration_timers consumer on `assettemplates.migration.timer.start` → `RunMigrationPlan`

## Driven Ports (outbound — what this module requires)
- `repositories.AssetTemplateRepository` — MongoDB persistence (`infrastructure/persistence/mongo`)
- `repositories.FieldVocabularyRepository` — read-only MongoDB access to the `field_vocabulary` collection (`infrastructure/persistence/mongo`)
- `ports.TemplateStoragePort` — MinIO writer for scripts (`infrastructure/storage/minio`)
- `common.AppCache` — Redis cache for `AvailableFields` (24h TTL) and counter (6h TTL)
- `natsModel.Bus` (name `core`) — consumer wiring for list-name sync
- `repositories.MigrationPlanRepository` / `repositories.MigrationExecutionRepository` — MongoDB persistence (`template_migration_plans` / `template_migration_executions`)
- `ports.MigrationSchedulerPort` — native NATS `@at` scheduling of the start timer; `CancelStart` is a no-op seam (no per-MsgId cancel primitive exists)
- `ports.TemplateSwitcherPort` → assets module `AssetServicePort.UpdateAssetById`, which owns the MinIO re-projection + FANOUT invalidate for the switched asset

## Invariants and Business Rules
- `DynamicField.FieldId` is assigned from `NextFieldId` and is IMMUTABLE — never reused, even after deletion
- Deletion of a DynamicField sets `Status=0` (deprecated); historical events keep resolving their field
- Maximum 200 active (`Status=1`) DynamicFields per template (per entity doc comment)
- System templates (`IsSystem=true`) have `OrgID=nil` and use MinIO prefix `mapexos_public`; private templates use `OrgID` hex
- `AvailableFields` cache is invalidated on create/update
- Classification names (`manufacturerName`/`modelName`/`categoryName`) are denormalized and kept in sync only via the `mapexos.lists.name_updated` consumer
- Counter cache (Redis) invalidated on create/delete
- Field vocabulary is read-only here: the effective list is `enabled:true` AND (`isSystem:true` OR org match); creating/editing org-scoped entries is out of scope for this module
- Migration TO template can be ANY template (no compatibility constraint with FROM)
- The asset ids are snapshotted at plan creation; assets added/removed afterward do not change the plan's scope
- Migration runs in batches of `BatchSize` (config `template_migration_batch_size`, default 100), continue-on-failure per asset — one asset's failure never aborts the plan
- A plan may be edited or cancelled only while `pending`/`scheduled`; any later attempt is a 409 conflict
- `RunMigrationPlan` is idempotent: a status guard plus a 5s stale-timer tolerance ensures each plan runs at most once even on a duplicate/late `@at` fire
- Finalization recomputes plan counters from the execution documents (self-healing, no batch transaction — the switch side effects live outside Mongo anyway)
- The migration never writes the asset collection directly; every template switch goes through the assets module

## Known Cross-Context Interactions
- Assets module (same service): reads templates to enrich asset responses with classification + `AssetIDPath`; migration switches an asset's template via `AssetServicePort.UpdateAssetById` (assets module owns the MinIO re-projection + FANOUT invalidate)
- JS-Executor service: consumes the MinIO L2 scripts payload via TieredCache, falls back to `/internal/templates` on miss
- Core `mapexos` service (lists module): publishes `mapexos.lists.name_updated` that this module consumes
- Events / edge: unaffected by migration — no new outbound events beyond the per-asset switch handled by the assets module
