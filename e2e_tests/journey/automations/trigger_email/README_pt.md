# Journey: Trigger Email

## O que esta journey prova

O trigger Email (SMTP) do serviço triggers funciona ponta a ponta.
Cada fase exercita um caminho de disparo diferente; ambas entregam
uma mensagem real para um servidor SMTP in-process e validam
envelope + subject.

| Fase | Caminho de disparo |
|---|---|
| [`phase1_connectivity`](./phase1_connectivity/) | Force-offline / force-online via healthmonitor → trigger dispara. |
| [`phase2_event_pipeline`](./phase2_event_pipeline/) | POST telemetria no gateway → js-executor → router → trigger dispara. |

## Como rodar

```bash
cd e2e_tests
./run-tests.sh saga trigger-email
```

## Requisitos

Cada README de fase lista requisitos específicos do protocolo. Comuns:

- Stack viva: `mapexos`, `assets`, `router`, `http_gateway`, `triggers` nas portas default.
- Porta `11025` livre no host (sobrescreva `SAGA_TRIGGER_SMTP_BIND_ADDR`).
