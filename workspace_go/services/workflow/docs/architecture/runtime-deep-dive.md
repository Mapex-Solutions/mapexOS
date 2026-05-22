# Workflow Runtime — Deep Dive da Arquitetura

## Modulos do Servico

O servico de workflow tem 8 modulos, inicializados em ordem fixa via DIG container:

```
1. Definitions    → CRUD para workflow definitions + armazenamento de scripts no MinIO
2. Plugins        → CRUD para manifests de plugins + TieredCache + invalidacao NATS Fanout
3. Credentials    → Armazenamento criptografado de credenciais (Envelope: Master Key → DEK → Data)
4. Engine         → Avaliador de condicoes + resolvedor de valores (logica pura, sem I/O)
5. Instances      → CRUD para workflow instances (templates com configuracao)
6. Runtime        → Motor de execucao DAG — O CORE
7. Archiver       → Consumer WORKFLOW-STATE → MongoDB BulkWrite + eventos ClickHouse
8. Reconciler     → Sweep periodico de timers para timeouts async
```

Cada modulo segue DDD + Hexagonal:
```
module/
  domain/          → entidades, ports, servicos (logica pura)
  application/     → servicos, ports, DTOs (orquestracao)
  infrastructure/  → persistencia, messaging (adaptadores)
  interfaces/      → http/routes + message/consumers (entry points)
```

---

## NATS Consumers (Entry Points)

### Modulo Runtime — 2 consumers

| Consumer | Stream | Subject | Handler | Proposito |
|----------|--------|---------|---------|-----------|
| **Execution** | `WORKFLOW-EXECUTION` | `workflow.execution.>` | `HandleExecution(msg)` | Entry point unico: newInstance, signal, signalOrStart, subworkflow |
| **Resume** | `WORKFLOW-RESUME` | `workflow.resume.>` | `HandleResume(msg)` | Retomar execucao pausada |

### Modulo Archiver — 1 consumer

| Consumer | Stream | Subject | Handler | Proposito |
|----------|--------|---------|---------|-----------|
| **State** | `WORKFLOW-STATE` | `workflow.state.>` | `ProcessStateBatch(msgs)` | Persistir MongoDB + ClickHouse |

### Padrao dos subjects de resume

```
workflow.resume.callback.*    → Callbacks async (code, plugin, subworkflow concluiu)
workflow.resume.signal.*      → Entrega de sinal externo (API HTTP)
workflow.resume.timer.*       → Expiracoes de timer (sweep do Reconciler)
workflow.resume.fasttimer.*   → Fast-path de timer curto (Archiver < 1min)
workflow.resume.reenqueue.*   → Re-enqueue apos MaxInlineSteps excedido
```

---

## Fluxo Core de Execucao

### Fase 1: Trigger → Criar Execucao

```
Mensagem WORKFLOW-EXECUTION chega (do Router ou de subworkflow)
         │
         ▼
HandleExecution(msg)
  │
  ├─ Unmarshal WorkflowExecutionMessage
  │    { mode, orgId, pathKey, event, data: { instanceId, workflowUUID, ... } }
  │
  ├─ Switch por mode:
  │    ├─ "newInstance"    → handleNewInstance()
  │    ├─ "signal"        → handleSignalMode()
  │    ├─ "signalOrStart" → tenta signal, fallback newInstance
  │    └─ "subworkflow"   → handleSubworkflow()
  │
  └─ handleNewInstance():
       │
       ├─ Carregar instance (TieredCache → MongoDB)
       │    Pegar definitionId, nomes, UUID config
       │
       ├─ Carregar definition (TieredCache → MongoDB)
       │    Validar: existe, enabled, status=valid
       │
       ├─ Construir ExecutionGraph da definition
       │    Parsear configs dos nodes, construir lista de adjacencia, achar __start__
       │
       ├─ Criar entidade WorkflowExecution
       │    State = InitializeState(def.States)
       │    ExternalInputs = InitializeExternalInputs(def.ExternalInputs, msg.ExternalInputs)
       │    Status = Running
       │    ActiveNodeIDs = [__start__]
       │
       ├─ Persistir no NATS KV  →  key: "exec.{workflowUUID}"
       │
       ├─ Publicar StateEvent "created"  →  Archiver insere stub no MongoDB
       │
       └─ execute(ctx, execution, graph, __start__)  ←── DAG Walker inicia
```

### Fase 2: Loop do DAG Walker

