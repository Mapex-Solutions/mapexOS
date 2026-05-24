# Module e2e: mapexIam / organizations

Este pacote cobre o módulo de organizations do IAM em duas
superfícies:

- O subpacote `e2e/` é a suite canônica de module e2e — CRUD, list
  paginada, filtros de query e cursor walks de tree contra o serviço
  `mapexIam` ao vivo.
- As pastas irmãs `steps/`, `payloads/` e `asserts/` são blocos de
  saga que outras journeys (e a própria suite e2e) importam para
  provisionar uma org customer scratch com shape determinístico.

Espelha o módulo organizations em
`workspace_go/services/mapexIam/src/modules/organizations/`.

## Endpoints exercitados

A suite `e2e/` exercita toda rota pública do módulo:

- `POST /api/v1/organizations` — create (happy path, validação,
  auth).
- `GET /api/v1/organizations/{id}` — read por id (happy path, 404).
- `PATCH /api/v1/organizations/{id}` — update parcial (happy path,
  404).
- `DELETE /api/v1/organizations/{id}` — delete + 404 pós-delete
  (happy path, 404).
- `GET /api/v1/organizations` — list paginada com filtros `page`,
  `perPage`, `name`, `type`, `enabled`; também valida `totalItems`,
  páginas fora do range e walks de página forward + backward.
- `GET /api/v1/organizations/tree` — navegação cursor-based sob o
  `X-Org-Context` ativo.

## Fixtures

Nenhuma em disco. A suite e2e e os blocos de saga montam payload em
Go a partir de `c.RunID` via `payloads.SagaTestCustomerOrg`, que
espelha o DTO de contrato em
`packages/contracts/services/mapexIam/organizations`.

## Blocos de saga

Importados pelas saga journeys (e pela suite e2e para construção do
payload):

- `steps/create_organization.go` — `CreateOrganization` faz POST do
  payload canônico de customer org, publica `iam.organizationID` e
  `iam.organizationPathKey` no bag e faz cascade-delete dos filhos
  (users, groups, roles, memberships) no `Compensate`.
- `steps/keys.go` — exporta `BagKeyOrgID` e `BagKeyOrgPathKey`.
- `payloads/saga_test_org.go` — `SagaTestCustomerOrg(runID)`, builder
  fluente; defaults `Type=customer`, `ParentOrgID` = id da seed root,
  `Enabled=true`, IDP interno, role policy `strict`, default scope
  `local`. `WithName` e `WithParentOrgID` são os overrides usuais.
- `asserts/assert_organization_exists.go` —
  `AssertOrganizationExists` busca a org pelo id do bag e confirma
  que a API devolve com `enabled=true`.

## Como rodar

A suite executável vive no subpacote `e2e/`:

```bash
cd e2e_tests

# Suite completa
go test ./services/mapexIam/organizations/e2e -v

# Um teste só
go test ./services/mapexIam/organizations/e2e -v -run TestCreate_201

# Todos os list / tree
go test ./services/mapexIam/organizations/e2e -v -run 'TestList_|TestTree_'
```

Os blocos de saga em `steps/`, `payloads/`, `asserts/` não têm
funções `Test*` — executam quando uma saga journey os importa
(`go test -tags=saga ./journey/...`).

## O que passar comprova

Passar a suite `e2e/` prova:

- Toda rota CRUD pública devolve o status code esperado em happy
  path, validação, auth e cenários de not-found.
- Paginação por página é estável: walks forward (1 -> 15) e backward
  (15 -> 1) com `perPage=1` visitam cada fixture exatamente uma vez;
  `totalItems` bate com o tamanho do universo independente de
  `perPage`; páginas fora do range devolvem `items=[]` em vez de
  clampar ou retornar 404.
- Filtros (`name`, `type`, `enabled`) compõem corretamente sob AND.
- O cursor walk de `/tree` cobre cada fixture dentro do
  `X-Org-Context` ativo, parando direito quando `hasNext=false`.

## Requisitos

- `mapexIam` acessível em `MAPEXOS_URL` (default `http://localhost:5000`).
- Mongo populado pelo `mongodb-init` (o seed admin precisa conseguir
  logar e o id da seed root org `0000000000000000000aa001` precisa
  existir para parentear as orgs criadas pela saga).
- A suite de list cria 15 fixtures por run; o cleanup é automático
  via `t.Cleanup`, mas a run só deixa zero resíduo se a stack estiver
  acessível durante o teardown.

## Notas

- A suite de list isola cada run carimbando o `runID` no nome da
  fixture e filtrando as queries seguintes com
  `?name=<orgNamePrefix>-<runID>`. Isso mantém orgs pré-populadas e
  runs paralelas fora do universo das asserções.
- `listFixtureCount = 15` foi escolhido de propósito para que walks
  com `perPage=1` cubram bordas de paginação (página 1, páginas do
  meio, última página) que fixture sets menores não cobririam.
- O PATCH em `TestUpdate_200` tolera 200 e 201 para resistir a
  mudanças de status no handler que não afetam semântica; o GET que
  vem em seguida é a asserção real.
- O cursor walk de `/tree` limita iteração a 200 páginas como
  salvaguarda contra bugs no backend que looparam para sempre; em
  comportamento correto o walk sai por `hasNext=false` bem antes
  desse limite.
- Os blocos de saga moram ao lado da suite e2e para que o builder de
  payload seja compartilhado — os testes e2e são o contract test mais
  confiável para o payload da saga em si.
