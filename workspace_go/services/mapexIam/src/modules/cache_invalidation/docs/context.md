# Bounded Context: Cache Invalidation

**Service:** mapexIam
**Module path:** `src/modules/cache_invalidation/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-04

## Purpose
NATS consumer module that centralizes reaction to IAM domain change events (role permissions/deletion, organization access policy/hierarchy, membership create/update/delete, group create/update/delete) and translates them into cache invalidations against `authorization_cache` and the coverage cache. Exists so that producer modules publish one event and do not need to know which users/org pairs must be invalidated — this module resolves the fan-out (including group → member expansion).

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| UserOrgPair | Deduplicated `(userId, orgId)` used to bump per-user/per-org auth cache | Membership, which may be assigned to a group rather than a user |
| Fan-out | Expanding a group-scoped change into all of its member users | Coverage expansion, which is about ancestor org descendants |
| Hierarchy change | `organization.hierarchy.changed` — created/deleted org with ancestor list, used to recompute coverage of users with recursive memberships on ancestors | Access policy change, which only affects permission cache, not coverage |

## Published Events (outbound)
Not applicable — this module is a consumer, not a publisher.

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| RolePermissionsChanged | `mapexos.cache.invalidation.role.{roleId}.permissions.changed` | `application/events.RolePermissionsChangedEvent` | roles (this service) |
| RoleDeleted | `mapexos.cache.invalidation.role.{roleId}.deleted` | `application/events.RoleDeletedEvent` | roles |
| OrgAccessPolicyChanged | `mapexos.cache.invalidation.organization.{orgId}.access_policy.changed` | `application/events.OrgAccessPolicyChangedEvent` | organizations |
| OrgHierarchyChanged | `mapexos.cache.invalidation.organization.{orgId}.hierarchy.changed` | `application/events.OrgHierarchyChangedEvent` | organizations |
| MembershipChanged | `mapexos.cache.invalidation.membership.{membershipId}.changed` | `application/events.MembershipChangedEvent` | memberships |
| MembershipDeleted | `mapexos.cache.invalidation.membership.{membershipId}.deleted` | `application/events.MembershipDeletedEvent` | memberships |
| GroupChanged | `mapexos.cache.invalidation.group.{groupId}.changed` | `application/events.GroupChangedEvent` | groups |
| GroupDeleted | `mapexos.cache.invalidation.group.{groupId}.deleted` | `application/events.GroupDeletedEvent` | groups |

Stream: `MAPEXOS_CACHE_INVALIDATION`. Wildcard subject: `mapexos.cache.invalidation.>`. Durable: `mapexos-cache-invalidation-consumer`. Queue group: `mapexIam-CACHE-INVALIDATION-GROUP`.

## Driving Ports (inbound)
- NATS consumer on `mapexos.cache.invalidation.>` (registered via `InitListeners` in phase 4, after repositories and services of other modules).

## Driven Ports (outbound)
- `authCacheRepos.AuthCacheRepository` (authorization_cache module)
- `authRepos.CoverageCacheRepository` (auth module)
- `membershipPorts.MembershipServicePort` (memberships module) — to enumerate memberships by role/org/scope
- `groupRepos.GroupMemberRepository` (groups module) — to expand group memberships into user lists

## Invariants and Business Rules
- Unparseable messages are ACKed (not retried) to avoid infinite redelivery of legacy payloads.
- DLQ is enabled (`ServiceType: mapex-iam`, `EventType: cache.invalidation`).
- Group membership fan-out MUST dedupe via `UserOrgPair` to avoid bumping the same key twice.
- On `organization.hierarchy.changed` with `action=deleted`, direct memberships on the deleted org are ALSO invalidated (in addition to recursive ancestors).
- No business state lives here — this module only invalidates caches; rebuilding is lazy via the `auth` module.

## Known Cross-Context Interactions
- Reads memberships via the `memberships` service port (not the repository) — preserves bounded context.
- Reads group members directly from `GroupMemberRepository` — acceptable because `groups` is in the same service.
- Downstream of every write in `roles`, `memberships`, `groups`, and `organizations`.