```
execute(ctx, execution, graph, startNodeID)
  │
  │  FOR step = 0 ATE MaxInlineSteps (300):
  │    │
  │    ├─ executeStep(currentNodeID, graph, nodeContext)
  │    │    │
  │    │    ├─ Buscar node no grafo
  │    │    ├─ Buscar executor no ExecutorRegistry
  │    │    │    core/* → lookup no mapa
  │    │    │    plugin/* → PluginExecutor (rota direta)
  │    │    ├─ Construir NodeExecutionContext
  │    │    ├─ executor.Execute(ctx, execCtx)
  │    │    └─ Retorna (result, pathEntry, error)
  │    │
  │    ├─ SE error → failExecution() → PARA
  │    │
  │    ├─ SE result.NodeState["waitType"] != nil
  │    │    └─ SUSPENDER (ver Fase 3) → PARA
  │    │
  │    ├─ Aplicar NodeState (contador do loop, count do merge)
  │    ├─ Rastrear loop stack (push em "body", pop em "done")
  │    ├─ Aplicar StatePatch no execution.State
  │    ├─ Aplicar NodeOutput no execution.NodeOutputs
  │    │
  │    ├─ SE sem OutputHandles:
  │    │    ├─ Verificar loop stack → voltar ao node loop
  │    │    └─ Senao → completeOrResuspend()
  │    │
  │    ├─ Resolver proximos nodes pelas edges do grafo
  │    │
  │    ├─ SE 1 proximo node:
  │    │    checkpoint() → avancar → CONTINUA
  │    │
  │    └─ SE multiplos proximos nodes:
  │         ├─ core/switch → executeSwitchBranches (sequencial)
  │         └─ outro → executeFanout (goroutines paralelas)
  │
  └─ MaxInlineSteps excedido → checkpoint → re-enqueue via NATS → PARA
```

### Fase 3: Suspensao (Node Async)

Quando um executor retorna `NodeState["waitType"]`:

```
suspendExecution(execution, nodeID, nodeType, nodeState)
  │
  ├─ 1. Setar execution.Status = WAITING
  │     Setar execution.ActiveNodeIDs = [nodeID]
  │     Setar execution.NodeStates[nodeID] = nodeState
  │
  ├─ 2. CHECKPOINT → KV.Put("exec.{UUID}", execution)
  │
  ├─ 3. Publicar StateEvent "waiting" → Archiver
  │
  └─ 4. DISPATCH por tipo do node:
         │
         ├─ core/delay       → SEM dispatch (Reconciler faz sweep)
         ├─ core/wait_signal  → SEM dispatch (endpoint HTTP escuta)
         ├─ core/wait_for     → SEM dispatch (Reconciler avalia)
         │
         ├─ core/code         → DispatchCodeExecution
         │                      → Publica "workflow.js.code" para WORKFLOW-JS-CODE
         │                      → js-workflow-executor executa no V8
         │
         ├─ core/subworkflow  → Publica para WORKFLOW-EXECUTION mode=subworkflow
         │
         └─ plugin/*          → dispatchPluginByActionType
              ├─ http/mqtt/nats/email → DispatchWorkflowTrigger
              │                         → Publica "trigger.WORKFLOW.execute" para Triggers service
              └─ script               → DispatchCodeExecution
                                        → Publica para js-workflow-executor
```

### Fase 4: Resume (Callback Chega)

```
Mensagem WORKFLOW-RESUME chega
         │
         ▼
HandleResume(msg)
  │
  ├─ Unmarshal ResumeMessage
  │    { instanceId, nodeId, status, outputHandle, output, statePatch, error }
  │
  ├─ GET execution do KV → validar status WAITING ou RUNNING
  │
  ├─ SE resume.IsTimeout:
  │    ├─ enableOutput=true → rotear para handle "timeout"
  │    └─ enableOutput=false → failExecution(TIMEOUT_EXCEEDED)
  │
  ├─ SE resume.Error != nil:
  │    ├─ Enriquecer erro (nodeId, nodeType, timestamp)
  │    └─ failExecution(resume.Error)
  │
  ├─ Aplicar dados do resume:
  │    ├─ StatePatch → merge no execution.State
  │    ├─ Output → setar execution.NodeOutputs[nodeId]
  │    ├─ Remover nodeId dos ActiveNodeIDs
  │    ├─ Deletar NodeStates[nodeId] (limpar estado de espera)
  │    ├─ Setar Status = RUNNING
  │    └─ Atualizar ExecutionPath (waiting → completed)
  │
  ├─ CHECKPOINT → KV.Put
  │
  ├─ Publicar StateEvent "resumed" → Archiver limpa timer
  │
  ├─ Resolver proximo node do node retomado + outputHandle
  │    SE nenhum proximo:
  │      ├─ Verificar loop stack → voltar ao loop
  │      └─ Senao → completeOrResuspend()
  │
  └─ execute(ctx, execution, graph, nextNodeID)  ←── DAG Walker retoma
```

