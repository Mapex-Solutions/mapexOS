# phase3_sensor_presence

Ciclo de presença do sensor em UDP + Basics Station: provisiona route group +
template de codec + gateway + sensor OTAA, conecta e faz join OTAA, dispara um
uplink real que traz o sensor asset online, depois força offline.

O outcome completo (PASS / FAIL) está no godoc do package em `journey.go`.

```bash
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase3_sensor_presence/
```
