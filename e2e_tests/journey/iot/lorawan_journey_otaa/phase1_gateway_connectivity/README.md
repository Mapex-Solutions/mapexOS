# phase1_gateway_connectivity

Gateway connectivity lifecycle over UDP + Basics Station. Provisions a gateway
asset (no template / route group — a gateway needs neither), asserts its
LNS-facing L3 auth projection, connects the simulated gateway to mapexLNS to go
online, then forces it offline.

Full outcome (PASS / FAIL) is in the `journey.go` package godoc.

```bash
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase1_gateway_connectivity/
```
