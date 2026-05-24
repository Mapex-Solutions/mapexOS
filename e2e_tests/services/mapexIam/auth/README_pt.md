# Module e2e: mapexIam / auth

Blocos de construção de saga para a superfície de autenticação do IAM
— as rotas que emitem e validam o bearer que todo outro teste de
módulo depende. Espelha o módulo auth em
`workspace_go/services/mapexIam/src/modules/auth/`. Não existe suite
CRUD por endpoint aqui porque autenticação não tem estado próprio
para listar ou atualizar; em vez disso, este pacote expõe os steps e
asserts que as saga journeys (e testes de módulo em outros pacotes)
compõem para logar, trocar de identidade e verificar o JWT resultante.

## Endpoints exercitados

Os building blocks chamam a superfície auth exposta pelo `mapexIam`:

- `POST /auth/login` — troca email + senha por access token; a
  resposta também carrega refresh token. Montado no root do serviço
  (sem prefixo `/api/v1`), diferente do resto do IAM.
- `GET /auth/users/me/coverage` — devolve as orgs que o chamador pode
  atuar; usado para verificar wiring de membership pós-login.

## Fixtures

Nenhuma. Credenciais vêm de `common/constants` (seed admin) ou do
email publicado no bag pelo step de onboarding, somado à constante de
senha compartilhada em `../onboarding/payloads`.

## Blocos de saga

- `steps/seed_admin_login.go` — `SeedAdminLogin` loga como o seed
  admin que o `mongodb-init` provisiona, atribui bearer +
  `X-Org-Context` em cada client por serviço, e publica
  `iam.userJWT` mais `iam.organizationID` no bag. Toda journey começa
  por aqui, a menos que precise de um ator não-admin.
- `steps/authenticate_user.go` — `AuthenticateUser` troca o seed
  admin pelo usuário criado mais cedo na journey, lendo
  `iam.userEmail` + `iam.organizationID` do bag.
- `steps/keys.go` — exporta `BagKeyUserJWT` (`iam.userJWT`);
  importado pelos consumidores downstream para que renames quebrem o
  build em vez de falharem em runtime.
- `asserts/assert_jwt_valid.go` — `AssertJwtValid` verifica que o
  bearer é um JWT bem-formado de três segmentos;
  `AssertJwtHasOrgContext` chama `/auth/users/me/coverage` e
  confirma que a saga org está na lista.

## Como rodar

Não há funções `Test*` de módulo neste pacote, então
`go test ./services/mapexIam/auth/...` é só um check de compilação.
Os steps e asserts executam quando uma saga journey que os importa
roda:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap
```

Para exercitar o fluxo de auth dentro de uma journey:

```bash
cd e2e_tests
go test -tags=saga ./journey/iot/mqtt_broker_auth/... -v -run TestPhase0
```

## O que passar comprova

Uma run verde que importa estes blocos prova:

- O seed admin provisionado pelo `mongodb-init` consegue autenticar
  contra a stack ao vivo.
- Um usuário criado pela saga (via o orchestrator de onboarding)
  consegue autenticar com as próprias credenciais.
- O bearer devolvido é estruturalmente válido e a coverage cache
  reflete o membership do usuário na organização da saga.

## Requisitos

- `mapexIam` acessível em `MAPEXOS_URL` (default `http://localhost:5000`).
- Mongo populado pelo `mongodb-init` para as credenciais do seed
  admin funcionarem.
- Para `AuthenticateUser`: um step prévio criando o usuário
  (`mapexIam/onboarding`) e a org pai (`mapexIam/organizations`).

## Notas

- `/auth/login` é o único endpoint IAM fora de `/api/v1`. Os demais
  módulos IAM (organizations, roles, onboarding) ficam todos sob o
  prefixo versionado.
- `AssertJwtValid` é estrutural de propósito — a verificação de
  assinatura acontece implicitamente na próxima vez que a saga chama
  um endpoint protegido e recebe algo diferente de `401`.
- `AssertJwtHasOrgContext` é o canário canônico de frescor do cache
  de membership: verde aqui prova que o grant da role dentro da saga
  org propagou para o cache.
