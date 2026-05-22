# Phase 1 — MQTT connectivity-driven Workflow execution

## What this test proves

The healthmonitor → router → workflow chain fires end-to-end for an
MQTT-protocol asset. CONNECT and DISCONNECT events surface workflow
executions on the events service for both the offline and the online
route groups.

The phase:

1. Creates a workflow definition and one workflow instance.
2. Creates two `kind=workflow` route groups (one for `online`, one for `offline`) targeting the same instance.
3. Creates an asset template.
4. Creates an MQTT connectivity asset wired to both route groups.
5. CONNECT warm-up → asset settles to `online` (silent — first observation is unknown→online, no workflow fires).
6. DISCONNECT → asset transitions to `offline` → offline RG fires → events service surfaces **workflow execution 1** scoped after the disconnect timestamp.
7. CONNECT again → asset transitions back to `online` → online RG fires → events service surfaces **workflow execution 2** scoped after the reconnect timestamp.
8. Deletes the asset; Compensate chain rolls everything else back.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_mqtt/phase1_workflow/...
```

## Requirements

- Live stack with these services running (defaults): `mapexos:5000`, `assets:5002`, `router:5003`, `events:5004`, `workflow:5005`. Check with `./run-tests.sh check`.
- MQTT broker reachable on `tcp://localhost:1883` (password listener).
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
