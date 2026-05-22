# Phase 1 — Connectivity-driven Teams trigger smoke

## What this test proves

Teams trigger fires through the healthmonitor → router → triggers
chain and POSTs its MessageCard payload to an HTTP listener.

Teams webhooks are HTTP POSTs, so this phase reuses the same shared
HTTP sink. Two real health transitions → two POSTs received.

## How to run

```bash
cd e2e_tests
go test ./journey/automations/trigger_teams/phase1_connectivity/...
```

Or:

```bash
./run-tests.sh saga trigger-teams
```

## Requirements

- Live stack: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`.
- Free port `11010` on the host (shared HTTP sink — override `SAGA_TRIGGER_SINK_BIND_ADDR`).
- When the triggers service runs in Docker, set `SAGA_TRIGGER_SINK_URL=http://host.docker.internal:11010`.
- Seed admin user provisioned (`admin@mapex.local`).
