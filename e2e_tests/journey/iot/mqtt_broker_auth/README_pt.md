# Journey: Autenticação MQTT do broker — ciclos por senha + cert

Cobertura ponta a ponta do pipeline de autenticação MQTT da plataforma
contra a stack `services_required` viva. Cada fase roda um ciclo de
vida completo de asset — **criar → conectar → presença online →
publicar → events round-trip → apagar → reconectar deve ser negado** —
provando tanto o caminho feliz quanto a invalidação de cache por
FANOUT na camada L1 (Pebble) do plugin do broker.

## Contrato do wire (pós TKT-2026-0040)

- Username MQTT = `assetUUID` puro (sem o prefixo legado `{orgId}:`)
- Formato de tópico do device = `events/{assetUUID}/{eventType}`
- Subject.CN do cert = `assetUUID` puro; escopo de tenant flui via
  `asset.protocol.mqtt.orgId` do lado do servidor
- Listeners: `tcp://localhost:1883` (senha) e
  `ssl://localhost:8883` (mTLS)

## Fases

| Fase                   | O que cobre                                                                   |
|------------------------|-------------------------------------------------------------------------------|
| `phase0_iam_bootstrap` | Login do admin seed → validade do JWT → coverage de org-context               |
| `phase1_password_user` | Ciclo de vida completo da auth por senha (10 passos)                          |
| `phase2_cert_user`     | Ciclo de vida completo da auth por cert / mTLS (11 passos)                    |
| `phase3_cascade`       | Cascata TieredStore L1 -> L2 -> L3 + fanout invalidate                        |

### `phase1_password_user` — ordem dos passos

1. **CreateRouteGroup** — route group para o asset
2. **CreateTemplate** — asset template (temperatura)
3. **CreateAsset** — asset persistido com `authType=password`
4. **ConnectMqttPassword** — CONNECT MQTT com `(assetUUID, password)`
5. **AssertHealthStatusEventually(online)** — presença flui ponta a ponta
6. **PublishTelemetry** — publica em `events/{assetUUID}/temperature`
7. **AssertRawEventReceivedAfter** — events service expõe a linha
8. **DisconnectMqtt** — disconnect MQTT limpo
9. **DeleteAsset** — invalidação FANOUT dispara
10. **AssertConnectDeniedPassword** — CONNECT novo é negado

### `phase2_cert_user` — ordem dos passos

1. **CreateRouteGroup**
2. **CreateTemplate**
3. **CreateAssetWith(SagaMqttCertTemperatureSensor)** — `authType=cert`, `certTTL={1 day}`
4. **IssueCert** — POST `/api/v1/mqtt_certs`; bundle PEM no bag, asset.currentCert persistido
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

# Só a phase 0 (smoke do bootstrap IAM)
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap/

# Ciclo por senha
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase1_password_user/

# Ciclo por cert (mTLS)
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/phase2_cert_user/

# Tudo na journey
go test -tags=saga -v -count=1 ./journey/iot/mqtt_broker_auth/...
```

`-count=1` desabilita o cache de testes do Go para que cada execução
bata na stack viva. Sem isso, uma segunda execução imprime `(cached)`
e mostra o resultado anterior sem realmente contactar o broker.

## Pré-requisitos

1. Pré-construir PKI:
   ```sh
   ./scripts/prebuild/pki/generate-pki.sh
   ```
2. Stack no ar (`services_required` + `mapex_services` OU `standalone/`):
   ```sh
   ./scripts/mapex-deploy.sh --full
   ```
3. mapexVault saudável + assets MS com `caReady=true` (caso contrário
   a fase de cert falha em `IssueCert`).
4. Listeners MQTT do broker acessíveis em `:1883` e `:8883`.

### Overrides de ambiente

| Variável                | Default                       | Notas                                |
|-------------------------|-------------------------------|--------------------------------------|
| `MAPEXOS_URL`           | `http://localhost:5000`       | base URL do mapexIam                 |
| `ASSETS_URL`            | `http://localhost:5002`       | base URL do serviço assets           |
| `ROUTER_URL`            | `http://localhost:5003`       | base URL do serviço router           |
| `GATEWAY_URL`           | `http://localhost:5001`       | base URL do serviço http_gateway     |
| `EVENTS_URL`            | `http://localhost:5004`       | base URL do serviço events           |
| `MQTT_BROKER_URL`       | `tcp://localhost:1883`        | listener de senha                    |
| `MQTT_BROKER_TLS_URL`   | `ssl://localhost:8883`        | listener mTLS                        |

## Diagnóstico de falha

Cada step + assert publica um nome qualificado no log do saga
(ex.: `assets/assets.ConnectMqttCert`,
`events/events.AssertRawEventReceivedAfter`). Quando uma fase falha:

1. Localize o item que falhou pelo nome no output do teste.
2. Abra o arquivo correspondente em `services/<service>/<module>/{steps,asserts}/`.
3. Leia o comentário acima do item que falhou — cada função documenta
   o que lê do bag, o que escreve, e o contrato de produção que
   exercita.
4. Tail no log do serviço certo (`docker logs mapex-assets`,
   `docker logs mapex-broker-mqtt`, ...) para a janela de tempo
   correspondente.
