# lorawan_journey

Ciclo de vida completo do device LoRaWAN, ponta a ponta contra o stack vivo,
dirigindo tráfego REAL de gateway + sensor através do mapexLNS (overlay do The
Things Stack) nos DOIS transportes de rádio — Semtech UDP e Basics Station key
mode.

Toda phase primeiro compõe o bootstrap IAM compartilhado
(`common/journey/iam_bootstrap`): login do admin seed → validade do JWT →
coverage de org.

## Phases

| Phase | O que cobre |
|-------|-------------|
| `phase1_gateway_connectivity` | Provisiona um gateway (UDP + BS), valida a projeção L3 de auth que o LNS lê, conecta → online, força offline → offline. |
| `phase2_sensor_uplink` | Provisiona route group + template de codec + gateway + sensor OTAA; conecta, faz join OTAA, dispara um uplink real, e valida que o MapexOS ingere (bytes crus + fPort/fCnt/rxInfo) e o sensor fica online. |
| `phase3_sensor_presence` | Provisiona + join de um sensor; um uplink real o traz online, depois força offline. |

Roadmap (phases futuras, precisam de blocos novos): negativo "recusa se GW não
registrado", ativação ABP, classe B / classe C.

## Como rodar

A partir da raiz do pacote `e2e_tests`:

```bash
# Todas as phases desta journey
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/...

# Uma phase específica
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase1_gateway_connectivity/
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase2_sensor_uplink/
go test -tags=saga -v ./journey/iot/lorawan_journey_otaa/phase3_sensor_presence/
```

## Pré-requisitos

O stack vivo precisa estar de pé: assets, mapexIam, router, events e **mapexLNS**
(o LoRaWAN Network Server). O `phase2_sensor_uplink` valida a ingestão via
`GET /api/v1/events/raw`, então o events service + ClickHouse precisam estar
saudáveis.

## Limitação conhecida

O `phase2_sensor_uplink` está estruturalmente completo mas hoje falha no assert
de ingestão: os bytes do raw event ainda não batem com o payload disparado por
causa de um mismatch de decrypt do FRMPayload (session-key AppSKey) no mapexLNS.
Esse é um fix de plataforma em trilha separada — o wiring da journey está certo.
