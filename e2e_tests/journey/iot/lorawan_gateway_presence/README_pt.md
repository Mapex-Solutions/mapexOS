# Presença de gateway LoRaWAN — online via connect no LNS / offline via force

## O que este teste prova

O estado de conectividade de um asset de **gateway** LoRaWAN acompanha a realidade,
sobre **ambos** os transportes (Semtech UDP e Basics Station key mode). Um gateway
asset não exige template nem route group, então a journey provisiona só o gateway.

O online vem do caminho real: o gateway simulado conecta no mapexLNS e seu
keepalive/stats fazem o LNS Gateway Server publicar um advisory de presença
`connect`, que o healthmonitor dos assets consome para virar o asset online. O
offline é forçado pelo endpoint interno de ops (`force_offline`), o mesmo hook que
as journeys de conectividade usam.

A journey, por bloco de transporte:

1. Cria um gateway asset LoRaWAN (UDP=eui, Basics Station=key) — sem template, sem route group.
2. Assere a projeção L3 de auth do gateway (o shape que o LNS lê no connect).
3. Conecta o gateway simulado ao mapexLNS (`lorawansim`).
4. Assere `healthStatus=online` (advisory de connect do LNS consumido).
5. `ForceOfflineByAdmin` (interno `/internal/health_monitor/{uuid}/force_offline`).
6. Assere `healthStatus=offline`.
7. (compensação) Fecha o link do simulador e deleta o gateway asset.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/lorawan_gateway_presence/
```

## Requisitos

- Stack vivo: `mapexos`, `assets` nas portas padrão.
- **mapexLNS** rodando, acessível nos ingressos UDP + Basics Station:
  - `LNS_UDP_HOST` (padrão `127.0.0.1`), `LNS_UDP_PORT` (padrão `1700`)
  - `LNS_BSTATION_URI` (padrão `ws://127.0.0.1:8887`)
- Usuário admin de seed provisionado (`admin@mapex.local`) — a fase 0 (bootstrap IAM) loga como ele.
