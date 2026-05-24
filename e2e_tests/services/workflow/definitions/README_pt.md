# Module e2e: workflow / definitions

## Escopo

Blocos de construção de saga para o módulo `definitions` do serviço
workflow (`workspace_go/services/workflow/src/modules/definitions/`).
Este pacote cobre a metade de definição do ciclo de vida do workflow
— a descrição estática do DAG que o operador desenha no canvas
(nodes, edges, states, retry policy, external inputs e signals).
Combina o step `CreateDefinition` com o payload canônico
`SagaSimpleDefinition` (uma captura literal da UI de um fluxo Start
→ Set State → Code → End) e publica id + version no bag para que o
pacote irmão `instances` consiga vincular uma instância à definição
recém-criada. Não há funções `Test*` aqui: a suite vive em
`journey/iot/connectivity_actions_*/phase1_workflow/` e carrega
estes steps para dirigir o serviço workflow ponta-a-ponta.

## Endpoints exercitados

- `POST /api/v1/workflow_definitions` — registra uma definição de
  workflow. O step envia o payload literal canônico e parametriza
  apenas o `name` com `runID` para que execuções concorrentes da
  saga não colidam no índice único de nome do Mongo.
- `DELETE /api/v1/workflow_definitions/{id}` — chamado pelo
  Compensate de `CreateDefinition`; idempotente contra `404`.

## Fixtures

O pacote não tem fixtures JSON; o corpo do payload é um literal Go
capturado da UI (DevTools → Network) e parametrizado por execução
pelo builder em `payloads/`.

| Arquivo                                 | Propósito                                                                                |
|-----------------------------------------|------------------------------------------------------------------------------------------|
| `payloads/saga_simple_definition.go`    | Definição canônica `Device Status` — Start → Set State → Code → End, um campo de state, um external signal e um plugin instalado (`telegram`); o nome é reescrito para `saga-workflow-def-<runID>`. |

A pasta `steps/` traz `create_definition.go` (o step da saga com seu
par Compensate) e `keys.go` (as constantes de bag key
`BagKeyDefinitionID` e `BagKeyDefinitionVersion` que o pacote
instances lê para montar seu payload de create).

## Como rodar

Este pacote não tem arquivos `*_test.go`; rodar `go test
./services/workflow/definitions` vai imprimir `[no test files]`. O
step é consumido pelas journeys de IoT connectivity-action:

```bash
cd e2e_tests

# Phase 1 das journeys de connectivity-action HTTP e MQTT
go test -tags=saga ./journey/iot/connectivity_actions_http/phase1_workflow
go test -tags=saga ./journey/iot/connectivity_actions_mqtt/phase1_workflow
```

## Resultado em caso de PASS

Quando a journey consumidora passa, este pacote provou que o módulo
workflow definitions respeita `POST /api/v1/workflow_definitions`
ponta-a-ponta para o menor DAG realista que o runtime aceita: a
definição é persistida com `_id` e `definitionVersion`, ambos
retornados na resposta do create, e o `DELETE` do Compensate remove
o registro sem deixar resíduo.

## Requisitos

- `workflow` disponível na porta `5007` (override via
  `WORKFLOW_URL`).
- `mapexos` disponível na porta `5000` para gerar o token de admin.
- Stack iniciado pelo `mapexOSDeploy`; o usuário admin do seed
  `admin@mapex.local` e a organização root
  `0000000000000000000aa001` precisam existir (provisionados pelo
  `mongodb-init`).

## Notas

- O serviço workflow distingue definição (template) de instância
  (vínculo executável). Este pacote cobre só a metade de definição;
  o ciclo de vida da instância está no pacote irmão
  `workflow/instances/`.
- `definitionVersion` é publicado no bag como `int` e o step
  `CreateInstance` faz type-assert de volta — preserve o tipo ao
  trocar o decoder da resposta.
- O payload canônico é uma captura literal da UI de modo que a forma
  da requisição fique em lockstep com o que o operador envia do
  canvas; um erro de parse em runtime é tratado como erro de
  desenvolvedor (o step entra em panic) e não como falha
  recuperável pela saga.
- O payload referencia um plugin instalado (`telegram`); o serviço
  workflow precisa ter o plugin registrado, caso contrário o create
  falha na validação.
