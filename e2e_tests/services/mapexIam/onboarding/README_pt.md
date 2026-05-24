# Module e2e: mapexIam / onboarding

Blocos de saga para o orchestrator de onboarding do IAM — o
entrypoint HTTP único que cria user + group + membership de forma
atômica. Espelha o módulo onboarding em
`workspace_go/services/mapexIam/src/modules/onboarding/`. Saga
journeys importam este pacote sempre que precisam de um ator não-admin
novo dentro da org scratch da saga; o orchestrator junta o que seriam
três POSTs sequenciais (user, group, membership) em uma chamada
transacional, então o step da saga permanece uma unidade
Do/Compensate só.

## Endpoints exercitados

- `POST /api/v1/onboarding/users` — cria user + group + membership em
  uma chamada atômica; devolve o id e o email do novo usuário e os
  ids dos groups criados junto (modo NewGroup).

## Fixtures

Nenhuma. O payload é montado em Go a partir de `c.RunID` mais o role
id lido do bag — ver `payloads/saga_iot_admin_user.go`. A constante
de senha compartilhada (`SagaIoTAdminUserPassword`) é hard-coded
junto do builder para o step de auth conseguir logar sem passar o
segredo pelo bag.

## Blocos de saga

- `steps/create_user_with_memberships.go` —
  `CreateUserWithMemberships` faz POST do payload canônico de IoT
  admin e publica `iam.userID`, `iam.userEmail` e `iam.groupID` no
  bag. Lê `iam.roleID`, então uma role precisa ter sido criada antes
  na saga (ver `../roles/steps`).
- `steps/keys.go` — exporta `BagKeyUserID`, `BagKeyUserEmail`,
  `BagKeyGroupID`; consumidores importam as constantes em vez de
  usar string literals.
- `payloads/saga_iot_admin_user.go` — `SagaIoTAdminUser(runID, roleID)`,
  builder fluente que devolve um DTO `CreateUserWithMemberships` com
  email carimbado por runID e um group novo bound à role recebida. A
  constante `SagaIoTAdminUserPassword` é a senha determinística que
  todo usuário de teste de saga compartilha.

## Como rodar

Não há funções `Test*` de módulo neste pacote, então
`go test ./services/mapexIam/onboarding/...` é só um check de
compilação. O step executa quando uma saga journey que o importa
roda:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap
```

Para rodar uma única journey que exercita o step de onboarding:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/... -v -run TestPhase0
```

## O que passar comprova

Uma run verde que importa este step prova:

- O orchestrator de onboarding cria user + group + membership em uma
  chamada transacional.
- O usuário novo consegue depois autenticar via `POST /auth/login`
  usando `SagaIoTAdminUserPassword`.
- O wiring de membership está visível na coverage cache (validado
  pelo assert correspondente no pacote auth).

## Requisitos

- `mapexIam` acessível em `MAPEXOS_URL` (default `http://localhost:5000`).
- Uma role provisionada antes na saga via
  `mapexIam/roles.CreateRole` para que `iam.roleID` esteja no bag.
- Um `X-Org-Context` setado no client HTTP (o runner da saga seta a
  partir do bag depois de `CreateOrganization` ou `SeedAdminLogin`).

## Notas

- O orchestrator é atômico: uma falha no meio (ex.: group já existe
  mas a escrita do membership falha) faz rollback dos três writes do
  lado server. O step da saga tem `Compensate` no-op porque o
  Compensate de org em `mapexIam/organizations` faz cascade-delete
  do user, group e membership na limpeza; rodar um delete por step
  aqui faria race com o cascade.
- O builder default usa modo `NewGroup` (cria um group novo junto do
  user). Testes que precisem de um shape diferente devem adicionar
  um override fluente no builder em vez de montar um DTO na mão.
- O email carimbado por runID evita que runs paralelas de saga colidam
  na constraint de email único.
