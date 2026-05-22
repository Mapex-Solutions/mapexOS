# Bounded Context: Onboarding Orchestrator

**Service:** mapexIam
**Module path:** `src/modules/onboarding_orchestrator/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-21

## Purpose
Application service that composes `users`, `memberships`, `organizations`, and `groups` into two atomic flows used by the admin UI: "create a user with their initial access" and "update a user and replace their access". Wraps the write sequence in a MongoDB transaction so partial onboardings do not leak. Owns no domain state itself.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Onboarding | Atomic user-create + membership(s) create in one transaction | Sign-up / self-registration — this is admin-driven |
| Access configuration | Either direct role assignment (creates a `user` membership) OR group assignment (creates a `group` membership) | Role — one access config may reference multiple roles |
| Default scope | Taken from the target org's `AccessPolicy.DefaultScope` (`"local"` or `"recursive"`) | Group scope or membership scope (set by the orchestrator from org policy) |

## Published Events (outbound)
Not applicable — this module does not publish directly. The downstream `users` / `memberships` / `groups` service calls publish their own events.

## Consumed Events (inbound)
None.

## Driving Ports (inbound)
- HTTP routes under `/api/v1/onboarding/users`:
  - `POST /users` — create user with memberships
  - `PATCH /users/:userId` — update user with access replacement
- `ports.UserOnboardingServicePort` — consumed only by its HTTP handlers (inferred).

## Driven Ports (outbound)
- `userPorts.UserServicePort`
- `membershipPorts.MembershipServicePort`
- `orgPorts.OrganizationServicePort`
- `groupPorts.GroupServicePort`
- `mongoManager.MongoManager` — required for `RunTransaction()` (ACID across modules)

Not applicable — this module is an orchestrator; it has no repository, no entity, no domain folder.

## Invariants and Business Rules
- Both endpoints run inside a MongoDB transaction; a failure at any step rolls back the full onboarding.
- Target `orgId` comes from `RequestContext.OrgContext`; `scope` is NOT accepted from the request body — it is read from `org.AccessPolicy.DefaultScope`.
- Update flow removes pre-existing memberships and group memberships for the user in the current org, then recreates them — there is no in-place patch.
- Access config is mutually exclusive per call: either direct roles OR a group, not both (inferred from the port comments).

## Known Cross-Context Interactions
- Pure orchestrator — every persistence path delegates to another module's service port.
- The cache invalidation events fired by downstream writes propagate to the `cache_invalidation` consumer as usual.
