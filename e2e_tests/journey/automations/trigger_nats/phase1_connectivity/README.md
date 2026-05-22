# Phase 1 — Connectivity-driven NATS trigger smoke

## What this test proves

NATS trigger fires through the healthmonitor → router → triggers
chain and publishes to a NATS subject. Validation uses the
events_trigger oracle: success=true is set after the NATS publish
call returns. Two real health transitions → two successful executions.

## How to run

```bash
cd e2e_tests
go test ./journey/automations/trigger_nats/phase1_connectivity/...
```

Or:

```bash
./run-tests.sh saga trigger-nats
```

## Requirements

- Live stack: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`, `events:5004`.
- NATS server reachable from the triggers service at `nats://localhost:4222` (override `server` in `SagaNatsTrigger` for other targets).
- Seed admin user provisioned (`admin@mapex.local`).
