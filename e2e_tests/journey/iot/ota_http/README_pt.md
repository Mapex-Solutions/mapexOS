# Journey: atualização de firmware OTA — HTTP (poll)

Cobertura ponta a ponta do rollout de firmware OTA para um device **HTTP** contra
o stack vivo: o device **consulta** o gateway pelo seu job, baixa + verifica +
aplica o firmware, reporta o progresso via `POST /api/v1/ota/status`, e o plano
chega a COMPLETED com o asset migrado para o template alvo.

É SOMENTE o transporte HTTP. O transporte MQTT é um journey **separado**
(`journey/iot/ota_mqtt/`); os dois nunca são misturados. Eles compartilham os
blocos de construção de Camada 1 em `services/assets/ota/` (essa reutilização é
permitida — só os pacotes de journey é que não podem se importar).

## O que este teste prova

O caminho de poll completo funciona: upload de firmware (PUT presignado) → plano →
o device HTTP consulta `GET /api/v1/ota/jobs` pelo gateway (auth por data source)
→ o gateway relaya pro endpoint interno de pending-job do assets, que gera um GET
presignado fresco → o device baixa, verifica checksum+size, e faz POST de
`downloading→downloaded→verified→updating→updated` em `POST /api/v1/ota/status` →
o gateway normaliza cada um em `ota.status.advisory` → o consumer de status do
assets avança a execução e, no `updated`, troca o template do asset source→target
→ o plano fecha COMPLETED. Os asserts leem o resultado pela **API pública apenas**.

## Fatos-chave

- Ações de operador usam só `/api/v1/ota/*`; o device sim é um cliente HTTP real
  (autenticado por data source via `?ds=`).
- Bytes do firmware nunca cruzam o assets MS: **PUT** presignado sobe, **GET**
  presignado baixa.
- Execução: `QUEUED → INITIATED → DOWNLOADING → DOWNLOADED → VERIFIED → UPDATING →
  UPDATED`; o assets é dono de `QUEUED`/`INITIATED`, o device dita o resto.
- Plano: `SCHEDULED → IN_PROGRESS → COMPLETED`.
- **HTTP NÃO tem presence-gate** — o device aparece ao pollar; sem pré-condição de
  online, sem corrida de push.

## O fluxo — ping / pong

Legenda: `OP` operador (API pública) · `AS` assets MS (`ota`) · `S3` MinIO ·
`GW` http_gateway · `DEV` device HTTP (sim). `⇄` via NATS.

```
# ── setup do operador ──
OP → AS   POST /api/v1/ota/firmware/init  {targetTemplateId, version, filename, size, sha256}
AS → OP   200  {firmwareId, uploadUrl}                       # PUT presignado
OP → S3   PUT  uploadUrl  <bytes do firmware>   → 200
OP → AS   POST /api/v1/ota/firmware/:firmwareId/complete     → 200 (READY)
OP → AS   POST /api/v1/ota/plans  {firmwareId, sourceTemplateId, assetIds, startAt=agora, maxTime, rolloutConfig}
AS → OP   200  {id, status:SCHEDULED, counters.total:1}      # StartAt ⇒ IN_PROGRESS

# ── poll HTTP (device puxa o job; sem push, sem presence gate) ──
DEV → GW  GET /api/v1/ota/jobs?ds={dataSourceId}&assetUUID={assetUUID}
GW → AS   (relay interno)  GET /jobs?assetUUID={assetUUID}
AS → GW   200  OTAUpdateCommand{ downloadUrl(GET presignado), checksum, size,
              reportTarget:"/api/v1/ota/status" }             # exec QUEUED → INITIATED (no primeiro serve)
GW → DEV  200  <ota_update>
DEV → S3  GET downloadUrl  → 200 <bytes>                       # DEV verifica sha256 + size
DEV → GW  POST /api/v1/ota/status?ds={dataSourceId}  {status:"downloading", progress:20}
GW ⇄ AS   ota.status.advisory  OTAStatusAdvisory              # exec → DOWNLOADING
DEV → GW  … POST "downloaded" → "verified" → "updating" → "updated"
GW ⇄ AS   ota.status.advisory  (um por report)               # → UPDATED (counters.succeeded++, template source→target)
          … todas as execuções terminais ⇒ early-close ⇒ plano COMPLETED
OP → AS   GET /api/v1/ota/plans/:planId             → {status:COMPLETED, counters.succeeded:1}
OP → AS   GET /api/v1/ota/plans/:planId/executions  → {state:UPDATED, percentage:100}
OP → AS   GET /api/v1/assets/:assetId               → {templateId: target}
```

## Ordenação (poll — sem corrida)

O device HTTP nunca é empurrado, então não há corrida de subscribe/presença; ele
só pola depois que o plano existe:

```
CreateTemplate(source) → CreateTemplate(target) → CreateDataSource
→ CreateConnectivityAsset(http, no source, ligado ao data source)
→ InitFirmware → UploadFirmwareBinary → CompleteFirmware → CreatePlan
→ RunHttpOtaDevice → AssertPlanStatus(COMPLETED) → AssertExecutionState(UPDATED,100)
→ AssertAssetTemplateSwitched
```

## Blocos de construção

Blocos de Camada 1 compartilhados em `services/assets/ota/{steps,payloads,asserts}`
(steps de firmware/plano, o sim de device HTTP, os asserts OTA) + `assettemplates`
(templates source/target) + `http_gateway/datasources` (CreateDataSource) +
`common/journey/iam_bootstrap`. O journey só sequencia.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/ota_http/
```

A build tag `saga` gateia o teste; `go test ./...` (sem tag) o pula.

## Requisitos

- Stack vivo: `mapexos`, `assets`, `http_gateway` nas portas padrão; MinIO + NATS
  acessíveis.
- Admin seed provisionado (`admin@mapex.local`) — fase 0 (IAM bootstrap) faz login.
- Determinismo: plano criado com `startAt=agora` + `ratePerMinute` alto; os
  timeouts de poll do job + poll do plano são generosos pra cobrir um scan.
