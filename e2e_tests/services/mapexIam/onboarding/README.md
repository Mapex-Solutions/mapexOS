# Module e2e: mapexIam / onboarding

Saga building blocks for the IAM onboarding orchestrator — the single
HTTP entrypoint that creates user + group + membership atomically.
Mirrors the onboarding module under
`workspace_go/services/mapexIam/src/modules/onboarding/`. Saga journeys
import this package whenever they need a fresh non-admin actor inside
the saga's scratch organization; the orchestrator collapses what would
otherwise be three sequential POSTs (user, group, membership) into one
transactional call, so the saga step stays a single Do/Compensate
unit.

## Endpoints exercised

- `POST /api/v1/onboarding/users` — creates user + group + membership
  in one atomic call; returns the new user id, email, and the ids of
  any groups created alongside it (NewGroup mode).

## Fixtures

None. The payload is built in Go from `c.RunID` plus the role id read
from the bag — see `payloads/saga_iot_admin_user.go`. The shared
password constant (`SagaIoTAdminUserPassword`) is hard-coded next to
the builder so the auth step can sign in without threading the secret
through the bag.

## Saga building blocks

- `steps/create_user_with_memberships.go` —
  `CreateUserWithMemberships` POSTs the canonical IoT admin payload
  and publishes `iam.userID`, `iam.userEmail`, and `iam.groupID` on
  the bag. Reads `iam.roleID` so a role must be created earlier in
  the saga (see `../roles/steps`).
- `steps/keys.go` — exports `BagKeyUserID`, `BagKeyUserEmail`,
  `BagKeyGroupID`; consumers import the constants instead of using
  string literals.
- `payloads/saga_iot_admin_user.go` — `SagaIoTAdminUser(runID, roleID)`
  fluent builder that returns a `CreateUserWithMemberships` DTO with
  a runID-stamped email and a freshly-created group bound to the
  given role. The `SagaIoTAdminUserPassword` constant is the
  deterministic password every saga test user shares.

## How to run

There are no module-level `Test*` functions in this package, so
`go test ./services/mapexIam/onboarding/...` is a no-op compile
check. The step executes when a saga journey that imports it runs:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap
```

To run a single journey that exercises the onboarding step:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/... -v -run TestPhase0
```

## Outcome on pass

A green run that imports this step proves:

- The onboarding orchestrator creates user + group + membership in a
  single transactional call.
- The new user can subsequently authenticate via `POST /auth/login`
  using `SagaIoTAdminUserPassword`.
- The membership wiring is visible to the coverage cache (verified by
  the matching auth-package assert).

## Requirements

- `mapexIam` reachable on `MAPEXOS_URL` (default `http://localhost:5000`).
- A role provisioned earlier in the saga via `mapexIam/roles.CreateRole`
  so `iam.roleID` is on the bag.
- An `X-Org-Context` set on the HTTP client (the saga runner sets it
  from the bag after `CreateOrganization` or `SeedAdminLogin`).

## Notes

- The orchestrator is atomic: a failure mid-way (e.g. group exists
  but membership write fails) rolls back all three writes server-side.
  The saga step itself has a no-op `Compensate` because the org-level
  Compensate in `mapexIam/organizations` cascade-deletes the user,
  group, and membership during teardown; running a per-step delete
  here would race the cascade.
- The default builder uses `NewGroup` mode (creates a brand-new group
  alongside the user). Tests that need a different shape should add a
  fluent override to the builder rather than hand-crafting a DTO.
- The runID-stamped email keeps parallel saga runs from colliding on
  the unique-email constraint.
