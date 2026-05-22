# Bounded Context: Users

**Service:** mapexIam
**Module path:** `src/modules/users/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Owns the user record: identity (email + password hash + external auth provider), profile fields (names, phone, job title, avatar, tour flag), and lifecycle (enabled, created/updated). Exposes `/api/v1/users/` CRUD, `/me` self-read/update, a counter endpoint, and `GetUserByEmail` used by the login flow. Users are tenant-agnostic here; tenancy comes from their memberships.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| User | Platform account — email, password/external-auth, profile | Membership assignee — a user becomes tenant-scoped only via memberships |
| AuthProvider | Per-user identity provider info (`type`, `externalId`, `metadata`) | Organization `AuthConfig` which is per-tenant provider settings |
| Myself (`/me`) | Shortcuts for the authenticated user (`GET /me`, `PATCH /me`, `PATCH /me/tour`) | Admin-driven user CRUD on `/users/:id` |
| StartTour | Boolean flag controlling first-login product tour | Not an onboarding state — only UI hint |

## Published Events (outbound)
None currently — there is no `NatsBus.Publish` call in this module (verified via Grep on `users/application/services/`).

## Consumed Events (inbound)
None.

## Driving Ports (inbound)
- HTTP routes under `/api/v1/users/`:
  - `/me`, `PATCH /me`, `PATCH /me/tour` — self-service
  - `GET /`, `POST /`, `GET /counter`, `GET /:userId`, `PATCH /:userId`, `DELETE /:userId` — admin CRUD
- `ports.UserServicePort` — consumed by `auth` (login/refresh), `groups`, `onboarding_orchestrator`.

## Driven Ports (outbound)
- `repositories.UserRepository` (MongoDB)
- `membershipPorts.MembershipServicePort` (coverage / tenant-scoping of user listing)
- `groupQueryPorts.GroupQueryServicePort` (group enrichment and PATH 2 listing)
- `orgPorts.OrganizationServicePort`
- `rolePorts.RoleServicePort`
- `common.AppCache` (Redis DB 0 — counter cache)

## Invariants and Business Rules
- `GetUserByEmail` returns the FULL entity including password hash and is reserved for auth flows; all HTTP responses use DTOs that strip sensitive fields.
- Users are not scoped by `orgId`; listing uses PATH 1 (direct memberships) and PATH 2 (group memberships) via the membership and group-query ports to derive visibility from the caller's `RequestContext`.
- `PATCH /me/tour` bypasses the permission check — only `AuthMiddleware` is required.
- Counter endpoint implements cache-aside via `AppCache`.

## Known Cross-Context Interactions
- Consumed by `auth` during login (password check) and during coverage/authorization cache build.
- Consumed by `onboarding_orchestrator` inside the create/update transactions.
- Consumed by `groups` (member add/remove flows).
- Does NOT publish invalidation events itself — upstream writes (memberships, groups) take care of auth cache invalidation when user access changes.