---

## Registro de Executors

18 executors registrados na inicializacao do servico:

### Inline (sincronos, sem I/O)
| Tipo | Executor | Comportamento |
|------|----------|---------------|
| `core/start` | StartExecutor | Passthrough → `["out"]` |
| `core/end` | EndExecutor | Terminal → `[]` ou ExecutionError |
| `core/condition` | ConditionExecutor | Avalia → `["true"]` ou `["false"]` |
| `core/switch` | SwitchExecutor | Multi-caso → `["case_X"]` ou `["default"]` |
| `core/set_state` | SetStateExecutor | Modifica estado → `["out"]` |
| `core/log` | LogExecutor | Cria log entry → `["out"]` |
| `core/goto` | GotoExecutor | Portal → `["out"]` (grafo roteia ao receiver) |

### Async (suspendem com waitType, retomam via callback)
| Tipo | Executor | waitType | Destino do dispatch |
|------|----------|----------|---------------------|
| `core/code` | CodeExecutor | `callback` | js-workflow-executor (V8) |
| `core/delay` | DelayExecutor | `timer` | Nenhum (Reconciler faz sweep) |
| `core/wait_signal` | WaitSignalExecutor | `signal` | Nenhum (endpoint HTTP) |
| `core/subworkflow` | SubworkflowExecutor | `callback` | WORKFLOW-EXECUTION (child) |
| `core/trigger_event` | TriggerEventExecutor | `callback` | Triggers service |

### Controle de fluxo (gerenciam estrutura de execucao)
| Tipo | Executor | Comportamento |
|------|----------|---------------|
| `core/fanout` | FanoutExecutor | Fork N branches → `["out_1", ..., "out_N"]` |
| `core/merge` | MergeExecutor | Join branches → `["out"]` quando count atingido |
| `core/sequence` | SequenceExecutor | Passo a passo → `["step_N"]` ou `["done"]` |
| `core/loop` | LoopExecutor | Iterar array → `["body"]` ou `["done"]` |
| `core/wait_for` | WaitForExecutor | Condicao → `["matched"]` ou suspende |

### Plugin (marketplace generico)
| Tipo | Executor | Comportamento |
|------|----------|---------------|
| `*/*` (non-core) | PluginExecutor | Resolver manifest + credentials + templates → suspende |

---

## Persistencia de Estado — 3 Camadas

```
                    ┌─────────────┐
                    │   Runtime    │
                    │  (execute)   │
                    └──────┬──────┘
                           │ checkpoint() apos cada step
                           ▼
                    ┌─────────────┐
                    │   NATS KV   │  ◄── ESTADO QUENTE
                    │  "exec.{id}"│      JSON completo da execution
                    └──────┬──────┘      Escrito: cada checkpoint
                           │             Deletado: no estado terminal
                           │ Publica StateEvent
                           ▼
                    ┌─────────────┐
                    │  Archiver   │  ◄── ESTADO MORNO
                    │  (consumer) │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌───────────┐
        │ MongoDB  │ │ MongoDB  │ │ClickHouse │  ◄── ESTADO FRIO
        │  Stub    │ │  Full    │ │  Evento   │
        │ ~200B    │ │ ~5-25KB  │ │  Analytics│
        │ created  │ │ terminal │ │  terminal │
        └──────────┘ └──────────┘ └───────────┘
```

### Timeline de escritas:

| Evento | KV | MongoDB | ClickHouse |
|--------|------|---------|------------|
| **created** | Criar execution completa | Inserir stub leve (~200B) | — |
| **cada step** | Atualizar (checkpoint) | — | — |
| **waiting** | Atualizar (estado de suspensao) | Atualizar status + timers | — |
| **resumed** | Atualizar (limpar espera) | Limpar timers | — |
| **completed/failed** | **Deletar** | Upsert execution completa (~5-25KB) + TTL 3 dias | Inserir evento analytics |

### Ciclo do Archiver (batch):

