# lorawan_journey

Full LoRaWAN device lifecycle end to end against the live stack, driving REAL
gateway + sensor traffic through mapexLNS (The Things Stack overlay) over BOTH
radio transports — Semtech UDP and Basics Station key mode.

Every phase first composes the shared IAM bootstrap
(`common/journey/iam_bootstrap`): seed admin login → JWT validity → org coverage.

## Phases

| Phase | What it covers |
|-------|----------------|
| `phase1_gateway_connectivity` | Provision a gateway (UDP + BS), assert its LNS-facing L3 auth projection, connect → online, force offline → offline. |
| `phase2_sensor_uplink` | Provision route group + codec template + gateway + OTAA sensor; connect, OTAA-join, fire a real uplink, assert MapexOS ingests it (raw bytes + fPort/fCnt/rxInfo) and the sensor goes online. |
| `phase3_sensor_presence` | Provision + join a sensor; a real uplink drives it online, then force it offline. |

Roadmap (future phases, need new building blocks): refuse-if-gateway-not-registered
(negative), ABP activation, class B / class C.

## How to run

From the `e2e_tests` package root:

```bash
# All phases of this journey
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/...

# A single phase
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase1_gateway_connectivity/
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase2_sensor_uplink/
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase3_sensor_presence/
```

## Prerequisites

The live stack must be up: assets, mapexIam, router, events, and **mapexLNS**
(the LoRaWAN Network Server). `phase2_sensor_uplink` asserts ingestion via
`GET /api/v1/events/raw`, so the events service + ClickHouse must be healthy.

## Known limitation

`phase2_sensor_uplink` is structurally complete but currently fails at the
ingestion assert: the raw event's bytes do not yet match the fired payload
because of a LoRaWAN session-key (AppSKey) FRMPayload decrypt mismatch in
mapexLNS. That is a platform fix tracked separately — the journey wiring is
correct.
