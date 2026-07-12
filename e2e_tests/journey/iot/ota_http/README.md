# Journey: OTA firmware update — HTTP (poll)

> **Status: SPEC — not implemented.** This describes the intended OTA journey flow.
> There is no `journey.go` or `journey/suite` registry entry yet (tracked under the
> OTA e2e ticket). Until then this is a design doc, not a runnable journey.

End-to-end coverage of the OTA firmware rollout for an **HTTP** device against the
live stack: the device **polls** the gateway for its job, downloads + verifies +
applies the firmware, reports its progress via `POST /api/v1/ota/status`, and the
plan reaches COMPLETED with the asset migrated to the target template.

This is the HTTP transport ONLY. The MQTT transport is a **separate** journey
(`journey/iot/ota_mqtt/`); the two are never mixed. They share the Layer-1
building blocks under `services/assets/ota/` (that reuse is allowed — only
journey packages may not import each other).

## What this test proves

The full poll path works: firmware upload (presigned PUT) → plan → the HTTP device
polls `GET /api/v1/ota/jobs` through the gateway (data-source auth) → the gateway
relays to the assets internal pending-job endpoint, which mints a fresh presigned
GET → the device downloads, verifies checksum+size, and POSTs
`downloading→downloaded→verified→updating→updated` to `POST /api/v1/ota/status` →
the gateway normalizes each into `ota.status.advisory` → the assets status
consumer advances the execution and, on `updated`, switches the asset's template
source→target → the plan closes COMPLETED. Asserts read the outcome through the
**public API only**.

## Key facts

- Operator actions use `/api/v1/ota/*` only; the sim device is a real HTTP client
  (data-source-authenticated via `?ds=`).
- Firmware bytes never cross the assets MS: presigned **PUT** up, presigned
  **GET** down.
- Execution: `QUEUED → INITIATED → DOWNLOADING → DOWNLOADED → VERIFIED → UPDATING
  → UPDATED`; assets owns `QUEUED`/`INITIATED`, the device drives the rest.
- Plan: `SCHEDULED → IN_PROGRESS → COMPLETED`.
- **HTTP is NOT presence-gated** — the device shows up by polling; no online
  precondition, no push race.

## The flow — ping / pong

Legend: `OP` operator (public API) · `AS` assets MS (`ota`) · `S3` MinIO ·
`GW` http_gateway · `DEV` HTTP device (sim). `⇄` over NATS.

```
# ── operator setup ──
OP → AS   POST /api/v1/ota/firmware/init  {targetTemplateId, version, filename, size, sha256}
AS → OP   200  {firmwareId, uploadUrl}                       # presigned PUT
OP → S3   PUT  uploadUrl  <firmware bytes>   → 200
OP → AS   POST /api/v1/ota/firmware/:firmwareId/complete     → 200 (READY)
OP → AS   POST /api/v1/ota/plans  {firmwareId, sourceTemplateId, assetIds, startAt=now, maxTime, rolloutConfig}
AS → OP   200  {id, status:SCHEDULED, counters.total:1}      # StartAt ⇒ IN_PROGRESS

# ── HTTP poll (device pulls its job; no push, no presence gate) ──
DEV → GW  GET /api/v1/ota/jobs?ds={dataSourceId}&assetUUID={assetUUID}
GW → AS   (internal relay)  GET /jobs?assetUUID={assetUUID}
AS → GW   200  OTAUpdateCommand{ downloadUrl(presigned GET), checksum, size,
              reportTarget:"/api/v1/ota/status" }             # exec QUEUED → INITIATED (on first serve)
GW → DEV  200  <ota_update>
DEV → S3  GET downloadUrl  → 200 <bytes>                       # DEV verifies sha256 + size
DEV → GW  POST /api/v1/ota/status?ds={dataSourceId}  {status:"downloading", progress:20}
GW ⇄ AS   ota.status.advisory  OTAStatusAdvisory              # exec → DOWNLOADING
DEV → GW  … POST "downloaded" → "verified" → "updating" → "updated"
GW ⇄ AS   ota.status.advisory  (one per report)               # → UPDATED (counters.succeeded++, template source→target)
          … all executions terminal ⇒ early-close ⇒ plan COMPLETED
OP → AS   GET /api/v1/ota/plans/:planId             → {status:COMPLETED, counters.succeeded:1}
OP → AS   GET /api/v1/ota/plans/:planId/executions  → {state:UPDATED, percentage:100}
OP → AS   GET /api/v1/assets/:assetId               → {templateId: target}
```

## Ordering (poll — no race)

The HTTP device is never pushed, so there is no subscribe/presence race; it simply
polls after the plan exists:

```
CreateTemplate(source) → CreateTemplate(target) → CreateDataSource
→ CreateConnectivityAsset(http, on source, bound to the data source)
→ InitFirmware → UploadFirmwareBinary → CompleteFirmware → CreatePlan
→ RunHttpOtaDevice → AssertPlanStatus(COMPLETED) → AssertExecutionState(UPDATED,100)
→ AssertAssetTemplateSwitched
```

## Building blocks

Shared Layer-1 blocks under `services/assets/ota/{steps,payloads,asserts}`
(firmware/plan steps, the HTTP device sim, the OTA asserts) + `assettemplates`
(source/target templates) + `http_gateway/datasources` (CreateDataSource) +
`common/journey/iam_bootstrap`. The journey only sequences them.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/ota_http/
```

The `saga` build tag gates the test; `go test ./...` (no tag) skips it.

## Requirements

- Live stack: `mapexos`, `assets`, `http_gateway` on default ports; MinIO + NATS
  reachable.
- Seed admin provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in.
- Determinism: plan created with `startAt=now` + high `ratePerMinute`; the job
  poll + plan poll use generous timeouts to cover one reconciler scan.
