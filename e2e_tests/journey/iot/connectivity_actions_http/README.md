# Journey: HTTP connectivity actions

## What this journey proves

For HTTP-protocol assets, the healthmonitor → router chain fires the
two router kinds permitted on the HealthMonitor surface
(`kind=workflow` and `kind=trigger`) end-to-end when the asset health
transitions between `online` and `offline`. Each phase exercises one
kind on the same shape of asset and the same route-group mechanics.

| Phase | What it covers |
|---|---|
| [`phase1_workflow`](./phase1_workflow/) | Healthmonitor transitions → workflow execution surfaces on the events service. |
| [`phase2_trigger`](./phase2_trigger/) | Healthmonitor transitions → trigger fires and lands on the in-process HTTP sink. |

## How to run

```bash
cd e2e_tests

# Every phase
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_http/...

# One phase only
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_http/phase1_workflow/
```

## Requirements

Each phase README lists the protocol-specific requirements. Common:

- Live stack: `mapexos`, `assets`, `router`, `http_gateway` on default ports.
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
- Asset uses HealthMonitor explicit mode; heartbeat reaches the asset directly via the gateway.
