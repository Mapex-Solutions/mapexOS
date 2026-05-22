# Bounded Context: Auth

**Service:** mapexIam
**Module path:** `src/modules/auth/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Owns authentication and session lifecycle for platform users: login, refresh-token rotation, logout, and the read-side endpoints that power the frontend context switcher (`/users/me/coverage`, `/me/permissions`). Coverage and authorization caches are read here on the hot path and rebuilt on demand via internal endpoints protected by API key.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Session | Refresh-token record stored in Redis keyed per user | HTTP Fiber session; Redis TTL-managed, not a stateful cookie |
| Coverage | List of organizations the authenticated user can reach through memberships (direct or recursive) | Permission set — coverage only answers "which orgs", not "which actions" |
| Authorization Cache | Per-user/per-org resolved permission list with a versioned key (`auth:org:{orgId}:user:{userId}:v{n}`) | Role cache (`role:{roleId}`) owned by authorization_cache module |
| Build endpoint | `/internal/auth/build-*` — rebuilds a cache entry. Called by other services, never by the browser | Public `/auth/*` endpoints which are user-facing |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| — | — | — | — |

None. The auth module does not publish domain events; cache invalidation is handled by the publishing modules (roles, memberships, groups, organizations).

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| — | — | — | — |

None directly (inferred — no consumer registered under `auth/`; invalidation listens in `cache_invalidation/`).

## Driving Ports (inbound)
- HTTP public routes under `/auth/` (Login, Logout, Refresh, `users/me/coverage`, `me/permissions`)
- HTTP internal routes under `/internal/auth/` (`build-authorization`, `build-coverage`) — API-Key protected
- `ports.AuthServicePort` — consumed by HTTP handlers; not exported cross-module (inferred)

## Driven Ports (outbound)
- `repositories.AuthRepository` (MongoDB — legacy entity lookups)
- `repositories.SessionRepository` (Redis — refresh token storage)
- `repositories.AuthorizationCacheRepository` (Redis — permission read/build)
- `repositories.CoverageCacheRepository` (Redis — coverage read/build)
- `userPorts.UserServicePort`, `membershipPorts.MembershipServicePort`, `rolePorts.RoleServicePort` (cross-module ports used during login and cache build)

## Invariants and Business Rules
- Refresh flow rotates tokens: the old refresh token is invalidated when a new pair is issued.
- `GetMyCoverage` and `GetMyPermissions` must read cache first; a miss triggers a synchronous build.
- Coverage cache key: `user:{userId}:orgs`. Authorization cache uses a versioned pointer so invalidation is O(1).
- Internal build endpoints require `X-API-Key`; they MUST NOT be exposed publicly.
- `Logout` extracts session info from the access token to invalidate the matching refresh token.

## Known Cross-Context Interactions
- Reads memberships/roles through their service ports during coverage and permission cache builds.
- Shares `AuthorizationCacheRepository`/`CoverageCacheRepository` with the `authorization_cache` and `cache_invalidation` modules (repositories registered in auth's `module.go`, consumed elsewhere).
- Coverage payload (`OrganizationCoverage`, `UserAccess` in `domain/repositories/types.go`) is the shape returned to the frontend shell.
