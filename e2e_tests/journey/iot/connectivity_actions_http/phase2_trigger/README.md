# Phase 2 — HTTP connectivity-driven Trigger execution

## What this test proves

The healthmonitor → router → triggers chain fires end-to-end for an
HTTP-protocol asset. Forced health transitions on the asset hit an
in-process HTTP sink — one hit per real transition.

The phase:

1. Starts an in-process HTTP sink server; every trigger fire is captured as a sink hit.
2. Creates a trigger pointing at the sink.
3. Creates two `kind=trigger` route groups (one for `online`, one for `offline`).
4. Creates an HTTP data source (push-mode + apiKey auth) and an asset template.
5. Creates an HTTP connectivity asset wired to both route groups.
6. Sends a warm-up heartbeat → asset settles to `online` (silent — no trigger fires).
7. Force-offline by admin → offline RG fires → sink captures **1 hit**.
8. Sends a new heartbeat → asset goes `online` → online RG fires → sink captures **2 hits**.
9. Deletes the asset; Compensate chain rolls everything else back.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_http/phase2_trigger/...
```

## Requirements

- Live stack with these services running (defaults): `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`. Check with `./run-tests.sh check`.
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
- The sink listens on a saga-owned host:port; when the triggers service runs in Docker, the trigger's target host must resolve back to the host running the saga.
