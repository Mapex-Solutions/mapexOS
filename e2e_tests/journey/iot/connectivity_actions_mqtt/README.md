# Journey: MQTT connectivity actions

## What this journey proves

For MQTT-protocol assets, the healthmonitor → router chain fires the
two router kinds permitted on the HealthMonitor surface
(`kind=workflow` and `kind=trigger`) end-to-end when the asset health
transitions between `online` and `offline`. Each phase exercises one
kind on the same shape of asset and the same route-group mechanics.

| Phase | What it covers |
|---|---|
| [`phase1_workflow`](./phase1_workflow/) | CONNECT / DISCONNECT transitions → workflow execution surfaces on the events service. |
| [`phase2_trigger`](./phase2_trigger/) | CONNECT / DISCONNECT transitions → trigger fires and lands on the in-process HTTP sink. |

## How to run

```bash
cd e2e_tests

# Every phase
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_mqtt/...

# One phase only
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_mqtt/phase1_workflow/
```

## Requirements

Each phase README lists the protocol-specific requirements. Common:

- Live stack: `mapexos`, `assets`, `router` on default ports.
- MQTT broker reachable on `tcp://localhost:1883` (password listener).
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
