# Module e2e: assets / assettemplates

## Escopo

Suite de e2e do módulo `assettemplates` do serviço assets
(`workspace_go/services/assets/src/modules/assettemplates/`). Asset
templates descrevem como parsear um par fabricante/modelo: um
`assetIdPath` para extrair o identificador do dispositivo e um par de
scripts JS inline (validator + conversion) executados sobre a
telemetria recebida. A suite exercita o CRUD completo mais a listagem
paginada com filtro por status.

## Endpoints exercitados

- `POST /api/v1/asset_templates` — cria template (payloads completo e
  mínimo).
- `GET /api/v1/asset_templates/{id}` — busca um template por id (cobre
  também 404).
- `PATCH /api/v1/asset_templates/{id}` — atualização parcial; coberta
  tanto para os campos de metadado (name, description, version) quanto
  para os scripts JS inline.
- `DELETE /api/v1/asset_templates/{id}` — deleta e confirma que o
  recurso saiu.
- `GET /api/v1/asset_templates` — listagem paginada com `page`,
  `perPage` e filtro `enabled`.

## Funções de teste

- `TestCreateAssetTemplate_Valid`
- `TestCreateAssetTemplate_Minimal`
- `TestGetAssetTemplateById`
- `TestGetAssetTemplateById_NotFound`
- `TestUpdateAssetTemplate_Scripts`
- `TestUpdateAssetTemplate_Metadata`
- `TestDeleteAssetTemplate`
- `TestListAssetTemplates`
- `TestListAssetTemplates_FilterByStatus`

## Fixtures

| Arquivo                | Cenário                                                                                  |
|------------------------|------------------------------------------------------------------------------------------|
| `create_template.json` | Template completo — Acme TS-2000 com os quatro slots de script e `assetIdPath`.          |
| `create_minimal.json`  | Template mínimo válido — só os obrigatórios e stubs inline de validator/conversion.      |
| `update_metadata.json` | PATCH parcial alterando `name`, `description`, `version` (mantém `assetIdPath`).         |
| `update_scripts.json`  | PATCH parcial trocando `scriptValidator` e `scriptConversion` por variantes marcadas.    |

## Como rodar

```bash
cd e2e_tests

# Suite completa do módulo
go test ./services/assets/assettemplates -v

# Um teste específico
go test ./services/assets/assettemplates -v -run TestUpdateAssetTemplate_Scripts
```

## Resultado em caso de PASS

Confirma que o módulo de asset templates respeita seu contrato HTTP
público: validação de campos obrigatórios, ciclo CRUD completo (create
→ read → patch → delete → 404), semântica de PATCH parcial tanto para
metadados quanto para o corpo dos scripts inline, e listagem paginada
com o filtro `enabled` — tudo sob a organização root do seed.

## Requisitos

- `assets` disponível na porta `5002` (override via `ASSETS_URL`).
- `mapexos` disponível na porta `5000` para gerar o token de admin.
- Stack iniciado pelo `mapexOSDeploy`; o usuário admin do seed
  `admin@mapex.local` e a organização root
  `0000000000000000000aa001` precisam existir (provisionados pelo
  `mongodb-init`).

## Notas

- Os testes de PATCH aceitam tanto `200` quanto `201` porque o serviço
  hoje devolve `201` em updates bem-sucedidos — o assert está
  intencionalmente flexível enquanto o serviço não é normalizado.
- Diferente da suite irmã `assets`, não há pré-requisitos de router
  nem de template: cada teste é dono do ciclo de vida do próprio
  template.
