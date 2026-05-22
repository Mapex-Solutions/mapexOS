# Bounded Context: MQTT Certificates

**Service:** assets
**Module path:** `src/modules/mqttcerts/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-11

## Purpose

Owns the device-cert lifecycle. On boot, fetches the intermediate CA from mapexVault (via HTTP) and caches it in RAM (`atomic.Pointer`) for high-throughput local signing. Issues device certs with `crypto/x509` (ECDSA P-256, 128-bit random serial). Tracks revoked cert metadata for 30 days (TTL on `revokedAt`); long-term audit migrates to ClickHouse in a future ticket. Cert/key PEM bytes are returned to the caller exactly once at issue time and NEVER persisted server-side.

## Ubiquitous Language

| Term | Meaning | Not to be confused with |
|---|---|---|
| Active cert | `asset.currentCert` subdoc (single, may be nil) | mqttRevokedCertificates row |
| Revoked cert | row in `mqttRevokedCertificates` (TTL 30d) | OAuth token revocation |
| Force flag | `IssueCertRequest.Force=true` accepts replacing an existing active cert | k8s force-delete |
| caReady flag | `MqttCertsService.IsCAReady()` — flips true after the RAM CA is loaded | k8s readiness probe |

## Published Events

None. Cert state changes ride on `mapexos.fanout.asset.invalidate` (owned by assets/assets).

## Consumed Events

None. All driving inputs are HTTP.

## Driving Ports

- HTTP external (JWT + permissions):
  - `POST /api/v1/mqtt_certs` (`MqttCertCreate`) — issue
  - `DELETE /api/v1/mqtt_certs/:serial` (`MqttCertRevoke`) — revoke
  - `GET /api/v1/mqtt_certs?assetUUID=` (`MqttCertRead`) — list revoked
- Lifecycle hook `OnMount` — kicks off RAM CA bootstrap (sync attempt + retry goroutine on failure).

## Driven Ports

- `RevokedRepository` (Mongo) — `mqttRevokedCertificates` with TTL 30d.
- `MapexVaultClientPort` — HTTP client for `GET /internal/pki/intermediate_ca_bundle`.
- `X509SignerPort` — local crypto signer (ECDSA P-256, 128-bit random serial, `ExtKeyUsageClientAuth`).
- `CAStorePort` — in-RAM atomic.Pointer holder; hot-swap ready.
- `AssetRepository` (cross-module via ports) — `UpdateCurrentCert`, `FindByCurrentCertSerial`.
- `AssetStoragePort` — existing asset module port that rewrites the `AssetReadModel` to MinIO so the broker plugin's L2 picks up the updated `CurrentCert.Serial` (or empty after revoke) on the next fanout-driven refresh.
- `FanoutPublisherPort` — existing asset module port for `mapexos.fanout.asset.invalidate`.

## Invariants and Business Rules

- Issue requires `caReady=true`; otherwise 502.
- Issue on asset with existing `currentCert` requires `force=true` (else 409).
- Cert/key PEM never persisted server-side; returned 1x in the issue response.
- Revoke MOVES the cert to `mqttRevokedCertificates`; the asset's `currentCert` becomes nil.
- Asset deletion is HARD DELETE: drops `currentCert` with the asset, drops `mqttRevokedCertificates` rows for the `assetUUID`, drops the L2 `mapex-assets/{orgId}/{assetUUID}.json` payload, publishes fanout. NO cascade revoke, NO audit retention.
- `mqttRevokedCertificates` has Mongo TTL 30d on `revokedAt`.

## Known Cross-Context Interactions

- `assets/assets` — writes asset.currentCert subdoc; HARD DELETE flow calls into `MqttCertsService.HardDeleteByAssetUUID`.
- `mapexVault/pki` — boot fetch of the intermediate CA bundle.
- `mapex-mqtt-broker` plugin — consumes `AssetReadModel.CurrentCert.Serial` (projected into the plugin's local `AuthEntry.CurrentCertSerial`) via its TieredCache (L1 Pebble → L2 MinIO → L3 GET /internal/assets/:uuid). Fanout invalidation evicts L1 on every issue/revoke. Plugin never talks to this module directly.
