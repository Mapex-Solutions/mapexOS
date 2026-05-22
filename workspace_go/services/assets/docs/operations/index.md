# Operations

## Start
```bash
./bin/assets
```

## Smoke checks
- `GET /health` — returns 200 healthy + `caReady: true` once `mqttcerts.OnMount` completes.
- `GET /api/v1/assets` with valid auth.
- `POST /api/v1/mqtt_certs` with valid auth + asset payload → expect `201` once the CA is mounted; expect `503 ca_not_ready` if mapexVault has not yet been reached.
- `GET /internal/assets/:assetUUID` with `X-API-Key: $INTERNAL_API_KEY` → `200` with `AssetReadModel` (carries `Protocol.Mqtt.PasswordHash` + `CurrentCert.Serial` for the broker plugin). MQTT CONNECT auth decisions happen inside the broker plugin — there is no separate auth callout endpoint to smoke.

## Dependencies
| Dependency | Required | Purpose |
|---|---|---|
| MongoDB | yes | Asset, template, cert and revoked-cert persistence |
| Redis (app) | yes | Application cache (MQTT permissions, counters) |
| Redis (shared) | yes | Shared cache (permissions, coverage) |
| NATS Core + JetStream | yes | Fanout invalidation, heartbeat / presence consumers, scan scheduling |
| MinIO | yes | L2 read-model storage for assets and templates |
| mapexVault | yes | Stateless CA source. `mqttcerts.OnMount` pulls the intermediate CA once at startup. |
| MapexOS API | yes | Permission and coverage middleware |
| Router service | optional | Route group assignment on asset create/update |
| mapex-mqtt-broker | indirect | Publishes presence advisories consumed by `healthmonitor`; reads `AssetReadModel` via its TieredCache (L1 Pebble → L2 MinIO → L3 `GET /internal/assets/:uuid`) to decide MQTT CONNECTs locally. Not a startup dependency of this service. |

## Startup ordering
1. Standard infra (Mongo, Redis, NATS, MinIO) up.
2. mapexVault up with the PKI module mounted (operator runs `scripts/prebuild/pki/generate-pki.sh` and points mapexVault at the seed dir on first start).
3. Assets service starts. `mqttcerts.OnMount` pulls the intermediate CA; on failure an exponential-backoff retry goroutine keeps trying while the rest of the service stays online. `caReady=true` flips the gauge and unlocks `/api/v1/mqtt_certs/*`.
4. mapex-mqtt-broker can start at any time — it does NOT block on Assets, but issuing new device certs requires `caReady=true`.

## Deploy
The end-to-end deploy is orchestrated by `deployment/docker-compose/scripts/mapex-deploy.sh` (see `deployment/docker-compose/scripts/README.md`). The `--pki` step generates the CA seed; `--infra` brings up Mongo/Redis/NATS/MinIO/broker; `--services` brings up the application services (this one included).

In **local dev**, mapexVault runs outside docker (`make run-service SERVICE=mapexVault` from `workspace_go/`), seeded from `.local/vault-pki-seed/`. This service points at `MAPEX_VAULT_URL=http://localhost:5010`.

## Benchmarks
See [Benchmarks](../benchmarks/index.md) for load-testing methodology, scripts, and results.
