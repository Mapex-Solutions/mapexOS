# phase3_sensor_presence

Sensor presence lifecycle over UDP + Basics Station: provision route group +
codec template + gateway + OTAA sensor, connect and OTAA-join, fire a real uplink
that drives the sensor asset online, then force it offline.

Full outcome (PASS / FAIL) is in the `journey.go` package godoc.

```bash
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase3_sensor_presence/
```
