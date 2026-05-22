# Contexto: iot

Saga journeys que exercitam o pipeline IoT — assets, route groups,
asset templates, callout de auth MQTT, telemetria, healthmonitor,
triggers, workflow e o sink de events. Toda journey neste contexto
assume que o usuário admin seed está autenticado (a fase 0 de cada
journey cuida disso).

## Journeys registradas

| Journey                    | Fases                                                | O que a história cobre                                                       |
|----------------------------|------------------------------------------------------|------------------------------------------------------------------------------|
| mqtt_broker_auth           | phase0..phase3                                       | Lifecycle MQTT por senha + cert + cascata TieredStore                        |
| connectivity_actions_http  | phase1_workflow, phase2_trigger                      | Healthmonitor de assets HTTP → workflow + trigger por route group            |
| connectivity_actions_mqtt  | phase1_workflow, phase2_trigger                      | Healthmonitor de assets MQTT → workflow + trigger por route group            |

Nomes de pasta de fase carregam o descritor; nomes de pacote dentro
delas são curtos (`phase0`, `phase1`) para que o alias de import deixe
a intenção óbvia.

## Como rodar

Todo comando roda a partir da raiz do pacote e2eTests.

```bash
cd workspace_go/packages/e2eTests

# Todas as fases de toda journey neste contexto
go test -tags=saga -v ./journey/iot/...

# Todas as fases de uma journey
go test -tags=saga -v ./journey/iot/mqtt_broker_auth/...

# Uma fase só
go test -tags=saga -v ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap/
go test -tags=saga -v ./journey/iot/mqtt_broker_auth/phase1_password_user/
```

A build tag `saga` controla esses testes: `go test ./...` (sem tag)
pula; só `go test -tags=saga` percorre as pastas de journey.

## Ambiente necessário

- mapexIam   acessível em `MAPEXOS_URL` (default `http://localhost:5000`)
- assets     acessível em `ASSETS_URL`  (default `http://localhost:5002`)
- router     acessível em `ROUTER_URL`  (default `http://localhost:5003`)
- http_gateway acessível em `GATEWAY_URL` (default `http://localhost:5001`)
- Usuário admin seed provisionado pelo seed canônico do mongodb-init
  (`admin@mapex.local` / `mapex@123`)

## Como adicionar uma nova journey neste contexto

1. `mkdir journey/iot/<journey_name>` (snake_case).
2. Crie `README.md` (e `README_pt.md`) com a narrativa, índice de
   fases e uma cópia do bloco "Como rodar" adaptada ao caminho da
   journey.
3. Adicione `phaseN_<descritor>/journey.go` e `journey_test.go` para
   cada fase. Reutilize a Phase 0 de uma journey existente como
   bootstrap quando a nova journey parte do mesmo ator; caso contrário
   escreva o próprio.
4. Adicione uma linha na tabela de journeys registradas acima.
