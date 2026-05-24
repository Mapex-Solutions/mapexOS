# Module e2e: http_gateway / datasources

## Scope

Module e2e suite for the `datasources` module of the http_gateway
service
(`workspace_go/services/http_gateway/src/modules/datasources/`).
DataSources describe how an inbound protocol (HTTP, MQTT, …) is
authenticated, rate-limited, and bound to assets at webhook ingest
time. The suite covers the CRUD round-trip plus the listing's filter
matrix (`name`, `enabled`, `mode`, `protocol`, projection, multi
filter) using both a ROOT-token client and an admin client scoped to
the seed organization.

## Endpoints exercised

- `POST /api/v1/data_sources` — create datasource (HTTP+apiKey and
  MQTT+jwt payloads).
- `GET /api/v1/data_sources/{id}` — fetch one datasource by id.
- `PATCH /api/v1/data_sources/{id}` — partial update of metadata.
- `DELETE /api/v1/data_sources/{id}` — delete and confirm 404 on the
  follow-up GET.
- `GET /api/v1/data_sources` — paginated listing with `page` /
  `perPage`, single-field filters (`name`, `enabled`, `mode`,
  `protocol`), `projection`, multi-filter composition, plus separate
  passes using the admin (org-scoped) and root client.

## Test functions

- `TestCreateDataSource`
- `TestGetDataSourceById`
- `TestUpdateDataSource`
- `TestDeleteDataSource`
- `TestListDataSources_BasicPagination`
- `TestListDataSources_FilterByName`
- `TestListDataSources_FilterByEnabled`
- `TestListDataSources_FilterByMode`
- `TestListDataSources_FilterByProtocol`
- `TestListDataSources_MultipleFilters`
- `TestListDataSources_Projection`
- `TestListDataSources_WithOrgContext`
- `TestListDataSources_RootUser`

## Fixtures

| File                          | Scenario                                                                                  |
|-------------------------------|-------------------------------------------------------------------------------------------|
| `create_datasource_http.json` | Pull-mode HTTP datasource with API-Key header auth and `uuidField` asset bind.            |
| `create_datasource_mqtt.json` | Push-mode MQTT datasource (initially disabled) with JWT auth and rate-limit configured.   |
| `update_datasource.json`      | Partial PATCH body changing `name` and `description`.                                     |
| `update_enabled.json`         | Partial PATCH body toggling `enabled` to false (reserved; not invoked by current tests).  |

## How to run

```bash
cd e2e_tests

# Full module suite
go test ./services/http_gateway/datasources -v

# A single test
go test ./services/http_gateway/datasources -v -run TestListDataSources_FilterByName
```

## Outcome on pass

Confirms the datasources module honors its public HTTP contract: full
CRUD round-trip (create → read → patch → delete → 404), pagination
metadata, the full server-side filter matrix on the listing
(case-insensitive `name`, exact `enabled`, exact `mode`, exact
`protocol`, multi-filter composition), the `projection` selector, and
the org-context resolution path — covered with both an admin client
(scoped to the seed organization) and a ROOT client.

## Requirements

- `http_gateway` reachable on port `5001` (override via `GATEWAY_URL`).
- `mapexos` reachable on port `5000` for the admin and ROOT token
  bootstrap (`utils.GetRootToken` / `utils.GetAdminToken`).
- Stack started from `mapexOSDeploy`; the seed admin user
  `admin@mapex.local` and the root organization
  `0000000000000000000aa001` must be present.

## Notes

- Two clients are constructed in `TestMain`: a ROOT client for tests
  that should see every datasource, and an org-scoped admin client
  with `X-Org-Context` set to the seed root org for the
  `WithOrgContext` test.
- The `name` filter is asserted with `strings.ToLower(...)` because the
  service performs a case-insensitive partial match — saga-created
  datasources with a lowercase `"http"` in the name are still valid
  matches.
- `PATCH` currently returns `201 Created` (not `200 OK`); the update
  test asserts on `201` accordingly.
