# LoRaWAN gateway presence — online via LNS connect / offline via force

## What this test proves

A LoRaWAN **gateway** asset's connectivity state tracks reality, over **both**
transports (Semtech UDP and Basics Station key mode). A gateway asset requires
neither a template nor a route group, so the journey provisions the gateway alone.

Online comes from the real path: the simulated gateway connects to mapexLNS and
its keepalive/stats make the LNS Gateway Server publish a `connect` presence
advisory, which the assets healthmonitor consumes to flip the asset online.
Offline is forced via the internal ops endpoint (`force_offline`), the same hook
the connectivity journeys use.

The journey, per transport block:

1. Creates a LoRaWAN gateway asset (UDP=eui, Basics Station=key) — no template, no route group.
2. Asserts the gateway's L3 auth projection (the shape the LNS reads on connect).
3. Connects the simulated gateway to mapexLNS (`lorawansim`).
4. Asserts `healthStatus=online` (LNS connect advisory consumed).
5. `ForceOfflineByAdmin` (internal `/internal/health_monitor/{uuid}/force_offline`).
6. Asserts `healthStatus=offline`.
7. (compensation) Closes the simulator link and deletes the gateway asset.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/lorawan_gateway_presence/
```

## Requirements

- Live stack: `mapexos`, `assets` on default ports.
- **mapexLNS** running, reachable at the UDP + Basics Station ingresses:
  - `LNS_UDP_HOST` (default `127.0.0.1`), `LNS_UDP_PORT` (default `1700`)
  - `LNS_BSTATION_URI` (default `ws://127.0.0.1:8887`)
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
