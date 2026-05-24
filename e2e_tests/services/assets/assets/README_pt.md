# Module e2e: assets / assets

## Escopo

Suite de e2e do módulo `assets` do serviço assets
(`workspace_go/services/assets/src/modules/assets/`). Exercita o
contrato CRUD usado por dispositivos IoT: criar um asset vinculado a um
template e a um route group, consultar, atualizar metadados, listar
com paginação e filtros, e deletar. O `TestMain` provisiona antes um
route group do tipo save-event no router e um asset template mínimo,
porque o payload `AssetCreate` exige ambos como chaves estrangeiras.

## Endpoints exercitados

- `POST /api/v1/assets` — cria asset (template + route group obrigatórios).
- `GET /api/v1/assets/{id}` — busca um asset por id (cobre também 404).
- `PATCH /api/v1/assets/{id}` — atualização parcial de nome e descrição.
- `DELETE /api/v1/assets/{id}` — deleta e confirma que o recurso saiu.
- `GET /api/v1/assets` — listagem paginada com `page`, `perPage`,
  `includeAll` e filtro por `category`.
- `POST /api/v1/asset_templates` / `DELETE /api/v1/asset_templates/{id}`
  — chamados pelo `TestMain` para criar e remover o template que os
  testes de asset dependem (o módulo template tem sua própria suite).
- `POST /api/v1/route_groups` / `DELETE /api/v1/route_groups/{id}` —
  chamados no serviço router pelo `TestMain` pelo mesmo motivo de
  pré-requisito.

## Funções de teste

- `TestCreateAsset_Valid`
- `TestCreateAsset_Minimal`
- `TestGetAssetById`
- `TestGetAssetById_NotFound`
- `TestUpdateAsset_Name`
- `TestDeleteAsset`
- `TestListAssets`
- `TestListAssets_FilterByCategory`

## Fixtures

| Arquivo               | Cenário                                                                                  |
|-----------------------|------------------------------------------------------------------------------------------|
| `create_asset.json`   | Dispositivo MQTT completo com auth por senha, geolocalização, descrição e route bind.    |
| `create_minimal.json` | Dispositivo HTTP mínimo válido — só os campos obrigatórios e as chaves injetadas.        |
| `update_name.json`    | Corpo de PATCH parcial alterando `name` e `description`.                                 |

Observação: `assetTemplateId` e `routeGroupIds` carregam o placeholder
`WILL_BE_INJECTED_BY_TEST` nos arquivos JSON e são sobrescritos em
tempo de execução com os ids que o `TestMain` provisionou.

## Como rodar

```bash
cd e2e_tests

# Suite completa do módulo
go test ./services/assets/assets -v

# Um teste específico
go test ./services/assets/assets -v -run TestCreateAsset_Valid
```

## Resultado em caso de PASS

Confirma que o módulo assets respeita seu contrato HTTP público
ponta-a-ponta: validação de campos obrigatórios, ligação por chave
estrangeira com asset templates e route groups do router, o ciclo CRUD
completo (create → read → patch → delete → 404), e a listagem paginada
com filtro por categoria — tudo sob a organização root do seed.

## Requisitos

- `assets` disponível na porta `5002` (override via `ASSETS_URL`).
- `mapexos` disponível na porta `5000` para gerar o token de admin.
- `router` disponível na porta `5003` para criar o route group
  pré-requisito.
- Stack iniciado pelo `mapexOSDeploy`; o usuário admin do seed
  `admin@mapex.local` e a organização root
  `0000000000000000000aa001` precisam existir (provisionados pelo
  `mongodb-init`).
- Variável `API_KEY` opcional — usa por padrão a chave interna
  canônica do stack.

## Notas

- O `TestMain` provisiona um template e um route group novos por
  execução e remove ambos no teardown, então a suite é autocontida e
  não deixa registros órfãos.
- O id de categoria `670a4cde48e006e3f95e8eb3` usado no teste de filtro
  vem do catálogo do seed; se o seed mudar, essa constante precisa ser
  ajustada.
- As pastas irmãs `payloads/`, `steps/` e `asserts/` são blocos de
  construção de saga consumidos pelas journeys de IoT; não fazem parte
  desta suite de e2e de módulo.
