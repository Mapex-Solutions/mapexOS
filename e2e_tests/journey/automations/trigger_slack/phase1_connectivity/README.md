# Phase 1 — Connectivity-driven Slack trigger smoke

## What this test proves

Slack trigger fires through the healthmonitor → router → triggers
chain and POSTs the webhook to an HTTP listener.

Slack webhooks are plain HTTP POSTs, so this phase reuses the same
in-process HTTP sink the HTTP trigger phase uses. Two real health
transitions → two webhook POSTs received.

## How to run

```bash
cd e2e_tests
go test ./journey/automations/trigger_slack/phase1_connectivity/...
```

Or:

```bash
./run-tests.sh saga trigger-slack
```

## Requirements

- Live stack: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`. Check with `./run-tests.sh check`.
- Free port `11010` on the host (shared HTTP sink — override `SAGA_TRIGGER_SINK_BIND_ADDR`).
- When the triggers service runs in Docker, set `SAGA_TRIGGER_SINK_URL=http://host.docker.internal:11010`.
- Seed admin user provisioned (`admin@mapex.local`).
