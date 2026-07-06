# LoRaWAN uplink ingestion — real GW + sensor → mapexLNS → MapexOS

## What this test proves

The full LoRaWAN ingestion path works end-to-end against the live stack, over
**both** radio transports (Semtech UDP and Basics Station key mode): a simulated
gateway + sensor send a *real* uplink — real OTAA crypto, real PHYPayload, real
gateway framing — through **mapexLNS**, and MapexOS ingests it.

What arrives is not just the hex: the on-the-wire uplink carries radio metadata
around the PHYPayload; mapexLNS validates the MIC, decrypts the FRMPayload, and
republishes it on `mapexos.lorawan.data.*`; js-executor runs the template codec
and records a raw event. The assert checks the **raw** event (undecoded bytes ==
the fired hex, plus `fPort`/`fCnt`/`rxInfo`) — the codec decode itself is a
downstream concern and is not asserted here.

Each transport uses its own labelled gateway + sensor (`gwUdp`/`snUdp`,
`gwBs`/`snBs`) so bag keys and Compensate teardown never collide.

The journey, per transport block:

1. (shared) Creates a route group and the LoRaWAN codec asset template.
2. Creates a LoRaWAN gateway asset (UDP=eui, Basics Station=key) and a sensor asset.
3. Asserts the gateway's L3 auth projection (the shape the LNS reads on connect).
4. Connects the simulated gateway to mapexLNS (`lorawansim`); OTAA-joins the sensor.
5. Fires one uplink (`0BB809F6025D0000000000`).
6. Asserts MapexOS ingested it — a raw event for the sensor's `threadId` whose
   `bytes` equal the fired hex, carrying `fPort`/`fCnt`/`rxInfo`.
7. Asserts the sensor asset transitions to `online` (uplink = presence signal).
8. (compensation) Tears down every asset, the template, and the route group; closes
   the simulator links.

## How to run

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/lorawan_uplink_ingest/
```

## Requirements

- Live stack: `mapexos`, `assets`, `router`, `events` on default ports.
- **mapexLNS** running, reachable at the UDP + Basics Station ingresses:
  - `LNS_UDP_HOST` (default `127.0.0.1`), `LNS_UDP_PORT` (default `1700`)
  - `LNS_BSTATION_URI` (default `ws://127.0.0.1:8887`)
- Seed admin user provisioned (`admin@mapex.local`) — phase 0 (IAM bootstrap) logs in as that.
