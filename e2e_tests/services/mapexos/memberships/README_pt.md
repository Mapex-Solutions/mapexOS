# Módulo e2e: mapexos / memberships

## Escopo

Cobre o módulo de memberships da plataforma mapexos — o vínculo entre
um assignee (usuário ou grupo) e uma ou mais roles dentro de uma
organização, com `scope` `local` (apenas esta org) ou `recursive`
(esta org + descendentes). Suportado por
`workspace_go/services/mapexIam/src/modules/memberships/`. A suíte
também exercita o endpoint `/api/v1/me/coverage`, que materializa a
coverage efetiva do chamador a partir das memberships ativas.

## Endpoints exercitados

- `POST /api/v1/memberships` — cria uma membership (user/local,
  user/recursive, group, multi-role).
- `GET /api/v1/memberships/{id}` — busca uma membership.
- `GET /api/v1/memberships?includeAll=true[&userId=]` — listagem
  paginada com filtro opcional `userId`.
- `PATCH /api/v1/memberships/{id}` — atualização parcial de `scope`,
  `enabled` ou `roleIds`.
- `DELETE /api/v1/memberships/{id}` — remove uma membership.
- `GET /api/v1/me/coverage` — coverage organizacional efetiva do
  usuário autenticado, recalculada a partir das memberships ativas.

## Fixtures

| Arquivo | Descrição |
|---|---|
| `create_user_membership_local.json` | Assignee user, `scope: local`, uma role. Fixture base do CRUD. |
| `create_user_membership_recursive.json` | Assignee user, `scope: recursive`. |
| `create_group_membership.json` | Assignee group, `scope: local`, uma role. |
| `create_multiple_roles.json` | Assignee user com duas roles (`{{ROLE_ID_1}}` + `{{ROLE_ID_2}}`). |
| `update_scope.json` | PATCH que promove `scope` para `recursive`. |
| `update_disable.json` | PATCH com `enabled: false`. |
| `update_roles.json` | PATCH que substitui `roleIds` por uma única nova role. |

Os placeholders `{{USER_ID}}`, `{{GROUP_ID}}`, `{{ORG_ID}}`,
`{{ROLE_ID}}`, `{{ROLE_ID_1}}`, `{{ROLE_ID_2}}` são substituídos na
carga pelos ids provisionados em runtime (org de teste + duas roles
de teste + um grupo de teste são criados no `TestMain` e destruídos
ao final).

## Como rodar

```bash
cd e2e_tests
go test ./services/mapexos/memberships -v

# Teste individual
go test ./services/mapexos/memberships -v -run TestUpdateMembership_Scope
```

## Resultado em caso de PASS

- CRUD completo nos quatro formatos de create (user/local,
  user/recursive, group, multi-role) → get → patch (scope, disable,
  roles) → delete → re-get retorna `404`.
- `TestCreateMembership_UserLocal` valida o shape completo da
  resposta: `assigneeType`, `assigneeId`, `orgId`, `scope`, `enabled`
  e o array `roleIds[]`.
- `TestCreateMembership_MultipleRoles` prova que atribuição
  multi-role persiste exatamente dois `roleIds`.
- `TestGetMembershipById_NotFound` prova que id desconhecido retorna
  `404`.
- `TestListMemberships` prova o envelope paginado e que tanto uma
  membership de usuário quanto uma de grupo são retornadas.
- `TestListMemberships_FilterByUser` prova que o filtro `userId=`
  reduz o resultado às memberships daquele usuário.
- `TestGetMeCoverage` prova que `/api/v1/me/coverage` retorna
  `{ userId, customers[] }` refletindo as memberships ativas.

## Requisitos

- Serviço mapexos / iam acessível em `http://localhost:5000`
  (sobrescrever via `MAPEXOS_URL`).
- Admin seed + role SuperAdmin + organização mapexos seed,
  provisionados pelo `mongodb-init`. Os ids determinísticos de
  usuário (`constants.RootUserID`, `constants.AdminUserID`) vêm do
  script de seed.
- `utils.SetupE2EEnvironment()` roda no `TestMain` (limpa DB + flush
  de cache + re-seed) — destrutivo contra a stack local.

## Notas

- O `TestMain` provisiona uma organização filha de teste, duas roles
  de teste e um grupo de teste, e apaga tudo ao final — assim cada
  execução é independente e não deixa rastro.
- Dois clientes são montados: `rootClient` (admin seed, wildcard) e
  `adminClient` (escopo org `admin_vendor.*`). O alias `client`
  padrão é o `rootClient` para a cobertura CRUD. Ambos carregam
  `X-Org-Context` fixado na org root seed.
- `scope=local` vs `scope=recursive` é uma distinção de primeira
  classe: memberships recursivas concedem ao assignee acesso à org e
  a todos os descendentes; `TestCreateMembership_UserRecursive` e
  `TestUpdateMembership_Scope` exercitam explicitamente as duas
  transições.
