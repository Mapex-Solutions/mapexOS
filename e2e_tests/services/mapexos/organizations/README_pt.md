# E2E de módulo: mapexos / organizations

## Escopo

Cobertura end-to-end do módulo de organizations — a espinha dorsal de
multi-tenancy da superfície mapexos servida pelo `mapexIam`
(`workspace_go/services/mapexIam/src/modules/organizations/`). É o maior
teste de módulo individual do repositório (~30 testes) e exercita o CRUD
completo, a hierarquia profunda `customer / site / building`, a paginação
em árvore por cursor e o middleware de cobertura que decide qual ator
enxerga qual organização. Todos os testes rodam contra a stack viva e
registram limpeza para que o dataset de seed permaneça intacto.

## Endpoints exercitados

- `POST   /api/v1/organizations`              — cria customer / site / building
- `GET    /api/v1/organizations/{id}`         — busca uma organização
- `GET    /api/v1/organizations`              — listagem paginada (page / perPage)
- `GET    /api/v1/organizations/tree`         — árvore paginada por cursor (next / previous)
- `PATCH  /api/v1/organizations/{id}`         — atualização parcial (name, enabled, full)
- `DELETE /api/v1/organizations/{id}`         — exclusão (idempotente em 404)
- `POST   /api/v1/roles`                      — cria role restrita para testes de negação
- `POST   /api/v1/onboarding/users`           — provisiona o admin restrito
- `POST   /api/v1/auth/login`                 — login como admin restrito

## Fixtures

| Arquivo                      | Propósito                                                                  |
|------------------------------|----------------------------------------------------------------------------|
| `create_customer.json`       | Organização `customer` no topo, filha do root semeado (ACME Corporation).  |
| `create_site.json`           | Organização `site` filha de um customer via placeholder `{{PARENT_ID}}`.   |
| `create_building.json`       | Organização `building` filha de um site via placeholder `{{PARENT_ID}}`.   |
| `create_minimal.json`        | Payload mínimo válido para customer — guia as asserts de campos obrigatórios. |
| `update_name.json`           | Corpo de `PATCH` que renomeia a organização.                               |
| `update_disable.json`        | Corpo de `PATCH` que alterna `enabled = false`.                            |
| `update_full.json`           | Corpo de `PATCH` que reescreve nome, endereço, telefone e access policy.   |

## Como executar

```bash
cd e2e_tests

# Pacote inteiro
go test ./services/mapexos/organizations -v

# Teste único
go test ./services/mapexos/organizations -v -run TestOrganizationHierarchy_PathKeyPropagation
```

## Resultado esperado em caso de sucesso

- Round-trip de CRUD: create / read / list / patch / delete retornam os
  status documentados e o documento persistido faz round-trip.
- Hierarquia profunda: `customer -> site -> building` é criada com sucesso
  e reporta o `parentOrgId` correto em cada nível.
- Propagação de `pathKey`: o `pathKey` de cada descendente é exatamente o
  do pai mais `/segment`; o `pathKey` de um building tem quatro segmentos.
- Herança de `customerId`: um customer é seu próprio `customerId`; sites e
  buildings abaixo dele herdam o mesmo valor.
- Paginação em árvore: `/tree` expõe cursores `next` / `previous` e flags
  `hasNext` / `hasPrevious`, e caminhadas para frente e para trás caem em
  páginas disjuntas.
- Middleware: ROOT (`mapex.*`) passa com ou sem `X-Org-Context`; um admin
  restrito sem o header recebe 403, e um admin restrito apontando para
  uma org fora da sua cobertura também recebe 403.

## Requisitos

- Stack do `mapexOSDeploy/` no ar (mongo, redis, NATS, mapexIam).
- `mapexos / iam` escutando em `:5000`.
- Admin semeado (`admin@mapex.local` / `mapex@123`), id do root semeado
  (`0000000000000000000aa001`) e id da role SuperAdmin semeada
  (`0000000000000000000aa201`) provisionados pelo `mongodb-init`.
- Go 1.25+.

## Notas

- O pacote contém um helper `provisionRestrictedAdmin` que monta, somente
  via API pública, um admin não-wildcard descartável: ele cria uma org
  customer, anexa uma role de escopo `local` com apenas
  `organization.read` + `organization.list`, em seguida faz onboarding de
  um usuário via `POST /api/v1/onboarding/users` e autentica para obter
  um JWT. O token resultante é o oposto do super-admin semeado — escopo
  de org, sem o wildcard `mapex.*` — e é o único caminho para exercitar
  as rotas de negação do middleware (`AdminWithoutOrgContext_Deny`,
  `AdminWithUnauthorizedOrgContext_Deny`).
- `coveragePropagationDelay` (4 segundos) é uma pausa intencional após
  qualquer escrita que mute organizações. O cache de cobertura do
  mapexIam é invalidado por um evento NATS, então requisições em
  sequência do mesmo cliente precisam dessa janela antes que a próxima
  chamada consiga apontar para o descendente recém-criado via
  `X-Org-Context`. Os helpers (`createTestOrganization`,
  `provisionRestrictedAdmin`) já dormem automaticamente; os testes de
  hierarquia dormem novamente entre níveis porque criam mais de um filho.
- Nenhum `t.Skip` permanece neste pacote — todos os testes rodam em toda
  execução.
