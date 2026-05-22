# Bounded Context: Groups

**Service:** mapexIam
**Module path:** `src/modules/groups/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Owns the lifecycle of IAM groups (named collections of users used as membership assignees) and of the `group_members` junction collection. Exposes full CRUD over `/api/v1/groups`, paginated member listing, add/remove member operations, and a cross-domain query port (`GroupQueryServicePort`) used by `users` and `memberships` to avoid circular dependencies.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Group | A named, multi-tenant, enabled/disabled collection used as a membership assignee | Organization — groups are scoped under one or are system-wide |
| GroupMember | Row in the junction collection linking `groupId` ↔ `userId` (used for scalable member queries) | Membership, which links an assignee (user or group) to an organization |
| Scope | `"global"` or `"local"` — inheritance behavior of the group across the org tree | Membership scope (`"local"`/`"recursive"`) which is different |
| PathKey | Denormalized hierarchical key for multi-tenant filtering | Organization `code` which is the local-only fragment |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| GroupChanged | `mapexos.cache.invalidation.group.{groupId}.changed` | `cache_invalidation/application/events.GroupChangedEvent` | cache_invalidation (this service) |
| GroupDeleted | `mapexos.cache.invalidation.group.{groupId}.deleted` | `cache_invalidation/application/events.GroupDeletedEvent` | cache_invalidation |

Published on create, update, delete, and on member add/remove (inferred — the service publishes `group.changed` at multiple points in `group_service.go`).

## Consumed Events (inbound)
None.

## Driving Ports (inbound)
- HTTP routes under `/api/v1/groups/` (list, create, counter, get, update, delete, members list/add/remove).
- `ports.GroupServicePort` — consumed by `onboarding_orchestrator` for onboarding flows.
- `ports.GroupQueryServicePort` — read-only port consumed by `users` and `memberships`, purposely narrow to break the circular dependency `GroupService → UserService/MembershipService`.

## Driven Ports (outbound)
- `repositories.GroupRepository` (MongoDB)
- `repositories.GroupMemberRepository` (MongoDB — junction table, scalable for 100K+ tenants)
- `orgPorts.OrganizationServicePort`
- `membershipPorts.MembershipServicePort`
- `userPorts.UserServicePort`
- `common.AppCache` (Redis DB 0) — service-private cache used for counters
- `natsModel.Publisher` — domain event publishing

## Invariants and Business Rules
- `AddMemberToGroup` / `RemoveMemberFromGroup` are idempotent (no error on duplicate/missing membership).
- Members are stored in the `group_members` junction collection, not in an embedded array on `Group`.
- Every write must publish the matching `cache.invalidation.group.*` event so the shared consumer can fan-out to all members.
- Multi-tenant fields (`OrgID`, `PathKey`, `Scope`) are populated from `RequestContext` at creation time.
- Counter endpoint implements cache-aside via `AppCache`.

## Known Cross-Context Interactions
- Consumed by `onboarding_orchestrator` (add/remove user to/from group during onboarding).
- Consumed by `cache_invalidation` consumer which expands `group.changed/deleted` into per-member cache invalidation using `GroupMemberRepository`.
- `GroupQueryServicePort` is consumed by `users` (group enrichment) and `memberships` (resolving group-scoped memberships).
