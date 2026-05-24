# Module e2e: assets / assettemplates

## Scope

Module e2e suite for the `assettemplates` module of the assets service
(`workspace_go/services/assets/src/modules/assettemplates/`). Asset
templates describe how a manufacturer/model is parsed: an `assetIdPath`
to extract the device identifier and a pair of inline JS scripts
(validator + conversion) executed against incoming telemetry. The
suite exercises the full CRUD plus the paginated listing with status
filtering.

## Endpoints exercised

- `POST /api/v1/asset_templates` — create template (full and minimal
  payloads).
- `GET /api/v1/asset_templates/{id}` — fetch one template by id (also
  covers 404).
- `PATCH /api/v1/asset_templates/{id}` — partial update; covered for
  both metadata fields (name, description, version) and the inline JS
  scripts.
- `DELETE /api/v1/asset_templates/{id}` — delete and confirm the
  resource is gone.
- `GET /api/v1/asset_templates` — paginated list with `page`,
  `perPage`, and `enabled` filter.

## Test functions

- `TestCreateAssetTemplate_Valid`
- `TestCreateAssetTemplate_Minimal`
- `TestGetAssetTemplateById`
- `TestGetAssetTemplateById_NotFound`
- `TestUpdateAssetTemplate_Scripts`
- `TestUpdateAssetTemplate_Metadata`
- `TestDeleteAssetTemplate`
- `TestListAssetTemplates`
- `TestListAssetTemplates_FilterByStatus`

## Fixtures

| File                    | Scenario                                                                                |
|-------------------------|-----------------------------------------------------------------------------------------|
| `create_template.json`  | Full template — Acme TS-2000 with all four script slots and `assetIdPath`.              |
| `create_minimal.json`   | Minimal valid template — only required fields and inline validator/conversion stubs.    |
| `update_metadata.json`  | Partial PATCH changing `name`, `description`, `version` (keeps `assetIdPath`).          |
| `update_scripts.json`   | Partial PATCH replacing `scriptValidator` and `scriptConversion` with marked variants.  |

## How to run

```bash
cd e2e_tests

# Full module suite
go test ./services/assets/assettemplates -v

# A single test
go test ./services/assets/assettemplates -v -run TestUpdateAssetTemplate_Scripts
```

## Outcome on pass

Confirms the asset templates module honors its public HTTP contract:
required-field validation, the CRUD round-trip (create → read → patch
→ delete → 404), partial PATCH semantics for both metadata fields and
the inline script body, and paginated listing with the `enabled`
filter — all under the seed root organization context.

## Requirements

- `assets` reachable on port `5002` (override via `ASSETS_URL`).
- `mapexos` reachable on port `5000` for the admin token bootstrap.
- Stack started from `mapexOSDeploy`; the seed admin user
  `admin@mapex.local` and the root organization
  `0000000000000000000aa001` must be present (provisioned by
  `mongodb-init`).

## Notes

- The PATCH tests accept both `200` and `201` because the service
  currently returns `201` on successful updates — the assert is loose
  on the status code on purpose, until the service is normalised.
- Unlike the sibling `assets` suite, no router or template
  prerequisites are needed: each test owns its own template lifecycle.
