# Bounded Context: Memberships

**Service:** mapexIam
**Module path:** `src/modules/memberships/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Owns the assignment of access: links an assignee (a user OR a group) to an organization with a set of roles and a scope (`local` or `recursive`). This is the core grant record the authorization pipeline reads to build per-user coverage and permissions. Also exposes `/api/v1/me/coverage`, the endpoint that answers "which customers/orgs am I a member of".

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Membership | `(assigneeType, assigneeId) → (orgId, roleIds, scope)` grant record | Group membership in `group_members` junction — different storage and purpose |
| AssigneeType | `"user"` or `"group"` | Not the user's role in the product |
| Scope | `"local"` (this org only) or `"recursive"` (this org + descendants) | Group scope (`"global"`/`"local"`) which is unrelated |
| Coverage | Derived view: orgs reachable via the user's memberships (direct + group + recursive expansion) | Authorization cache (permissions per org), different key space |
| CustomerID | Denormalized tenant anchor on the membership; absent for vendor-org memberships | Organization `customerId`, the source of the denormalization |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| MembershipChanged | `mapexos.cache.invalidation.membership.{membershipId}.changed` | `cache_invalidation/application/events.MembershipChangedEvent` | cache_invalidation (this service) |
| MembershipDeleted | `mapexos.cache.invalidation.membership.{membershipId}.deleted` | `cache_invalidation/application/events.MembershipDeletedEvent` | cache_invalidation |

## Consumed Events (inbound)
None.

## Driving Ports (inbound)
- HTTP routes under `/api/v1/memberships/` (list, create, get, update, delete).
- HTTP routes under `/api/v1/me/coverage` (authenticated user's coverage).
- `ports.MembershipServicePort` — consumed by `auth`, `cache_invalidation`, `groups`, `onboarding_orchestrator`, `users`.

## Driven Ports (outbound)
- `repositories.MembershipRepository` (MongoDB)
- `common.CacheGetOrSetEx` (for coverage read cache)
- `natsModel.Publisher` (domain events)
- `orgPorts.OrganizationServicePort`
- `groupQueryPorts.GroupQueryServicePort` (read-only group lookup — avoids circular dep)

## Invariants and Business Rules
- Every create/update/delete MUST publish the matching `cache.invalidation.membership.*` event; the consumer reacts by bumping the authorization cache version for the affected user+org and invalidating coverage.
- `GetAllMemberships` (internal use) bypasses coverage filtering; `GetMemberships` (external) applies it via `RequestContext`.
- `GetAssigneeIdsByOrgIds` is paginated internally to handle large datasets.
- Group-scoped memberships do NOT immediately enumerate members at write time — the consumer expands them lazily using `GroupMemberRepository`.
- `OrgPathKey` is denormalized from `Organization` to enable range queries without a join.

## Known Cross-Context Interactions
- Upstream of the `cache_invalidation` consumer for every write.
- Consumed by `auth` (coverage + authorization cache build paths).
- Consumed by `onboarding_orchestrator` (atomic user creation with memberships).
- Consumed by `users` (PATH 2 in user listing — direct memberships and group-inherited memberships).
- Reads groups via `GroupQueryServicePort`, never via the group repo directly.
