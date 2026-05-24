# Module e2e: mapexIam / roles

Saga building blocks for the IAM roles module — the surface that
defines named permission bundles attached to users and groups.
Mirrors the roles module under
`workspace_go/services/mapexIam/src/modules/roles/`. The canonical
fixture here bundles every permission a saga test user needs to drive
the IoT pipeline end-to-end (org, user, group, asset, route, trigger,
workflow CRUD plus events read), so journeys do not have to hand-pick
permission strings each time they provision an actor.

## Endpoints exercised

- `POST /api/v1/roles` — creates a role with the given name, scope,
  flags, and permission list.
- `DELETE /api/v1/roles/{id}` — used by `Compensate` to clean up the
  role; the role is created in the seed parent org (where the
  bootstrap actor has coverage), so the child-org cascade does not
  reach it and an explicit delete is required.

## Fixtures

None. The payload is built in Go from `c.RunID` plus permission
constants imported from `permissions/{mapexos,assets,router}` — see
`payloads/saga_iot_admin_role.go`.

## Saga building blocks

- `steps/create_role.go` — `CreateRole` POSTs the canonical IoT admin
  role and publishes `iam.roleID` on the bag. Owns its own
  `Compensate` that calls `DELETE /api/v1/roles/{id}` (idempotent on
  404).
- `steps/keys.go` — exports `BagKeyRoleID` (`iam.roleID`); imported by
  the onboarding step so renames are caught at compile time.
- `payloads/saga_iot_admin_role.go` — `SagaIoTAdminRole(runID)` fluent
  builder returning a `RoleCreate` DTO with `Scope=local`,
  `IsSystem=false`, `IsTemplate=false`, and the canonical IoT-admin
  permission set. `WithName` is the only override most callers need.

## How to run

There are no module-level `Test*` functions in this package, so
`go test ./services/mapexIam/roles/...` is a no-op compile check. The
step executes when a saga journey that imports it runs:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap
```

To target a specific journey:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/... -v -run TestPhase0
```

## Outcome on pass

A green run that imports this step proves:

- A role with `Scope=local` can be created under the seed parent org.
- The role's permissions land in the response (and, once bound to a
  user via the onboarding step, expand correctly into the coverage
  cache).
- The `Compensate` delete cleans the role even when the org cascade
  cannot reach it.

## Requirements

- `mapexIam` reachable on `MAPEXOS_URL` (default `http://localhost:5000`).
- The caller must be authenticated with coverage on the seed parent
  org — the saga runner sets this via `SeedAdminLogin` before the
  step runs.

## Notes

- The role is intentionally created against the seed parent org, not
  the saga's scratch child org. The bootstrap actor (seed admin) has
  coverage on the seed root, so role creation always succeeds; the
  trade-off is that the saga's `CreateOrganization` cascade does not
  reach the role and we need the explicit `Compensate` delete.
- The permission bundle is the union of what every IoT-pipeline saga
  needs. Journeys that exercise a narrower surface can still reuse
  this role — coverage strictly above what the test calls is benign.
- Permission constants come from `permissions/{mapexos,assets,router}`
  in the contracts module so renames in the source of truth break the
  payload at compile time.
