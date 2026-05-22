# Phase 1 — Connectivity-driven HTTP trigger smoke

## What this test proves

HTTP trigger fires through the healthmonitor → router → triggers
chain and POSTs to an in-process HTTP server.

The phase starts a local HTTP sink, creates an HTTP trigger pointing
at it, wires two route groups (`online` / `offline`) to the trigger,
then drives an asset through two real health transitions and checks
the sink received one POST per transition.

## How to run

```bash
cd e2e_tests
go test ./journey/automations/trigger_http/phase1_connectivity/...
```

Or:

```bash
./run-tests.sh saga trigger-http
```

## Requirements

- Live stack: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`. Check with `./run-tests.sh check`.
- Port `11010` free on the host (override `SAGA_TRIGGER_SINK_BIND_ADDR`).
- When the triggers service runs in Docker, set `SAGA_TRIGGER_SINK_URL=http://host.docker.internal:11010` so the trigger can reach the host.
- Seed admin user provisioned (`admin@mapex.local`).
