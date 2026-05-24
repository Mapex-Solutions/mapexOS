# Module e2e: mapexos / lists

## Scope

Covers the lists module of the mapexos platform — typed key/value
catalogs (e.g. `assetGroup`, `assetType`) consumed by the rest of the
platform as enumerations. Backed by
`workspace_go/services/mapexIam/src/modules/lists/`. Lists carry a
`type`, a `name`, a `value`, an `isSystem` flag (true for built-in
catalogs, false for org-defined ones), and an `orgId` for scope.

## Endpoints exercised

- `POST /api/v1/lists` — create a list entry (org-scoped, system, and
  minimal payloads).
- `GET /api/v1/lists/{id}` — fetch one entry.
- `GET /api/v1/lists?includeAll=true&page=&perPage=` — paginated list,
  also exercised with `type=` and `name=` filters.
- `PATCH /api/v1/lists/{id}` — partial update of `name` + `value`.
- `DELETE /api/v1/lists/{id}` — delete an entry.

## Fixtures

| File | Description |
|---|---|
| `create_minimal.json` | `assetGroup` "Workstations", `isSystem: false`, seed root org. |
| `create_org_list.json` | `assetGroup` "Servers", `isSystem: false` — main CRUD fixture. |
| `create_system_list.json` | `assetType` "Physical Server", `isSystem: true`. |
| `update_name.json` | PATCH body that renames to "Updated Name" / `updated_value`. |

## How to run

```bash
cd e2e_tests
go test ./services/mapexos/lists -v

# Single test
go test ./services/mapexos/lists -v -run TestListLists_FilterByType
```

## Outcome on pass

- Full CRUD round-trip on the three payload shapes (org / system /
  minimal): create → get → patch (name + value) → delete → re-get
  returns either `404` or `200` with `nil data`.
- `TestGetListById_NotFound` proves missing ids return `404` or
  `200`+`nil` consistently.
- `TestListLists` proves pagination envelope shape (`items[]` +
  `pagination{page, perPage}`).
- `TestListLists_FilterByType` proves the `type` query filter — every
  returned item satisfies `type == assetGroup`.
- `TestListLists_FilterByName` exercises the partial-name filter and
  logs the count (tolerant of org-scope visibility).

## Requirements

- mapexos / iam service reachable on `http://localhost:5000`
  (override via `MAPEXOS_URL`).
- Seed mapexos organization (id `0000000000000000000aa001`),
  provisioned by `mongodb-init`.
- `utils.SetupE2EEnvironment()` runs in `TestMain` (clean DB + flush
  cache + re-seed) — destructive against the local stack.

## Notes

- Two clients are wired (`rootClient` wildcard + `adminClient`
  org-scoped `admin_vendor.*`); the default `client` alias is
  `rootClient` for CRUD coverage. Both carry `X-Org-Context` pinned to
  the seed root org.
- `system` vs org distinction is data-level (`isSystem` boolean), not
  a separate endpoint — the same `POST /api/v1/lists` handles both.
- The `TestGetListById_NotFound` and `TestDeleteList` assertions
  accept both `404` and `200` with `nil data`, because the service may
  return either shape for a missing entry.
- `TestListLists_FilterByName` is intentionally tolerant — it logs
  rather than asserts the count because results depend on org
  visibility.
