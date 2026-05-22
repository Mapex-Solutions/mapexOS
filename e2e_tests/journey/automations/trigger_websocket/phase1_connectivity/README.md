# Phase 1 — Connectivity-driven WebSocket trigger smoke

## What this test proves

WebSocket trigger fires through the healthmonitor → router → triggers
chain and connects to an in-process WS server, writing a frame.

The phase starts a local WS sink on `WsSinkBindAddr` (`/ws`), creates
a WebSocket trigger pointing at it, wires two route groups
(`online` / `offline`) to the trigger, then drives an asset through
two real health transitions. The events_trigger oracle confirms the
trigger published successfully on each transition.

## How to run

```bash
cd e2e_tests
go test ./journey/automations/trigger_websocket/phase1_connectivity/...
```

Or:

```bash
./run-tests.sh saga trigger-websocket
```

## Requirements

- Live stack: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`, `events:5004`.
- Port `11026` free on the host (override `SAGA_TRIGGER_WS_BIND_ADDR`).
- When the triggers service runs in Docker, set `SAGA_TRIGGER_WS_URL=ws://host.docker.internal:11026/ws`.
- Seed admin user provisioned (`admin@mapex.local`).
