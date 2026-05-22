# Journey: Ações de conectividade MQTT

## O que esta journey prova

Para assets de protocolo MQTT, a cadeia healthmonitor → router dispara
ponta a ponta os dois kinds de router permitidos na superfície do
HealthMonitor (`kind=workflow` e `kind=trigger`) quando a saúde do
asset transita entre `online` e `offline`. Cada fase exercita um kind
no mesmo formato de asset e na mesma mecânica de route group.

| Fase | O que cobre |
|---|---|
| [`phase1_workflow`](./phase1_workflow/) | Transições de CONNECT / DISCONNECT → execução de workflow aparece no serviço de events. |
| [`phase2_trigger`](./phase2_trigger/) | Transições de CONNECT / DISCONNECT → trigger dispara e cai no sink HTTP in-process. |

## Como rodar

```bash
cd e2e_tests

# Todas as fases
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_mqtt/...

# Uma fase
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_mqtt/phase1_workflow/
```

## Requisitos

Cada README de fase lista os requisitos específicos. Comuns:

- Stack viva: `mapexos`, `assets`, `router` nas portas default.
- Broker MQTT acessível em `tcp://localhost:1883` (listener de senha).
- Usuário admin seed provisionado (`admin@mapex.local`) — a phase 0 (IAM bootstrap) faz login como ele.
