# Fase 2 — Ciclo de vida completo de auth MQTT por certificado (mTLS)

## O que este teste prova

O pipeline de autenticação MQTT por certificado (mTLS) funciona ponta
a ponta contra o broker vivo. O fluxo espelha a fase de senha, mas o
device apresenta um certificado emitido na conexão contra o listener
mTLS do broker.

A fase:

1. Cria um route group para o asset.
2. Cria um asset template (schema de temperatura).
3. Cria o asset com `protocol=mqtt` + `authType=cert` usando a variante de payload de cert (carrega `certTTL`).
4. Emite um cert via `POST /api/v1/mqtt_certs`; o bundle PEM cai no bag e `asset.currentCert` é persistido; invalidação FANOUT dispara.
5. CONNECT mTLS no listener `:8883` do broker com o cert recém-emitido.
6. Verifica `healthStatus=online` (advisory de presença consumido, status persistido).
7. PUBLISH em `events/{assetUUID}/temperature`; o ACL do broker precisa aceitar o tópico com assetUUID puro.
8. Verifica que o serviço de events expõe a linha via `/api/v1/events/raw`.
9. Disconnect MQTT limpo.
10. DELETE no asset; invalidação FANOUT chega ao plugin do broker.
11. Verifica que um CONNECT mTLS novo com o MESMO cert é negado — asset removido, serial de `currentCert` removido, auth modo cert falha.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/mqtt_broker_auth/phase2_cert_user/
```

## Requisitos

- Stack viva: `mapexos`, `assets`, `router`, `events` nas portas default, mais `mapexVault` saudável e assets MS reportando `caReady=true`.
- Broker MQTT acessível em `ssl://localhost:8883` (`MQTT_BROKER_TLS_URL`).
- PKI pré-construída: `./scripts/prebuild/pki/generate-pki.sh`.
- Usuário admin seed provisionado (`admin@mapex.local`) — a phase 0 (IAM bootstrap) faz login como ele.
- `mqtt_broker_auth/README.md` lista todas as overrides de URL por serviço.
