# Fase 1 — Execução de Workflow pela conectividade HTTP

## O que este teste prova

A cadeia healthmonitor → router → workflow dispara ponta a ponta para
um asset de protocolo HTTP. Transições de saúde forçadas no asset
fazem aparecer uma execução de workflow no serviço de events tanto
para o route group de offline quanto para o de online.

A fase:

1. Cria uma workflow definition e uma workflow instance.
2. Cria dois route groups `kind=workflow` (um para `online`, outro para `offline`) apontando para a mesma instance.
3. Cria um data source HTTP (modo push + auth apiKey) e um asset template.
4. Cria um asset HTTP de conectividade ligado aos dois route groups.
5. Envia um heartbeat de aquecimento → asset assenta em `online` (silencioso — sem workflow).
6. Force-offline do admin → RG de offline dispara → events service expõe a **execução de workflow 1** filtrada após o timestamp do force-offline.
7. Envia um novo heartbeat → asset volta para `online` → RG de online dispara → events service expõe a **execução de workflow 2** filtrada após o timestamp do heartbeat.
8. Apaga o asset; a cadeia de Compensate desfaz o resto.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/connectivity_actions_http/phase1_workflow/...
```

## Requisitos

- Stack viva com estes services rodando (defaults): `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `events:5004`, `workflow:5005`. Verifique com `./run-tests.sh check`.
- Usuário admin seed provisionado (`admin@mapex.local`) — a phase 0 (IAM bootstrap) faz login como ele.
- O asset usa o modo explicit do HealthMonitor; o heartbeat chega ao asset direto pelo gateway, sem precisar do scan agendado.
