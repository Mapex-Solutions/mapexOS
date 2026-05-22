# Refactor: Eliminar Stream WORKFLOW-TRIGGER

## Problema

Duas NATS streams fazem a mesma coisa com um hop desnecessário:

```
Router → WORKFLOW-EXECUTION → HandleExecution
           → mode=newInstance
              → Carrega instance (buscar definitionId)     ← 1o load da instance
              → PublishTriggerStart → WORKFLOW-TRIGGER
                 → HandleTrigger
                    → Carrega definition
                    → Carrega instance (buscar nomes/config)  ← 2o load da instance (DUPLICADO)
                    → Cria execution
                    → execute()
```

**Problemas:**
1. Round-trip extra NATS publish/consume em **toda** execução
2. Instance carregada **duas vezes** (uma no HandleExecution, outra no HandleTrigger)
3. Duas streams para manter, monitorar e debugar
4. Subworkflow também usa WORKFLOW-TRIGGER desnecessariamente

## Arquitetura Atual

### Streams envolvidas

| Stream | Subject | Produtor | Consumidor | Propósito |
|--------|---------|----------|------------|-----------|
| `WORKFLOW-EXECUTION` | `workflow.execution.>` | Router, HTTP API | `HandleExecution` | Entry point: dispatch por modo |
| `WORKFLOW-TRIGGER` | `workflow.trigger.>` | `HandleExecution` (newInstance), Subworkflow executor | `HandleTrigger` | Cria e executa o workflow |
| `WORKFLOW-RESUME` | `workflow.resume.>` | Callbacks externos, Reconciler | `HandleResume` | Retoma execução pausada |

### Quem publica para WORKFLOW-TRIGGER hoje

1. **`handleNewInstance()`** em `signal.go:79` — publica `workflow.trigger.start`
   - Depois de resolver instanceId → definitionId
2. **`DispatchSubworkflowTrigger()`** em `runtime_publisher.go:129` — publica `workflow.trigger.subworkflow.*`
   - Quando um node `core/subworkflow` suspende

### Modos do HandleExecution (signal.go)

```go
switch execMsg.Mode {
case "newInstance":    → handleNewInstance → PublishTriggerStart → WORKFLOW-TRIGGER
case "signal":        → handleSignalMode → deliverSignal → PublishSignalResume → WORKFLOW-RESUME
case "signalOrStart": → tenta signal, fallback para handleNewInstance
}
```

### HandleTrigger (runtime_service.go:52-179)

Recebe `TriggerMessage` com:
- `WorkflowID` (definitionId)
- `InstanceID`
- `WorkflowUUID`
- `EventPayload`
- `ExternalInputs`
- `Depth`, `ParentInstanceID`, `ParentNodeID`, `CallbackSubject`

Faz:
1. Carrega definition pelo WorkflowID
2. Constrói grafo
3. Encontra node start
4. Carrega instance (para nomes e UUID config) ← DUPLICADO
5. Cria WorkflowExecution
6. Persiste no KV
7. Publica state event "created"
8. execute() DAG walker

## Arquitetura Alvo

Eliminar `WORKFLOW-TRIGGER` por completo. `HandleExecution` absorve toda a lógica de trigger. Subworkflow usa a mesma stream WORKFLOW-EXECUTION.

### Antes (3 streams, 2 hops)

```
Router ──publish──► WORKFLOW-EXECUTION ──consume──► HandleExecution
                                                       │
                                                       ▼
                                              handleNewInstance()
                                                       │
                                                       │ publish (HOP EXTRA)
                                                       ▼
                                              WORKFLOW-TRIGGER ──consume──► HandleTrigger
                                                                              │
                                                                              ▼
                                                                        Cria execution
                                                                        execute() DAG
```

### Depois (2 streams, 0 hops extras)

```
Router ──publish──► WORKFLOW-EXECUTION ──consume──► HandleExecution
                                                       │
                                                       ├─ mode=newInstance
                                                       │    Carrega instance + definition
                                                       │    Cria execution
                                                       │    execute() DAG         ← DIRETO
                                                       │
                                                       ├─ mode=signal
                                                       │    deliverSignal → WORKFLOW-RESUME
                                                       │
                                                       ├─ mode=signalOrStart
                                                       │    Tenta signal, fallback newInstance
                                                       │
                                                       └─ mode=subworkflow (NOVO)
                                                            Cria child execution
                                                            execute() DAG
```

