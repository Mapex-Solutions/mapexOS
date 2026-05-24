# Module e2e: mapexIam / roles

Blocos de saga para o módulo de roles do IAM — a superfície que
define bundles nomeados de permissões atribuídos a usuários e groups.
Espelha o módulo roles em
`workspace_go/services/mapexIam/src/modules/roles/`. A fixture
canônica aqui junta toda permissão que um usuário de teste de saga
precisa para tocar o pipeline IoT ponta a ponta (CRUD de org, user,
group, asset, route, trigger, workflow mais leitura de events), então
as journeys não precisam escolher strings de permissão na mão toda
vez que provisionam um ator.

## Endpoints exercitados

- `POST /api/v1/roles` — cria role com nome, scope, flags e lista de
  permissões dados.
- `DELETE /api/v1/roles/{id}` — usado pelo `Compensate` na limpeza;
  a role é criada na seed parent org (onde o ator de bootstrap tem
  coverage), então o cascade da org-filha não chega nela e o delete
  explícito é necessário.

## Fixtures

Nenhuma. O payload é montado em Go a partir de `c.RunID` mais
constantes de permissão importadas de
`permissions/{mapexos,assets,router}` — ver
`payloads/saga_iot_admin_role.go`.

## Blocos de saga

- `steps/create_role.go` — `CreateRole` faz POST da role canônica de
  IoT admin e publica `iam.roleID` no bag. Tem `Compensate` próprio
  que chama `DELETE /api/v1/roles/{id}` (idempotente em 404).
- `steps/keys.go` — exporta `BagKeyRoleID` (`iam.roleID`); importado
  pelo step de onboarding para que renames quebrem o build.
- `payloads/saga_iot_admin_role.go` — `SagaIoTAdminRole(runID)`,
  builder fluente que devolve um DTO `RoleCreate` com
  `Scope=local`, `IsSystem=false`, `IsTemplate=false` e o conjunto
  canônico de permissões de IoT admin. `WithName` é o único override
  que a maioria dos callers precisa.

## Como rodar

Não há funções `Test*` de módulo neste pacote, então
`go test ./services/mapexIam/roles/...` é só um check de compilação.
O step executa quando uma saga journey que o importa roda:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap
```

Para mirar em uma journey específica:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/... -v -run TestPhase0
```

## O que passar comprova

Uma run verde que importa este step prova:

- Uma role com `Scope=local` pode ser criada sob a seed parent org.
- As permissões da role aparecem na resposta (e, uma vez bound a um
  user via o step de onboarding, expandem corretamente na coverage
  cache).
- O delete do `Compensate` limpa a role mesmo quando o cascade da
  org não consegue alcançá-la.

## Requisitos

- `mapexIam` acessível em `MAPEXOS_URL` (default `http://localhost:5000`).
- O caller precisa estar autenticado com coverage na seed parent org
  — o runner da saga garante isso via `SeedAdminLogin` antes do step
  rodar.

## Notas

- A role é criada de propósito contra a seed parent org, não contra
  a org-filha scratch da saga. O ator de bootstrap (seed admin) tem
  coverage na seed root, então a criação da role sempre passa; o
  trade-off é que o cascade do `CreateOrganization` da saga não
  alcança a role, e precisamos do delete explícito no `Compensate`.
- O bundle de permissões é a união do que cada saga de pipeline IoT
  precisa. Journeys com superfície mais estreita ainda podem reusar
  essa role — coverage acima do que o teste chama é inofensivo.
- As constantes de permissão vêm de
  `permissions/{mapexos,assets,router}` no módulo de contracts, então
  renames na fonte da verdade quebram o payload em compile time.
