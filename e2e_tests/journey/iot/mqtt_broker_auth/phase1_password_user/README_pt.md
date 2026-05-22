# Fase 1 — Ciclo de vida completo de auth MQTT por senha

## O que este teste prova

O pipeline de autenticação MQTT por senha funciona ponta a ponta
contra o broker vivo. A fase executa o fluxo completo que o operador
faz à mão para um asset com `authType=password`: cria → conecta →
presença flui → publica → events round-trip → apaga → reconectar deve
ser negado. A negação prova a invalidação L1 disparada por FANOUT no
plugin do broker.

A fase:

1. Cria um route group para o asset.
2. Cria um asset template (schema de temperatura).
3. Cria o asset com `protocol=mqtt` + `authType=password`; a senha em claro vai para o bag.
4. CONNECT MQTT no listener de senha com `(assetUUID, password)`.
5. Verifica `healthStatus=online` (advisory de presença consumido, status persistido, read model expõe).
6. PUBLISH em `events/{assetUUID}/temperature`; o ACL do broker precisa aceitar o tópico com assetUUID puro.
7. Verifica que o serviço de events expõe a linha via `/api/v1/events/raw`.
8. Disconnect MQTT limpo.
9. DELETE no asset; invalidação FANOUT chega ao plugin do broker.
10. Verifica que um CONNECT novo com as mesmas credenciais é negado.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/mqtt_broker_auth/phase1_password_user/
```

## Requisitos

- Stack viva: `mapexos`, `assets`, `router`, `events` nas portas default.
- Broker MQTT acessível em `tcp://localhost:1883` (`MQTT_BROKER_URL`).
- Usuário admin seed provisionado (`admin@mapex.local`) — a phase 0 (IAM bootstrap) faz login como ele.
- `mqtt_broker_auth/README.md` lista todas as overrides de URL por serviço.
