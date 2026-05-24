# Module e2e: mapexos / auth

## Scope

Covers the authentication surface of the mapexos platform — login,
logout, refresh token, and the "who am I" coverage endpoint exposed by
`workspace_go/services/mapexIam/src/modules/auth/`. Unlike the other
mapexos modules, this suite does not need a bootstrap JWT in
`TestMain`: it is exactly the surface that mints those JWTs, so each
test drives the login flow itself.

## Endpoints exercised

- `POST /auth/login` — credential exchange for an access + refresh token pair.
- `POST /auth/logout` — invalidates the current access token.
- `POST /auth/refresh` — rotates tokens; refresh token is passed in the
  `X-Refresh-Token` header.
- `GET /auth/users/me/coverage` — returns the authenticated user's
  organization coverage.

## Fixtures

| File | Description |
|---|---|
| `login_valid.json` | Seed admin credentials (`admin@mapex.local` / `mapex@123`), `keepConnected: true`. |
| `login_invalid_email.json` | Malformed email — triggers 400. |
| `login_short_password.json` | Password under the 8-char minimum — triggers 400. |
| `login_wrong_password.json` | Valid email, wrong password — triggers 401. |

## How to run

```bash
cd e2e_tests
go test ./services/mapexos/auth -v

# Single test
go test ./services/mapexos/auth -v -run TestRefreshToken
```

## Outcome on pass

- `TestLogin_Valid` proves the seed admin can log in and that the
  response carries `access_token`, `refresh_token`, and a `user`
  object with `id` and `email`.
- `TestLogin_InvalidEmail` and `TestLogin_ShortPassword` prove input
  validation rejects malformed payloads with `400 Bad Request`.
- `TestLogin_WrongPassword` proves wrong credentials return
  `401 Unauthorized`.
- `TestLogout` proves a freshly-issued access token is accepted by
  `/auth/logout`.
- `TestRefreshToken` proves the refresh flow rotates both tokens when
  the refresh token is supplied in `X-Refresh-Token`.
- `TestGetMyCoverage` and `TestGetMyCoverage_Unauthorized` prove the
  coverage endpoint returns the caller's accessible orgs and rejects
  unauthenticated requests with `401`.

## Requirements

- mapexos / iam service reachable on `http://localhost:5000`
  (override via `MAPEXOS_URL`).
- Seed admin user, role, organization, and recursive membership
  provisioned by `mongodb-init` on first boot of the stack.

## Notes

- `TestMain` deliberately does NOT call `SetupE2EEnvironment` and does
  NOT pre-acquire a token — this suite is the source of truth for the
  login/refresh/coverage contract itself.
- Refresh token transport: header `X-Refresh-Token`, not body. The
  access token still goes in `Authorization: Bearer`.
- No `X-Org-Context` header is set here; the coverage endpoint resolves
  scope from the bearer alone.
