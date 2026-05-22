# Bounded Context: App (Service Bootstrap Orchestrator)

**Service:** assets
**Module path:** `src/modules/app/`
**Owner:** assets team
**Last reviewed:** 2026-05-11

## Purpose

The `app` module is NOT a business bounded context. It is the **service bootstrap
orchestrator** for the `assets` service. Its sole responsibility is to coordinate
the deterministic initialization of every other module in the service, in a fixed
three-phase order, so that dependency wiring is complete and consistent before
any driving adapter (HTTP server, NATS consumer) starts accepting work.

Everything under `app/` is infrastructure-of-assembly. It owns no entities,
enforces no invariants of its own, and emits/consumes no domain events. If it
were removed, the service would have no entry point — but no business rule
would change.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Module | A business bounded context folder (`assets`, `assettemplates`, `healthmonitor`, etc.) registered in `configuration/modules` | Go package or npm module |
| Phase | One of the three ordered stages of initialization (Repositories, Services, Interfaces) | Workflow runtime phase |
| Lazy module | A module flagged `Lazy = true` that opts out of eager phase execution | Go lazy evaluation |

## Published Events (driven — outbound)

None. The bootstrap orchestrator does not publish NATS events, HTTP responses,
or any other outbound messages.

## Consumed Events (driving — inbound)

None. The orchestrator is invoked exactly once at process start via
`InitModule(c *fiber.App)` from the service's `main.go`. It does not subscribe
to NATS, listen on HTTP, or react to external signals.

## Driving Ports (what can call this module)

- `InitModule(*fiber.App)` — called once from `main.go` during service startup.

## Driven Ports (what this module requires)

- `configuration/modules.Modules` — ordered list of module descriptors, each
  exposing optional `InitRepositories`, `InitServices`, `InitInterfaces` hooks
  and a `Lazy` flag.
- `microservices/logger` — structured logging for the bootstrap progress trace.

## Invariants and Business Rules

The `app` module enforces no business invariants. It does enforce two
**operational** invariants that every business module relies on:

1. **Three-phase strict ordering.** Phase 1 (Repositories) completes for all
   non-lazy modules before Phase 2 (Services) begins; Phase 2 completes before
   Phase 3 (Interfaces). A Phase-2 service may therefore safely resolve any
   repository registered in Phase 1, and a Phase-3 HTTP route / NATS consumer
   may safely resolve any service registered in Phase 2.
2. **Lazy modules are skipped at bootstrap.** Any module with `Lazy = true` is
   excluded from all three phases here; its initialization is the caller's
   responsibility (typically deferred to first-use).

## Known Cross-Context Interactions

- Invokes `InitRepositories` / `InitServices` / `InitInterfaces` on every
  non-lazy business module registered under `shared/configuration/modules`
  (currently: `assettemplates`, `assets`, `healthmonitor`, `mqttcerts`).
  The orchestrator has no knowledge of what those hooks do — it only
  guarantees the order in which they run.
- Does not cross the service boundary. All interaction is local process wiring.
