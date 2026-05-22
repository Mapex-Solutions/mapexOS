# Fase 1 — Smoke do trigger HTTP pela conectividade

## O que este teste prova

Trigger HTTP dispara pela cadeia healthmonitor → router → triggers e
faz POST para um servidor HTTP in-process.

A fase sobe um sink HTTP local, cria um trigger HTTP apontando para
ele, liga dois route groups (`online` / `offline`) ao trigger, e então
guia um asset por duas transições reais de saúde, verificando que o
sink recebeu um POST por transição.

## Como rodar

```bash
cd e2e_tests
go test ./journey/automations/trigger_http/phase1_connectivity/...
```

Ou:

```bash
./run-tests.sh saga trigger-http
```

## Requisitos

- Stack viva: `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`. Verifique com `./run-tests.sh check`.
- Porta `11010` livre no host (sobrescreva `SAGA_TRIGGER_SINK_BIND_ADDR`).
- Quando o serviço triggers roda em Docker, defina `SAGA_TRIGGER_SINK_URL=http://host.docker.internal:11010` para o trigger alcançar o host.
- Usuário admin seed provisionado (`admin@mapex.local`).
