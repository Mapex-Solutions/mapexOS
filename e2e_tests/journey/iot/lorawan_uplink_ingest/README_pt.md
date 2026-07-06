# Ingestão de uplink LoRaWAN — GW + sensor reais → mapexLNS → MapexOS

## O que este teste prova

O caminho completo de ingestão LoRaWAN funciona ponta a ponta contra o stack vivo,
sobre **ambos** os transportes de rádio (Semtech UDP e Basics Station key mode): um
gateway + sensor simulados enviam um uplink *real* — crypto OTAA real, PHYPayload
real, framing de gateway real — através do **mapexLNS**, e o MapexOS o ingere.

O que chega não é só o hex: o uplink na parede carrega metadados de rádio em volta
do PHYPayload; o mapexLNS valida o MIC, descriptografa o FRMPayload e republica em
`mapexos.lorawan.data.*`; o js-executor roda o codec do template e grava um raw
event. O assert confere o raw event **cru** (bytes não decodificados == hex
disparado, mais `fPort`/`fCnt`/`rxInfo`) — o decode do codec é downstream e não é
asserido aqui.

Cada transporte usa seu próprio gateway + sensor rotulados (`gwUdp`/`snUdp`,
`gwBs`/`snBs`) para as bag keys e o teardown do Compensate nunca colidirem.

A journey, por bloco de transporte:

1. (compartilhado) Cria um route group e o asset template com codec LoRaWAN.
2. Cria um gateway asset LoRaWAN (UDP=eui, Basics Station=key) e um sensor asset.
3. Assere a projeção L3 de auth do gateway (o shape que o LNS lê no connect).
4. Conecta o gateway simulado ao mapexLNS (`lorawansim`); faz o join OTAA do sensor.
5. Dispara um uplink (`0BB809F6025D0000000000`).
6. Assere que o MapexOS ingeriu — um raw event no `threadId` do sensor cujos `bytes`
   são iguais ao hex disparado, carregando `fPort`/`fCnt`/`rxInfo`.
7. Assere que o sensor asset vira `online` (uplink = sinal de presença).
8. (compensação) Derruba todos os assets, o template e o route group; fecha os links
   do simulador.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/lorawan_uplink_ingest/
```

## Requisitos

- Stack vivo: `mapexos`, `assets`, `router`, `events` nas portas padrão.
- **mapexLNS** rodando, acessível nos ingressos UDP + Basics Station:
  - `LNS_UDP_HOST` (padrão `127.0.0.1`), `LNS_UDP_PORT` (padrão `1700`)
  - `LNS_BSTATION_URI` (padrão `ws://127.0.0.1:8887`)
- Usuário admin de seed provisionado (`admin@mapex.local`) — a fase 0 (bootstrap IAM) loga como ele.
