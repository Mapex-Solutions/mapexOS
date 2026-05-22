# Bounded Context: App

**Service:** workflow
**Module path:** `src/modules/app/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-22

## Purpose

Bootstrap/orchestration module. Not a domain context — this module has no entities, events, or business rules of its own. It wires the workflow service startup by iterating the module list from `shared/configuration/modules` and invoking each module's three lifecycle hooks in order: `InitRepositories` → `InitServices` → `InitInterfaces`. All dependency construction happens through the shared DIG container; this module only drives phase ordering and emits `[MODULE:*]` boot logs.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Module | A bounded context registered in `configMod.Modules` with optional `InitRepositories` / `InitServices` / `InitInterfaces` callbacks | Go package (much smaller granularity) |
| Lazy | Flag that excludes a module from eager startup (it self-registers on demand) | Runtime idle state |
| Phase | One of the three boot stages (repositories, services, interfaces) | Workflow execution lifecycle |

## Published Events (driven — outbound)

_None._

## Consumed Events (driving — inbound)

_None._

## Driving Ports (what can call this module)

- `InitModule(app *fiber.App)` — invoked once from the service entrypoint (`cmd/main.go`) after DIG container and MongoDB/NATS infrastructure are ready.

## Driven Ports (what this module requires)

- `shared/configuration/modules.Modules` — ordered list of module descriptors (name, lazy flag, three init callbacks).
- `microservices/logger` — boot-sequence logging.

## Invariants and Business Rules

- Phases MUST execute in order: repositories first, services next, interfaces last. Reordering breaks DIG resolution.
- Only non-lazy modules are eagerly initialized; lazy modules are responsible for their own registration elsewhere.
- This module owns no state, persistence, or messaging — it is pure bootstrap.

## Known Cross-Context Interactions

- Invokes every other module's `Init*` hooks. Any module added to `shared/configuration/modules.Modules` is automatically picked up by this module on next boot.
