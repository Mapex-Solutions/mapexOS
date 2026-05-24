# Módulo e2e: mapexos / auth

## Escopo

Cobre a superfície de autenticação da plataforma mapexos — login,
logout, refresh de token e o endpoint "quem sou eu" de coverage,
expostos por `workspace_go/services/mapexIam/src/modules/auth/`. Ao
contrário dos outros módulos do mapexos, este pacote não precisa de
um JWT bootstrap no `TestMain`: ele é justamente a superfície que
emite esses JWTs, então cada teste executa o fluxo de login sozinho.

## Endpoints exercitados

- `POST /auth/login` — troca credenciais por um par access + refresh token.
- `POST /auth/logout` — invalida o access token corrente.
- `POST /auth/refresh` — rotaciona tokens; o refresh token vai no
  header `X-Refresh-Token`.
- `GET /auth/users/me/coverage` — retorna a coverage organizacional do
  usuário autenticado.

## Fixtures

| Arquivo | Descrição |
|---|---|
| `login_valid.json` | Credenciais do admin seed (`admin@mapex.local` / `mapex@123`), `keepConnected: true`. |
| `login_invalid_email.json` | Email malformado — dispara 400. |
| `login_short_password.json` | Senha abaixo dos 8 caracteres mínimos — dispara 400. |
| `login_wrong_password.json` | Email válido, senha errada — dispara 401. |

## Como rodar

```bash
cd e2e_tests
go test ./services/mapexos/auth -v

# Teste individual
go test ./services/mapexos/auth -v -run TestRefreshToken
```

## Resultado em caso de PASS

- `TestLogin_Valid` prova que o admin seed consegue logar e que a
  resposta traz `access_token`, `refresh_token` e um objeto `user`
  com `id` e `email`.
- `TestLogin_InvalidEmail` e `TestLogin_ShortPassword` provam que a
  validação de entrada rejeita payloads malformados com
  `400 Bad Request`.
- `TestLogin_WrongPassword` prova que credencial inválida retorna
  `401 Unauthorized`.
- `TestLogout` prova que um access token recém-emitido é aceito por
  `/auth/logout`.
- `TestRefreshToken` prova que o fluxo de refresh rotaciona os dois
  tokens quando o refresh token vem em `X-Refresh-Token`.
- `TestGetMyCoverage` e `TestGetMyCoverage_Unauthorized` provam que o
  endpoint de coverage retorna as orgs acessíveis do chamador e
  rejeita requisições não autenticadas com `401`.

## Requisitos

- Serviço mapexos / iam acessível em `http://localhost:5000`
  (sobrescrever via `MAPEXOS_URL`).
- Usuário admin seed, role, organização e membership recursiva
  provisionados pelo `mongodb-init` no primeiro boot da stack.

## Notas

- O `TestMain` propositalmente NÃO chama `SetupE2EEnvironment` e NÃO
  obtém token previamente — este pacote é a fonte da verdade do
  próprio contrato de login/refresh/coverage.
- Transporte do refresh token: header `X-Refresh-Token`, não no body.
  O access token continua em `Authorization: Bearer`.
- Nenhum header `X-Org-Context` é enviado aqui; o endpoint de coverage
  resolve escopo somente pelo bearer.
