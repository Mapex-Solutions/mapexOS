# phase2_sensor_uplink

Full LoRaWAN ingestion path over UDP + Basics Station: provision route group +
codec template + gateway + OTAA sensor, connect, OTAA-join, fire a real uplink
through mapexLNS, and assert MapexOS ingests it (raw bytes + fPort/fCnt/rxInfo)
and the sensor goes online.

Full outcome (PASS / FAIL) is in the `journey.go` package godoc.

Note: currently fails at the ingestion assert due to a LoRaWAN AppSKey FRMPayload
decrypt mismatch in mapexLNS (platform fix tracked separately) — the wiring is
correct.

```bash
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase2_sensor_uplink/
```
