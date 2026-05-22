# Bounded Context: App (mapexIam Bootstrap)

**Service:** mapexIam
**Module path:** `src/modules/app/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-22

## Purpose

The `app` module is NOT a business bounded context — it is the **service bootstrap
orchestrator** for the `mapexIam` service. Its sole responsibility is to iterate
the module registry declared in `src/shared/configuration/modules/config.go` and
invoke each module's lifecycle hooks in a fixed, deterministic order so that
dependency wiring is complete and consistent before any driving adapter
(HTTP server, NATS consumer, NATS listener) starts accepting work.

Everything under `app/` is infrastructure-of-assembly. It owns no entities,
enforces no business invariants of its own, and emits/consumes no domain events.
If it were removed, the service would have no entry point — but no business rule
would change. It is documented here because §1 of `/go-arch` mandates a
`docs/context.md` for every folder under `src/modules/`.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Module | One entry in `configMod.Modules` (a `common.ModuleConfig` with a name, lazy flag, and lifecycle hooks) | A Go package or DDD bounded context |
| Phase | One of the four ordered init passes: Repositories → Services → Interfaces → Listeners | Auth/session lifecycle phases |
| Lazy module | A module flagged `Lazy = true` that opts out of eager phase execution; its hooks are invoked on-demand by another owner | Go lazy evaluation / DIG lazy resolve |

## Published Events (driven — outbound)

None. The bootstrap orchestrator does not publish NATS events, HTTP responses,
or any other outbound messages.

## Consumed Events (driving — inbound)

None. The orchestrator is invoked exactly once at process start via
`InitModule(c *fiber.App)` from the service's `main.go`. It does not subscribe
to NATS, listen on HTTP, or react to external signals.

## Driving Ports (what can call this module)

- `InitModule(*fiber.App)` — called once from `main.go` during service startup,
  after the Fiber app and the DIG container have been constructed.

## Driven Ports (what this module requires)

- `src/shared/configuration/modules.Modules` — the ordered registry of
  `common.ModuleConfig` values, each exposing optional `InitRepositories`,
  `InitServices`, `InitInterfaces`, `InitListeners` hooks and a `Lazy` flag.
- `github.com/Mapex-Solutions/mapexGoKit/microservices/logger` — structured
  logging for phase-by-phase boot progress.

## Invariants and Business Rules

The `app` module enforces no business invariants. It does enforce the following
**operational** invariants that every business module relies on:

1. **Four-phase strict ordering.** Phase 1 (Repositories) completes for all
   non-lazy modules before Phase 2 (Services) begins; Phase 2 completes before
   Phase 3 (Interfaces); Phase 3 completes before Phase 4 (Listeners).
   A Phase-2 service may therefore safely resolve any repository registered in
   Phase 1; a Phase-3 HTTP route / NATS consumer may safely resolve any service
   registered in Phase 2; a Phase-4 NATS event listener may safely rely on all
   prior phases being complete.
2. **Lazy modules are skipped at bootstrap.** Any module with `Lazy = true` is
   excluded from all four phases here; its initialization is the caller's
   responsibility (typically deferred to first-use).
3. **Nil hooks are silently skipped.** Each lifecycle hook is optional; a `nil`
   hook means "nothing to do for this phase" and the bootstrap simply advances.
4. **Per-phase iteration, not per-module.** Modules are iterated four times
   (once per phase), NOT initialized one at a time end-to-end. Cross-module DI
   resolution in later phases relies on all modules having completed earlier
   phases.
5. **No domain logic in this module.** No types, no ports, no DI struct. Any
   business logic inside `app/` is a violation of its bootstrap-only role.

## Known Cross-Context Interactions

- Reads the module registry owned by `src/shared/configuration/modules/config.go`.
- Invokes `InitRepositories` / `InitServices` / `InitInterfaces` / `InitListeners`
  on every non-lazy business bounded context in this service: `auth`,
  `authorization_cache`, `cache_invalidation`, `groups`, `lists`, `memberships`,
  `onboarding_orchestrator`, `organizations`, `roles`, `users`, plus any future
  siblings. The orchestrator has no knowledge of what those hooks do — it only
  guarantees the order in which they run.
- Does not import from any module's `domain/` or `application/` — only from the
  module's top-level package (which re-exports the `Init*` functions).
- Does not cross the service boundary. All interaction is local process wiring.
