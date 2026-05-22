# Journey: MQTT Broker Auth — ciclos password + certificado

Cobertura end-to-end do pipeline de autenticação MQTT da plataforma
contra o stack `services_required` em execução. Cada fase roda o ciclo
completo de um asset — **criar → conectar → presença online →
publicar → ida e volta no events → deletar → reconectar tem que ser
negado** — provando o caminho feliz e também a invalidação de cache
disparada pelo FANOUT no tier L1 (Pebble) do plugin do broker.

## Contrato de fio (pós TKT-2026-0040)

- Username MQTT = `assetUUID` puro (sem o prefixo legado `{orgId}:`)
- Formato de tópico = `events/{assetUUID}/{eventType}`
- Subject.CN do certificado = `assetUUID` puro; o scoping de tenant
  flui do `asset.protocol.mqtt.orgId` server-side
- Listeners: `tcp://localhost:1883` (password) e
  `ssl://localhost:8883` (mTLS)

## Fases

| Fase                   | O que cobre                                                                |
|------------------------|----------------------------------------------------------------------------|
| `phase0_iam_bootstrap` | Login do seed admin → validade do JWT → cobertura do org-context          |
| `phase1_password_user` | Ciclo completo com password (10 passos)                                   |
| `phase2_cert_user`     | Ciclo completo com cert (mTLS) (11 passos)                                |
| `phase3_cascade`       | _(esqueleto — fora do escopo deste ticket)_                               |

### `phase1_password_user` — ordem dos passos

1. **CreateRouteGroup** — route group do asset
2. **CreateTemplate** — asset template (temperatura)
3. **CreateAsset** — asset persistido com `authType=password`
4. **ConnectMqttPassword** — CONNECT MQTT com `(assetUUID, password)`
5. **AssertHealthStatusEventually(online)** — presença flui ponta a ponta
6. **PublishTelemetry** — publish em `events/{assetUUID}/temperature`
7. **AssertRawEventReceivedAfter** — events expõe a linha
8. **DisconnectMqtt** — teardown MQTT limpo
9. **DeleteAsset** — fanout dispara a invalidação
10. **AssertConnectDeniedPassword** — CONNECT novo é negado

### `phase2_cert_user` — ordem dos passos

1. **CreateRouteGroup**
2. **CreateTemplate**
3. **CreateAssetWith(SagaMqttCertTemperatureSensor)** — `authType=cert`, `certTTL={1 dia}`
4. **IssueCert** — POST `/api/v1/mqtt_certs`; PEM bundle no bag, `asset.currentCert` persistido
5. **ConnectMqttCert** — handshake mTLS contra `:8883`
6. **AssertHealthStatusEventually(online)**
7. **PublishTelemetry**
8. **AssertRawEventReceivedAfter**
9. **DisconnectMqtt**
10. **DeleteAsset**
11. **AssertConnectDeniedCert** — CONNECT mTLS novo é negado

## Como rodar

```sh
cd workspace_go/packages/e2eTests

# Só a Fase 0 (smoke do bootstrap IAM)
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap/

# Ciclo password
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase1_password_user/

# Ciclo cert (mTLS)
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase2_cert_user/

# Tudo da journey
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/...
```

`-count=1` desabilita o cache de testes do Go pra que cada invocação
realmente bata no stack. Sem ele, uma segunda execução imprime
`(cached)` e devolve o veredito anterior sem tocar no broker.

## Pré-requisitos

1. PKI pré-construída:
   ```sh
   ./scripts/prebuild/pki/generate-pki.sh
   ```
2. Stack no ar (ou pelos composes `services_required` + `mapex_services`,
   ou pelo `standalone/`):
   ```sh
   ./scripts/mapex-deploy.sh --full
   ```
3. mapexVault saudável + assets MS com `caReady=true` (a fase do cert
   falha no `IssueCert` se isso não tiver subido)
4. Listeners MQTT do broker acessíveis em `:1883` e `:8883`

### Overrides via env

| Variável                | Default                       | Notas                                |
|-------------------------|-------------------------------|--------------------------------------|
| `MAPEXOS_URL`           | `http://localhost:5000`       | base do mapexIam                     |
| `ASSETS_URL`            | `http://localhost:5002`       | base do assets                       |
| `ROUTER_URL`            | `http://localhost:5003`       | base do router                       |
| `GATEWAY_URL`           | `http://localhost:5001`       | base do http_gateway                 |
| `EVENTS_URL`            | `http://localhost:5004`       | base do events                       |
| `MQTT_BROKER_URL`       | `tcp://localhost:1883`        | listener password                    |
| `MQTT_BROKER_TLS_URL`   | `ssl://localhost:8883`        | listener mTLS                        |

## Diagnóstico de falha

Cada step + assert publica um nome qualificado no log da saga (ex:
`assets/assets.ConnectMqttCert`,
`events/events.AssertRawEventReceivedAfter`). Quando uma fase falha:

1. Localize o item que falhou pelo nome na saída do teste.
2. Abra o arquivo correspondente em
   `services/<service>/<module>/{steps,asserts}/`.
3. Leia o comentário acima do item — cada função documenta o que lê
   do bag, o que escreve, e o contrato de produção que exercita.
4. Acompanhe o log do serviço certo (`docker logs mapex-assets`,
   `docker logs mapex-broker-mqtt`, ...) na janela correspondente.

## Bag keys tocadas por essa journey

| Chave                              | Fase  | Escritor                        |
|------------------------------------|-------|---------------------------------|
| `iam.userJWT`                      | 0     | `authSteps.SeedAdminLogin`      |
| `iam.organizationID`               | 0     | `authSteps.SeedAdminLogin`      |
| `router.routeGroupID`              | 1, 2  | `rgSteps.CreateRouteGroup`      |
| `assets.assetTemplateID`           | 1, 2  | `templateSteps.CreateTemplate`  |
| `assets.assetID`                   | 1, 2  | `assetSteps.CreateAsset`        |
| `assets.assetUUID`                 | 1, 2  | `assetSteps.CreateAsset`        |
| `assets.assetMqttPassword`         | 1     | `assetSteps.CreateAsset`        |
| `assets.assetCertPEM` (+ key, ca)  | 2     | `assetSteps.IssueCert`          |
| `assets.assetCertSerial`           | 2     | `assetSteps.IssueCert`          |
| `assets.mqttClient`                | 1, 2  | step `Connect*`                 |
| `assets.mqttConnectedAt`           | 1, 2  | step `Connect*`                 |
| `assets.telemetrySentAt`           | 1, 2  | `assetSteps.PublishTelemetry`   |
| `assets.assetDeleted`              | 1, 2  | `assetSteps.DeleteAsset`        |
