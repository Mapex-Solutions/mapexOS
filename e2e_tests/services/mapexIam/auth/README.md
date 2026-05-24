# Module e2e: mapexIam / auth

Saga building blocks for the IAM authentication surface — the routes
that issue and verify the bearer token every other module test relies
on. Mirrors the auth module under `workspace_go/services/mapexIam/src/modules/auth/`.
There is no per-endpoint CRUD suite here because authentication has no
state of its own to list or update; instead, this package exposes the
steps and asserts that saga journeys (and module tests in other
packages) compose to log in, switch identity, and verify the resulting
JWT.

## Endpoints exercised

The building blocks call the auth surface exposed by `mapexIam`:

- `POST /auth/login` — exchanges email + password for an access token;
  the response also carries a refresh token. Mounted at the service
  root (no `/api/v1` prefix), unlike the rest of IAM.
- `GET /auth/users/me/coverage` — returns the orgs the caller can act
  on; used to verify membership wiring after login.

## Fixtures

None. Credentials come from `common/constants` (seed admin) or from
the bag-published email created by the onboarding step plus the shared
saga password constant in `../onboarding/payloads`.

## Saga building blocks

- `steps/seed_admin_login.go` — `SeedAdminLogin` logs in as the seed
  admin shipped by `mongodb-init`, sets the bearer + `X-Org-Context`
  on every per-service client, and publishes `iam.userJWT` plus
  `iam.organizationID` on the bag. Every journey starts here unless it
  needs a non-admin actor.
- `steps/authenticate_user.go` — `AuthenticateUser` swaps the seed
  admin for the saga test user created earlier in the journey, reading
  `iam.userEmail` + `iam.organizationID` from the bag.
- `steps/keys.go` — exports `BagKeyUserJWT` (`iam.userJWT`); imported
  by downstream consumers so renames break the build instead of
  failing at runtime.
- `asserts/assert_jwt_valid.go` — `AssertJwtValid` checks the bearer
  is a well-formed three-segment JWT; `AssertJwtHasOrgContext` hits
  `/auth/users/me/coverage` and verifies the saga org is in the list.

## How to run

There are no module-level `Test*` functions in this package, so
`go test ./services/mapexIam/auth/...` is a no-op compile check. The
steps and asserts execute when a saga journey that imports them runs:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap
```

To exercise the auth flow specifically inside a journey:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/... -v -run TestPhase0
```

## Outcome on pass

A green run that imports these building blocks proves:

- The seed admin user shipped by `mongodb-init` can authenticate
  against the live stack.
- A saga-provisioned user (created via the onboarding orchestrator)
  can authenticate with its own credentials.
- The returned bearer is structurally valid and the coverage cache
  reflects the user's membership in the saga organization.

## Requirements

- `mapexIam` reachable on `MAPEXOS_URL` (default `http://localhost:5000`).
- Mongo seeded by `mongodb-init` so the seed admin credentials work.
- For `AuthenticateUser`: a prior step that created the user
  (`mapexIam/onboarding`) and the parent org (`mapexIam/organizations`).

## Notes

- `/auth/login` is the only IAM endpoint outside `/api/v1`. The other
  IAM modules (organizations, roles, onboarding) all live under the
  versioned prefix.
- `AssertJwtValid` is intentionally structural — signature verification
  happens implicitly the next time the saga calls a protected endpoint
  and gets back something other than `401`.
- `AssertJwtHasOrgContext` is the canonical canary for membership
  cache freshness: a green here proves the role-grant inside the saga
  org propagated to the coverage cache.
