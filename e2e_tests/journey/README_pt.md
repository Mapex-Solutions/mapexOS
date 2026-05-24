# Hierarquia de Journeys

Os testes escalam aos milhares conforme a plataforma cresce. Para se
manterem encontráveis e para manter o reuso entre journeys honesto, a
árvore de diretórios é organizada assim:

```
journey/
├── README.md                     # este arquivo (hierarquia global + regras)
└── {context}/                    # domínio de negócio (iot, automations, workflow, vault, platform, ...)
    ├── README.md                 # narrativa do contexto + comandos de run de toda journey aqui
    └── {journey_name}/           # uma história ponta a ponta nomeada (snake_case)
        ├── README.md             # o que a journey cobre + outcome resumido
        └── {phaseN_descriptor}/  # uma pasta por fase do journey
            ├── journey.go        # helpers locais + Run + (opcional) ItemsForCompose
            └── journey_test.go   # //go:build saga; uma func TestPhaseN_<Descriptor>_Saga
```

## Regras

- **Context** agrupa journeys por domínio de negócio. Exemplos: `iot/`,
  `automations/`, `workflow/`, `vault/`, `iam/`, `platform/`
  (cross-cutting). Cria um contexto novo só quando um existente
  enganaria quem busca a journey.

- **Journey** é a história nomeada (snake_case). Uma journey = uma
  narrativa coesa ponta a ponta. Exemplos: `mqtt_full_pipeline`,
  `trigger_http`, `mqtt_broker_auth`, `connectivity_actions_http`.

- **Phase** é uma etapa nomeada da journey. As fases rodam em ordem
  (Phase 0 antes de Phase 1, etc). Cada pasta de fase é um package Go;
  o nome do package é curto (`phase0`, `phase1`) para os consumidores
  aliarem como `phase0 ".../phase0_iam_bootstrap"`. Pastas de fase
  carregam sufixo descritivo (`phase0_iam_bootstrap`) para o nome da
  pasta sozinho dizer o que a fase faz.

- **Reuso cross-journey** nunca acontece importando a fase de outra
  journey. O reuso vive em `services/{svc}/{mod}/{steps,asserts,payloads}`
  — os blocos de construção que toda fase compõe. Dentro de uma única
  journey, as fases PODEM importar uma à outra (PhaseN+1 tipicamente
  compõe `PhaseN.BootstrapItems`).

- **Cada fase DEVE carregar um bloco de Outcome** tanto em `journey.go`
  (godoc do package) quanto em `journey_test.go` (godoc da func de
  teste) descrevendo o que passar a fase prova e a que tipo de falha
  geralmente aponta. Os blocos de Outcome são como um dev descobre o
  que uma journey cobre sem abrir o source.

- **Todo README de contexto DEVE carregar um bloco "Como rodar"** com
  os comandos go test no escopo de contexto, journey e fase. Devs caem
  num contexto, veem os comandos, e rodam o escopo certo sem buscar
  docs.

- **Documentação nunca grava caminho absoluto específico de dev.**
  Comandos de run assumem cwd na raiz do repo e usam caminhos
  relativos (ex.: `cd e2e_tests`). Quem clonar o monorepo copia e cola
  sem editar.

## Journeys atualmente registradas

| Contexto    | Journey                       | Fases                                                                                  |
|-------------|-------------------------------|----------------------------------------------------------------------------------------|
| automations | trigger_email                 | phase1_connectivity, phase2_event_pipeline                                             |
| automations | trigger_http                  | phase1_connectivity, phase2_event_pipeline                                             |
| automations | trigger_mqtt                  | phase1_connectivity, phase2_event_pipeline                                             |
| automations | trigger_nats                  | phase1_connectivity, phase2_event_pipeline                                             |
| automations | trigger_rabbitmq              | phase1_connectivity, phase2_event_pipeline                                             |
| automations | trigger_slack                 | phase1_connectivity, phase2_event_pipeline                                             |
| automations | trigger_teams                 | phase1_connectivity, phase2_event_pipeline                                             |
| automations | trigger_websocket             | phase1_connectivity, phase2_event_pipeline                                             |
| iot         | connectivity_actions_http     | phase1_workflow, phase2_trigger                                                        |
| iot         | connectivity_actions_mqtt     | phase1_workflow, phase2_trigger                                                        |
| iot         | mqtt_broker_auth              | phase0_iam_bootstrap, phase1_password_user, phase2_cert_user, phase3_cascade (skeleton) |

## Como rodar

Todo comando roda a partir do diretório `e2e_tests`.

```bash
cd e2e_tests

# Todas as journeys de todos os contextos
go test -tags=saga -v ./journey/...

# Todas as journeys de um contexto (ver README do contexto pra escopos mais finos)
go test -tags=saga -v ./journey/automations/...
go test -tags=saga -v ./journey/iot/...
```

O build tag `saga` gateia esses testes: `go test ./...` (sem tag)
pula; só `go test -tags=saga` percorre as pastas de journey.

## Como adicionar uma nova journey

1. Escolha o contexto (ou crie um se nenhum existente serve).
2. Se criar contexto, adicionar `journey/{context}/README.md` com a
   narrativa, tabela de journeys registradas, e bloco "Como rodar".
3. Criar `journey/{context}/{journey_name}/README.md` com a narrativa
   e índice de fases.
4. Criar `phaseN_<descriptor>/journey.go` e `journey_test.go`. Copiar
   o layout de uma fase existente por consistência.
5. Atualizar a tabela de journeys registradas deste README e a do
   README do contexto.
6. PhaseN reusa Phase0..N-1 importando direto dentro da mesma pasta
   de journey. Reuso cross-journey só no nível de blocos de
   construção.
