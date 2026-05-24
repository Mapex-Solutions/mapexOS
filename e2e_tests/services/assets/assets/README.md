# Module e2e: assets / assets

## Scope

Module e2e suite for the `assets` module of the assets service
(`workspace_go/services/assets/src/modules/assets/`). It exercises the
CRUD contract used by IoT devices: create an asset bound to a template
and a route group, fetch it, update its metadata, list with pagination
and filters, and delete it. `TestMain` first provisions a save-event
route group on the router and a minimal asset template, because the
`AssetCreate` payload requires both as foreign keys.

## Endpoints exercised

- `POST /api/v1/assets` — create asset (template + route group required).
- `GET /api/v1/assets/{id}` — fetch one asset by id (also covers 404).
- `PATCH /api/v1/assets/{id}` — partial update of name and description.
- `DELETE /api/v1/assets/{id}` — delete and confirm the resource is gone.
- `GET /api/v1/assets` — paginated list with `page`, `perPage`,
  `includeAll`, and `category` filter.
- `POST /api/v1/asset_templates` / `DELETE /api/v1/asset_templates/{id}`
  — invoked from `TestMain` to provision and tear down the template the
  asset tests depend on (template module is covered by its own suite).
- `POST /api/v1/route_groups` / `DELETE /api/v1/route_groups/{id}` —
  invoked against the router service from `TestMain` for the same
  prerequisite reason.

## Test functions

- `TestCreateAsset_Valid`
- `TestCreateAsset_Minimal`
- `TestGetAssetById`
- `TestGetAssetById_NotFound`
- `TestUpdateAsset_Name`
- `TestDeleteAsset`
- `TestListAssets`
- `TestListAssets_FilterByCategory`

## Fixtures

| File                  | Scenario                                                                         |
|-----------------------|----------------------------------------------------------------------------------|
| `create_asset.json`   | Full MQTT device with password auth, geolocation, description, and route bind.   |
| `create_minimal.json` | Minimal valid HTTP device — only the required fields plus injected foreign keys. |
| `update_name.json`    | Partial PATCH body changing `name` and `description`.                            |

Note: `assetTemplateId` and `routeGroupIds` carry the placeholder string
`WILL_BE_INJECTED_BY_TEST` in the JSON files and are overwritten at
runtime with the ids `TestMain` provisioned.

## How to run

```bash
cd e2e_tests

# Full module suite
go test ./services/assets/assets -v

# A single test
go test ./services/assets/assets -v -run TestCreateAsset_Valid
```

## Outcome on pass

Confirms the assets module honors its public HTTP contract end-to-end:
required-field validation, foreign-key linkage to asset templates and
router groups, the CRUD round-trip (create → read → patch → delete →
404), and the paginated list with category filtering — all under the
seed root organization context.

## Requirements

- `assets` reachable on port `5002` (override via `ASSETS_URL`).
- `mapexos` reachable on port `5000` for the admin token bootstrap.
- `router` reachable on port `5003` to create the prerequisite route
  group.
- Stack started from `mapexOSDeploy`; the seed admin user
  `admin@mapex.local` and the root organization
  `0000000000000000000aa001` must be present (provisioned by
  `mongodb-init`).
- `API_KEY` env var optional — defaults to the canonical internal key
  used by the stack.

## Notes

- `TestMain` provisions a fresh template and route group per test run
  and removes both in the teardown, so the suite is self-contained and
  leaves no orphan records.
- The category id `670a4cde48e006e3f95e8eb3` used by the filter test
  comes from the seed catalog; if the seed changes that constant must
  be updated to match.
- The `payloads/`, `steps/`, and `asserts/` sibling folders are saga
  building blocks consumed by the IoT journeys; they are not part of
  this module e2e suite.
