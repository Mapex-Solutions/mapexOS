# Journey: atualização de firmware OTA — MQTT (push)

Cobertura ponta a ponta do rollout de firmware OTA para um device **MQTT** contra
o stack vivo: o broker **empurra** o comando `ota_update` para um device
conectado e com presence-gate, o device baixa + verifica + aplica o firmware,
reporta o progresso em `events/{assetUUID}/ota_status`, e o plano chega a
COMPLETED com o asset migrado para o template alvo.

É SOMENTE o transporte MQTT. O transporte HTTP é um journey **separado**
(`journey/iot/ota_http/`); os dois nunca são misturados. Eles compartilham os
blocos de construção de Camada 1 em `services/assets/ota/` (essa reutilização é
permitida — só os pacotes de journey é que não podem se importar).

## O que este teste prova

O caminho de push completo funciona: upload de firmware (PUT presignado) → plano →
o reconciler despacha o comando por `mqtt.downlink` para o broker **só quando o
device está online** → o broker entrega em `commands/{assetUUID}/ota_update` → o
device baixa via GET presignado, verifica checksum+size, e reporta
`downloading→downloaded→verified→updating→updated` → o consumer de status do
assets avança a execução e, no `updated`, troca o template do asset source→target
→ o plano fecha COMPLETED. Os asserts leem o resultado pela **API pública apenas**.

## Fatos-chave

- Ações de operador usam só `/api/v1/ota/*`; o device sim é um cliente MQTT real.
- Bytes do firmware nunca cruzam o assets MS: **PUT** presignado sobe, **GET**
  presignado baixa.
- Execução: `QUEUED → INITIATED → DOWNLOADING → DOWNLOADED → VERIFIED → UPDATING →
  UPDATED`; o assets é dono de `QUEUED`/`INITIATED`, o device dita o resto.
- Plano: `SCHEDULED → IN_PROGRESS → COMPLETED`.
- **MQTT tem presence-gate** — o device precisa estar `online` (e inscrito no
  tópico de comando) antes do plano despachar, ou o comando empurrado é perdido.

## O fluxo — ping / pong

Legenda: `OP` operador (API pública) · `AS` assets MS (`ota`) · `S3` MinIO ·
`BR` broker MQTT · `DEV` device MQTT (sim). `⇄` via NATS.

```
# ── setup do operador ──
OP → AS   POST /api/v1/ota/firmware/init  {targetTemplateId, version, filename, size, sha256}
AS → OP   200  {firmwareId, uploadUrl}                       # PUT presignado
OP → S3   PUT  uploadUrl  <bytes do firmware>   → 200
OP → AS   POST /api/v1/ota/firmware/:firmwareId/complete     → 200 (READY)
OP → AS   POST /api/v1/ota/plans  {firmwareId, sourceTemplateId, assetIds, startAt=agora, maxTime, rolloutConfig}
AS → OP   200  {id, status:SCHEDULED, counters.total:1}      # StartAt ⇒ IN_PROGRESS

# ── push MQTT (device já CONECTADO + inscrito, presença online) ──
AS ⇄ BR   mqtt.downlink  DownlinkEnvelope{ ota_update, OTAUpdateCommand{
              downloadUrl(GET presignado), checksum, size,
              reportTarget:"events/{assetUUID}/ota_status" } }  # exec QUEUED → INITIATED (presence-gated)
BR → DEV  MQTT publish  commands/{assetUUID}/ota_update  <ota_update>
DEV → S3  GET downloadUrl  → 200 <bytes>                        # DEV verifica sha256 + size
DEV → BR  MQTT publish  events/{assetUUID}/ota_status  {status:"downloading", progress:20}
BR ⇄ AS   ota.status.advisory  OTAStatusAdvisory               # exec → DOWNLOADING
DEV → BR  … "downloaded" → "verified" → "updating" → "updated"
BR ⇄ AS   ota.status.advisory  (um por report)                 # → UPDATED (counters.succeeded++, template source→target)
          … todas as execuções terminais ⇒ early-close ⇒ plano COMPLETED
OP → AS   GET /api/v1/ota/plans/:planId             → {status:COMPLETED, counters.succeeded:1}
OP → AS   GET /api/v1/ota/plans/:planId/executions  → {state:UPDATED, percentage:100}
OP → AS   GET /api/v1/assets/:assetId               → {templateId: target}
```

## Ordenação (o contrato de presença + subscribe)

O reconciler empurra **uma vez**, só para devices **online**, e avança a execução
para `INITIATED` após o dispatch. Então o sim PRECISA conectar, inscrever em
`commands/{assetUUID}/ota_update`, e estar com presença `online` **antes** do
`CreatePlan`:

```
CreateTemplate(source) → CreateTemplate(target) → CreateRouteGroup
→ CreateConnectivityAsset(mqtt, no source)
→ ConnectMqttPassword → SubscribeOtaCommand → AssertHealthStatusEventually("online")
→ InitFirmware → UploadFirmwareBinary → CompleteFirmware → CreatePlan
→ RunMqttOtaDevice → AssertPlanStatus(COMPLETED) → AssertExecutionState(UPDATED,100)
→ AssertAssetTemplateSwitched
```

## Blocos de construção

Blocos de Camada 1 compartilhados em `services/assets/ota/{steps,payloads,asserts}`
(steps de firmware/plano, o sim de device MQTT, os asserts OTA) + `assettemplates`
(templates source/target) + `common/journey/iam_bootstrap`. O journey só sequencia
(`Items()` + `Run()`).

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -run 'TestSuite/iot/ota_mqtt'
```

A journey roda pelo runner único (`journey/suite`); a build tag `saga` a gateia, e
`go test ./...` (sem tag) a pula.

## Requisitos

- Stack vivo: `mapexos`, `assets` nas portas padrão; MinIO + NATS acessíveis.
- `mapexMQTTBroker` rodando — listener `tcp://localhost:1883` (senha).
- Admin seed provisionado (`admin@mapex.local`) — fase 0 (IAM bootstrap) faz login.
- Determinismo: plano criado com `startAt=agora` + `ratePerMinute` alto (despacho
  no primeiro scan); os timeouts de espera do comando + poll do plano são
  generosos pra cobrir um scan do reconciler.
