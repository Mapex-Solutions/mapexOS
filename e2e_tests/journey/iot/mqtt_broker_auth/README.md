# Journey: MQTT Broker Auth — password + cert lifecycles

End-to-end coverage of the platform's MQTT authentication pipeline
against the live `services_required` stack. Each phase runs a complete
asset lifecycle — **create → connect → presence online → publish →
events round-trip → delete → reconnect must be denied** — proving both
the happy path and the FANOUT-driven cache invalidation in the broker
plugin's L1 (Pebble) tier.

## Wire contract (post TKT-2026-0040)

- MQTT username = bare `assetUUID` (no legacy `{orgId}:` prefix)
- Device topic shape = `events/{assetUUID}/{eventType}`
- Cert Subject.CN = bare `assetUUID`; tenant scoping flows from
  `asset.protocol.mqtt.orgId` server-side
- Listeners: `tcp://localhost:1883` (password) and
  `ssl://localhost:8883` (mTLS)

## Phases

| Phase                  | What it covers                                                                |
|------------------------|-------------------------------------------------------------------------------|
| `phase0_iam_bootstrap` | Seed admin login → JWT validity → org-context coverage                        |
| `phase1_password_user` | Password auth full lifecycle (10 steps)                                       |
| `phase2_cert_user`     | Cert (mTLS) auth full lifecycle (11 steps)                                    |
| `phase3_cascade`       | _(skeleton — out of scope for this ticket)_                                   |

### `phase1_password_user` — step ordering

1. **CreateRouteGroup** — route group for the asset
2. **CreateTemplate** — asset template (temperature)
3. **CreateAsset** — asset persisted with `authType=password`
4. **ConnectMqttPassword** — MQTT CONNECT with `(assetUUID, password)`
5. **AssertHealthStatusEventually(online)** — presence flowed end-to-end
6. **PublishTelemetry** — publish on `events/{assetUUID}/temperature`
7. **AssertRawEventReceivedAfter** — events service surfaces the row
8. **DisconnectMqtt** — clean MQTT teardown
9. **DeleteAsset** — fanout invalidation fires
10. **AssertConnectDeniedPassword** — fresh CONNECT now denied

### `phase2_cert_user` — step ordering

1. **CreateRouteGroup**
2. **CreateTemplate**
3. **CreateAssetWith(SagaMqttCertTemperatureSensor)** — `authType=cert`, `certTTL={1 day}`
4. **IssueCert** — POST `/api/v1/mqtt_certs`; PEM bundle on bag, asset.currentCert persisted
5. **ConnectMqttCert** — mTLS handshake against `:8883`
6. **AssertHealthStatusEventually(online)**
7. **PublishTelemetry**
8. **AssertRawEventReceivedAfter**
9. **DisconnectMqtt**
10. **DeleteAsset**
11. **AssertConnectDeniedCert** — fresh mTLS CONNECT now denied

## Run

```sh
cd e2e_tests

# Phase 0 only (smoke the IAM bootstrap)
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap/

# Password lifecycle
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase1_password_user/

# Cert (mTLS) lifecycle
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase2_cert_user/

# Everything in the journey
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/...
```

`-count=1` disables Go's test cache so each invocation hits the live
stack. Without it a second run prints `(cached)` and reports the prior
verdict without actually contacting the broker.

## Prerequisites

1. Pre-build PKI:
   ```sh
   ./scripts/prebuild/pki/generate-pki.sh
   ```
2. Stack up (either `services_required` + `mapex_services` compose stacks
   OR `standalone/`):
   ```sh
   ./scripts/mapex-deploy.sh --full
   ```
3. mapexVault healthy + assets MS `caReady=true` (the cert phase fails
   at `IssueCert` otherwise)
4. Broker MQTT listeners reachable on `:1883` and `:8883`

### Environment overrides

| Variable                | Default                       | Notes                                |
|-------------------------|-------------------------------|--------------------------------------|
| `MAPEXOS_URL`           | `http://localhost:5000`       | mapexIam base URL                    |
| `ASSETS_URL`            | `http://localhost:5002`       | assets service base URL              |
| `ROUTER_URL`            | `http://localhost:5003`       | router service base URL              |
| `GATEWAY_URL`           | `http://localhost:5001`       | http_gateway service base URL        |
| `EVENTS_URL`            | `http://localhost:5004`       | events service base URL              |
| `MQTT_BROKER_URL`       | `tcp://localhost:1883`        | password listener                    |
| `MQTT_BROKER_TLS_URL`   | `ssl://localhost:8883`        | mTLS listener                        |

## Failure diagnosis

Every step + assert publishes a qualified name to the saga log
(e.g. `assets/assets.ConnectMqttCert`,
`events/events.AssertRawEventReceivedAfter`). When a phase fails:

1. Locate the failing item by name in the test output.
2. Open the matching file under `services/<service>/<module>/{steps,asserts}/`.
3. Read the comment above the failing item — each function documents
   what it reads from the bag, what it writes, and the production
   contract it exercises.
4. Tail the right service log (`docker logs mapex-assets`,
   `docker logs mapex-broker-mqtt`, ...) for the matching window.

## Bag keys touched by this journey

| Key                                | Phase | Writer                          |
|------------------------------------|-------|---------------------------------|
| `iam.userJWT`                      | 0     | `authSteps.SeedAdminLogin`      |
| `iam.organizationID`               | 0     | `authSteps.SeedAdminLogin`      |
| `router.routeGroupID`              | 1, 2  | `rgSteps.CreateRouteGroup`      |
| `assets.assetTemplateID`           | 1, 2  | `templateSteps.CreateTemplate`  |
| `assets.assetID`                   | 1, 2  | `assetSteps.CreateAsset`        |
| `assets.assetUUID`                 | 1, 2  | `assetSteps.CreateAsset`        |
| `assets.assetMqttPassword`         | 1     | `assetSteps.CreateAsset`        |
| `assets.assetCertPEM` (+ key, ca)  | 2     | `assetSteps.IssueCert`          |
| `assets.assetCertSerial`           | 2     | `assetSteps.IssueCert`          |
| `assets.mqttClient`                | 1, 2  | `Connect*` step                 |
| `assets.mqttConnectedAt`           | 1, 2  | `Connect*` step                 |
| `assets.telemetrySentAt`           | 1, 2  | `assetSteps.PublishTelemetry`   |
| `assets.assetDeleted`              | 1, 2  | `assetSteps.DeleteAsset`        |
