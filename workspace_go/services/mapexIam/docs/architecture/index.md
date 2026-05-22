# Architecture

## Design
Modular architecture with clean separation of concerns.

## Project Structure
```
src/
├── modules/
│   ├── auth/                    # Authentication (JWT/OAuth2)
│   ├── authorization_cache/     # Shared auth + coverage cache
│   ├── cache_invalidation/      # Listens for invalidation events via NATS
│   ├── groups/                  # Group management
│   ├── lists/                   # Lookup lists + name sync
│   ├── memberships/             # User ↔ Group membership
│   ├── onboarding_orchestrator/ # User/org onboarding flows
│   ├── organizations/           # Organization hierarchy
│   ├── roles/                   # Custom role management
│   └── users/                   # User management
└── shared/
    └── configuration/            # Service configuration
```

## Module Initialization Order

Modules initialize in a specific order to ensure dependencies are satisfied.

| # | Module | Purpose |
|---|--------|---------|
| 1 | `lists` | Core — no dependencies |
| 2 | `organizations` | Core — no dependencies |
| 3 | `authorization_cache` | Shared auth cache used by multiple modules |
| 4 | `roles` | Authorization — custom role definitions |
| 5 | `groups` | Authorization — group management |
| 6 | `memberships` | Authorization — user-group links |
| 7 | `users` | Depends on memberships for multi-tenant filtering |
| 8 | `auth` | Authentication (JWT/OAuth2) |
| 9 | `cache_invalidation` | NATS consumer for invalidation events |
| 10 | `onboarding_orchestrator` | Orchestrator — **always last** |

## Module Responsibilities
- `lists`: lookup list CRUD, publishes name update events for downstream sync (e.g., asset templates)
- `organizations`: org hierarchy CRUD, publishes cache invalidation on access policy / hierarchy changes
- `authorization_cache`: shared Redis store for auth and coverage caches
- `roles`: custom role CRUD, publishes cache invalidation on permission changes
- `groups`: group CRUD + member management, publishes cache invalidation on group changes
- `memberships`: user-group membership CRUD, publishes cache invalidation on membership changes
- `users`: user CRUD + profile management
- `auth`: JWT authentication, token refresh, builds authorization and coverage caches via internal endpoints
- `cache_invalidation`: listens on `MAPEXOS_CACHE_INVALIDATION` stream and refreshes Redis caches
- `onboarding_orchestrator`: coordinates user creation with membership setup

## IAM Model
- Roles are **custom**, not fixed
- Groups aggregate users
- Memberships connect users and groups

## Organization Tree
- Supports multi‑level hierarchy (Vendor → Customer → …)
- Coverage and permissions are resolved by hierarchy

## Request Flow
1. Requests are authenticated via JWT middleware
2. Authorized requests are processed by the appropriate module
3. When data changes (roles, groups, memberships, organizations), cache invalidation events are published to NATS so other services stay in sync
