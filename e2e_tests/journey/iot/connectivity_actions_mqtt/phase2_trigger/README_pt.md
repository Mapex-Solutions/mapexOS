# Fase 2 — Execução de Trigger pela conectividade MQTT

## O que este teste prova

A cadeia healthmonitor → router → triggers dispara ponta a ponta para
um asset de protocolo MQTT. Eventos de CONNECT e DISCONNECT batem em
um sink HTTP in-process — um hit por transição real.

A fase:

1. Sobe um sink HTTP in-process; cada disparo do trigger é capturado como sink hit.
2. Cria um trigger apontando para o sink.
3. Cria dois route groups `kind=trigger` (um para `online`, outro para `offline`).
4. Cria um asset template.
5. Cria um asset MQTT de conectividade ligado aos dois route groups.
6. CONNECT de aquecimento → asset assenta em `online` (silencioso — primeira observação unknown→online, sem trigger).
7. DISCONNECT → asset transita para `offline` → RG de offline dispara → sink captura **1 hit**.
8. CONNECT de novo → asset volta para `online` → RG de online dispara → sink captura **2 hits**.
9. Apaga o asset; a cadeia de Compensate desfaz o resto.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_mqtt/phase2_trigger/...
```

## Requisitos

- Stack viva com estes services rodando (defaults): `mapexos:5000`, `assets:5002`, `router:5003`, `triggers:5006`. Verifique com `./run-tests.sh check`.
- Broker MQTT acessível em `tcp://localhost:1883` (listener de senha).
- Usuário admin seed provisionado (`admin@mapex.local`) — a phase 0 (IAM bootstrap) faz login como ele.
- O sink escuta em host:port do saga; quando o serviço triggers roda em Docker, o host configurado no trigger precisa resolver de volta para a máquina que está rodando o saga.