## O que muda

### 1. Absorver HandleTrigger no handleNewInstance

**Arquivos:** `signal.go` + `runtime_service.go`

O `handleNewInstance()` atualmente:
1. Carrega instance → pega definitionId
2. Publica para WORKFLOW-TRIGGER

Muda para:
1. Carrega instance → pega definitionId + nomes + UUID config
2. Carrega definition
3. Constrói grafo, encontra start node
4. Cria WorkflowExecution (com State, ExternalInputs, etc.)
5. Persiste no KV
6. Publica state event "created"
7. execute() DAG walker
8. ACK mensagem

Essencialmente move o corpo do `HandleTrigger` para dentro do `handleNewInstance`, removendo o publish intermediário.

**Benefício:** Instance carregada uma só vez. Zero overhead de serialização/deserialização NATS.

### 2. Subworkflow via WORKFLOW-EXECUTION

**Arquivos:** `runtime_publisher.go` + `lifecycle.go`

Atualmente `DispatchSubworkflowTrigger()` publica `SubworkflowTrigger` para `workflow.trigger.subworkflow.*`.

**Mudança:** Publicar para `WORKFLOW-EXECUTION` com `mode=subworkflow`.

- `HandleExecution` recebe e despacha para novo handler `handleSubworkflow()`
- Campos específicos de subworkflow vão no `data` map: `parentInstanceId`, `parentNodeId`, `callbackSubject`, `depth`, `definitionId`
- Mantém garantias NATS (at-least-once, retry, DLQ) sem o overhead da stream extra
- Subject muda de `workflow.trigger.subworkflow.*` para `workflow.execution.subworkflow.*`

### 3. Remover WORKFLOW-TRIGGER

**Arquivos para deletar:**
- `modules/runtime/interfaces/message/consumers/workflow_trigger/consumer.go`
- `modules/runtime/interfaces/message/consumers/workflow_trigger/constants.go`

**Arquivos para modificar:**
- `consumers.go` — remover `NewWorkflowTriggerConsumer`
- `module.go` — remover registro do consumer
- `runtime.constant.go` — remover TriggerStream/TriggerSubject se definidos
- `runtime_publisher.go` — remover `PublishTriggerStart`, modificar `DispatchSubworkflowTrigger`
- `runtime_publisher_port.go` — remover `PublishTriggerStart` da interface

### 4. Campos do TriggerMessage para WorkflowExecutionMessage

O `TriggerMessage` tem campos extras não presentes no `WorkflowExecutionMessage`:

| Campo | Usado por | Destino |
|-------|-----------|---------|
| `ParentInstanceID` | Subworkflow | `data["parentInstanceId"]` |
| `ParentNodeID` | Subworkflow | `data["parentNodeId"]` |
| `CallbackSubject` | Subworkflow | `data["callbackSubject"]` |
| `Depth` | Subworkflow | `data["depth"]` |
| `ExternalInputs` | newInstance | `data["externalInputs"]` ou campo direto |

Para `newInstance`: o `handleNewInstance` já tem acesso ao event payload e instance — não precisa de TriggerMessage.

Para `subworkflow`: os campos vão no `data` map da `WorkflowExecutionMessage`.

### 5. Remover HandleTrigger

- Deletar método `HandleTrigger` do `RuntimeService`
- A lógica foi absorvida em `handleNewInstance` e `handleSubworkflow`
- Deletar struct `TriggerMessage` de `types.go`

## Checklist de Implementação

### Fase 1: Absorver trigger no execution (modo newInstance)

- [ ] Extrair lógica core do `HandleTrigger` para uma função interna `createAndExecute()`
- [ ] Chamar `createAndExecute()` direto no `handleNewInstance()` — sem publicar para NATS
- [ ] Remover `PublishTriggerStart()` do publisher (port + implementação)
- [ ] Unificar carregamento da instance (carregar uma vez, usar em ambos os contextos)
- [ ] **Testar:** router → WORKFLOW-EXECUTION → newInstance → execução roda

### Fase 2: Subworkflow via WORKFLOW-EXECUTION

