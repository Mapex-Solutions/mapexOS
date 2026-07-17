# Journey: OTA firmware update — MQTT (push)

End-to-end coverage of the OTA firmware rollout for an **MQTT** device against the
live stack: the broker **pushes** the `ota_update` command to a connected,
presence-gated device, the device downloads + verifies + applies the firmware,
reports its progress on `events/{assetUUID}/ota_status`, and the plan reaches
COMPLETED with the asset migrated to the target template.

This is the MQTT transport ONLY. The HTTP transport is a **separate** journey
(`journey/iot/ota_http/`); the two are never mixed. They share the Layer-1
building blocks under `services/assets/ota/` (that reuse is allowed — only
journey packages may not import each other).

## What this test proves

The full push path works: firmware upload (presigned PUT) → plan → the reconciler
dispatches the command over `mqtt.downlink` to the broker **only when the device
is online** → the broker delivers it on `commands/{assetUUID}/ota_update` → the
device downloads via a presigned GET, verifies checksum+size, and reports
`downloading→downloaded→verified→updating→updated` → the assets status consumer
advances the execution and, on `updated`, switches the asset's template
source→target → the plan closes COMPLETED. Asserts read the outcome through the
**public API only**.

## Key facts

- Operator actions use `/api/v1/ota/*` only; the sim device is a real MQTT client.
- Firmware bytes never cross the assets MS: presigned **PUT** up, presigned
  **GET** down.
- Execution: `QUEUED → INITIATED → DOWNLOADING → DOWNLOADED → VERIFIED → UPDATING
  → UPDATED`; assets owns `QUEUED`/`INITIATED`, the device drives the rest.
- Plan: `SCHEDULED → IN_PROGRESS → COMPLETED`.
- **MQTT is presence-gated** — the device must be `online` (and subscribed to the
  command topic) before the plan dispatches, or the pushed command is missed.

## The flow — ping / pong

Legend: `OP` operator (public API) · `AS` assets MS (`ota`) · `S3` MinIO ·
`BR` MQTT broker · `DEV` MQTT device (sim). `⇄` over NATS.

```
# ── operator setup ──
OP → AS   POST /api/v1/ota/firmware/init  {targetTemplateId, version, filename, size, sha256}
AS → OP   200  {firmwareId, uploadUrl}                       # presigned PUT
OP → S3   PUT  uploadUrl  <firmware bytes>   → 200
OP → AS   POST /api/v1/ota/firmware/:firmwareId/complete     → 200 (READY)
OP → AS   POST /api/v1/ota/plans  {firmwareId, sourceTemplateId, assetIds, startAt=now, maxTime, rolloutConfig}
AS → OP   200  {id, status:SCHEDULED, counters.total:1}      # StartAt ⇒ IN_PROGRESS

# ── MQTT push (device already CONNECTed + subscribed, presence online) ──
AS ⇄ BR   mqtt.downlink  DownlinkEnvelope{ ota_update, OTAUpdateCommand{
              downloadUrl(presigned GET), checksum, size,
              reportTarget:"events/{assetUUID}/ota_status" } }  # exec QUEUED → INITIATED (presence-gated)
BR → DEV  MQTT publish  commands/{assetUUID}/ota_update  <ota_update>
DEV → S3  GET downloadUrl  → 200 <bytes>                        # DEV verifies sha256 + size
DEV → BR  MQTT publish  events/{assetUUID}/ota_status  {status:"downloading", progress:20}
BR ⇄ AS   ota.status.advisory  OTAStatusAdvisory               # exec → DOWNLOADING
DEV → BR  … "downloaded" → "verified" → "updating" → "updated"
BR ⇄ AS   ota.status.advisory  (one per report)                # → UPDATED (counters.succeeded++, template source→target)
          … all executions terminal ⇒ early-close ⇒ plan COMPLETED
OP → AS   GET /api/v1/ota/plans/:planId             → {status:COMPLETED, counters.succeeded:1}
OP → AS   GET /api/v1/ota/plans/:planId/executions  → {state:UPDATED, percentage:100}
OP → AS   GET /api/v1/assets/:assetId               → {templateId: target}
```

## Ordering (the presence + subscribe contract)

The reconciler pushes **once**, to **online** devices only, and advances the
execution to `INITIATED` after dispatch. So the sim MUST connect, subscribe to
`commands/{assetUUID}/ota_update`, and be presence-`online` **before**
`CreatePlan`:

```
CreateTemplate(source) → CreateTemplate(target) → CreateRouteGroup
→ CreateConnectivityAsset(mqtt, on source)
→ ConnectMqttPassword → SubscribeOtaCommand → AssertHealthStatusEventually("online")
→ InitFirmware → UploadFirmwareBinary → CompleteFirmware → CreatePlan
→ RunMqttOtaDevice → AssertPlanStatus(COMPLETED) → AssertExecutionState(UPDATED,100)
→ AssertAssetTemplateSwitched
```

## Building blocks

Shared Layer-1 blocks under `services/assets/ota/{steps,payloads,asserts}`
(firmware/plan steps, the MQTT device sim, the OTA asserts) + `assettemplates`
(source/target templates) + `common/journey/iam_bootstrap`. The journey only
sequences them (`Items()` + `Run()`).

## How to run

```bash
cd e2e_tests
go test -tags=saga -run 'TestSuite/iot/ota_mqtt'
```

The journey runs through the single suite runner (`journey/suite`); the `saga`
build tag gates it, and `go test ./...` (no tag) skips it.

## Requirements

- Live stack: `mapexos`, `assets` on default ports; MinIO + NATS reachable.
- `mapexMQTTBroker` running — listener `tcp://localhost:1883` (password).
- Seed admin provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in.
- Determinism: plan created with `startAt=now` + high `ratePerMinute` (dispatch on
  first scan); the command-wait + plan poll use generous timeouts to cover one
  reconciler scan.