```
ProcessStateBatch(messages):
  │
  ├─ BATCH 1: Stubs leves
  │    Subject: workflow.state.created
  │    Operacao: BulkInsertLightweight() — stub ~200B
  │
  ├─ BATCH 2: Updates de waiting + fast-path de timers curtos
  │    Subject: workflow.state.waiting
  │    Operacao: BulkUpdateWaiting() — status + timers
  │    Se timer < 1min: publica resume direto (fast-path)
  │
  ├─ BATCH 3: Updates de resumed
  │    Subject: workflow.state.resumed
  │    Operacao: BulkUpdateResumed() — limpa timers
  │
  └─ BATCH 4: Upserts terminais completos
       Subject: workflow.state.completed|failed|cancelled
       Sequencia:
         a) KV.Get("exec.{UUID}") → execution completa
         b) BulkUpsertFull() → MongoDB com TTL
         c) Publicar evento → ClickHouse para analytics
         d) KV.Delete("exec.{UUID}") → limpar estado quente
```

---

## Execucao de Fanout

```
executeFanout(branchStartNodes, mode)
  │
  ├─ Deep copy do estado para isolamento (cada branch recebe copia propria)
  │
  ├─ Spawn N goroutines → executeBranch() por branch
  │    Cada branch caminha independente ate:
  │    ├─ core/merge encontrado → branch termina, reporta mergeNodeID
  │    ├─ Node async → branch suspende, reporta BranchWaiting
  │    ├─ Terminal → branch completa
  │    └─ Erro → branch falha
  │
  ├─ WaitAll goroutines
  │
  └─ Processar resultados:
       │
       ├─ ALGUMA falhou → failExecution()
       │
       ├─ ALGUMA waiting:
       │    ├─ mode=waitAll: guardar todos waiting nodeIDs em ActiveNodeIDs
       │    │    Checkpoint + suspendFanoutExecution (dispatch cada waiting node)
       │    └─ mode=firstCompleted: guardar __fanout_meta
       │         Quando primeira branch completar → completeOrResuspend cancela as outras
       │
       └─ TODAS completaram:
            Aplicar state patches + outputs de todas as branches
            Setar branchCount no merge
            Retornar mergeNodeID → walker continua do merge
```

---

## Loop Stack (Suporte a Body Async)

```
NodeStates["__loop_stack"]["stack"] = ["loop_1", "loop_2"]  (LIFO)

PUSH: quando loop emite "body" → push loopNodeId
POP:  quando loop emite "done" → pop
CHECK: quando DAG walker chega a terminal:
       SE stack nao vazio → pop → currentNodeID = loop → continua
       SENAO → completeOrResuspend()

Persiste no KV entre suspensoes async:
  Loop iter 0 → body → code (suspende) → resume → end → CHECK → pop → Loop iter 1 → ...

Suporta loops aninhados:
  Loop externo → body → Loop interno → body → ... → done(interno) → pop → done(externo) → pop
```

---

## Logica do completeOrResuspend

Chamado quando uma branch/walker chega a um ponto terminal:

```
completeOrResuspend(execution):
  │
  ├─ Verificar __fanout_meta (modo firstCompleted):
  │    SE "firstCompleted" → cancelar waiting nodes restantes → completeExecution()
  │
  ├─ Verificar ActiveNodeIDs por nodes ainda esperando:
  │    SE algum node tem waitType no NodeStates:
  │      Re-suspender execution (Status=WAITING, checkpoint, publicar)
  │    SENAO:
  │      completeExecution()
  │
  └─ completeExecution():
       Status = COMPLETED
       CompletedAt = now
       ActiveNodeIDs = nil
       checkpoint → KV.Put
       Publicar StateEvent "completed" → Archiver arquiva + KV deleta
```

---

## Reconciler (Timer em Background)

```
A cada 30 segundos:
  │
  ├─ Query MongoDB: timerExpiresAt <= now AND status = "waiting"
  │
  ├─ PARA CADA execution expirada:
  │    ├─ GET do KV (ainda viva?)
  │    ├─ Publicar ResumeMessage com isTimeout=true
  │    │    Subject: workflow.resume.timer.{UUID}
  │    └─ HandleResume processa:
  │         ├─ enableOutput=true → rotear para handle "timeout"
  │         └─ enableOutput=false → failExecution(TIMEOUT_EXCEEDED)
  │
  └─ Fast-path: Archiver detecta timers < 1 minuto
       Publica direto para workflow.resume.fasttimer.{UUID}
       (evita delay de 30s do sweep para timers curtos)
```
