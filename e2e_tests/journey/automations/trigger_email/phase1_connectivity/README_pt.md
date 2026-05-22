# Fase 1 — Smoke do trigger Email pela conectividade

## O que este teste prova

O trigger Email dispara pela cadeia healthmonitor → router → triggers
e entrega uma mensagem real para um servidor SMTP in-process.

A fase:

1. Sobe um servidor SMTP in-process em `SmtpSinkBindAddr` (default `0.0.0.0:11025`).
2. Cria um trigger Email apontando para esse sink.
3. Cria dois route groups `kind=trigger` (um para `online`, outro para `offline`).
4. Cria um asset HTTP de conectividade ligado aos dois route groups.
5. Envia um heartbeat de aquecimento → asset assenta em `online` (silencioso — sem trigger).
6. Force-offline do admin → RG de offline dispara → sink SMTP captura **1 mensagem**.
7. Verifica que a mensagem capturada tem o `from`, `to` configurados e um fragmento do subject.
8. Envia um novo heartbeat → asset volta para `online` → RG de online dispara → sink SMTP captura **2 mensagens**.
9. Apaga o asset; a cadeia de Compensate desfaz o resto.

## Como rodar

```bash
cd e2e_tests
go test ./journey/automations/trigger_email/phase1_connectivity/...
```

Ou via helper:

```bash
./run-tests.sh saga trigger-email
```

## Requisitos

- Stack viva com estes services rodando (defaults): `mapexos:5000`, `assets:5002`, `router:5003`, `http_gateway:5001`, `triggers:5006`. Verifique com `./run-tests.sh check`.
- Porta `11025` livre no host (sobrescreva com `SAGA_TRIGGER_SMTP_BIND_ADDR`).
- Quando o serviço triggers roda em Docker, defina `SAGA_TRIGGER_SMTP_HOST=host.docker.internal` para o `smtpHost` do trigger alcançar o listener SMTP no host.
- Usuário admin seed provisionado (`admin@mapex.local`) — a phase 0 faz login como ele.
