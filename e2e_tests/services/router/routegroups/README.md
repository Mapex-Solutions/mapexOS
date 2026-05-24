# Module e2e: router / routegroups

## Scope

Module e2e suite for the `routegroups` module of the router service
(`workspace_go/services/router/src/modules/routegroups/`). It exercises
the public CRUD contract of a route group: create with one or more
inline routers, fetch by id, partial update, delete with a 404
follow-up, and the paginated list with several filters (name, enabled,
version, multi-filter, projection). The suite drives the router under
two identities — the wildcard root user and the seed admin under the
root organization context — to prove that both flows honour the
coverage middleware enforced by the service.

## Endpoints exercised

- `POST /api/v1/route_groups` — create a route group with inline
  `routers` (`kind=save_event` in the baseline fixture).
- `GET /api/v1/route_groups/{id}` — fetch one route group; also used to
  assert a 404 after delete.
- `PATCH /api/v1/route_groups/{id}` — partial update (name swap).
  Accepts both `200 OK` and `201 Created` for build skew.
- `DELETE /api/v1/route_groups/{id}` — remove a route group; cleanup
  helper tolerates `404`.
- `GET /api/v1/route_groups` — paginated list with combinations of
  `page`, `perPage`, `name`, `enabled`, `version`, multi-filter, and
  `projection`.

## Test functions

- `TestCreateRouteGroup`
- `TestGetRouteGroupById`
- `TestUpdateRouteGroup`
- `TestDeleteRouteGroup`
- `TestListRouteGroups_BasicPagination`
- `TestListRouteGroups_FilterByName`
- `TestListRouteGroups_FilterByEnabled`
- `TestListRouteGroups_FilterByVersion`
- `TestListRouteGroups_MultipleFilters`
- `TestListRouteGroups_Projection`
- `TestListRouteGroups_WithOrgContext`
- `TestListRouteGroups_RootUser`

## Fixtures

| File                              | Scenario                                                                                              |
|-----------------------------------|-------------------------------------------------------------------------------------------------------|
| `create_routegroup.json`          | Baseline route group `API Routes v1`, enabled, version `1.0.0`, bound to the seed root org `0000000000000000000aa001`, with a single `save_event` router carrying `metadata.source=api`. |
| `create_routegroup_versioned.json`| Companion route group `API Routes v2`, disabled, version `2.0.0`, same root org; used to populate the listing for filter tests so projections distinguish multiple records. |
| `update_routegroup.json`          | Partial PATCH body — renames the target to `API Routes v1 Updated`.                                   |
| `update_enabled.json`             | Partial PATCH body toggling `enabled=false`; available for ad-hoc runs (not loaded by the current `Test*` functions). |

All `orgId` references point at the canonical seed root organization
`0000000000000000000aa001`; the suite does not provision additional
orgs.

## How to run

```bash
cd e2e_tests

# Full module suite
go test ./services/router/routegroups -v

# A single test
go test ./services/router/routegroups -v -run TestCreateRouteGroup
```

## Outcome on pass

Confirms the routegroups module honours its public HTTP contract
end-to-end: required-field validation on create, the CRUD round-trip
(create → read → patch → delete → 404), and the paginated list with
every supported filter combination (name, enabled, version, multi
filter, projection). It also proves that the coverage middleware
correctly resolves the root organization context for both the wildcard
root user and the seed admin.

## Requirements

- `router` reachable on port `5003` (override via `ROUTER_URL`).
- `mapexos` reachable on port `5000` for the root/admin token
  bootstrap.
- Stack started from `mapexOSDeploy`; the seed admin user
  `admin@mapex.local` and the root organization
  `0000000000000000000aa001` must be present (provisioned by
  `mongodb-init`).

## Notes

- Every CRUD endpoint demands `X-Org-Context` even when the bearer
  carries the wildcard role; `TestMain` sets it on both clients to the
  root org id.
- The PATCH endpoint accepts `200 OK` or `201 Created` — the suite
  tolerates both to absorb build skew on older router images.
- The `payloads/` and `steps/` sibling folders are saga building blocks
  consumed by the IoT and automation journeys (save-event, trigger, and
  workflow route group variants); they are not part of this module e2e
  suite.
