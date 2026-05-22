# Journey: Ações de conectividade HTTP

## O que esta journey prova

Para assets de protocolo HTTP, a cadeia healthmonitor → router dispara
ponta a ponta os dois kinds de router permitidos na superfície do
HealthMonitor (`kind=workflow` e `kind=trigger`) quando a saúde do
asset transita entre `online` e `offline`. Cada fase exercita um kind
no mesmo formato de asset e na mesma mecânica de route group.

| Fase | O que cobre |
|---|---|
| [`phase1_workflow`](./phase1_workflow/) | Transições do healthmonitor → execução de workflow aparece no serviço de events. |
| [`phase2_trigger`](./phase2_trigger/) | Transições do healthmonitor → trigger dispara e cai no sink HTTP in-process. |

## Como rodar

```bash
cd e2e_tests

# Todas as fases
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_http/...

# Uma fase
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_http/phase1_workflow/
```

## Requisitos

Cada README de fase lista os requisitos específicos. Comuns:

- Stack viva: `mapexos`, `assets`, `router`, `http_gateway` nas portas default.
- Usuário admin seed provisionado (`admin@mapex.local`) — a phase 0 (IAM bootstrap) faz login como ele.
- Asset usa modo explicit do HealthMonitor; o heartbeat chega ao asset direto pelo gateway.
