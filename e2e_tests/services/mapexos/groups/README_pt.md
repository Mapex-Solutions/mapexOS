# Módulo e2e: mapexos / groups

## Escopo

Cobre o módulo de groups da plataforma mapexos — coleções de usuários
que agregam atribuições de role dentro de uma organização. Suportado
por `workspace_go/services/mapexIam/src/modules/groups/`. Um grupo é
sempre criado contra um `orgId` (escopo organizacional); a distinção
`isSystem` fica reservada para grupos seedados de fábrica, e a suíte
prova que ambos os fluxos chegam ao mesmo CRUD.

## Endpoints exercitados

- `POST /api/v1/groups` — cria grupo (escopo org, mínimo, system e um
  caso negativo sem `orgId`).
- `GET /api/v1/groups/{id}` — busca um grupo.
- `GET /api/v1/groups?includeAll=true` — listagem paginada dos grupos
  acessíveis.
- `PATCH /api/v1/groups/{id}` — atualização parcial de `name`,
  `description`, `enabled` ou todos juntos.
- `DELETE /api/v1/groups/{id}` — remove um grupo.

## Fixtures

| Arquivo | Descrição |
|---|---|
| `create_minimal.json` | Payload válido mínimo — name + enabled + `{{ORG_ID}}` + um roleId. |
| `create_org_group.json` | "Engineering Team" com descrição, usado pela maioria dos testes CRUD. |
| `create_system_group.json` | Payload de grupo de administradores do sistema (ainda escopo org). |
| `update_name.json` | PATCH que renomeia para "Engineering Team Updated". |
| `update_description.json` | PATCH que altera apenas a descrição. |
| `update_disable.json` | PATCH com `enabled: false`. |
| `update_full.json` | PATCH que atualiza name + description + enabled de uma vez. |

`{{ORG_ID}}` é substituído na carga pelo id da organização mapexos
seed vindo de `common/constants`.

## Como rodar

```bash
cd e2e_tests
go test ./services/mapexos/groups -v

# Teste individual
go test ./services/mapexos/groups -v -run TestUpdateGroup_Full
```

## Resultado em caso de PASS

- CRUD completo: create (org / minimal / system) → get → patch (name,
  description, disable, full) → delete → re-get retorna 404.
- `TestCreateGroup_NoOrgIdForNonSystem` prova que o validador rejeita
  criação de grupo sem `orgId` com `400 Bad Request`.
- `TestGetGroupById_NotFound` prova que id desconhecido retorna `404`.
- `TestListGroups` prova que `includeAll=true` retorna o envelope
  paginado (`items[]` + `pagination`) e que um grupo recém-criado
  aparece na listagem.

## Requisitos

- Serviço mapexos / iam acessível em `http://localhost:5000`
  (sobrescrever via `MAPEXOS_URL`).
- Admin seed + role SuperAdmin (id `0000000000000000000aa201`) +
  organização mapexos seed, todos provisionados pelo `mongodb-init`.
- `utils.SetupE2EEnvironment()` roda no `TestMain` (limpa DB + flush
  de cache + re-seed) — destrutivo contra a stack local.

## Notas

- Dois clientes são montados: `rootClient` (admin seed, wildcard) e
  `adminClient` (escopo org `admin_vendor.*`); ambos carregam
  `X-Org-Context` fixado na org root seed. O alias `client` padrão
  aponta para o `rootClient` para a cobertura CRUD.
- O middleware do mapexos exige `X-Org-Context` em todos os endpoints
  CRUD, mesmo para o bearer wildcard.
- Cada teste mutante registra `t.Cleanup` para apagar o grupo criado,
  aceitando tanto `200` quanto `404` no teardown.
