# Módulo e2e: mapexos / lists

## Escopo

Cobre o módulo de lists da plataforma mapexos — catálogos tipados de
chave/valor (por exemplo `assetGroup`, `assetType`) consumidos pelo
resto da plataforma como enumerações. Suportado por
`workspace_go/services/mapexIam/src/modules/lists/`. Cada list carrega
um `type`, um `name`, um `value`, um flag `isSystem` (true para
catálogos de fábrica, false para os criados pela org) e um `orgId` de
escopo.

## Endpoints exercitados

- `POST /api/v1/lists` — cria uma entrada (payloads escopo org,
  system e mínimo).
- `GET /api/v1/lists/{id}` — busca uma entrada.
- `GET /api/v1/lists?includeAll=true&page=&perPage=` — listagem
  paginada, também exercitada com filtros `type=` e `name=`.
- `PATCH /api/v1/lists/{id}` — atualização parcial de `name` + `value`.
- `DELETE /api/v1/lists/{id}` — remove uma entrada.

## Fixtures

| Arquivo | Descrição |
|---|---|
| `create_minimal.json` | `assetGroup` "Workstations", `isSystem: false`, org root seed. |
| `create_org_list.json` | `assetGroup` "Servers", `isSystem: false` — fixture principal do CRUD. |
| `create_system_list.json` | `assetType` "Physical Server", `isSystem: true`. |
| `update_name.json` | PATCH que renomeia para "Updated Name" / `updated_value`. |

## Como rodar

```bash
cd e2e_tests
go test ./services/mapexos/lists -v

# Teste individual
go test ./services/mapexos/lists -v -run TestListLists_FilterByType
```

## Resultado em caso de PASS

- CRUD completo nos três formatos de payload (org / system /
  minimal): create → get → patch (name + value) → delete → re-get
  retorna `404` ou `200` com `nil data`.
- `TestGetListById_NotFound` prova que id inexistente retorna `404`
  ou `200`+`nil` de forma consistente.
- `TestListLists` prova o formato do envelope paginado (`items[]` +
  `pagination{page, perPage}`).
- `TestListLists_FilterByType` prova o filtro de query `type` — todo
  item retornado satisfaz `type == assetGroup`.
- `TestListLists_FilterByName` exercita o filtro parcial de nome e
  loga a contagem (tolerante à visibilidade por escopo de org).

## Requisitos

- Serviço mapexos / iam acessível em `http://localhost:5000`
  (sobrescrever via `MAPEXOS_URL`).
- Organização mapexos seed (id `0000000000000000000aa001`),
  provisionada pelo `mongodb-init`.
- `utils.SetupE2EEnvironment()` roda no `TestMain` (limpa DB + flush
  de cache + re-seed) — destrutivo contra a stack local.

## Notas

- Dois clientes são montados (`rootClient` wildcard + `adminClient`
  escopo org `admin_vendor.*`); o alias `client` padrão é o
  `rootClient` para a cobertura CRUD. Ambos carregam `X-Org-Context`
  fixado na org root seed.
- A distinção `system` vs org é de dado (`isSystem` booleano), não um
  endpoint separado — o mesmo `POST /api/v1/lists` atende os dois.
- As asserções de `TestGetListById_NotFound` e `TestDeleteList`
  aceitam tanto `404` quanto `200` com `nil data`, pois o serviço
  pode retornar qualquer um dos formatos para entrada inexistente.
- `TestListLists_FilterByName` é tolerante de propósito — apenas loga
  ao invés de asseverar a contagem porque o resultado depende da
  visibilidade por org.
