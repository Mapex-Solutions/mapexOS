# Module e2e: mapexIam / organizations

This package covers the IAM organizations module on two surfaces:

- The `e2e/` subpackage is the canonical module e2e suite — CRUD,
  paginated list, query filters, and cursor-based tree walks against
  the live `mapexIam` service.
- The sibling `steps/`, `payloads/`, and `asserts/` folders are saga
  building blocks that other journeys (and the e2e suite itself)
  import to provision a scratch customer org with a deterministic
  shape.

Mirrors the organizations module under
`workspace_go/services/mapexIam/src/modules/organizations/`.

## Endpoints exercised

The `e2e/` suite exercises every public route of the organizations
module:

- `POST /api/v1/organizations` — create (happy path, validation,
  auth).
- `GET /api/v1/organizations/{id}` — read by id (happy path, 404).
- `PATCH /api/v1/organizations/{id}` — partial update (happy path,
  404).
- `DELETE /api/v1/organizations/{id}` — delete + post-delete 404
  (happy path, 404).
- `GET /api/v1/organizations` — paginated list with `page`, `perPage`,
  `name`, `type`, `enabled` filters; also asserts `totalItems`,
  out-of-range pages, and forward + backward page walks.
- `GET /api/v1/organizations/tree` — cursor-based tree navigation
  under the active `X-Org-Context`.

## Fixtures

None on disk. The e2e suite and the saga building blocks both build
payloads in Go from `c.RunID` via `payloads.SagaTestCustomerOrg`,
which mirrors the contract DTO from
`packages/contracts/services/mapexIam/organizations`.

## Saga building blocks

These are imported by saga journeys (and by the e2e suite for
payload construction):

- `steps/create_organization.go` — `CreateOrganization` POSTs the
  canonical customer-org payload, publishes `iam.organizationID` and
  `iam.organizationPathKey` on the bag, and cascade-deletes children
  (users, groups, roles, memberships) in `Compensate`.
- `steps/keys.go` — exports `BagKeyOrgID` and `BagKeyOrgPathKey`.
- `payloads/saga_test_org.go` — `SagaTestCustomerOrg(runID)` fluent
  builder; defaults to `Type=customer`, `ParentOrgID` =
  seed root id, `Enabled=true`, internal IDP, `strict` role policy,
  `local` default scope. `WithName` and `WithParentOrgID` are the
  fluent overrides callers reach for.
- `asserts/assert_organization_exists.go` — `AssertOrganizationExists`
  fetches the org by id from the bag and verifies the API returns it
  with `enabled=true`.

## How to run

The runnable suite lives in the `e2e/` subpackage:

```bash
cd e2e_tests

# Full suite
go test ./services/mapexIam/organizations/e2e -v

# Single test
go test ./services/mapexIam/organizations/e2e -v -run TestCreate_201

# All list / tree tests
go test ./services/mapexIam/organizations/e2e -v -run 'TestList_|TestTree_'
```

The saga building blocks in `steps/`, `payloads/`, `asserts/` carry no
`Test*` functions — they execute when imported by a saga journey
(`go test -tags=saga ./journey/...`).

## Outcome on pass

Passing the `e2e/` suite proves:

- Every public CRUD route returns the expected status code on happy
  path, validation, auth, and not-found scenarios.
- Page-based pagination is stable: forward (1 -> 15) and backward
  (15 -> 1) walks with `perPage=1` visit every fixture exactly once;
  `totalItems` matches the universe size regardless of `perPage`;
  out-of-range pages return `items=[]` rather than clamping or 404.
- Filters (`name`, `type`, `enabled`) compose correctly under AND.
- The `/tree` cursor walk covers every fixture in the active
  `X-Org-Context`, stopping cleanly when `hasNext=false`.

## Requirements

- `mapexIam` reachable on `MAPEXOS_URL` (default `http://localhost:5000`).
- Mongo seeded by `mongodb-init` (the seed admin must be able to log
  in, and the seed root org id `0000000000000000000aa001` must exist
  to parent the saga-created orgs).
- The list tests create 15 fixtures per run; cleanup is automatic via
  `t.Cleanup`, but the run leaves no residue only if the stack is
  reachable through teardown.

## Notes

- The list suite isolates each run by stamping `runID` into the
  fixture name and filtering subsequent queries with
  `?name=<orgNamePrefix>-<runID>`. This keeps pre-seeded orgs and
  parallel test runs out of the assertion universe.
- `listFixtureCount = 15` is deliberately picked so `perPage=1` walks
  exercise pagination edges (page 1, mid pages, last page) that
  smaller fixture sets would miss.
- The `TestUpdate_200` PATCH tolerates both 200 and 201 to stay
  resilient against handler-level status changes that do not affect
  semantics; the GET that follows is the real assertion.
- The `/tree` cursor walk caps iteration at 200 pages as a safety net
  against backend bugs that would otherwise loop forever; under
  correct behaviour the walk exits via `hasNext=false` long before
  that bound.
- The saga building blocks live alongside the e2e suite so the
  payload builder is shared — the e2e tests are the most reliable
  contract test for the saga payload itself.
