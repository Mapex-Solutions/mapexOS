# Bounded Context: Roles

**Service:** mapexIam
**Module path:** `src/modules/roles/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Owns definition of roles — named bundles of permission strings referenced by memberships. Supports multi-tenant visibility (system, template, local) and hierarchical inheritance through the org tree via `PathKey` + `Scope`. Writes publish events so the cache_invalidation consumer can bump the per-role cache and the auth cache for every affected user.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Role | Named collection of permission strings (e.g., `["read_users", "write_devices"]`) | Membership — a role is a template, a membership is an assignment |
| Permission | String literal under `Role.Permissions` (see `packages/permissions/`) | HTTP permission middleware key (same strings, different layer) |
| isSystem | MAPEX-global role visible to all tenants | `isTemplate` (vendor/customer-level template) |
| Global vs local scope | `"global"` cascades down the org tree; `"local"` is limited to the owning org | Membership scope, which is different (`"local"` vs `"recursive"`) |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| RolePermissionsChanged | `mapexos.cache.invalidation.role.{roleId}.permissions.changed` | `cache_invalidation/application/events.RolePermissionsChangedEvent` | cache_invalidation (this service) |
| RoleDeleted | `mapexos.cache.invalidation.role.{roleId}.deleted` | `cache_invalidation/application/events.RoleDeletedEvent` | cache_invalidation |

## Consumed Events (inbound)
None.

## Driving Ports (inbound)
- HTTP routes under `/api/v1/roles/` (list, create, get, update, delete).
- `ports.RoleServicePort` — consumed by `auth` (permission resolution during cache build) and `users` (role name enrichment).

## Driven Ports (outbound)
- `repositories.RoleRepository` (MongoDB)
- `orgPorts.OrganizationServicePort` (for hierarchical role scoping)
- `natsModel.Publisher`

## Invariants and Business Rules
- Permission updates MUST publish `RolePermissionsChanged`; the consumer deletes `role:{roleId}` and bumps auth cache for every membership carrying the role.
- Deletion MUST publish `RoleDeleted`; memberships still exist at publish time so the consumer can enumerate affected user-org pairs.
- Listing applies hierarchical inheritance: system + MAPEX-exclusive + local + global roles from ancestor orgs.
- Multi-tenant fields (`OrgID`, `PathKey`, `Scope`) are populated from `RequestContext` at creation.

## Known Cross-Context Interactions
- Upstream of `cache_invalidation` on every write.
- `auth` reads role permissions during authorization cache population.
- `users` reads role names for display enrichment.
