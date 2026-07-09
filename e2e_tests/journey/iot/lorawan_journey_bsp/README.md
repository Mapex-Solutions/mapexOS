# lorawan_journey_bsp

ABP (Activation By Personalization) LoRaWAN device lifecycle end to end against
the live stack, driving REAL gateway + sensor traffic through mapexLNS over BOTH
radio transports — Semtech UDP and Basics Station key mode.

Unlike OTAA (`lorawan_journey_otaa`), ABP needs no join handshake: the sensor is
provisioned with a fixed session (DevAddr + NwkSKey + AppSKey) and activated in
place. The gateway that carries the uplink is stood up by the shared
`common/journey/lorawan_gateway` building block; every phase first composes the
shared IAM bootstrap (`common/journey/iam_bootstrap`).

## Phases

| Phase | What it covers |
|-------|----------------|
| `phase1_sensor_uplink` | Provision route group + codec template + gateway + ABP sensor; connect, activate (no join), fire a real uplink, assert MapexOS ingests it (raw bytes + fPort/fCnt/rxInfo) and the sensor goes online. |
| `phase2_sensor_presence` | Provision + activate an ABP sensor; a real uplink drives it online, then force it offline. |

## How to run

From the `e2e_tests` package root:

```bash
# All phases of this journey
go test -tags=saga -v ./journey/iot/lorawan_journey_bsp/...

# A single phase
go test -tags=saga -v ./journey/iot/lorawan_journey_bsp/phase1_sensor_uplink/
go test -tags=saga -v ./journey/iot/lorawan_journey_bsp/phase2_sensor_presence/
```

## Prerequisites

The live stack must be up: assets, mapexIam, router, events, and **mapexLNS**.
`phase1_sensor_uplink` asserts ingestion via `GET /api/v1/events/raw`, so the
events service + ClickHouse must be healthy.

## Known limitation

`phase1_sensor_uplink` shares the OTAA journey's ingestion assert, so it fails
for the same reason until the LoRaWAN session-key (AppSKey) FRMPayload decrypt
mismatch in mapexLNS is fixed (platform fix tracked separately) — the journey
wiring is correct.
