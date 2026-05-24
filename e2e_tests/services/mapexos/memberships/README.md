# Module e2e: mapexos / memberships

## Scope

Covers the memberships module of the mapexos platform — the binding
between an assignee (a user or a group) and one or more roles inside
an organization, with a `scope` of `local` (this org only) or
`recursive` (this org + descendants). Backed by
`workspace_go/services/mapexIam/src/modules/memberships/`. The suite
also exercises the `/api/v1/me/coverage` endpoint that materializes
the caller's effective customer coverage from active memberships.

## Endpoints exercised

- `POST /api/v1/memberships` — create a membership (user/local,
  user/recursive, group, multi-role).
- `GET /api/v1/memberships/{id}` — fetch one membership.
- `GET /api/v1/memberships?includeAll=true[&userId=]` — paginated list
  with optional `userId` filter.
- `PATCH /api/v1/memberships/{id}` — partial update of `scope`,
  `enabled`, or `roleIds`.
- `DELETE /api/v1/memberships/{id}` — delete a membership.
- `GET /api/v1/me/coverage` — the authenticated user's effective
  organization coverage rebuilt from active memberships.

## Fixtures

| File | Description |
|---|---|
| `create_user_membership_local.json` | User assignee, `scope: local`, single role. Base CRUD fixture. |
| `create_user_membership_recursive.json` | User assignee, `scope: recursive`. |
| `create_group_membership.json` | Group assignee, `scope: local`, single role. |
| `create_multiple_roles.json` | User assignee with two roles (`{{ROLE_ID_1}}` + `{{ROLE_ID_2}}`). |
| `update_scope.json` | PATCH body that promotes `scope` to `recursive`. |
| `update_disable.json` | PATCH body with `enabled: false`. |
| `update_roles.json` | PATCH body that replaces `roleIds` with a single new role. |

Placeholders `{{USER_ID}}`, `{{GROUP_ID}}`, `{{ORG_ID}}`, `{{ROLE_ID}}`,
`{{ROLE_ID_1}}`, `{{ROLE_ID_2}}` are substituted at load time with
runtime-provisioned ids (test org + two test roles + test group are
created in `TestMain` and torn down at the end).

## How to run

```bash
cd e2e_tests
go test ./services/mapexos/memberships -v

# Single test
go test ./services/mapexos/memberships -v -run TestUpdateMembership_Scope
```

## Outcome on pass

- Full CRUD round-trip across the four create flavors (user/local,
  user/recursive, group, multi-role) → get → patch (scope, disable,
  roles) → delete → re-get returns `404`.
- `TestCreateMembership_UserLocal` validates the full response shape:
  `assigneeType`, `assigneeId`, `orgId`, `scope`, `enabled`, and the
  `roleIds[]` array.
- `TestCreateMembership_MultipleRoles` proves multi-role assignment
  persists exactly two `roleIds`.
- `TestGetMembershipById_NotFound` proves unknown ids return `404`.
- `TestListMemberships` proves the paginated envelope and that both a
  user-based and a group-based membership are returned.
- `TestListMemberships_FilterByUser` proves the `userId=` filter
  scopes results to memberships of that user.
- `TestGetMeCoverage` proves `/api/v1/me/coverage` returns
  `{ userId, customers[] }` reflecting active memberships.

## Requirements

- mapexos / iam service reachable on `http://localhost:5000`
  (override via `MAPEXOS_URL`).
- Seed admin + SuperAdmin role + seed mapexos organization,
  provisioned by `mongodb-init`. Deterministic user ids
  (`constants.RootUserID`, `constants.AdminUserID`) come from the
  seed script.
- `utils.SetupE2EEnvironment()` runs in `TestMain` (clean DB + flush
  cache + re-seed) — destructive against the local stack.

## Notes

- `TestMain` provisions a child test organization, two test roles,
  and a test group, then deletes everything at the end — so each run
  is independent and leaves no trail.
- Two clients are wired: `rootClient` (seed admin, wildcard) and
  `adminClient` (org-scoped `admin_vendor.*`). The default `client`
  alias is `rootClient` for CRUD coverage. Both carry `X-Org-Context`
  pinned to the seed root org.
- `scope=local` vs `scope=recursive` is a first-class distinction:
  recursive memberships grant the assignee access to the org and all
  descendants; `TestCreateMembership_UserRecursive` and
  `TestUpdateMembership_Scope` exercise both transitions explicitly.
