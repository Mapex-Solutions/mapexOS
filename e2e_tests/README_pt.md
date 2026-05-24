# Testes E2E MapexOS

Suite de testes end-to-end da plataforma MapexOS. Dois tipos coexistem
aqui:

- **Module e2e** (`services/<svc>/<module>/`) — testes de CRUD e
  contrato de um módulo só. Cada endpoint é exercitado com fixtures
  reais contra a stack rodando.
- **Saga journeys** (`journey/<context>/<journey_name>/phaseN_*/`) —
  fluxos cross-módulo organizados em fases ordenadas. Cada fase é uma
  saga com steps + asserts + rollback. Gated pelo build tag `saga`.

Toda a suite roda contra a stack canônica do
[`mapexOSDeploy`](../../mapexOSDeploy/) e usa o seed admin
(`admin@mapex.local` / `mapex@123`) como ator de bootstrap.

## Pré-requisitos

1. Stack no ar via docker-compose canônico:
   ```bash
   cd ../mapexOSDeploy
   docker compose up -d
   ```
2. Serviços Go rodando nas portas default — ver [Portas](#portas)
   abaixo.
3. Go 1.25+ na máquina que vai rodar os testes.

O seed admin (user, role, organization, membership recursivo) é
provisionado pelo container `mongodb-init` no primeiro boot do mongo.
Testes que precisam de usuários adicionais provisionam em runtime via
o orchestrator público (`POST /api/v1/onboarding/users`) — o seed é
para dados reais/produção e nunca é modificado por teste.

## Layout

```
e2e_tests/
├── common/                     # Código compartilhado
│   ├── constants/              # URLs, ids, credenciais
│   ├── types/                  # StandardResponse, ErrorResponse, ...
│   └── utils/                  # Login, asserções, setup
│
├── services/                   # Module e2e — um pacote por módulo
│   ├── assets/{assets,assettemplates}/
│   ├── events/events/
│   ├── http_gateway/datasources/
│   ├── mapexIam/{auth,onboarding,organizations,roles}/
│   ├── mapexos/{auth,groups,lists,memberships,organizations,roles,users}/
│   ├── router/routegroups/
│   ├── triggers/triggers/
│   └── workflow/{definitions,instances}/
│
└── journey/                    # Saga journeys (build tag: saga)
    ├── automations/            # 8 trigger journeys (http, email, websocket,
    │   │                       #   slack, teams, mqtt, nats, rabbitmq)
    │   └── trigger_<type>/{phase1_connectivity,phase2_event_pipeline}/
    └── iot/
        ├── connectivity_actions_{http,mqtt}/{phase1_workflow,phase2_trigger}/
        └── mqtt_broker_auth/{phase0_iam_bootstrap..phase3_cascade}/
```

## Como rodar

Todos os comandos rodam a partir de `e2e_tests/`.

```bash
cd e2e_tests

# Todos os module e2e (saga tag NÃO setado)
go test ./services/...

# Um módulo específico
go test ./services/mapexos/organizations -v

# Um teste específico
go test ./services/mapexos/organizations -v -run TestCreateOrganization_Customer

# Todas as saga journeys (saga tag OBRIGATÓRIO)
go test -tags=saga ./journey/...

# Contexto / journey / fase específicos
go test -tags=saga ./journey/automations/...
go test -tags=saga ./journey/automations/trigger_http/...
go test -tags=saga ./journey/automations/trigger_http/phase1_connectivity
```

Para suites longas, ajustar `-timeout`:

```bash
go test -tags=saga -timeout 15m ./journey/...
```

## Portas

| Serviço        | Porta | Necessário para                    |
|----------------|-------|------------------------------------|
| mapexos / iam  | 5000  | Todos os testes (auth + org CRUD)  |
| http_gateway   | 5001  | datasources + saga phase 2         |
| assets         | 5002  | assets / assettemplates / saga IoT |
| router         | 5003  | routegroups + saga                 |
| events         | 5004  | fases event-pipeline               |
| triggers       | 5006  | trigger journeys                   |
| workflow       | 5007  | workflow tests + IoT actions       |

Saga journeys também sobem sinks in-process: `11010` (HTTP),
`11025` (SMTP), `11026` (WebSocket). Essas portas precisam estar
livres.

## Variáveis de ambiente

Os defaults batem com a stack canônica. Override só se a sua stack
está em outros hosts/portas:

```bash
export MAPEXOS_URL=http://localhost:5000
export GATEWAY_URL=http://localhost:5001
export ASSETS_URL=http://localhost:5002
export ROUTER_URL=http://localhost:5003
export EVENTS_URL=http://localhost:5004
export TRIGGERS_URL=http://localhost:5006
export WORKFLOW_URL=http://localhost:5007
```

## Convenções

- **Testes usam só a API pública.** Rotas internas (`/internal/*`) são
  fallback de rebuild de cache e nunca são chamadas por teste.
- **O seed admin é o ator de bootstrap.** Qualquer usuário adicional
  que um teste precise é provisionado em runtime via o orchestrator
  (`POST /api/v1/onboarding/users`). O seed JSON nunca é modificado
  por teste.
- **Fixtures ficam ao lado do teste** em pasta `fixtures/` e usam os
  ids canônicos do seed (`0000000000000000000aa001` para root org,
  `0000000000000000000aa201` para a role SuperAdmin).
- **Cleanup é obrigatório.** Todo teste mutante registra `t.Cleanup`
  ou `defer` para apagar o que criou.

## Documentação

Cada pacote module e2e e cada saga journey carrega
`README.md` + `README_pt.md` descrevendo escopo, fixtures e como rodar
isolado. Comece pelo diretório que te interessa — os READMEs por
pasta são a fonte da verdade para cada superfície de teste.

Ver também: [`journey/README.md`](./journey/README.md) para as regras
de hierarquia das saga journeys.
