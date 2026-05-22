# Fase 1 — Execução de Workflow pela conectividade MQTT

## O que este teste prova

A cadeia healthmonitor → router → workflow dispara ponta a ponta para
um asset de protocolo MQTT. Eventos de CONNECT e DISCONNECT fazem
aparecer execuções de workflow no serviço de events tanto para o route
group de offline quanto para o de online.

A fase:

1. Cria uma workflow definition e uma workflow instance.
2. Cria dois route groups `kind=workflow` (um para `online`, outro para `offline`) apontando para a mesma instance.
3. Cria um asset template.
4. Cria um asset MQTT de conectividade ligado aos dois route groups.
5. CONNECT de aquecimento → asset assenta em `online` (silencioso — primeira observação unknown→online, sem workflow).
6. DISCONNECT → asset transita para `offline` → RG de offline dispara → events service expõe a **execução de workflow 1** filtrada após o timestamp do disconnect.
7. CONNECT de novo → asset volta para `online` → RG de online dispara → events service expõe a **execução de workflow 2** filtrada após o timestamp do reconnect.
8. Apaga o asset; a cadeia de Compensate desfaz o resto.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_mqtt/phase1_workflow/...
```

## Requisitos

- Stack viva com estes services rodando (defaults): `mapexos:5000`, `assets:5002`, `router:5003`, `events:5004`, `workflow:5005`. Verifique com `./run-tests.sh check`.
- Broker MQTT acessível em `tcp://localhost:1883` (listener de senha).
- Usuário admin seed provisionado (`admin@mapex.local`) — a phase 0 (IAM bootstrap) faz login como ele.
