# Presença de sensor LoRaWAN — online por dado / offline via force

## O que este teste prova

O estado de conectividade de um asset de **sensor** LoRaWAN acompanha a realidade,
sobre **ambos** os transportes (Semtech UDP e Basics Station key mode). Um sensor
vira online por dado real — um uplink é um sinal de presença (modo implícito do
healthmonitor) — e é forçado offline pelo endpoint interno de ops.

A journey, por bloco de transporte:

1. (compartilhado) Cria um route group e o asset template com codec LoRaWAN.
2. Cria um gateway asset (para carregar o uplink) e um sensor asset.
3. Conecta o gateway simulado ao mapexLNS (`lorawansim`); faz o join OTAA do sensor.
4. Dispara um uplink (`0BB809F6025D0000000000`).
5. Assere `healthStatus=online` (uplink = sinal de presença).
6. `ForceOfflineByAdmin` (interno `/internal/health_monitor/{uuid}/force_offline`).
7. Assere `healthStatus=offline`.
8. (compensação) Derruba o sensor + gateway, o template e o route group; fecha os
   links do simulador.

Cada transporte usa seu próprio gateway + sensor rotulados para as bag keys e o
teardown do Compensate nunca colidirem.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/lorawan_sensor_presence/
```

## Requisitos

- Stack vivo: `mapexos`, `assets`, `router`, `events` nas portas padrão.
- **mapexLNS** rodando, acessível nos ingressos UDP + Basics Station:
  - `LNS_UDP_HOST` (padrão `127.0.0.1`), `LNS_UDP_PORT` (padrão `1700`)
  - `LNS_BSTATION_URI` (padrão `ws://127.0.0.1:8887`)
- Usuário admin de seed provisionado (`admin@mapex.local`) — a fase 0 (bootstrap IAM) loga como ele.
