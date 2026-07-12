# Quick Start — E2E Tests

Fast path to running the MapexOS e2e tests. Everything runs from `e2e_tests/` with
plain `go test`. See `README.md` for the full layout and conventions.

## Prerequisites

The stack must be reachable. Run the app services locally (`go run` via the CLI) or
bring them up with the deploy compose. The saga runner also brings up any missing
infra itself (`infra.EnsureAll`), reusing whatever is already running and tearing
down only what it started — so a fully-up local stack needs no extra setup.

## TL;DR

```bash
cd e2e_tests

# Module e2e — HTTP contract of one module (no build tag)
go test ./services/...
go test ./services/mapexos/organizations -v

# Saga journeys — ONE runner brings the stack up once and runs every journey as a
# parallel subtest on ephemeral ports.
go test -tags=saga ./journey/suite/... -parallel 4

# A single journey — filter by its registry name (see journey/suite/suite_test.go)
go test -tags=saga ./journey/suite/... -run 'TestSuite/iot/gateway_provisioning'
```

## Two test families

| Family | Command | Notes |
|---|---|---|
| **Module e2e** | `go test ./services/{svc}/{mod}` | one module's HTTP contract; own `TestMain` login |
| **Saga journey** | `go test -tags=saga ./journey/suite/... -run TestSuite/<name>` | cross-service flow with rollback, run via the single runner |

There is **no per-journey test file** — journeys are registered once in
`journey/suite/suite_test.go`. Adding a journey is one aliased import + one registry
line.

## Config knobs

| Env | Default | Use |
|---|---|---|
| `MAPEXOS_URL`, `ASSETS_URL`, … | `localhost:500x` | service base URLs |
| `SAGA_SINK_HOST` | `localhost` | set to `host.docker.internal` when a Dockerized service must reach a host-side sink |
| `SAGA_TIMEOUT_MULTIPLIER` | `1.0` | stretch poll timeouts on a slow stack (only extends) |
| `MAPEX_COMPOSE_FILE` | discovered | point EnsureAll at the deploy compose |

## Notes

- **Single-tenant by design** — every journey runs as the seed admin; `runID`
  isolates data across parallel journeys. Not a defect.
- **Ephemeral ports** — sinks/servers bind OS-assigned free ports, so nothing needs a
  fixed port free.
- Longer runs: add `-timeout 20m`.
