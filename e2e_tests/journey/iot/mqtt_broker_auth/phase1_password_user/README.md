# Phase 1 — MQTT password auth full lifecycle

## What this test proves

The MQTT password auth pipeline works end-to-end against the live
broker. The phase drives the full operator-by-hand workflow for an
asset using `authType=password`: create → connect → presence flowed →
publish → events round-trip → delete → reconnect must be denied. The
deny-probe proves the FANOUT-driven L1 invalidation in the broker
plugin.

The phase:

1. Creates a route group for the asset.
2. Creates an asset template (temperature schema).
3. Creates the asset with `protocol=mqtt` + `authType=password`; the plaintext password lands on the bag.
4. MQTT CONNECT against the password listener with `(assetUUID, password)`.
5. Asserts `healthStatus=online` (presence advisory consumed, status persisted, read model surfaces it).
6. PUBLISH on `events/{assetUUID}/temperature`; broker ACL must accept the bare-assetUUID topic.
7. Asserts the events service surfaces the row via `/api/v1/events/raw`.
8. Clean MQTT disconnect.
9. DELETE the asset; FANOUT invalidation reaches the broker plugin.
10. Asserts a fresh CONNECT with the same credentials is denied.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/mqtt_broker_auth/phase1_password_user/
```

## Requirements

- Live stack: `mapexos`, `assets`, `router`, `events` on default ports.
- MQTT broker reachable on `tcp://localhost:1883` (`MQTT_BROKER_URL`).
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
- `mqtt_broker_auth/README.md` lists the full per-service URL overrides.
