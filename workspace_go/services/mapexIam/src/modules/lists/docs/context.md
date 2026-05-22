# Bounded Context: Lists

**Service:** mapexIam
**Module path:** `src/modules/lists/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Owns generic value-typed reference data ("lists") used across the product as enumerations and hierarchical taxonomies (e.g., manufacturers, asset types, job titles). Supports multi-tenant visibility via `isSystem` / `isTemplate` flags plus `orgId` / `pathKey` / `scope`, and a `parentId` self-link for hierarchies such as manufacturer → model.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| List | A single named value belonging to a typed collection (`type`), optionally with a parent and metadata | Not a page of items; a list is ONE entry with a `type` discriminator |
| Type | String discriminator grouping items (e.g., `manufacturer`, `jobTitle`) | Organization `type` (vendor/customer/etc.), unrelated |
| isSystem | True = MAPEX-global, visible to all tenants | `isTemplate` = vendor/customer-shared template but not global |
| ParentId | Self-reference to build hierarchies (manufacturer → model) | Organization `parentOrgId`, unrelated |

## Published Events (outbound)
None currently published (the DI injects `natsModel.Publisher` but the service has no `Publish(` call — inferred from `Grep`). A "list name update" is mentioned in DI comments as intended future use.

## Consumed Events (inbound)
None.

## Driving Ports (inbound)
- HTTP routes under `/api/v1/lists/` (list, create, get, update, delete).
- `ports.ListServicePort` — not known to be consumed by other modules (inferred).

## Driven Ports (outbound)
- `repositories.ListRepository` (MongoDB)
- `natsModel.Publisher` (injected, publish not observed in service — inferred)

## Invariants and Business Rules
- `type` acts as the discriminator and MUST be set at creation time.
- Multi-tenant visibility is controlled by `isSystem`, `isTemplate`, `OrgID`, `PathKey`, `Scope`.
- Hierarchical items (with `parentId`) expect the parent to exist before children are created (inferred).
- `GetListByEmail` exists on the port but the entity has no `email` field; this is likely a legacy or specialized lookup (inferred — worth verifying).

## Known Cross-Context Interactions
- Consumed by the frontend for selector population (job titles, manufacturers, etc.).
- No direct runtime dependency on other IAM modules.
