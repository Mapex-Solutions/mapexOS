# Tests

## Prerequisites
- Go 1.25+
- No running infrastructure required for unit tests (use mocks)
- MongoDB, Redis, NATS, MinIO, and mapexVault for integration tests

## Run
```bash
go test ./... -count=1
```

## Test Files

| File | Module | Description |
|---|---|---|
| `src/modules/assets/application/services/asset_service_test.go` | `assets` | AssetService unit tests (CRUD, L2 read-model carrying `PasswordHash` + `CurrentCert`, fanout, HARD delete cascade to cert state) |
| `src/modules/assettemplates/application/services/assettemplate_service_test.go` | `assettemplates` | AssetTemplateService unit tests (CRUD, L2, fanout, list name sync) |
| `src/modules/mqttcerts/application/services/mqttcerts_service_test.go` | `mqttcerts` | Cert issuance / revocation unit tests + `OnMount` retry behaviour against a fake mapexVault client |
| `src/modules/healthmonitor/application/services/healthmonitor_handler_heartbeat_test.go` | `healthmonitor` | Heartbeat handler unit tests (online flip, threshold, dedup) |
| `src/modules/healthmonitor/infrastructure/messaging/nats/alert_publisher_test.go` | `healthmonitor` | Alert publisher unit tests against a fake `natsModel.Bus` |

## Mocks

- `src/modules/assets/mocks/` — `asset_repository_mock`, `asset_storage_port_mock`, `routegroup_port_mock`, `mqtt_cert_port_mock`
- `src/modules/assettemplates/mocks/` — `assettemplate_repository_mock`, `template_storage_port_mock`
- `src/modules/mqttcerts/mocks/` — `mqttcerts_repository_mock`, `vault_client_mock`, `revoked_cert_repository_mock`
- `src/modules/healthmonitor/mocks/` — `health_repository_mock`, `alert_publisher_mock`, `presence_consumer_mock`
- `src/shared/mocks/` — `nats`, `app_cache_mock`, `tiered_cache_mock`

## End-to-end coverage

Cross-service flows that include this service are covered under `workspace_go/packages/e2eTests/`:

- `journey/iot/mqtt_full_pipeline/` (saga, `-tags=saga`) — IAM bootstrap → asset create → cert issue → broker connect → presence advisory → state flip. See [saga journey doc](../../../../packages/e2eTests/journey/iot/mqtt_full_pipeline/README.md).
- `services/assets/...` — module-level contract tests for CRUD endpoints.
