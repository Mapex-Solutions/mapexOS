# phase1_gateway_connectivity

Ciclo de conectividade do gateway em UDP + Basics Station. Provisiona um gateway
asset (sem template / route group — gateway não precisa), valida a projeção L3
de auth que o LNS lê, conecta o gateway simulado ao mapexLNS para ficar online,
depois força offline.

O outcome completo (PASS / FAIL) está no godoc do package em `journey.go`.

```bash
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase1_gateway_connectivity/
```
