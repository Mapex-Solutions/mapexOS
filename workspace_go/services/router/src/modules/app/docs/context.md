# Bounded Context: App (Router Bootstrap)

**Service:** router
**Module path:** `src/modules/app/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-22

## Purpose

This is NOT a business bounded context — it is the **service bootstrap orchestrator** for the `router` service. Its single responsibility is to iterate the module registry declared in `src/shared/configuration/modules/config.go` and invoke each module's lifecycle hooks (`InitRepositories`, `InitServices`, `InitInterfaces`) in the correct order. It carries no domain logic, no ports, no entities, and no state of its own — it only wires the DI container and starts the consumer/HTTP surfaces of the real modules (`routegroups`, `events`).

It is documented here because §1 of `/go-arch` mandates a `docs/context.md` for every folder under `src/modules/`. This file formalizes that `app/` exists solely to own the boot sequence.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Module | One entry in `configMod.Modules` (a `common.ModuleConfig` with a name and three lifecycle hooks) | A Go package or DDD bounded context |
| Lifecycle phase | One of the three ordered init passes: Repositories → Services → Interfaces | Service runtime phases (scan, publish, etc.) |
| Lazy module | `ModuleConfig{Lazy: true}` — skipped by the bootstrap; its hooks are expected to be invoked on-demand by some other owner | Lazy DI provider (dig only instantiates on first resolve) |

## Published Events (driven — outbound)

None. This module does not publish.

## Consumed Events (driving — inbound)

None. This module does not consume messages. It is invoked exactly once from `main.go` after the Fiber app and the DI container have been constructed.

## Driving Ports (what can call this module)

- `InitModule(c *fiber.App)` — called from `main.go` during service startup.

## Driven Ports (what this module requires)

- `configMod.Modules` — the ordered registry of `common.ModuleConfig` values declared in `src/shared/configuration/modules/config.go`.
- Each listed module's `InitRepositories` / `InitServices` / `InitInterfaces` functions (optional per module).
- `github.com/Mapex-Solutions/mapexGoKit/microservices/logger` — for phase-by-phase progress logging.

## Invariants and Business Rules

- The three phases MUST run in this exact order for every non-lazy module: **Repositories → Services → Interfaces**. This guarantees that when interfaces (HTTP routes, NATS consumers) start, all services and their repository dependencies are already registered in the DIG container.
- A module with `Lazy: true` is skipped entirely here; bootstrap never calls its hooks.
- Each hook is optional (`nil` means "nothing to do for this phase") — the bootstrap silently skips nil hooks rather than erroring.
- No domain types, no ports, no DI struct — the module MUST remain a thin orchestration loop. Any business logic in this file is a violation.
- Modules are iterated per phase (three passes over the slice), NOT per-module (one module fully initialized before the next). Cross-module DI resolution during Phase 3 relies on ALL modules having completed Phase 2.

## Known Cross-Context Interactions

- Reads the module registry owned by `src/shared/configuration/modules/config.go`.
- Calls lifecycle hooks of the two real bounded contexts in this service: **routegroups** (Phase 1+2+3) and **events** (Phase 2+3, no repositories).
- Does not import from any module's `domain/` or `application/` — only from the module's top-level package (which re-exports the `Init*` functions).
