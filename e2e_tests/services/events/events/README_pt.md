# Module e2e: events / events

## Escopo

Esta pasta está reservada para a suite de e2e de módulo do módulo
`events` do serviço events
(`workspace_go/services/events/src/modules/events/`), mas o pacote de
testes em si ainda não foi escrito. O que existe aqui hoje é o pacote
compartilhado `asserts/`: oráculos de saga que consultam a API HTTP
pública do serviço events em nome de journeys cross-module (pipelines
IoT, smokes de trigger, verificação de eventos de workflow). Os testes
CRUD de módulo para as três listagens de eventos chegarão ao lado
desses asserts em uma iteração futura.

## Endpoints exercitados

Indiretamente — apenas via saga journeys, através dos oráculos em
`asserts/`:

- `GET /api/v1/events/raw` — lista registros raw de ingestão (filtra
  por `threadId` + `startTime`); usado por
  `AssertRawEventReceivedAfter` para provar que a ingestão pelo
  gateway chegou ao store de eventos.
- `GET /api/v1/events/trigger` — lista execuções de trigger (filtra
  por `triggerId`); usado por `AssertTriggerEventReceivedAfter`,
  `AssertTriggerExecutedSuccessfullyEventually` e
  `AssertLastTriggerRequestDataContains` para verificar entrega do
  trigger e inspecionar o `requestData` resolvido.
- `GET /api/v1/events/workflow` — lista execuções de workflow (filtra
  por `instanceId` + `startTime`); usado por
  `AssertWorkflowEventReceivedAfter` para confirmar que um workflow
  rodou para o asset cuja telemetria o disparou.

## Funções de teste

Nenhuma na camada de e2e de módulo. Os asserts desta pasta são
consumidos pelas saga journeys (rodadas com a build tag `saga`), não
por `go test ./services/events/events`.

## Fixtures

Sem fixtures externas — os asserts montam suas query strings a partir
de valores do bag produzidos por steps de saga anteriores.

## Como rodar

A pasta ainda não tem arquivo de teste, então `go test
./services/events/events` é no-op. Os asserts são exercitados por toda
saga journey que verifica entrega de evento downstream, por exemplo:

```bash
cd e2e_tests

# Exercita AssertTriggerExecutedSuccessfullyEventually +
# AssertLastTriggerRequestDataContains via os smokes de trigger
go test -tags=saga ./journey/automations/... -v

# Exercita AssertRawEventReceivedAfter + AssertWorkflowEventReceivedAfter
go test -tags=saga ./journey/iot/... -v
```

## Resultado em caso de PASS

Quando as saga journeys consumidoras passam, esses asserts em conjunto
provam que o serviço events está acessível em seus paths públicos de
leitura e que os consumers NATS downstream populam corretamente as
tabelas ClickHouse raw, trigger e workflow expostas por esses
endpoints.

## Requisitos

- `events` disponível na porta `5004` (override via `EVENTS_URL`).
- ClickHouse, NATS e os produtores upstream (`http_gateway`,
  `triggers`, `workflow`) acessíveis para que o serviço events tenha
  registros para devolver.
- Stack iniciado pelo `mapexOSDeploy`; o usuário admin do seed
  `admin@mapex.local` e a organização root
  `0000000000000000000aa001` precisam existir.

## Notas

- Os asserts vão deliberadamente apenas pela API HTTP pública — nunca
  lêem ClickHouse, Mongo ou subjects NATS diretamente, para que o
  contrato testado continue sendo o exposto ao usuário.
- Cada assert subtrai 2 segundos de folga do `startTime` aplicado,
  para absorver desvio de relógio entre o runner e o serviço events.
- O pipeline de insert no ClickHouse tem buffer em batch, então a
  primeira observação de um registro pode levar de dez a trinta
  segundos em DEV; os budgets de poll em `assert_trigger_success.go`
  refletem isso.
