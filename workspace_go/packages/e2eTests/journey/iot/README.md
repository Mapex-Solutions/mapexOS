# Context: iot

Saga journeys that exercise the IoT pipeline — assets, route groups, asset
templates, MQTT auth callout, telemetry, healthmonitor, triggers, workflow
and the events sink. Every journey in this context assumes the seed
admin user is signed in (Phase 0 of every journey performs that).

## Registered journeys

| Journey                | Phases                                  | What the story covers                                                |
|------------------------|-----------------------------------------|----------------------------------------------------------------------|
| mqtt_full_pipeline     | phase0_iam_bootstrap, phase1_iot_setup  | MQTT happy path: tenant scaffolding -> route group + template + asset |

Phase folder names carry the descriptor; package names inside them are
short (`phase0`, `phase1`) so the import alias makes intent obvious.

## How to run

Every command runs from the e2eTests package root.

```bash
cd workspace_go/packages/e2eTests

# All phases of every journey in this context
go test -tags=saga -v ./journey/iot/...

# Every phase of one journey
go test -tags=saga -v ./journey/iot/mqtt_full_pipeline/...

# A single phase of one journey
go test -tags=saga -v ./journey/iot/mqtt_full_pipeline/phase0_iam_bootstrap/
go test -tags=saga -v ./journey/iot/mqtt_full_pipeline/phase1_iot_setup/
```

The `saga` build tag gates these tests: `go test ./...` (no tag) skips
them; only `go test -tags=saga` walks the journey folders.

## Required environment

- mapexIam   reachable at `MAPEXOS_URL` (default `http://localhost:5000`)
- assets     reachable at `ASSETS_URL`  (default `http://localhost:5002`)
- router     reachable at `ROUTER_URL`  (default `http://localhost:5003`)
- http_gateway reachable at `GATEWAY_URL` (default `http://localhost:5001`)
- Seed admin user provisioned by the canonical mongodb-init seed
  (`admin@mapex.local` / `mapex@123`)

## How to add a new journey to this context

1. `mkdir journey/iot/<journey_name>` (snake_case).
2. Create `README.md` with the narrative, phase index, and a copy of the
   "How to run" block tailored to the journey path.
3. Add `phaseN_<descriptor>/journey.go` and `journey_test.go` for each
   phase. Reuse Phase 0 of an existing journey as the bootstrap when the
   new journey starts from the same actor; otherwise write its own.
4. Add a row to the registered-journeys table above.