- [ ] Adicionar `mode=subworkflow` no switch do `HandleExecution`
- [ ] Criar `handleSubworkflow()` — extrai campos do `data` map e chama `createAndExecute()`
- [ ] Modificar `DispatchSubworkflowTrigger()` para publicar em `WORKFLOW-EXECUTION` com `mode=subworkflow`
- [ ] Atualizar subject de `workflow.trigger.subworkflow.*` para `workflow.execution.subworkflow.*`
- [ ] **Testar:** node subworkflow → WORKFLOW-EXECUTION → child executa → callback para parent

### Fase 3: Remover WORKFLOW-TRIGGER

- [ ] Deletar `workflow_trigger/consumer.go` e `workflow_trigger/constants.go`
- [ ] Remover `NewWorkflowTriggerConsumer` do `consumers.go`
- [ ] Remover registro do consumer no `module.go`
- [ ] Deletar método `HandleTrigger` do `RuntimeService`
- [ ] Deletar struct `TriggerMessage` de `types.go`
- [ ] Remover `PublishTriggerStart` do port + implementação do publisher
- [ ] Limpar referências restantes (grep por "WORKFLOW-TRIGGER", "HandleTrigger", "TriggerMessage")

### Fase 4: Verificar que nada quebrou

- [ ] Workflows existentes executam corretamente via router
- [ ] Node subworkflow funciona end-to-end (parent → child → callback → parent continua)
- [ ] Modo signalOrStart funciona (entrega de signal + fallback newInstance)
- [ ] Modo signal funciona
- [ ] Caminhos de erro funcionam (instance missing, definition missing, workflow disabled)
- [ ] Archiver recebe state events corretos (created, waiting, resumed, completed/failed)
- [ ] Reconciler recupera timeouts
- [ ] Execution viewer mostra dados corretos
- [ ] Loop com body async continua funcionando
- [ ] Fanout com branches async continua funcionando

## Avaliação de Risco

| Risco | Mitigação |
|-------|-----------|
| Quebrar workflows em execução | Dados do KV não mudam. Apenas o fluxo de mensagens muda. Workflows em execução usam RESUME stream (inalterada). |
| Mismatch no callback do subworkflow | O `callbackSubject` não depende da stream — é um subject NATS arbitrário. Continua funcionando. |
| Deletar stream com mensagens pendentes | Drenar WORKFLOW-TRIGGER antes de remover. Verificar zero mensagens pendentes. |
| Quebra na integração com router | Router só publica para WORKFLOW-EXECUTION — inalterado. |
| Outros serviços publicando para WORKFLOW-TRIGGER | Grep confirma que apenas o workflow service publica. Zero produtores externos. |
| Função interna em vez de NATS para newInstance | A mensagem WORKFLOW-EXECUTION já garante at-least-once via NATS. O processamento interno é idempotente (KV create falha se UUID já existe). |

## Resumo de Arquivos

| Arquivo | Ação |
|---------|------|
| `signal.go` | **Major:** absorver lógica do HandleTrigger em handleNewInstance + criar handleSubworkflow |
| `runtime_service.go` | Extrair `createAndExecute()` do HandleTrigger, depois deletar HandleTrigger |
| `runtime_publisher.go` | Remover PublishTriggerStart, modificar DispatchSubworkflowTrigger |
| `runtime_publisher_port.go` | Remover PublishTriggerStart da interface |
| `types.go` | Remover TriggerMessage e SubworkflowTrigger |
| `workflow_trigger/consumer.go` | **Deletar** |
| `workflow_trigger/constants.go` | **Deletar** |
| `consumers.go` | Remover NewWorkflowTriggerConsumer |
| `module.go` | Remover registro do consumer trigger |
| `contracts/.../executions/` | Verificar se WorkflowExecutionMessage suporta campos de subworkflow no data map |

## Resultado Final

### Streams após refactor

| Stream | Subject | Propósito |
|--------|---------|-----------|
| `WORKFLOW-EXECUTION` | `workflow.execution.>` | **Único entry point:** newInstance, signal, signalOrStart, subworkflow |
| `WORKFLOW-RESUME` | `workflow.resume.>` | Retomar execução pausada (callbacks, timers, signals, re-enqueue) |
| `WORKFLOW-STATE` | `workflow.state.>` | Archiver: persistência MongoDB + ClickHouse |
| `WORKFLOW-JS-CODE` | `workflow.js.code` | Dispatch para js-workflow-executor (V8) |
| ~~`WORKFLOW-TRIGGER`~~ | ~~`workflow.trigger.>`~~ | **REMOVIDA** |
