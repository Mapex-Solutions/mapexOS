# Module e2e: mapexos / groups

## Scope

Covers the groups module of the mapexos platform — collections of
users that aggregate role assignments inside an organization. Backed
by `workspace_go/services/mapexIam/src/modules/groups/`. A group is
always created against an `orgId` (organization-scoped); the
`isSystem` distinction is reserved for built-in seeded groups and the
suite proves both flows reach the same CRUD shape.

## Endpoints exercised

- `POST /api/v1/groups` — create a group (org-scoped, minimal, system,
  and a negative case without `orgId`).
- `GET /api/v1/groups/{id}` — fetch one group.
- `GET /api/v1/groups?includeAll=true` — paginated list of accessible
  groups.
- `PATCH /api/v1/groups/{id}` — partial update of `name`,
  `description`, `enabled`, or all of them.
- `DELETE /api/v1/groups/{id}` — delete a group.

## Fixtures

| File | Description |
|---|---|
| `create_minimal.json` | Smallest valid payload — name + enabled + `{{ORG_ID}}` + one roleId. |
| `create_org_group.json` | "Engineering Team" with description, used by most CRUD tests. |
| `create_system_group.json` | System administrators group payload (still org-scoped). |
| `update_name.json` | PATCH body that renames to "Engineering Team Updated". |
| `update_description.json` | PATCH body that changes only the description. |
| `update_disable.json` | PATCH body with `enabled: false`. |
| `update_full.json` | PATCH body that updates name + description + enabled together. |

`{{ORG_ID}}` is replaced at load time with the seed mapexos
organization id from `common/constants`.

## How to run

```bash
cd e2e_tests
go test ./services/mapexos/groups -v

# Single test
go test ./services/mapexos/groups -v -run TestUpdateGroup_Full
```

## Outcome on pass

- Full CRUD round-trip: create (org / minimal / system) → get → patch
  (name, description, disable, full) → delete → re-get returns 404.
- `TestCreateGroup_NoOrgIdForNonSystem` proves the validator rejects
  group creation without `orgId` with `400 Bad Request`.
- `TestGetGroupById_NotFound` proves an unknown id returns `404`.
- `TestListGroups` proves `includeAll=true` returns a paginated
  envelope (`items[]` + `pagination`) and that a freshly-created group
  appears in the list.

## Requirements

- mapexos / iam service reachable on `http://localhost:5000`
  (override via `MAPEXOS_URL`).
- Seed admin + SuperAdmin role (id `0000000000000000000aa201`) + seed
  mapexos organization, all provisioned by `mongodb-init`.
- `utils.SetupE2EEnvironment()` runs in `TestMain` (clean DB + flush
  cache + re-seed) — destructive against the local stack.

## Notes

- Two clients are wired: `rootClient` (seed admin, wildcard) and
  `adminClient` (org-scoped `admin_vendor.*`); both carry
  `X-Org-Context` pinned to the seed root org. The default `client`
  alias points at `rootClient` for CRUD coverage.
- The mapexos middleware requires `X-Org-Context` on every CRUD
  endpoint, even for the wildcard bearer.
- Every mutating test registers `t.Cleanup` to delete the created
  group, accepting both `200` and `404` during teardown.
