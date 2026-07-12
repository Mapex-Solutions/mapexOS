# infra — the e2e Environment layer

`infra.Ensure(...)` lets a saga journey **declare the live-stack services it
needs** and have them provisioned before the run and released after — so e2e is
self-provisioning instead of assuming an already-running stack. It sits on the
saga `Fixture` hook.

## The saga Fixture

`core/saga` has a saga-level lifecycle pair, run around the whole item walk:

```go
type Fixture struct {
    Setup   func(*Context) error   // once, before ANY item (even before iam_bootstrap)
    CleanUp func(*Context) error   // once, after ALL items and ALL step Compensates; ALWAYS
}
```

It is named distinctly from `Step{Do, Compensate}` on purpose: a Step's
`Compensate` undoes a mutation to the system under test; `CleanUp` disposes of
resources acquired FOR the run. Different concept, different name.

Wire it with `RunWithFixture` (the plain `Run` is this with a zero fixture, so
existing journeys are unchanged):

```go
saga.RunWithFixture(t, ctx, runID, clients,
    infra.Ensure("mapexos", "assets", "broker", "minio", "nats"),
    items...)
```

### Ordering guarantee

```
Setup  →  items (Do/Check)  →  step Compensates (reverse)  →  CleanUp
  │                                                             │
  └ infra up, before everything            infra down, after everything, ALWAYS ┘
```

`CleanUp` runs on every exit path — success, an item failure, or a `Setup`
failure (partial acquisition still gets released). The runner guarantees it via
deferred-registration order (CleanUp registered first → runs last; rollback
registered second → runs before CleanUp).

## `infra.Ensure` — idempotent provisioning

Per requested service, `Setup`:

- **probes readiness** (HTTP `GET {url}/health` for HTTP services; a TCP dial for
  `broker`/`nats`; the health endpoint for `minio`);
- if **already up** → reuses it (marks nothing);
- if **down** → `docker compose up -d <service>` and waits until healthy, marking
  the service as **owned by this run**.

`CleanUp` stops **only the services this run started** (`docker compose stop`, in
reverse). A service that was already running is never touched. When the whole
stack is already up, `Setup` changes nothing and `CleanUp` is a no-op — so the
same journey runs unchanged against an externally-managed stack (CI) or a bare
machine. The owned set is held in the fixture's closure (Setup fills it, CleanUp
reads it), so ownership is per-run and never leaks across journeys.

**Ownership caveats (honest limits):**

- **Sequential-only.** Ownership is per-run with no cross-process coordination, so
  reuse is safe only when journeys run sequentially or within one runner process (the
  suite runner). Two independent `go test` processes could stop a service the other is
  using; a compose-file lock + refcount is required before parallel *processes* are
  safe. (Parallel journeys *within* the single runner are fine — one owner.)
- **Transitive deps are NOT owned.** `docker compose up -d <service>` also starts that
  service's `depends_on` chain (mongo, nats, …), which is not tracked as owned, so
  `CleanUp` leaves those running. "Stops only what it started" is scoped to the
  directly-requested services, not their dependencies.

## `infra.EnsureAll` — the suite entry point (two stacks)

The runner's `TestMain` calls `EnsureAll()`, which provisions **two independent
stacks** once and returns a teardown that stops only what it started:

| Stack | Services | Compose file | Default mode |
|---|---|---|---|
| **infra** | nats, broker, minio | `mapexOSDeploy/infra/docker-compose.yml` | `docker` |
| **mapex** | mapex-iam, http-gateway, assets, router, events, triggers, workflow, lns, js-executor | `mapexOSDeploy/services/docker-compose.yml` | `local` |

Each stack has a mode, env-overridable:

- `MAPEX_E2E_INFRA_MODE` / `MAPEX_E2E_MAPEX_MODE` = `docker` | `local`.
- **A service already answering its health probe is ALWAYS reused (skipped)** —
  regardless of mode.
- A **down** service is handled by the mode: `docker` → `docker compose up` it via that
  stack's compose file (and teardown stops it); `local` → it is your job to run it
  (`go run`), so a down one is reported as a clear "start it" error, never brought up.

So the common dev flow — infra in Docker, the app services from source (`go run`) —
is the default: EnsureAll brings up any missing infra and reuses your running apps.
CI sets both modes to `docker`. Because it runs from `TestMain` (no `*testing.T`),
EnsureAll logs a per-stack `reused / brought-up / failed` summary via the std logger.
Compose files: `MAPEX_COMPOSE_FILE` (infra) / `MAPEX_SERVICES_COMPOSE_FILE` (services)
override the runtime-resolved defaults.

## Configuration

- **Compose file**: `MAPEX_COMPOSE_FILE` env; otherwise resolved at runtime from the
  module root (walk up to `go.mod`) → sibling
  `mapexOSDeploy/infra/docker-compose.yml`. Anchored at the module root, not the
  process working directory, so it resolves to the same absolute path from any package.
- **Supported services** (logical name → compose service): `mapexos`, `assets`,
  `http_gateway`, `router`, `events`, `triggers`, `workflow`,
  `broker`→`mapex-broker-mqtt`, `nats`→`nats-core`, `minio`. Service URLs come
  from `common/constants`.

Requires `docker` + `docker compose` on PATH only when a service actually needs
bringing up; the unit tests inject fakes for the readiness probe and the compose
runner, so they need neither docker nor the network.
