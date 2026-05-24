# Module e2e: router / routegroups

## Escopo

Suite de e2e do módulo `routegroups` do serviço router
(`workspace_go/services/router/src/modules/routegroups/`). Exercita o
contrato CRUD público de um route group: criar com um ou mais routers
inline, consultar por id, atualizar parcialmente, deletar e validar
404 em seguida, e a listagem paginada com vários filtros (name,
enabled, version, múltiplos filtros combinados, projection). A suite
dirige o router em duas identidades — o usuário root wildcard e o
admin do seed sob o contexto da organização root — para provar que
ambos os fluxos respeitam o coverage middleware aplicado pelo serviço.

## Endpoints exercitados

- `POST /api/v1/route_groups` — cria um route group com `routers`
  inline (`kind=save_event` na fixture base).
- `GET /api/v1/route_groups/{id}` — busca um route group; também serve
  para verificar 404 após delete.
- `PATCH /api/v1/route_groups/{id}` — atualização parcial (troca de
  nome). Aceita `200 OK` ou `201 Created` por compatibilidade de
  build.
- `DELETE /api/v1/route_groups/{id}` — remove um route group; o helper
  de cleanup tolera `404`.
- `GET /api/v1/route_groups` — listagem paginada com combinações de
  `page`, `perPage`, `name`, `enabled`, `version`, múltiplos filtros e
  `projection`.

## Funções de teste

- `TestCreateRouteGroup`
- `TestGetRouteGroupById`
- `TestUpdateRouteGroup`
- `TestDeleteRouteGroup`
- `TestListRouteGroups_BasicPagination`
- `TestListRouteGroups_FilterByName`
- `TestListRouteGroups_FilterByEnabled`
- `TestListRouteGroups_FilterByVersion`
- `TestListRouteGroups_MultipleFilters`
- `TestListRouteGroups_Projection`
- `TestListRouteGroups_WithOrgContext`
- `TestListRouteGroups_RootUser`

## Fixtures

| Arquivo                            | Cenário                                                                                              |
|------------------------------------|-------------------------------------------------------------------------------------------------------|
| `create_routegroup.json`           | Route group base `API Routes v1`, habilitado, versão `1.0.0`, vinculado à org root do seed `0000000000000000000aa001`, com um único router `save_event` carregando `metadata.source=api`. |
| `create_routegroup_versioned.json` | Route group companheiro `API Routes v2`, desabilitado, versão `2.0.0`, mesma org root; usado para popular a listagem nos testes de filtro de modo que as projections distingam vários registros. |
| `update_routegroup.json`           | Corpo de PATCH parcial — renomeia o alvo para `API Routes v1 Updated`.                                |
| `update_enabled.json`              | Corpo de PATCH parcial alternando `enabled=false`; disponível para execuções ad-hoc (não carregado pelas funções `Test*` atuais). |

Todas as referências de `orgId` apontam para a organização root
canônica do seed `0000000000000000000aa001`; a suite não provisiona
organizações extras.

## Como rodar

```bash
cd e2e_tests

# Suite completa do módulo
go test ./services/router/routegroups -v

# Um teste específico
go test ./services/router/routegroups -v -run TestCreateRouteGroup
```

## Resultado em caso de PASS

Confirma que o módulo routegroups respeita seu contrato HTTP público
ponta-a-ponta: validação de campos obrigatórios no create, o ciclo
CRUD completo (create → read → patch → delete → 404) e a listagem
paginada com todas as combinações de filtro suportadas (name, enabled,
version, múltiplos filtros, projection). Também prova que o coverage
middleware resolve corretamente o contexto da organização root tanto
para o usuário root wildcard quanto para o admin do seed.

## Requisitos

- `router` disponível na porta `5003` (override via `ROUTER_URL`).
- `mapexos` disponível na porta `5000` para gerar os tokens de root e
  admin.
- Stack iniciado pelo `mapexOSDeploy`; o usuário admin do seed
  `admin@mapex.local` e a organização root
  `0000000000000000000aa001` precisam existir (provisionados pelo
  `mongodb-init`).

## Notas

- Todo endpoint CRUD exige `X-Org-Context` mesmo quando o bearer
  carrega o role wildcard; o `TestMain` configura o header em ambos os
  clients apontando para a org root.
- O endpoint PATCH aceita `200 OK` ou `201 Created` — a suite tolera
  os dois para absorver desvios de build em imagens mais antigas do
  router.
- As pastas irmãs `payloads/` e `steps/` são blocos de construção de
  saga consumidos pelas journeys de IoT e automation (variantes de
  route group save-event, trigger e workflow); não fazem parte desta
  suite de e2e de módulo.
