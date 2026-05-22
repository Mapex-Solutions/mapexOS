# Assets Service Documentation

## Overview
The Assets service is the **source of truth** for Assets and Asset Templates. It writes read-models to **L2 (MinIO/S3)** for distributed caching and exposes a single internal fallback endpoint (`GET /internal/assets/:assetUUID`) that repopulates L2 on cache misses, enabling CQRS-style reads across the platform with low latency.

It also owns the device-side MQTT control plane: per-asset X.509 certificate lifecycle (`mqttcerts`) and health monitoring driven by broker presence advisories + heartbeats (`healthmonitor`). MQTT auth decisions (bcrypt-password and cert-serial-match) happen entirely INSIDE the mapex-mqtt-broker plugin off the AssetReadModel — this service does not run an auth callout endpoint.

## Responsibilities
- CRUD for Assets and Asset Templates
- EVA field configuration and template scripting
- Publish a single L2 read-model (`AssetReadModel`) carrying everything every consumer needs — including `Protocol.Mqtt.PasswordHash` and `CurrentCert.Serial` consumed by the broker plugin for local CONNECT decisions
- Issue and revoke MQTT device certificates (`mqttcerts` module) using a RAM-cached intermediate CA pulled from mapexVault on startup
- Ingest device heartbeats (HTTP + NATS) and broker presence advisories to drive online/offline state (`healthmonitor` module)

## Non-Responsibilities
- Event ingestion (HTTP Gateway)
- Rule evaluation (RuleEngine)
- Event persistence (Events service)
- MQTT CONNECT auth decisions (the mapex-mqtt-broker plugin owns this end-to-end via its TieredCache; this service is just the read-model writer)
- PKI material storage (mapexVault — stateless CA storage; this service holds the intermediate CA only in RAM)

## Primary Data Flow
1. Asset/Template created or updated in MongoDB
2. Service writes `AssetReadModel` to L2 (MinIO) — single payload, all fields
3. Fanout invalidation notifies every consumer (Router / JS-Executor / Events / mapex-mqtt-broker plugin) to drop their local L1 entry
4. On next access, consumer reads L2 (or falls back to `GET /internal/assets/:assetUUID` which repopulates L2 on the way out)

## CQRS + Distributed Cache
- **Write model**: MongoDB in Assets service
- **Read model**: a single `AssetReadModel` written to L2 MinIO/S3 (bucket `mapex-assets`, key `{orgId}/{assetUUID}.json`). It carries every field every consumer needs — including `Protocol.Mqtt.PasswordHash` and `CurrentCert.Serial` so the broker plugin decides MQTT CONNECTs locally without a callout.
- **L3 fallback**: `GET /internal/assets/:assetUUID` returns the same shape and repopulates L2 inline on the way out, so the next reader hits the cache.

## MQTT Control Plane
- **Issuance** — `mqttcerts` signs device certs locally using the intermediate CA cached in RAM via `common.Mountable.OnMount` (pulled from mapexVault `pki/` at startup, exponential-backoff retry on failure). Device private key never persists server-side — it is downloaded once at issuance time.
- **Revocation** — HARD delete on revoke + 30-day TTL audit row in `mqttRevokedCertificates` (Mongo TTL on `revokedAt`).
- **Auth decisions** — the mapex-mqtt-broker plugin reads the `AssetReadModel` from its TieredCache (L1 Pebble → L2 MinIO → L3 HTTP) and decides every CONNECT locally: bcrypt-compare `Protocol.Mqtt.PasswordHash` for password mode, equality-check `CurrentCert.Serial` against the device's cert serial for cert mode. There is NO HTTP auth callout — the only HTTP path between broker and assets is the L3 read-model fallback.
- **Presence** — broker plugin publishes connect/disconnect advisories on `${env}.mapexos.mqtt.presence.advisory` (single subject, `Event` field discriminates). `healthmonitor` consumes both with separate durables and updates the asset's online state.

## Docs Map
- [Architecture](architecture/index.md)
- [Endpoints](endpoints/index.md)
- [Configuration](configuration/index.md)
- [Operations](operations/index.md)
- [Observability](observability/index.md)
- [Tests](tests/index.md)
- [Benchmarks](benchmarks/index.md)
