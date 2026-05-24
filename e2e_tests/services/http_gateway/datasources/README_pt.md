# Module e2e: http_gateway / datasources

## Escopo

Suite de e2e do módulo `datasources` do serviço http_gateway
(`workspace_go/services/http_gateway/src/modules/datasources/`).
DataSources descrevem como um protocolo de entrada (HTTP, MQTT, …) é
autenticado, rate-limitado e ligado a assets no momento da ingestão
via webhook. A suite cobre o ciclo CRUD completo mais a matriz de
filtros da listagem (`name`, `enabled`, `mode`, `protocol`,
projection, multi filter) usando tanto um cliente com token ROOT
quanto um cliente admin escopado na organização do seed.

## Endpoints exercitados

- `POST /api/v1/data_sources` — cria datasource (payloads HTTP+apiKey
  e MQTT+jwt).
- `GET /api/v1/data_sources/{id}` — busca uma datasource por id.
- `PATCH /api/v1/data_sources/{id}` — atualização parcial de metadado.
- `DELETE /api/v1/data_sources/{id}` — deleta e confirma 404 no GET
  seguinte.
- `GET /api/v1/data_sources` — listagem paginada com `page` /
  `perPage`, filtros por campo (`name`, `enabled`, `mode`, `protocol`),
  `projection`, composição multi-filter, mais passes separados usando
  o cliente admin (escopado por org) e o cliente ROOT.

## Funções de teste

- `TestCreateDataSource`
- `TestGetDataSourceById`
- `TestUpdateDataSource`
- `TestDeleteDataSource`
- `TestListDataSources_BasicPagination`
- `TestListDataSources_FilterByName`
- `TestListDataSources_FilterByEnabled`
- `TestListDataSources_FilterByMode`
- `TestListDataSources_FilterByProtocol`
- `TestListDataSources_MultipleFilters`
- `TestListDataSources_Projection`
- `TestListDataSources_WithOrgContext`
- `TestListDataSources_RootUser`

## Fixtures

| Arquivo                       | Cenário                                                                                       |
|-------------------------------|-----------------------------------------------------------------------------------------------|
| `create_datasource_http.json` | Datasource HTTP modo pull com auth por API-Key em header e bind por `uuidField`.              |
| `create_datasource_mqtt.json` | Datasource MQTT modo push (inicialmente desabilitada) com auth JWT e rate-limit configurado.  |
| `update_datasource.json`      | Corpo de PATCH parcial alterando `name` e `description`.                                      |
| `update_enabled.json`         | Corpo de PATCH parcial mudando `enabled` para false (reservado; não usado pelos testes hoje). |

## Como rodar

```bash
cd e2e_tests

# Suite completa do módulo
go test ./services/http_gateway/datasources -v

# Um teste específico
go test ./services/http_gateway/datasources -v -run TestListDataSources_FilterByName
```

## Resultado em caso de PASS

Confirma que o módulo datasources respeita seu contrato HTTP público:
ciclo CRUD completo (create → read → patch → delete → 404), metadados
de paginação, toda a matriz de filtros server-side da listagem (`name`
case-insensitive, `enabled` exato, `mode` exato, `protocol` exato e
composição multi-filter), o seletor `projection` e o caminho de
resolução de org-context — coberto tanto com cliente admin (escopado
na organização root do seed) quanto com cliente ROOT.

## Requisitos

- `http_gateway` disponível na porta `5001` (override via
  `GATEWAY_URL`).
- `mapexos` disponível na porta `5000` para gerar os tokens admin e
  ROOT (`utils.GetRootToken` / `utils.GetAdminToken`).
- Stack iniciado pelo `mapexOSDeploy`; o usuário admin do seed
  `admin@mapex.local` e a organização root
  `0000000000000000000aa001` precisam existir.

## Notas

- Dois clientes são montados no `TestMain`: um ROOT para testes que
  precisam enxergar todas as datasources e um admin escopado por org
  com `X-Org-Context` apontado para a org root do seed, usado no
  teste `WithOrgContext`.
- O filtro `name` é validado com `strings.ToLower(...)` porque o
  serviço faz match parcial case-insensitive — datasources criadas
  por saga com `"http"` em minúsculo também contam como match válido.
- O `PATCH` hoje devolve `201 Created` (não `200 OK`); o teste de
  update afirma sobre `201` por causa disso.
