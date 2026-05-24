# Module e2e: mapexos / users

## Scope

End-to-end coverage of the users module — the human-identity surface of
mapexos served by `mapexIam`
(`workspace_go/services/mapexIam/src/modules/users/`). The suite proves
that users created through the public onboarding orchestrator (atomic
user + membership in one call) round-trip through fetch, list, patch and
delete, that both `internal` and `google` auth providers are accepted,
and that update payloads (name, password, full profile, disable) flow
through to the persisted document. Every user is provisioned with a
local-scope membership in the seed Mapexos org under the SuperAdmin
role, then cleaned up via `t.Cleanup`.

## Endpoints exercised

- `POST   /api/v1/onboarding/users`     — atomic user + membership creation
- `POST   /api/v1/users`                — direct create (used only for the invalid-email guard)
- `GET    /api/v1/users/{id}`           — fetch a single user
- `GET    /api/v1/users`                — paginated list (page / perPage)
- `PATCH  /api/v1/users/{id}`           — partial update (name, password, full, disable)
- `DELETE /api/v1/users/{id}`           — delete (idempotent on 404)

## Fixtures

| File                    | Purpose                                                                |
|-------------------------|------------------------------------------------------------------------|
| `create_internal.json`  | Full internal-auth user with password, phone, job title and avatar.    |
| `create_google.json`    | Google-OAuth user with `externalId` + provider metadata.               |
| `create_minimal.json`   | Minimum internal-auth user (email, password, first / last name).       |
| `update_name.json`      | `PATCH` body that renames `firstName` and `lastName`.                  |
| `update_password.json`  | `PATCH` body that rotates the password and forces next-login change.   |
| `update_full.json`      | `PATCH` body that rewrites email, name, phone, job title and avatar.   |
| `update_disable.json`   | `PATCH` body that toggles `enabled = false`.                           |

## How to run

```bash
cd e2e_tests

# Whole package
go test ./services/mapexos/users -v

# Single test
go test ./services/mapexos/users -v -run TestCreateUser_Internal
```

## Outcome on pass

- Onboarding round-trip: posting an internal, Google, or minimal user
  payload to `/api/v1/onboarding/users` returns `{user, memberships}`
  and the user is immediately retrievable by id.
- Validation: a malformed email on the direct `/api/v1/users` endpoint
  yields 400 without persisting anything.
- Reads: `GET /api/v1/users/{id}` returns the persisted email, first
  name and last name; an unknown id yields 404.
- Updates: name, password (with `changePasswordNextLogin = true`), full
  profile, and `enabled = false` all surface on the next `GET`.
- Delete: a deleted user vanishes from the read path with a 404.
- List: paginated listing surfaces the freshly-created user in the
  `items` array and reports `totalItems / page / perPage` metadata.

## Requirements

- Stack from `mapexOSDeploy/` running (mongo, redis, NATS, mapexIam).
- `mapexos / iam` listening on `:5000`.
- Seed admin (`admin@mapex.local` / `mapex@123`), seed Mapexos org id
  (`constants.MapexosOrgID`) and seed SuperAdmin role id
  (`constants.SuperAdminRoleID`) provisioned by `mongodb-init` — every
  test user is onboarded into that org under that role.
- Go 1.25+.

## Notes

- The package never POSTs to `/api/v1/users` for happy-path creation —
  it always goes through the orchestrator at
  `POST /api/v1/onboarding/users` so the membership is wired atomically.
  The direct `/api/v1/users` endpoint is only used by
  `TestCreateUser_InvalidEmail` to confirm validation rejects bad
  payloads.
- `loadFixture` returns only the base user fields; `createTestUser` and
  the explicit `TestCreateUser_*` cases wrap them with a `memberships`
  array pinned to `MapexosOrgID` + `SuperAdminRoleID` with `local`
  scope before posting.
- A ROOT client (`mapex.*` with `X-Org-Context = MapexosOrgID`) drives
  every test; the `adminClient` is wired in `TestMain` for parity with
  sibling packages but the users suite has no middleware deny cases.
