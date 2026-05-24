# Module e2e: mapexos / organizations

## Scope

End-to-end coverage of the organizations module — the multi-tenant hierarchy
backbone of the mapexos surface served by `mapexIam`
(`workspace_go/services/mapexIam/src/modules/organizations/`). The suite is
the largest single-module test in the repo (~30 tests) and exercises full
CRUD, the deep `customer / site / building` hierarchy, cursor-based tree
pagination, and the coverage middleware that decides which actor sees which
org. Every test runs against the live stack and registers cleanup so the
seed dataset is left untouched.

## Endpoints exercised

- `POST   /api/v1/organizations`              — create customer / site / building
- `GET    /api/v1/organizations/{id}`         — fetch a single org
- `GET    /api/v1/organizations`              — paginated list (page / perPage)
- `GET    /api/v1/organizations/tree`         — cursor-paginated tree (next / previous)
- `PATCH  /api/v1/organizations/{id}`         — partial update (name, enabled, full)
- `DELETE /api/v1/organizations/{id}`         — delete (idempotent on 404)
- `POST   /api/v1/roles`                      — create restricted role for deny tests
- `POST   /api/v1/onboarding/users`           — provision restricted admin actor
- `POST   /api/v1/auth/login`                 — log in as the restricted admin

## Fixtures

| File                         | Purpose                                                                    |
|------------------------------|----------------------------------------------------------------------------|
| `create_customer.json`       | Top-level `customer` org under the seed root (ACME Corporation).           |
| `create_site.json`           | `site` org parented to a customer via `{{PARENT_ID}}` placeholder.         |
| `create_building.json`       | `building` org parented to a site via `{{PARENT_ID}}` placeholder.         |
| `create_minimal.json`        | Smallest valid customer payload — drives required-field assertions.        |
| `update_name.json`           | `PATCH` body that renames the org.                                         |
| `update_disable.json`        | `PATCH` body that toggles `enabled = false`.                               |
| `update_full.json`           | `PATCH` body that rewrites name, address, phone, and access policy.        |

## How to run

```bash
cd e2e_tests

# Whole package
go test ./services/mapexos/organizations -v

# Single test
go test ./services/mapexos/organizations -v -run TestOrganizationHierarchy_PathKeyPropagation
```

## Outcome on pass

- CRUD round-trip: create / read / list / patch / delete each return the
  documented status codes and the persisted document round-trips.
- Deep hierarchy: `customer -> site -> building` creates succeed and report
  the correct `parentOrgId` at every level.
- `pathKey` propagation: each descendant's `pathKey` is exactly the parent
  `pathKey` plus a `/segment`; a building's `pathKey` has four segments.
- `customerId` inheritance: a customer is its own `customerId`; sites and
  buildings under it inherit the same value.
- Tree pagination: `/tree` exposes `next` / `previous` cursors and
  `hasNext` / `hasPrevious` flags, and forward + backward walks land on
  disjoint pages.
- Middleware: ROOT (`mapex.*`) passes with or without `X-Org-Context`; a
  restricted admin without the header gets 403, and a restricted admin
  pointed at an org outside their coverage also gets 403.

## Requirements

- Stack from `mapexOSDeploy/` running (mongo, redis, NATS, mapexIam).
- `mapexos / iam` listening on `:5000`.
- Seed admin (`admin@mapex.local` / `mapex@123`), seed root org id
  (`0000000000000000000aa001`), and seed SuperAdmin role id
  (`0000000000000000000aa201`) provisioned by `mongodb-init`.
- Go 1.25+.

## Notes

- The package contains a `provisionRestrictedAdmin` helper that builds a
  scratch non-wildcard admin actor purely through the public API: it
  creates a customer org, attaches a `local`-scope role with just
  `organization.read` + `organization.list`, then onboards a user via
  `POST /api/v1/onboarding/users` and logs in to obtain a JWT. The
  resulting token is the opposite of the seed super-admin — org-scoped,
  no `mapex.*` wildcard — and is the only way to exercise the
  middleware deny paths (`AdminWithoutOrgContext_Deny`,
  `AdminWithUnauthorizedOrgContext_Deny`).
- `coveragePropagationDelay` (4 seconds) is a deliberate pause after any
  org-mutating write. The mapexIam coverage cache is invalidated through
  a NATS event, so back-to-back requests from the same client need this
  window before the next call can target the freshly created descendant
  via `X-Org-Context`. The helpers (`createTestOrganization`,
  `provisionRestrictedAdmin`) sleep automatically; the hierarchy tests
  sleep again between layers because they create more than one child.
- No `t.Skip` calls remain in this package — every test runs on every
  invocation.
