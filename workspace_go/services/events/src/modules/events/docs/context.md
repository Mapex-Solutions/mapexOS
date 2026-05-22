# Bounded Context: Events

**Service:** events
**Module path:** `src/modules/events/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose
Terminal ClickHouse sink for every operational event in the platform: processed events (with EVA field resolution), raw ingestion payloads, JS-executor debug logs, DLQ entries, and execution history for router / business rule / trigger / workflow pipelines. Exposes cursor-paginated HTTP read APIs over each stream. Owns no business decisions — this BC is about durable, query-efficient storage and retrieval of things that already happened.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Event (processed) | Row in the `events` table with resolved EVA fields (eva_number/string/bool/date MAPs keyed by `fieldId`) | `RawEvent` (pre-processing payload) or `event_type` string |
| EVA | Entity-Value-Attribute storage pattern using `MAP<UInt16, T>` keyed by `fieldId` from `AssetTemplate.DynamicFields` | The legacy `Array(Tuple(String, T))` layout it replaced |
| EventTrackerId | UUID propagated end-to-end across services for one logical event | ClickHouse row id / batch index |
| AssetTemplateOrgId | Owner org of the template, `"mapexos_public"` or an orgId — used for cache key `{templateOrgId}/{templateId}` | `OrgId` (tenant that owns the event) |
| TemplateCachePort | Tiered-cache port (L0/L1/L2 + Fallback HTTP) for `CachedTemplate` used during EVA resolution | Generic HTTP client for assets service |
| DLQ | Events that failed processing elsewhere, captured from `MAPEXOS-DLQ` for analysis | NATS retry / Nack — DLQ is the terminal destination |
| RetentionDays | Per-row TTL value resolved from the retention module and stamped at write time | NATS stream retention / ClickHouse partition TTL |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| — | — | — | — (terminal sink; does not publish) |

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| Processed event (EVA) | `mapexos.events.save` (stream `EVENTS`) | `domain/entities.Event` | router |
| Raw ingestion payload | `mapexos.events.raw` (stream `EVENTS-RAW`) | `domain/entities.RawEvent` | http_gateway, MQTT gateway (inferred) |
| JS executor debug | `mapexos.events.logs.jsexecutor` (stream `EVENTS-JSEXEC`) | `domain/entities.JsExecEvent` | js-executor |
| Router execution history | `mapexos.events.router` (stream `EVENTS-ROUTER`) | `domain/entities.RouterEvent` | router |
| Business rule execution history | `mapexos.events.businessrule` (stream `EVENTS-BUSINESSRULE`) | `domain/entities.BusinessRuleEvent` | ruleengine (inferred) |
| Trigger execution history | `mapexos.events.trigger` (stream `EVENTS-TRIGGER`) | `domain/entities.TriggerEvent` | triggers |
| Workflow execution history | `mapexos.events.workflow` (stream `EVENTS-WORKFLOW`) | `domain/entities.WorkflowEvent` | workflow |
| Dead Letter Queue | `mapexos.dlq` (stream `MAPEXOS-DLQ`) | `domain/entities.DLQEvent` | any service's DLQ policy |
| TemplateInvalidate | `mapexos.fanout.template.invalidate` (stream `FANOUT`, ephemeral) | `contracts/services/assets/assettemplates/types.go::TemplateInvalidatePayload` | assets |

## Driving Ports (inbound)
- 8 NATS batch consumers (see table above). All use `DefaultRetryPolicy` and DLQ policy except the DLQ consumer itself, which ACKs unconditionally to avoid redelivery loops.
- HTTP under `/api/v1/events`, all protected by `AuthMiddleware` + `InjectRequestContext` + per-route permission:
  - `GET /raw`, `GET /jsexec`, `GET /router`, `GET /businessrule`, `GET /trigger`, `GET /workflow`, `GET /dlq`, `GET /dlq/counts` — cursor-paginated lists.
  - `GET /workflow/execution/:executionId` — single workflow event by Mongo hex id.
  - `POST /store/query` — processed-event list with optional EVA filters (POST to carry `EvaFilters` in body).
  - `GET /store/:eventTrackerId` — single processed event with resolved EVA field names in `advancedSearch`.

## Driven Ports (outbound)
- `domain/repositories.EventRepository` + per-stream repos (raw, jsexec, router, businessrule, trigger, workflow, dlq, eventstore) — implemented in `infrastructure/persistence/clickhouse`.
- `application/ports.TemplateCachePort` — implemented by `infrastructure/cache/tieredcache`, wraps the DI-injected `TieredCache` named `"templates"` (L0/L1/L2 + fallback HTTP to assets service).
- `domain/services.eva_mapper` — pure mapping helper from incoming fields to EVA MAPs using the template's `DynamicFields`.

## Invariants and Business Rules
- Batch processors handle Ack/Nack/Reject per message and always return `nil`: invalid JSON → `Reject` (DLQ); bulk-insert failure → `Nack` (retry with backoff); success → `Ack`.
- The DLQ consumer never Nacks (would loop); it ACKs even on storage failure.
- EVA resolution uses `fieldId` (uint16) as the map key — never field name. Keys come from `AssetTemplate.DynamicFields` resolved via `TemplateCachePort`.
- Human-readable fields (`AssetName`, `TemplateName`, etc.) are denormalized at write-time by the router, not re-joined at read-time.
- `RetentionDays` is stamped per row on write from the retention module and cached 24h; it drives ClickHouse TTL eviction per-org.
- `GetEventStoreDetail` resolves `DynamicFields` from different sources depending on `source`: `"asset"` → `AssetTemplate` via cache; `"rule"` → `BusinessRule` (future).

## Known Cross-Context Interactions
- Consumes from **router**, **triggers**, **workflow**, **ruleengine** (inferred), **js-executor**, **http_gateway**, and any service's DLQ overflow.
- Depends on **retention** to source per-org TTLs stamped into `RetentionDays`.
- Depends on **assets** (AssetTemplate source) via `TemplateCachePort` fallback HTTP for EVA field resolution.
- Read path feeds the frontend's events / DLQ / workflow-execution explorer pages.
- Consumes `mapexos.fanout.template.invalidate` (FANOUT) to invalidate the local TieredCache (L0+L1) when assets edits/renames an AssetTemplate. Without this, EVA fieldId mappings would drift from the assets-side template.
