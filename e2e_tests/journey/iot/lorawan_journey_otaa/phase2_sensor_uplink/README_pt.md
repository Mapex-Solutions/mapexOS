# phase2_sensor_uplink

Caminho completo de ingestão LoRaWAN em UDP + Basics Station: provisiona route
group + template de codec + gateway + sensor OTAA, conecta, faz join OTAA,
dispara um uplink real pelo mapexLNS, e valida que o MapexOS ingere (bytes crus +
fPort/fCnt/rxInfo) e o sensor fica online.

O outcome completo (PASS / FAIL) está no godoc do package em `journey.go`.

Nota: hoje falha no assert de ingestão por um mismatch de decrypt do FRMPayload
(AppSKey) no mapexLNS (fix de plataforma em trilha separada) — o wiring está
certo.

```bash
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase2_sensor_uplink/
```
