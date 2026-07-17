# Bounded Context: OTA (Over-The-Air firmware updates)

**Service:** assets
**Module path:** `src/modules/ota/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-07-17

## Purpose
Owns the remote firmware update lifecycle for IoT assets. An OTA update is modeled as a **template-to-template migration**: an asset moves from a source `AssetTemplate` to a target template (the version of record), and only switches template once its device reports `UPDATED`. The module manages firmware artifacts (upload to MinIO, checksum, retention), OTA plans (scheduled rollouts with two timers — start and maxTime), and per-device executions (one document per asset, scaling to thousands). It dispatches update commands to devices over MQTT/LoRaWAN (via the edge) or serves them to HTTP pollers, ingests device status advisories, and reconciles executions toward the target template at a paced, presence-gated rate.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Firmware | Immutable binary artifact stored in MinIO with a SHA-256 checksum; keyed `{orgId}/{targetTemplateId}/{firmwareId}.bin` | `AssetTemplate` (the version-of-record schema the firmware targets) |
| FirmwareStatus | `PENDING_UPLOAD → READY → ACTIVE → DEPRECATED/REVOKED/ABANDONED/PURGED` | `PlanStatus` (the rollout lifecycle) |
| OTAPlan | A rollout of one firmware from a source template to a target template, bounded by `StartAt` and `MaxTime` timers | `OTAExecution` (the per-device record) |
| OTAExecution | One asset's progress through a plan; carries state, attempts, and the resolved `AssetUUID`/protocol | `OTAPlan` (holds only aggregate counters) |
| ExecutionState | `QUEUED → INITIATED → DOWNLOADING → DOWNLOADED → VERIFIED → UPDATING → UPDATED`; terminal: `UPDATED`, `FAILED`, `TIMED_OUT` | device-reported status strings (mapped in the status handler) |
| Reconciler | Engine that dispatches the next paced batch of actionable executions to online devices between the two timers | ReconcilerTimers (drives the start/close timers; its `RunScan` sweep is driven by an elected leader, not a self-scheduled timer) |
| Presence gate | A device is dispatched only when the health monitor reports it online | — |
| Abandon check | A per-firmware timer that GCs an artifact never finalized after upload | maxTime close (the per-plan close routine) |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| Edge downlink command | `mqtt.downlink` / `lorawan.downlink` (STATIC; org/asset in payload) | `contracts/services/assets/downlink::DownlinkEnvelope` (commandType `ota_update` → `OTAUpdateCommand`) | mapexMQTTBroker / mapexLNS downlink consumers |
| OTA status history | `events.ota.status` (stream `events`) | `contracts/services/ota/events/ota_status.go::OTAStatusAdvisory` | Events MS (ClickHouse history) |
| Plan/firmware timers | `ota.timer.start` / `ota.timer.close` / `ota.timer.abandon` (STATIC; id in payload + MsgId dedup) | scheduled self-publish (native NATS `@at`) | this module's timers consumer |

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| Device status advisory | `ota.status.advisory` (stream `ota`) | `contracts/services/ota/events/ota_status.go::OTAStatusAdvisory` | HTTP gateway (device status endpoint) / mqtt-broker forwarder |
| Plan/firmware timers | `ota.timer.>` (stream `ota`, `AllowMsgSchedules`) | id in payload | this module (self-scheduled) |

## Driving Ports (inbound — who calls this module)
- HTTP `/api/v1/ota` (JWT auth): firmware upload init/complete, OTA plan CRUD, execution listing
- NATS consumer on `ota.status.advisory` → `StatusHandler` (device status ingestion)
- NATS consumer on `ota.timer.>` → `ReconcilerTimers` + firmware abandon (start/close/abandon)
- In-process leader ticker → `ReconcilerTimers.RunScan` (the paced scan; only the elected pod ticks)

## Driven Ports (outbound — what this module requires)
- `repositories.{Firmware,OTAPlan,OTAExecution}Repository` — MongoDB persistence (`infrastructure/persistence/mongo`)
- `ports.FirmwareStorePort` — MinIO artifact store: presigned PUT/GET, stat, delete (`infrastructure/storage`)
- `ports.OTASchedulerPort` / `ports.FirmwareSchedulerPort` — native NATS `@at` scheduling for the per-plan start/close + firmware abandon timers (`infrastructure/messaging/nats`)
- `natsModel.LeaderElection` over a `natsModel.KeyValueStore` (bucket `LEADER-assets`) — elects the single pod that runs the paced scan; wired at the composition root (`module.go`)
- `ports.EdgeDispatchPort` — publishes downlink commands to the edge (`infrastructure/messaging/nats`)
- `ports.StatusHistoryPublisherPort` — forwards advisories to Events MS (`infrastructure/messaging/nats`)
- `ports.LiveStatePort` — Redis live execution state (`infrastructure/persistence/redis`)
- `ports.PresenceReaderPort` — health-monitor presence (cross-module adapter)
- `ports.AssetReaderPort` / `ports.AssetTemplateSwitcherPort` — asset read + template switch (cross-module adapter)

## Invariants and Business Rules
- Firmware artifacts are immutable: the object key embeds a fresh `firmwareId`, so no upload overwrites another.
- `CompleteUpload` returns success ONLY after the object is confirmed in storage (size + checksum via HeadObject).
- A plan can only be created against a `READY` firmware (`Firmware.CanReference()`).
- An asset switches template ONLY when its device reports `UPDATED`; the switch targets the plan's `TargetTemplateID`.
- Dispatch is presence-gated (online only) and paced at `RatePerMinute`; offline devices stay `QUEUED` for a later tick.
- The paced scan is a network-wide singleton: exactly one pod, elected via a NATS KV lease, runs `RunScan` on an in-process ticker at `ota_scan_interval`; another pod takes over on leader death. It replaces the former self-rescheduling `ota.timer.scan` message (which could be silently dedup-dropped, stalling dispatch).
- HTTP-protocol devices are never pushed — they poll; the reconciler only records the resolved protocol/UUID for them.
- Status advisories are idempotent and ownership-checked: unknown, already-terminal, or non-owned executions are no-ops.
- A plan auto-aborts when the failure rate crosses `AbortThresholdPct` after `AbortMinExecuted` devices have executed.
- Retention uses two timers: the per-plan `maxTime` close (finalizes non-terminal executions + deletes the `.bin`) and the per-firmware abandon check (GCs never-finalized uploads).

## Known Cross-Context Interactions
- HealthMonitor module (same service): presence source for the dispatch gate.
- Assets module (same service): asset read + template switch on `UPDATED`.
- AssetTemplates module (same service): source/target templates are the versions of record.
- Events MS: consumes `events.ota.status` for the ClickHouse status history.
- Edge (mapex-mqtt-broker / LNS): forwards downlink commands to devices and relays device status back as advisories.
