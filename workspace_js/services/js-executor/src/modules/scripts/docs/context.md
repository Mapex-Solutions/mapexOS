# Bounded Context: Scripts

**Service:** js-executor
**Module path:** `src/modules/scripts/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose
Orchestrates IoT event processing end-to-end: resolves an event's asset from cache, fetches the asset's template + scripts, dispatches the batch to the engine, and then publishes the transformed event plus optional debug/heartbeat payloads to NATS. Also exposes HTTP endpoints (public + internal) for testing scripts and generating sample payloads from templates. This is the only module that fans data out to the rest of the platform; `events/` owns ingress and `engine/` owns pure execution.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
| --- | --- | --- |
| AssetScripts | `{ decode, validation, transform }` scripts resolved from a template | Engine's `ScriptSet` (same shape, different layer) |
| AssetReadModel | CQRS read projection of an asset stored in MinIO (L2), hydrated into L0/L1 on demand | Asset entity in the assets service |
| TemplateReadModel | Read projection carrying `scriptProcessor`, `scriptValidator`, `scriptConversion`, `scriptTest` | Template aggregate in the assets service |
| SourceType | `'http' \| 'mqtt' \| 'lorawan'` — origin gateway, drives assetUUID resolution strategy | NATS subject / stream name |
| AssetBind | HTTP-only config describing how to extract assetUUID from the event body (`fixedAssetId`, `uuidField`) | SourceType |
| EventTrackerId | End-to-end correlation ID propagated across services; fallback is `seq-{streamSequence}` | NATS `Nats-Msg-Id` (used for dedup, derived from it) |
| StandardizedPayload | Zod-validated output of the transform phase; contract shared via `@mapexos/schemas` | Raw event payload / decoded payload |
| Debug mode | `asset.debugEnabled === true` — triggers `events.raw` + success `events.logs.jsexecutor` publishes | Health monitoring flag |
| Public org | `mapexos_public` — pseudo-orgId used for system templates (`isSystem=true`) | Real tenant org IDs |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
| --- | --- | --- | --- |
| Route execute (success only) | `mapexos.route.execute` | `{ eventSource: 'assetEvent', assetUUID, assetId, orgId, pathKey, eventTrackerId, dataSource, event }` — `PublishResultParams`, dedup `${eventTrackerId}-route` | router service |
| Raw event (debug) | `mapexos.events.raw` | `{ eventTrackerId, threadId, orgId, pathKey, event, source: {http\|mqtt\|lorawan}_gateway, created, name, description, success:true, error:'' }` — dedup `${eventTrackerId}-raw` | events service (`events_raw` ClickHouse) |
| JS execution log | `mapexos.events.logs.jsexecutor` | `{ eventTrackerId, created, threadId, orgId, pathKey, name, description, execution: { success, failedAt, totalExecutionTime, error }, event }` — dedup `${eventTrackerId}-jslog` | events service (`events_jsexecutor` ClickHouse) |
| Asset heartbeat (implicit) | `mapexos.asset.heartbeat.{orgId}` | `{ orgId, assetUUID, pathKey, ts }` — fire-and-forget, no dedup. ONLY emitted when `asset.healthMonitor.enabled=true` AND `heartbeatMode='implicit'`. Same subject is also published by http_gateway for explicit-mode HTTP-protocol assets; the assets/healthmonitor consumer is origin-agnostic. Explicit-mode MQTT-protocol assets do NOT pass through this subject — their liveness is captured by NATS broker presence advisories and consumed directly by the assets/healthmonitor module. | assets service (healthmonitor) |

All four use core NATS (`publishCore`) and are flushed together in one TCP roundtrip per batch via `EventPublisherPort.flush()` before the consumer ACKs.

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
| --- | --- | --- | --- |
| (none directly) | — | — | — |

This module does not subscribe to NATS. Inbound events arrive through the `events/` module, which calls `handleHttpBatch` / `handleMqttBatch`.

## Driving Ports (inbound — who calls this module)
- `ScriptServicePort` (`application/ports/script_service_port.ts`)
  - `handleHttpBatch(messages)` / `handleMqttBatch(messages)` — batch NATS pipelines called by consumers in `events/`.
  - `executeScripts(message)` — single-message path that resolves + executes + enriches; *(declared on the port but currently not invoked by consumers, which use the batch paths)* (inferred).
  - `scripsTest(payload, scripts)` — exercised by `POST /api/v1/scripts/test`.
  - `getScriptTest(orgId, templateId)` / `getSamplePayload(orgId, templateId)` — exercised by public and internal HTTP routes.
  - `fetchAssetScripts`, `resolveAssetUUID` — exposed for the batch handlers.
- HTTP interfaces (`interfaces/http/`)
  - Public router under `/api/v1/scripts`: `POST /test`, `GET /templates/:templateId/script_test`, `GET /templates/:templateId/sample_payload`.
  - Internal router under `/internal/templates` (API-key guarded): `GET /:orgId/:templateId/script_test`, `GET /:orgId/:templateId/sample_payload`.

## Driven Ports (outbound — what this module requires)
- `ScriptEngineServicePort` (`@modules/engine`) — `runBatch` for NATS pipelines, `runScriptPipeline` for HTTP test/sample paths.
- `AssetCachePort` — `get(orgId/assetUUID)` → `AssetReadModel`; `invalidate` used by the FANOUT adapter in `events/`. Backed by `AssetCacheAdapter` → `TieredCacheClient('AssetCache')`.
- `TemplateCachePort` — `get(orgId/templateId)` → `TemplateReadModel`; same adapter pattern. Backed by `TemplateCacheAdapter` → `TieredCacheClient('TemplateCache')`.
- `EventPublisherPort` — `publishResult`, `publishRawEvent`, `publishExecutionLog`, `publishHeartbeat`, `flush`. Backed by `NatsEventPublisherAdapter` → `NatsBus.publishCore` + `flush`.
- `Logger` (`@mapexos/microservices`).
- `@mapexos/schemas` — `ZodStandardizedPayloadSchema`, `ScriptTest` DTO, `StandardizedPayload` type.
- `@mapexos/validations` — `zodValidationError` for error normalization.
- `@mapexos/utils` — `getByPath` for assetBind path resolution.
- Prometheus `JsExecutorMetrics` — `eventsProcessed`, `eventDuration`, `payloadSize`, `assetCache`, `templateCache` (several are wired through DI but not observed yet) (inferred).

## Invariants and Business Rules
- `orgId` MUST be present in `dataSource` (HTTP) or parsed from the MQTT subject; missing orgId → permanent error.
- HTTP assetUUID resolution: `assetBind.type === 'fixedAssetId'` returns the fixed ID path; `uuidField` tries each configured path in order and uses the first match. MQTT/LoRaWAN require `assetUUID` to be pre-provided.
- An asset with no `assetTemplateId` fails processing permanently.
- Missing transform script (`scriptConversion`) is fatal — decode and validation are optional, transform is not.
- Template org resolution: if `asset.assetTemplateOrgId` is empty, falls back to `PUBLIC_ORG_ID = 'mapexos_public'` for system templates.
- `dataSource.pathKey`, `name`, `description` are overwritten from `AssetReadModel` — the asset cache is the source of truth, not the incoming message.
- Heartbeat publishes are gated by `asset.healthMonitor.enabled === true` AND `heartbeatMode === 'implicit'` (missing/undefined defaults to `'implicit'`). When the gate skips, the metric `heartbeatsSkipped{reason}` is incremented with `reason='disabled'` (gate off) or `reason='explicit_mode'` (delegated to broker presence for MQTT or the dedicated HTTP endpoint for HTTP-protocol assets).
- Debug behavior: `debugEnabled=true` → publish `events.raw` before execution and `events.logs.jsexecutor` on success; failures always publish `events.logs.jsexecutor` regardless of debug flag.
- Only the final transform phase's output is validated against `ZodStandardizedPayloadSchema` (in `executeScripts`; batch handler validates structurally via success flag only).
- `OOMError` thrown from the engine is surfaced as `BatchMessageResult.isOOM` so consumers NACK for retry; other errors return `success:false` with the sanitized message.
- `eventTrackerId` is assigned from `seq-{streamSequence}` when absent, guaranteeing a correlation ID downstream.
- All per-batch publishes are queued and flushed exactly once (`eventPublisher.flush()`) before results return — this is what makes the consumer's ACK safe even though individual publishes are fire-and-forget.

## Known Cross-Context Interactions
- Engine module: in-process dispatch target (`ScriptEngineServicePort.runBatch` / `runScriptPipeline`); `BytecodeCacheContext` is built here from the loaded template.
- Events module: inbound adapter that calls this module's batch handlers and bridges results to NATS ack/nack/reject.
- Assets service (external): upstream producer of both the `AssetReadModel` / `TemplateReadModel` entries in MinIO L2 and the `mapexos.fanout.asset.invalidate` / `mapexos.fanout.template.invalidate` events consumed by `events/` to keep L0/L1 fresh.
- Router service (external): sole consumer of `mapexos.route.execute`. The `eventSource: 'assetEvent'` discriminator tells the router to use `asset.RouteGroupIds` (vs. healthmonitor's `'healthStatus'`).
- Events service (external): ClickHouse ingest for `events.raw` and `events.logs.jsexecutor`.
- Healthmonitor submodule of assets service: consumer of `mapexos.asset.heartbeat.{orgId}`; absence of heartbeats after `thresholdMinutes` × `requiredMisses` marks the asset stale.
- http_gateway service: upstream producer of the HTTP datasource messages ultimately routed here via NATS.
- Plugin marketplace / Rule Test Runner UI: consumers of the HTTP `sample_payload` and `script_test` endpoints to hydrate editor UIs.
