# E2E de módulo: mapexos / roles

## Escopo

Cobertura end-to-end do módulo de roles — a primitiva de concessão de
permissões da superfície mapexos servida pelo `mapexIam`
(`workspace_go/services/mapexIam/src/modules/roles/`). A suíte prova que
roles de sistema e roles com escopo de organização fazem round-trip por
CRUD, que permissões wildcard e namespaced (`mapex.*`, `user.*`,
`admin.*`, `asset.read`) sobrevivem ao create + read intactas, e que a
validação rejeita payloads malformados. Uma org customer descartável é
provisionada no `TestMain` para hospedar todas as roles com escopo de
organização; o dataset semeado nunca é modificado.

## Endpoints exercitados

- `POST   /api/v1/roles`                — cria roles de sistema / de org
- `GET    /api/v1/roles/{id}`           — busca uma role
- `GET    /api/v1/roles`                — listagem paginada (page / perPage)
- `PATCH  /api/v1/roles/{id}`           — atualização parcial (name, permissions, full)
- `DELETE /api/v1/roles/{id}`           — exclusão (idempotente em 404)
- `POST   /api/v1/organizations`        — bootstrap da org pai descartável no `TestMain`
- `DELETE /api/v1/organizations/{id}`   — derruba a org ao final da suíte

## Fixtures

| Arquivo                      | Propósito                                                                |
|------------------------------|--------------------------------------------------------------------------|
| `create_system_role.json`    | Role de sistema global com `mapex.*` (`isSystem: true`, `scope: global`).|
| `create_org_role.json`       | Role `Site Manager` com escopo de org e permissões de user / asset.      |
| `create_minimal.json`        | Menor role válida de org (`Viewer` com duas permissões de read).         |
| `update_name.json`           | Corpo de `PATCH` que renomeia a role.                                    |
| `update_permissions.json`    | Corpo de `PATCH` que substitui o array de permissões.                    |
| `update_full.json`           | Corpo de `PATCH` que reescreve name, description e permissions.          |
| `update_disable.json`        | Corpo de `PATCH` que alterna `enabled = false` (paridade; v1 ignora).    |

## Como executar

```bash
cd e2e_tests

# Pacote inteiro
go test ./services/mapexos/roles -v

# Teste único
go test ./services/mapexos/roles -v -run TestCreateRole_SystemRole
```

## Resultado esperado em caso de sucesso

- Round-trip de CRUD: roles de sistema e de org podem ser criadas,
  buscadas, listadas (com metadados de paginação), atualizadas e
  excluídas.
- O `orgId` é resolvido pelo serviço a partir do `X-Org-Context` em vez
  de ser confiado no payload — a resposta carrega um id preenchido
  mesmo quando o teste posta um placeholder.
- Permissões wildcard: `mapex.*`, `user.*`, `admin.*` e `asset.*`
  sobrevivem ao create + read sem alteração.
- Validação: `name` ausente, `orgId` ausente em roles não-sistema e
  arrays de `permissions` vazios são rejeitados com 400.

## Requisitos

- Stack do `mapexOSDeploy/` no ar (mongo, redis, NATS, mapexIam).
- `mapexos / iam` escutando em `:5000`.
- Admin semeado (`admin@mapex.local` / `mapex@123`) e id da org Mapexos
  semeada (`constants.MapexosOrgID`) provisionados pelo `mongodb-init`
  — a org descartável dos testes é criada como filha dela.
- Go 1.25+.

## Notas

- O `TestMain` cria uma única org customer compartilhada (`Test
  Organization for Roles`) e a remove após o `m.Run()`. As fixtures com
  escopo de organização usam `{{CUSTOMER_ID}}` como placeholder, que o
  `loadFixture` reescreve com o id em tempo de execução.
- O cliente de CRUD é um cliente com token ROOT (`mapex.*`) e
  `X-Org-Context` fixado em `MapexosOrgID`. Um `adminClient` também é
  montado no `TestMain` por paridade com pacotes irmãos, mas os testes
  de roles em si dependem do cliente ROOT porque cada assert é um
  resultado de CRUD e não uma negação de middleware.
- O campo `enabled` nos payloads de role é um resquício de v1; o
  serviço o ignora, então `update_disable.json` é mantido apenas por
  simetria e não é objeto de assert.
