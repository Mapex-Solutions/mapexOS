# Phase 2 — MQTT cert (mTLS) auth full lifecycle

## What this test proves

The MQTT cert (mTLS) auth pipeline works end-to-end against the live
broker. The flow mirrors the password phase but the device presents
an issued cert on the wire against the broker's mTLS listener.

The phase:

1. Creates a route group for the asset.
2. Creates an asset template (temperature schema).
3. Creates the asset with `protocol=mqtt` + `authType=cert` using the cert payload variant (carries `certTTL`).
4. Issues a cert via `POST /api/v1/mqtt_certs`; the PEM bundle lands on the bag and `asset.currentCert` is persisted; FANOUT invalidation fires.
5. mTLS CONNECT against the broker's `:8883` listener with the freshly issued cert.
6. Asserts `healthStatus=online` (presence advisory consumed, status persisted).
7. PUBLISH on `events/{assetUUID}/temperature`; broker ACL must accept the bare-assetUUID topic.
8. Asserts the events service surfaces the row via `/api/v1/events/raw`.
9. Clean MQTT disconnect.
10. DELETE the asset; FANOUT invalidation reaches the broker plugin.
11. Asserts a fresh mTLS CONNECT with the SAME cert is denied — asset gone, `currentCert` serial gone, cert-mode auth fails.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/mqtt_broker_auth/phase2_cert_user/
```

## Requirements

- Live stack: `mapexos`, `assets`, `router`, `events` on default ports, plus `mapexVault` healthy and assets MS reporting `caReady=true`.
- MQTT broker reachable on `ssl://localhost:8883` (`MQTT_BROKER_TLS_URL`).
- PKI prebuilt: `./scripts/prebuild/pki/generate-pki.sh`.
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
- `mqtt_broker_auth/README.md` lists the full per-service URL overrides.
