# lorawan_journey_bsp

Ciclo de vida de device LoRaWAN em ABP (Activation By Personalization), ponta a
ponta contra o stack vivo, dirigindo tráfego REAL de gateway + sensor através do
mapexLNS nos DOIS transportes de rádio — Semtech UDP e Basics Station key mode.

Diferente do OTAA (`lorawan_journey_otaa`), ABP não faz join: o sensor é
provisionado com sessão fixa (DevAddr + NwkSKey + AppSKey) e ativado no lugar. O
gateway que carrega o uplink é levantado pela peça compartilhada
`common/journey/lorawan_gateway`; toda phase primeiro compõe o bootstrap IAM
compartilhado (`common/journey/iam_bootstrap`).

## Phases

| Phase | O que cobre |
|-------|-------------|
| `phase1_sensor_uplink` | Provisiona route group + template de codec + gateway + sensor ABP; conecta, ativa (sem join), dispara um uplink real, e valida que o MapexOS ingere (bytes crus + fPort/fCnt/rxInfo) e o sensor fica online. |
| `phase2_sensor_presence` | Provisiona + ativa um sensor ABP; um uplink real o traz online, depois força offline. |

## Como rodar

A partir da raiz do pacote `e2e_tests`:

```bash
# Todas as phases desta journey
go test -tags=saga -v ./journey/iot/lorawan_journey_bsp/...

# Uma phase específica
go test -tags=saga -v ./journey/iot/lorawan_journey_bsp/phase1_sensor_uplink/
go test -tags=saga -v ./journey/iot/lorawan_journey_bsp/phase2_sensor_presence/
```

## Pré-requisitos

O stack vivo precisa estar de pé: assets, mapexIam, router, events e **mapexLNS**.
O `phase1_sensor_uplink` valida a ingestão via `GET /api/v1/events/raw`, então o
events service + ClickHouse precisam estar saudáveis.

## Limitação conhecida

O `phase1_sensor_uplink` usa o mesmo assert de ingestão da journey OTAA, então
falha pelo mesmo motivo até o mismatch de decrypt do FRMPayload (AppSKey) no
mapexLNS ser corrigido (fix de plataforma em trilha separada) — o wiring está
certo.
