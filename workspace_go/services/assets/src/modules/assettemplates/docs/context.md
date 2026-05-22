# Bounded Context: AssetTemplates

**Service:** assets
**Module path:** `src/modules/assettemplates/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-11

## Purpose
Owns the definition of an `AssetTemplate` — the reusable classification (manufacturer/model/category), schema (AvailableFields, DynamicFields with EVA `fieldId`) and scripts (validator, conversion, test, processor) that an Asset is bound to. Provides CRUD over HTTP, publishes scripts to MinIO (L2) so JS-Executor can consume them via TieredCache, and keeps denormalized classification names in sync by listening to list name-change events.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| AssetTemplate | Blueprint (classification + scripts + field schema) assigned to many Assets | `Asset` (the actual device instance) |
| System template | `IsSystem=true`, OrgID nil, visible to everyone (MAPEX global) | `IsTemplate=true` shared/vendor template scoped to an org |
| DynamicField | Field mapping with immutable `fieldId` (uint16) for EVA storage in ClickHouse | `AvailableFields` (string list used for Rule autocomplete only) |
| FieldId | Immutable numeric key per field; never reused, soft-deleted via `Status=0` | `NextFieldId` (auto-increment counter) |
| AvailableFields | Flat list of field names for rule/UI autocomplete, cached 24h in Redis | `DynamicFields` which carries typing and EVA info |
| Script | One of `ScriptValidator`, `ScriptConversion`, `ScriptTest`, `ScriptProcessor` — JS source executed by JS-Executor | — |
| List classification | External "list" entities (manufacturer/model/category) whose names are denormalized here | MongoDB `_id` references in `ManufacturerId`/`ModelId`/`CategoryId` |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| Template scripts write | MinIO bucket key `{orgId|mapexos_public}/{templateId}.json` (object storage, not NATS) | Scripts bundle (inferred from `TemplateStoragePort.WriteScripts`) | JS-Executor (via TieredCache L2) |
| Template cache invalidate | `mapexos.fanout.template.invalidate` (stream `FANOUT`) | `contracts/services/assets/assettemplates/types.go::TemplateInvalidatePayload` | Router, Events, JS-Executor |

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| List name updated | `mapexos.lists.name_updated` (stream `MAPEXOS-LISTS`) | `Event{ListId, ListType, NewName, OrgId}` (`list_name_updated/types.go`) | Core `mapexos` service (lists module) |

Handled `ListType` values: `asset_manufacturer`, `asset_model`, `asset_category` — any other type is acked and ignored. Consumer uses NATS Core connection with queue group `{service}-LIST-NAME-GROUP` and DLQ policy.

## Driving Ports (inbound — who calls this module)
- HTTP `/api/v1/asset_templates` (JWT auth): CRUD + `GET /counter` + `GET /:id/available_fields`
- HTTP `/internal/templates` (API-Key auth): TieredCache fallback endpoints for JS-Executor
- NATS consumer on `mapexos.lists.name_updated` for denormalized-name sync

## Driven Ports (outbound — what this module requires)
- `repositories.AssetTemplateRepository` — MongoDB persistence (`infrastructure/persistence/mongo`)
- `ports.TemplateStoragePort` — MinIO writer for scripts (`infrastructure/storage/minio`)
- `common.AppCache` — Redis cache for `AvailableFields` (24h TTL) and counter (6h TTL)
- `natsModel.Bus` (name `core`) — consumer wiring for list-name sync

## Invariants and Business Rules
- `DynamicField.FieldId` is assigned from `NextFieldId` and is IMMUTABLE — never reused, even after deletion
- Deletion of a DynamicField sets `Status=0` (deprecated); historical events keep resolving their field
- Maximum 200 active (`Status=1`) DynamicFields per template (per entity doc comment)
- System templates (`IsSystem=true`) have `OrgID=nil` and use MinIO prefix `mapexos_public`; private templates use `OrgID` hex
- `AvailableFields` cache is invalidated on create/update
- Classification names (`manufacturerName`/`modelName`/`categoryName`) are denormalized and kept in sync only via the `mapexos.lists.name_updated` consumer
- Counter cache (Redis) invalidated on create/delete

## Known Cross-Context Interactions
- Assets module (same service): reads templates to enrich asset responses with classification + `AssetIDPath`
- JS-Executor service: consumes the MinIO L2 scripts payload via TieredCache, falls back to `/internal/templates` on miss
- Core `mapexos` service (lists module): publishes `mapexos.lists.name_updated` that this module consumes
