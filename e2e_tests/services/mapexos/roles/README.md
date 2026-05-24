# Module e2e: mapexos / roles

## Scope

End-to-end coverage of the roles module — the permission grant primitive
of the mapexos surface served by `mapexIam`
(`workspace_go/services/mapexIam/src/modules/roles/`). The suite proves
that system roles and org-scoped roles round-trip through CRUD, that
wildcard and namespaced permissions (`mapex.*`, `user.*`, `admin.*`,
`asset.read`) survive create + read intact, and that validation rejects
malformed payloads. A throwaway customer org is provisioned in `TestMain`
to host every org-scoped role; the seed dataset is never modified.

## Endpoints exercised

- `POST   /api/v1/roles`                — create system / org-scoped roles
- `GET    /api/v1/roles/{id}`           — fetch a single role
- `GET    /api/v1/roles`                — paginated list (page / perPage)
- `PATCH  /api/v1/roles/{id}`           — partial update (name, permissions, full)
- `DELETE /api/v1/roles/{id}`           — delete (idempotent on 404)
- `POST   /api/v1/organizations`        — bootstrap throwaway parent org in `TestMain`
- `DELETE /api/v1/organizations/{id}`   — tear it down on suite exit

## Fixtures

| File                         | Purpose                                                                |
|------------------------------|------------------------------------------------------------------------|
| `create_system_role.json`    | Global system role with `mapex.*` (`isSystem: true`, `scope: global`). |
| `create_org_role.json`       | Org-scoped `Site Manager` role with user / asset permissions.          |
| `create_minimal.json`        | Smallest valid org role (`Viewer` with two read permissions).          |
| `update_name.json`           | `PATCH` body that renames the role.                                    |
| `update_permissions.json`    | `PATCH` body that replaces the permissions array.                      |
| `update_full.json`           | `PATCH` body that rewrites name, description and permissions.          |
| `update_disable.json`        | `PATCH` body toggling `enabled = false` (kept for parity; v1 ignores). |

## How to run

```bash
cd e2e_tests

# Whole package
go test ./services/mapexos/roles -v

# Single test
go test ./services/mapexos/roles -v -run TestCreateRole_SystemRole
```

## Outcome on pass

- CRUD round-trip: system and org-scoped roles can be created, fetched,
  listed (with pagination metadata), updated, and deleted.
- `orgId` is resolved from `X-Org-Context` by the service rather than
  trusted from the payload — the response carries a populated id even
  when the test posts a placeholder.
- Wildcard permissions: `mapex.*`, `user.*`, `admin.*` and `asset.*`
  survive create + read unchanged.
- Validation: missing `name`, missing `orgId` on non-system roles, and
  empty `permissions` arrays are all rejected with 400.

## Requirements

- Stack from `mapexOSDeploy/` running (mongo, redis, NATS, mapexIam).
- `mapexos / iam` listening on `:5000`.
- Seed admin (`admin@mapex.local` / `mapex@123`) and seed Mapexos org id
  (`constants.MapexosOrgID`) provisioned by `mongodb-init` — the
  throwaway test org is created as a child of it.
- Go 1.25+.

## Notes

- `TestMain` creates a single shared customer org (`Test Organization
  for Roles`) and deletes it after `m.Run()`. The org-scoped fixtures
  use `{{CUSTOMER_ID}}` as a placeholder that `loadFixture` rewrites
  with the runtime id.
- The CRUD client is a ROOT-token client (`mapex.*`) with
  `X-Org-Context` pinned to `MapexosOrgID`. An `adminClient` is also
  wired in `TestMain` for parity with sibling packages but the role
  tests themselves rely on the ROOT client because every assertion is a
  CRUD outcome rather than a middleware deny.
- The `enabled` field in role payloads is a v1 holdover; the service
  ignores it, so `update_disable.json` is kept for symmetry only and is
  not asserted against.
