# Bounded Context: Authorization Cache

**Service:** mapexIam
**Module path:** `src/modules/authorization_cache/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-29

## Purpose
Shared Redis adapter that centralizes invalidation of the IAM authorization data set (per-user/per-org permissions, per-user coverage, per-role permissions). Exposes only invalidation primitives — reads live in the `auth` module. Exists so that `memberships`, `roles`, `groups`, `organizations`, and the `cache_invalidation` consumer can all invalidate the same keys without duplicating cache code.

## Module Layout (intentional §1 exemption)
This is a repository-only adapter, not a service. The module deliberately ships only `domain/` (port) and `infrastructure/` (Redis adapter) — `application/` and `interfaces/` are absent because there is no orchestration to host. DI wiring lives in `module.go::InitRepositories()` via `container.Provide()`, so a `dig.In` struct in `application/di/` would be empty boilerplate. A compile-time port check (`var _ repositories.AuthCacheRepository = (*redis.AuthCacheRepository)(nil)`) in `infrastructure/cache/redis/types.go` covers the safety concern that an `application/di/` struct would otherwise enforce.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Versioning strategy | Invalidation by bumping `auth:org:{orgId}:user:{userId}:ver` (1-100 round robin). Old payloads expire via TTL. | Simple `DEL` used for coverage and role keys |
| SharedCache | Redis DB 5 — cross-service authorization data | AppCache (DB 0) which is service-private |
| Role cache | Key `role:{roleId}` storing resolved permission list for a role | Per-user permission cache, which is keyed by user+org+version |

## Published Events (outbound)
Not applicable — this module is a cache adapter. It does not publish.

## Consumed Events (inbound)
Not applicable — invalidation is triggered by direct method calls from the `cache_invalidation` consumer (and indirectly by other modules via that consumer).

## Driving Ports (inbound)
- `repositories.AuthCacheRepository` — only public surface. Methods: `InvalidateUserAuth`, `InvalidateCoverage`, `InvalidateRole`.

## Driven Ports (outbound)
- `common.SharedCache` (Redis DB 5) — via `Get`, `Set`, `Del`.

## Invariants and Business Rules
- Version pointers have NO TTL; they are written forever and round-robin 1 → 100.
- Versioned payload keys (`...:v{n}`) are expected to have a TTL (managed by the builder side in the `auth` module); this module only bumps the pointer.
- Coverage and role invalidation are plain `DEL`; there is no versioning for those keys.
- Implementation logs every invalidation with `[REPO:AuthCache]` — callers must not log duplicates.

## Known Cross-Context Interactions
- Registered in the DI container by `authorization_cache.InitRepositories()`; consumed by the `cache_invalidation` consumer and the `auth` module.
- The `auth` module owns the READ/BUILD side of the same keys; this module owns only the WRITE-delete side. Changes to key format must be coordinated across both.
