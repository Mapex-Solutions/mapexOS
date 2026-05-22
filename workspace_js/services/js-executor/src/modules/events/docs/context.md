# Bounded Context: Events

**Service:** js-executor
**Module path:** `src/modules/events/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose
Interface-layer module that owns all NATS JetStream subscriptions for js-executor. Runs four consumers: two queue consumers that ingest raw device events (HTTP datasource payloads and MQTT telemetry) and fan them into `ScriptService` batch handlers, and two FANOUT consumers that react to asset and template cache-invalidation events so every replica drops stale L0/L1 entries. Translates the NATS message lifecycle (ack / nack / reject / DLQ) into script-service outcomes; it does not decode or execute payloads itself.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
| --- | --- | --- |
| Queue consumer | JetStream durable consumer with a queue group for load-balanced batch pull | FANOUT consumer (broadcast, ephemeral) |
| FANOUT consumer | Ephemeral per-replica subscription on the `FANOUT` stream for cache invalidation | JetStream durable queue consumer |
| Batch | Array of `Message` pulled via `batchMessageHandlerV2`; results are returned index-aligned for ACK/Nack | Piscina batch inside the engine |
| Permanent error | Parse / schema failure that must never retry → `msg.reject(...)` + DLQ | Transient error → `msg.nack()` for redelivery |
| OOM result | `BatchMessageResult.isOOM === true` — engine ran out of V8 memory, NACK for retry | Permanent error |
| DLQ policy | `{ serviceName, serviceType, eventType }` appended on terminal failures by the NatsBus | NATS stream retention policy |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
| --- | --- | --- | --- |
| (none directly) | — | — | — |

This module does not publish. Downstream publishes (`mapexos.route.execute`, `mapexos.events.raw`, `mapexos.events.logs.jsexecutor`, `mapexos.asset.heartbeat.{orgId}`) are emitted from the `scripts` module's `NatsEventPublisherAdapter` during the batch pipeline.

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
| --- | --- | --- | --- |
| JS execute (HTTP datasource) | `mapexos.processor.js.execute` (stream `PROCESSOR-JS-EXECUTE`, durable `processor-js-execute`) | `ScriptProcessorMessage` JSON — `scripts/application/types/message.types.ts` | http_gateway service |
| MQTT telemetry | `mapexos.mqtt.data.>` (stream `MQTT-DATA`, durable `mqtt-data-processor`) | Raw device JSON; orgId + assetUUID parsed from subject `mapexos.mqtt.data.{orgId}.{assetUUID}.>` | NATS MQTT leaf → core republish |
| Asset cache invalidate | `mapexos.fanout.asset.invalidate` (stream `FANOUT`, ephemeral) | `{ orgId, assetUUID }` | assets service |
| Template cache invalidate | `mapexos.fanout.template.invalidate` (stream `FANOUT`, ephemeral) | `{ orgId, templateId }` | assets service (templates) |

`FANOUT` stream is ensured at startup with `maxAge=5min`, `maxMsgs=10_000`, subjects `mapexos.fanout.>`.

## Driving Ports (inbound — who calls this module)
None. `module.ts::initListeners()` is invoked by the bootstrap phase 4; from then on the NATS runtime drives this module.

## Driven Ports (outbound — what this module requires)
- `NatsBus` (`@mapexos/infrastructure`) — `startConsumer`, `subscribeFanout`, `ensureFanoutStream`.
- `ScriptServicePort` (`@modules/scripts/application/ports`) — `handleHttpBatch(messages)` and `handleMqttBatch(messages)` return per-message `BatchMessageResult[]` that drive ack/nack/reject decisions.
- `AssetCachePort` / `TemplateCachePort` (`@modules/scripts/application/ports`) — `invalidate(key)` called from FANOUT consumers.
- `ConfigModule` (`@mapexos/microservices`) — `service_name` + tuning knobs resolved via `resolveConsumerConfig` / `resolveAllTuning` (batch size, fetch timeout, maxAckPending).
- `Logger` (`@mapexos/microservices`).
- Prometheus `JsExecutorMetrics` — observes `batchSize`, increments `eventsProcessed{consumer,status}`; `payloadSize` is wired but unused (inferred).
- Shared constants: `SERVICE_NAME`, `SERVICE_TYPE`, `DEFAULT_RETRY_POLICY`.

## Invariants and Business Rules
- Consumers are `void`-launched (fire-and-forget) from `initListeners`; they must register themselves with NatsBus but must not block module init.
- For each message in a queue batch, exactly one of `msg.ack()`, `msg.reject(reason)`, or `msg.nack(error)` is called based on `BatchMessageResult`:
  - `success` → `ack`, status=`success`.
  - `isPermanent` → `reject` (sent to DLQ per `dlqPolicy`), status=`rejected`.
  - `isOOM` → `nack` for retry, status=`failure`.
  - otherwise → `nack`, status=`failure`.
- The service returns results BEFORE acks happen, which means publishes in `NatsEventPublisherAdapter` must have flushed — ACK safety depends on that flush ordering.
- FANOUT consumers must be idempotent and best-effort: invalid payloads (missing orgId / assetUUID / templateId) are logged and dropped silently; parsing errors are caught and never crash the handler.
- MQTT subject MUST match `mapexos.mqtt.data.{orgId}.{assetUUID}.*` — otherwise parse fails upstream and the HTTP/MQTT batch handler flags the message as permanent.
- FANOUT stream is created by this module if absent; other services may also declare it (idempotent).
- Durable names and stream names are constants — changing them is a breaking migration, not a config tweak.

## Known Cross-Context Interactions
- Scripts module: consumes its `ScriptServicePort`, `AssetCachePort`, `TemplateCachePort`. This module lives in `interfaces/` of js-executor and depends on `application/` of scripts — that's the DDD inbound-adapter direction.
- Engine module: indirectly, via `ScriptService` → `ScriptEngineServicePort`. OOM propagation round-trip is: engine worker → `OOMError` → `ScriptService` batch handler → `BatchMessageResult.isOOM` → this module's consumer → `msg.nack`.
- http_gateway service: upstream publisher of `mapexos.processor.js.execute`.
- NATS MQTT leaf (external): upstream source of `mapexos.mqtt.data.>` (republished from `dt/{orgId}/{assetUUID}/telemetry`).
- assets service: upstream publisher of both FANOUT invalidation subjects after mutating its read models in MinIO (L2).
- Router service: indirect downstream — receives `mapexos.route.execute` published by the scripts module after this module's consumers complete.
