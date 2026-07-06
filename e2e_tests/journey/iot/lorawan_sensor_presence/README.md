# LoRaWAN sensor presence — online by data / offline via force

## What this test proves

A LoRaWAN **sensor** asset's connectivity state tracks reality, over **both**
transports (Semtech UDP and Basics Station key mode). A sensor comes online from
real data — an uplink is a presence signal (implicit healthmonitor mode) — and is
forced offline via the internal ops endpoint.

The journey, per transport block:

1. (shared) Creates a route group and the LoRaWAN codec asset template.
2. Creates a gateway asset (to carry the uplink) and a sensor asset.
3. Connects the simulated gateway to mapexLNS (`lorawansim`); OTAA-joins the sensor.
4. Fires one uplink (`0BB809F6025D0000000000`).
5. Asserts `healthStatus=online` (uplink = presence signal).
6. `ForceOfflineByAdmin` (internal `/internal/health_monitor/{uuid}/force_offline`).
7. Asserts `healthStatus=offline`.
8. (compensation) Tears down the sensor + gateway, the template, and the route
   group; closes the simulator links.

Each transport uses its own labelled gateway + sensor so bag keys and Compensate
teardown never collide.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/lorawan_sensor_presence/
```

## Requirements

- Live stack: `mapexos`, `assets`, `router`, `events` on default ports.
- **mapexLNS** running, reachable at the UDP + Basics Station ingresses:
  - `LNS_UDP_HOST` (default `127.0.0.1`), `LNS_UDP_PORT` (default `1700`)
  - `LNS_BSTATION_URI` (default `ws://127.0.0.1:8887`)
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
