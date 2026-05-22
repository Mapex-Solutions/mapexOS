# Phase 1 — HTTP connectivity-driven Workflow execution

## What this test proves

The healthmonitor → router → workflow chain fires end-to-end for an
HTTP-protocol asset. Forced health transitions on the asset surface a
workflow execution on the events service for both the offline and the
online route groups.

The phase:

1. Creates a workflow definition and one workflow instance.
2. Creates two `kind=workflow` route groups (one for `online`, one for `offline`) targeting the same instance.
3. Creates an HTTP data source (push-mode + apiKey auth) and an asset template.
4. Creates an HTTP connectivity asset wired to both route groups.
5. Sends a warm-up heartbeat → asset settles to `online` (silent — no workflow fires).
6. Force-offline by admin → offline RG fires → events service surfaces **workflow execution 1** scoped after the force-offline timestamp.
7. Sends a new heartbeat → asset goes `online` → online RG fires → events service surfaces **workflow execution 2** scoped after the heartbeat timestamp.
8. Deletes the asset; Compensate chain rolls everything else back.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_http/phase1_workflow/...
```

## Requirements

- Live stack with these services running (defaults): `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `events:5004`, `workflow:5005`. Check with `./run-tests.sh check`.
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
- The asset uses HealthMonitor explicit mode; the heartbeat reaches the asset directly via the gateway, no scheduled scan needed.
