# Module e2e: workflow / instances

## Escopo

Blocos de construção de saga para o módulo `instances` do serviço
workflow (`workspace_go/services/workflow/src/modules/instances/`).
Este pacote cobre a metade de instância do ciclo de vida do workflow
— o vínculo executável que combina uma definição já existente (por
id + version) com inputs externos concretos e a política de execução
escolhida pelo operador. Combina o step `CreateInstance` com o
payload canônico `SagaSimpleInstance` (uma captura literal da UI) e
injeta `definitionId`, `definitionVersion` e `definitionName` a
partir das bag keys publicadas pelo pacote irmão
`workflow/definitions`, e depois publica o id da instância
resultante para que route groups de `kind=workflow` possam
referenciar. Não há funções `Test*` aqui: a suite vive em
`journey/iot/connectivity_actions_*/phase1_workflow/` e carrega
estes steps para dirigir o serviço workflow ponta-a-ponta.

## Endpoints exercitados

- `POST /api/v1/workflow_instances` — registra uma instância de
  workflow vinculada a uma definição existente. O step envia o
  payload literal canônico com `definitionId`, `definitionVersion`,
  `definitionName` e `name` sobrescritos por execução.
- `DELETE /api/v1/workflow_instances/{id}` — chamado pelo Compensate
  de `CreateInstance`; idempotente contra `404`.

## Fixtures

O pacote não tem fixtures JSON; o corpo do payload é um literal Go
capturado da UI e parametrizado por execução pelo builder em
`payloads/`.

| Arquivo                              | Propósito                                                                                |
|--------------------------------------|------------------------------------------------------------------------------------------|
| `payloads/saga_simple_instance.go`   | Instância canônica `Device Status` — `externalInputs` vazio, `isTemplate=false`, `uniqueExecution=false`, sem `pathKey` nem `workflowUUID`; `name` é reescrito para `saga-workflow-inst-<runID>` e `definitionId` / `definitionVersion` / `definitionName` são sobrescritos por `CreateInstance` a partir do bag. |

A pasta `steps/` traz `create_instance.go` (o step da saga com seu
par Compensate) e `keys.go` (a constante de bag key
`BagKeyInstanceID` que o payload de route group `kind=workflow`
lê).

## Como rodar

Este pacote não tem arquivos `*_test.go`; rodar `go test
./services/workflow/instances` vai imprimir `[no test files]`. O
step é consumido pelas journeys de IoT connectivity-action:

```bash
cd e2e_tests

# Phase 1 das journeys de connectivity-action HTTP e MQTT
go test -tags=saga ./journey/iot/connectivity_actions_http/phase1_workflow
go test -tags=saga ./journey/iot/connectivity_actions_mqtt/phase1_workflow
```

## Resultado em caso de PASS

Quando a journey consumidora passa, este pacote provou que o módulo
workflow instances respeita `POST /api/v1/workflow_instances`
ponta-a-ponta: uma instância vinculada a uma definição + version
reais é persistida com `_id` retornado na resposta do create, a
saga consegue publicar esse id para um route group downstream de
kind=workflow referenciar, e o `DELETE` do Compensate remove o
registro sem deixar resíduo.

## Requisitos

- `workflow` disponível na porta `5007` (override via
  `WORKFLOW_URL`).
- `mapexos` disponível na porta `5000` para gerar o token de admin.
- Uma definição de workflow precisa estar no bag — o step
  `CreateInstance` falha imediatamente se `BagKeyDefinitionID` ou
  `BagKeyDefinitionVersion` estiver ausente, então as journeys
  precisam rodar `definitions.CreateDefinition` antes.
- Stack iniciado pelo `mapexOSDeploy`; o usuário admin do seed
  `admin@mapex.local` e a organização root
  `0000000000000000000aa001` precisam existir (provisionados pelo
  `mongodb-init`).

## Notas

- O serviço workflow distingue definição (template) de instância
  (vínculo executável). Este pacote cobre só a metade de instância;
  o ciclo de vida da definição está no pacote irmão
  `workflow/definitions/`.
- `definitionName` é mantido alinhado com o nome derivado de runID
  que o pacote definitions escreve, de modo que o operador
  listando instâncias de uma execução de saga veja um rótulo
  consistente no par definição + instância.
- O payload canônico é uma captura literal da UI; um erro de parse
  em runtime é tratado como erro de desenvolvedor (o step entra em
  panic) e não como falha recuperável pela saga.
- Os valores placeholder em `sagaSimpleInstanceJSON`
  (`definitionId`, `definitionVersion`, `definitionName`) são
  sempre sobrescritos em runtime — existem só para manter o
  literal como um POST body válido em isolamento.
