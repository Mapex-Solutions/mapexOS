# Bounded Context: Organizations

**Service:** mapexIam
**Module path:** `src/modules/organizations/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Owns the org hierarchy used as the multi-tenant backbone: vendors → customers → sites/buildings/floors/zones. Maintains hierarchical `PathKey` (dot/slash-delimited), local `Code`, `Depth`, denormalized `CustomerID` anchor, and per-org `AccessPolicy` + `AuthConfig`. Exposes flat and tree list endpoints and a pathKey-based descendant lookup used by authorization flows.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| PathKey | Full hierarchical key (e.g., `000001/000001/0001`) used for range/prefix queries | `Code` which is only the local fragment |
| Depth | Integer level in the tree (0=vendor, 1=customer, 2=site, …) | Organization `type`, which is the semantic label |
| CustomerID | Denormalized tenant anchor pointing to the customer org in the ancestor chain; null for vendor-level orgs | Organization ID itself |
| AccessPolicy | `RolePolicy` (`merge`/`strict`) + `DefaultScope` (`local`/`recursive`) — governs how permissions inherit | `AuthConfig` which governs identity-provider integration |
| AuthConfig | Per-org identity provider settings (providerType, issuer, client, JWT claim mappings) | Service-level auth middleware config |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| OrgAccessPolicyChanged | `mapexos.cache.invalidation.organization.{orgId}.access_policy.changed` | `cache_invalidation/application/events.OrgAccessPolicyChangedEvent` | cache_invalidation (this service) |
| OrgHierarchyChanged | `mapexos.cache.invalidation.organization.{orgId}.hierarchy.changed` | `cache_invalidation/application/events.OrgHierarchyChangedEvent` | cache_invalidation |

Hierarchy events are published for both create and delete with the full ancestor list so the consumer can invalidate recursive-scope coverage caches.

## Consumed Events (inbound)
None.

## Driving Ports (inbound)
- HTTP routes under `/api/v1/organizations/` (list, tree, create, get, update, delete).
- HTTP internal routes under `/api/internal/v1/organizations/` — currently a no-op group; the previous retention-policies endpoint has been moved to the Events service (confirmed by route file).
- `ports.OrganizationServicePort` — consumed by `memberships`, `groups`, `onboarding_orchestrator`, `roles`, `users`.

## Driven Ports (outbound)
- `repositories.OrganizationRepository` (MongoDB)
- `natsModel.Publisher`

## Invariants and Business Rules
- `Code` is generated from the parent's `childCount` to guarantee local uniqueness (inferred from port comments).
- `PathKey` and `Depth` are computed from the parent chain and MUST stay in sync with `parentOrgId`.
- `AccessPolicy.RolePolicy` changes MUST publish `OrgAccessPolicyChanged` so all memberships under the org get their auth cache invalidated.
- Create and delete MUST publish `OrgHierarchyChanged` with the full ancestor list for recursive-coverage invalidation.
- `AccessPolicy.AllowDirectPermissions` was removed — V1 is pure role-based.
- Filtering respects `RequestContext`: `OrgContext + includeChildren` uses PathKey range; `OrgContext` alone uses direct id; no context uses `$in` over accessible orgs.

## Known Cross-Context Interactions
- Read by almost every other IAM module via `OrganizationServicePort` (membership scoping, role scoping, onboarding, user visibility).
- Upstream of coverage cache invalidation when the hierarchy changes.
- `AuthConfig` is the contract surface read by login flows that need per-tenant identity providers.
