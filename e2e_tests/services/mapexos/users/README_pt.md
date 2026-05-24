# E2E de módulo: mapexos / users

## Escopo

Cobertura end-to-end do módulo de users — a superfície de identidade
humana do mapexos servida pelo `mapexIam`
(`workspace_go/services/mapexIam/src/modules/users/`). A suíte prova que
usuários criados pelo orquestrador público de onboarding (usuário +
membership atômicos em uma única chamada) fazem round-trip por fetch,
list, patch e delete, que tanto o provedor `internal` quanto o `google`
são aceitos, e que os payloads de update (nome, senha, perfil completo,
desabilitar) chegam ao documento persistido. Cada usuário é provisionado
com membership de escopo `local` na org Mapexos semeada sob a role
SuperAdmin, e depois removido via `t.Cleanup`.

## Endpoints exercitados

- `POST   /api/v1/onboarding/users`     — criação atômica de usuário + membership
- `POST   /api/v1/users`                — create direto (usado apenas no guard de email inválido)
- `GET    /api/v1/users/{id}`           — busca um usuário
- `GET    /api/v1/users`                — listagem paginada (page / perPage)
- `PATCH  /api/v1/users/{id}`           — atualização parcial (name, password, full, disable)
- `DELETE /api/v1/users/{id}`           — exclusão (idempotente em 404)

## Fixtures

| Arquivo                 | Propósito                                                                |
|-------------------------|--------------------------------------------------------------------------|
| `create_internal.json`  | Usuário completo com auth interna: senha, telefone, cargo e avatar.      |
| `create_google.json`    | Usuário Google OAuth com `externalId` + metadata do provedor.            |
| `create_minimal.json`   | Usuário mínimo com auth interna (email, senha, primeiro / último nome).  |
| `update_name.json`      | Corpo de `PATCH` que renomeia `firstName` e `lastName`.                  |
| `update_password.json`  | Corpo de `PATCH` que rotaciona a senha e força troca no próximo login.   |
| `update_full.json`      | Corpo de `PATCH` que reescreve email, nome, telefone, cargo e avatar.    |
| `update_disable.json`   | Corpo de `PATCH` que alterna `enabled = false`.                          |

## Como executar

```bash
cd e2e_tests

# Pacote inteiro
go test ./services/mapexos/users -v

# Teste único
go test ./services/mapexos/users -v -run TestCreateUser_Internal
```

## Resultado esperado em caso de sucesso

- Round-trip de onboarding: postar um payload internal, Google ou
  mínimo em `/api/v1/onboarding/users` retorna `{user, memberships}` e
  o usuário fica imediatamente recuperável pelo id.
- Validação: um email malformado no endpoint direto `/api/v1/users`
  resulta em 400 sem persistir nada.
- Leituras: `GET /api/v1/users/{id}` retorna o email, primeiro e último
  nome persistidos; um id desconhecido resulta em 404.
- Updates: nome, senha (com `changePasswordNextLogin = true`), perfil
  completo e `enabled = false` aparecem no `GET` seguinte.
- Delete: um usuário excluído desaparece do caminho de leitura com 404.
- List: a listagem paginada expõe o usuário recém-criado no array
  `items` e reporta os metadados `totalItems / page / perPage`.

## Requisitos

- Stack do `mapexOSDeploy/` no ar (mongo, redis, NATS, mapexIam).
- `mapexos / iam` escutando em `:5000`.
- Admin semeado (`admin@mapex.local` / `mapex@123`), id da org Mapexos
  semeada (`constants.MapexosOrgID`) e id da role SuperAdmin semeada
  (`constants.SuperAdminRoleID`) provisionados pelo `mongodb-init` —
  cada usuário de teste é onboardado nessa org com essa role.
- Go 1.25+.

## Notas

- O pacote nunca faz POST em `/api/v1/users` no caminho feliz de
  criação — sempre passa pelo orquestrador em
  `POST /api/v1/onboarding/users` para que a membership seja conectada
  atomicamente. O endpoint direto `/api/v1/users` só é usado por
  `TestCreateUser_InvalidEmail` para confirmar que a validação rejeita
  payloads inválidos.
- `loadFixture` retorna apenas os campos base do usuário; tanto
  `createTestUser` quanto os casos `TestCreateUser_*` explícitos os
  envelopam com um array `memberships` fixado em `MapexosOrgID` +
  `SuperAdminRoleID` com escopo `local` antes de postar.
- Um cliente ROOT (`mapex.*` com `X-Org-Context = MapexosOrgID`)
  dirige todos os testes; o `adminClient` é montado no `TestMain` por
  paridade com pacotes irmãos, mas a suíte de users não possui casos
  de negação de middleware.
