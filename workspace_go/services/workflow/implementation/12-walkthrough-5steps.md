# Walkthrough Interno: 5 Nodes — O que REALMENTE acontece no Go

Rastreio EXATO do que acontece internamente, com referencias ao codigo real.

---

## O Workflow

```
[1:start] → [2:set_state] → [3:condition] ──true──→ [4:log] → [5:end]
                                            └─false─→ [5:end]
```

Nodes:
```
node-1  core/start       config: nil
node-2  core/set_state   config: {operation: "set", targetField: "discount_rate", valueSource: {type: "literal", value: "10"}}
node-3  core/condition   config: {condition: {logic: "AND", items: [{field: {type: "state", value: "discount_rate"}, operator: "greater_than", value: {type: "literal", value: "5"}}]}}
node-4  core/log         config: {message: "Discount: ${state.discount_rate}%", level: "info"}
node-5  core/end         config: {terminateWithError: false}
```

Edges:
```
node-1 "out"   → node-2
node-2 "out"   → node-3
node-3 "true"  → node-4
node-3 "false" → node-5
node-4 "out"   → node-5
```

Variables: `[{field: "discount_rate", type: "number", defaultValue: 0}]`
Event:     `{"orderId": "ORD-001", "total": 200}`

---

## FASE 0: TRIGGER — HandleTrigger()

O Trigger Service publica no NATS subject `workflow.trigger.{workflowId}`.
O consumer WORKFLOW-TRIGGER recebe e chama `HandleTrigger(msg)`.

### runtime_service.go:50-158

```
Linha 54: json.Unmarshal(msg.Data, &trigger)
  trigger = {
    WorkflowID:      "wf-abc",
    OrgID:           "org-123",
    EventPayload:    {"orderId": "ORD-001", "total": 200},
    Depth:           0,
    ParentInstanceID: "",
    CallbackSubject:  "",
    ParentNodeID:     "",
  }

Linha 62: def = s.deps.DefinitionRepo.FindById(ctx, "wf-abc")
  → Busca WorkflowDefinition no MongoDB (ou TieredCache)
  → def.Enabled = true ✓
  → def.Variables = [{field: "discount_rate", type: "number", defaultValue: 0}]
  → def.Nodes = [node-1..node-5]
  → def.Edges = [5 edges]

Linha 73: graph = domainServices.BuildGraph(def)
  graph.Adjacency = {
    "node-1": {"out":   "node-2"},
    "node-2": {"out":   "node-3"},
    "node-3": {"true":  "node-4", "false": "node-5"},
    "node-4": {"out":   "node-5"},
  }
  graph.ParsedConfigs = {
    "node-1": nil,
    "node-2": &SetStateNodeConfig{Operation: "set", TargetField: "discount_rate", ValueSource: {Type: "literal", Value: "10"}},
    "node-3": &ConditionNodeConfig{Condition: {Logic: "AND", Items: [...]}},
    "node-4": &LogNodeConfig{Message: "Discount: ${state.discount_rate}%", Level: "info"},
    "node-5": &EndNodeConfig{TerminateWithError: false},
  }

Linha 77: startNodeID = "node-1"  (loop nos nodes procurando type == "core/start")

Linha 90: Cria WorkflowInstance:
  instance = {
    ID:            ObjectID("inst-001"),
    WorkflowID:    ObjectID("wf-abc"),
    WorkflowName:  "Discount Workflow",
    OrgID:         ObjectID("org-123"),
    Version:       1,
    Status:        "running",
    ActiveNodeIDs: ["node-1"],
    NodeStates:    {},
    State:         {"discount_rate": 0},          ← initializeState(variables)
    StateDefaults: {"discount_rate": 0},          ← defaults para code executor context
    EventPayload:  {"orderId": "ORD-001", "total": 200},
    ExecutionPath: [],
    NodeOutputs:   {},
    Depth:         0,
  }

Linha 133: kvKey = "inst:inst-001"
Linha 134: instanceData = json.Marshal(instance)       ← serializa ~1KB
Linha 139: s.deps.KVStore.Create("inst:inst-001", instanceData)
  → NATS KV Create — primeira escrita no disco. Se crashar aqui, instancia existe no KV.

Linha 145: s.publishStateEvent(instance, "workflow.state.created")
  → Publica no WORKFLOW-STATE stream
  → Archiver vai fazer InsertOne leve no MongoDB (~200 bytes)

Linha 148: s.executeInline(ctx, instance, graph, "node-1")
  → ENTRA NO LOOP DO DAG WALKER ↓
```

---

## executeInline — O LOOP

### runtime_service.go:275-500

```go
func (s *RuntimeService) executeInline(ctx, instance, graph, startNodeID) error {
    kvKey = "inst:inst-001"
    currentNodeID = "node-1"

    for step := 0; step < MaxInlineSteps; step++ {   // MaxInlineSteps = 500
        // ... corpo do loop (ver iteracoes abaixo)
    }
}
```

O loop roda ate:
- OutputHandles vazio (terminal)
- NodeState com waitType presente (async)
- Erro
- MaxInlineSteps atingido

---

## ITERACAO 0: start (node-1)

### Entrada do loop

```
step = 0
currentNodeID = "node-1"
```

### Linhas 285-286: Busca node no grafo

```
node = graph.GetNode("node-1")
  → node = {ID: "node-1", Type: "core/start", Label: "Start", Config: nil}
```

### Linhas 301-302: Busca executor no registry

```
executor = s.registry.Get("core/start")
  → executor = &StartExecutor{}
```

### Linhas 318-330: Monta NodeExecutionContext

```go
execCtx = &NodeExecutionContext{
    InstanceID:   ObjectID("inst-001"),
    State:        {"discount_rate": 0},
    EventPayload: {"orderId": "ORD-001", "total": 200},
    NodeOutputs:  {},
    NodeStates:   {},
    Depth:        0,
    NodeID:       "node-1",
    NodeType:     "core/start",
    ParsedConfig: nil,
    Label:        "Start",
    Graph:        graph,
    Timezone:     "",
}
```

### Linhas 332-334: Executa

```
enteredAt = time.Now()                    // T+0.000ms
result, err = executor.Execute(ctx, execCtx)
exitedAt = time.Now()                     // T+0.000ms
```

**Dentro do StartExecutor.Execute():**
```go
func (e *StartExecutor) Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
    return &NodeExecutionResult{
        OutputHandles: []string{"out"},
    }, nil
}
```
Passthrough. Retorna `["out"]` e nada mais.

### Linhas 336-343: Cria PathEntry

```
pathEntry = {
    NodeID:     "node-1",
    NodeType:   "core/start",
    Status:     "completed",
    EnteredAt:  T+0.000ms,
    ExitedAt:   T+0.000ms,
    DurationMs: 0,
}
```

### Linha 345: err == nil? SIM → pula bloco de erro

### Linha 367: OutputHandles[0] = "out" → pathEntry.OutputHandle = "out"

### Linha 371: result.Error == nil → pula bloco de erro

### Linha 385: result.NodeState nil (sem waitType) → pula bloco de async

### Linha 406: ExecutionPath append
```
instance.ExecutionPath = [{nodeId: "node-1", status: "completed", outputHandle: "out", durationMs: 0}]
```

### Linha 408: result.StatePatch == nil → pula

### Linha 414: result.NodeOutput == nil → pula

### Linha 418: len(OutputHandles) == 0? NAO (tem "out") → pula terminal

### Linha 431: Resolve proximo node
```
nextNodes = graph.ResolveNextNodes("node-1", ["out"])
  → adjacency["node-1"]["out"] = "node-2"
  → nextNodes = ["node-2"]
```

### Linha 433: len(nextNodes) == 0? NAO → pula completed

### Linha 448: len(nextNodes) == 1? SIM → caminho linear
```
currentNodeID = "node-2"
instance.ActiveNodeIDs = ["node-2"]
instance.Updated = time.Now()
```

### Linha 452: KV CHECKPOINT
```
s.kvCheckpoint("inst:inst-001", instance)
  → json.Marshal(instance) → KV Put
  → Se crashar AGORA: restart do node-2 (node-1 NAO re-executa)
```

### Linha 455: `continue` → volta ao topo do loop

### Estado da instancia apos iteracao 0

```
Status:        "running"
ActiveNodeIDs: ["node-2"]
State:         {"discount_rate": 0}       ← sem mudanca
NodeOutputs:   {}
ExecutionPath: [
  {nodeId: "node-1", status: "completed", outputHandle: "out"}
]
```

---

## ITERACAO 1: set_state (node-2)

### Entrada do loop

```
step = 1
currentNodeID = "node-2"
```

### Linhas 285-302: Busca node + executor

```
node = {ID: "node-2", Type: "core/set_state"}
executor = &SetStateExecutor{resolver: valueResolver}
```

### Linhas 318-330: Monta execCtx

```go
execCtx = &NodeExecutionContext{
    State:        {"discount_rate": 0},      ← valor ATUAL
    EventPayload: {"orderId": "ORD-001", "total": 200},
    NodeOutputs:  {},
    NodeID:       "node-2",
    NodeType:     "core/set_state",
    ParsedConfig: &SetStateNodeConfig{
        Operation:   "set",
        TargetField: "discount_rate",
        ValueSource: FieldValue{Type: "literal", Value: "10"},
    },
}
```

### Linhas 332-334: Executa

**Dentro do SetStateExecutor.Execute():**
```go
func (e *SetStateExecutor) Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
    cfg := execCtx.ParsedConfig.(*SetStateNodeConfig)   // type assertion
    if cfg == nil {
        return &NodeExecutionResult{OutputHandles: []string{"out"}}, nil
    }

    // 1. Resolve o valor via ValueResolver
    value, err := e.resolver.Resolve(
        &cfg.ValueSource,            // FieldValue{Type: "literal", Value: "10"}
        execCtx.State,               // {"discount_rate": 0}
        execCtx.EventPayload,        // {"orderId": "ORD-001", "total": 200}
        execCtx.NodeOutputs,         // {}
    )
    // value = float64(10)   ← literal "10" convertido para numero

    // 2. Aplica operacao
    // cfg.Operation == "set"
    patch := map[string]interface{}{
        cfg.TargetField: value,      // "discount_rate": float64(10)
    }

    return &NodeExecutionResult{
        OutputHandles: []string{"out"},
        StatePatch:    patch,           // {"discount_rate": 10}
    }, nil
}
```

**result:**
```
OutputHandles: ["out"]
StatePatch:    {"discount_rate": float64(10)}
NodeState:     nil
NodeOutput:    nil
LogEntries:    nil
Error:         nil
NodeState:     nil
```

### Linha 367: pathEntry.OutputHandle = "out"

### Linhas 371, 385: Error nil, NodeState nil → pula

### Linha 406: ExecutionPath append

### Linha 408-412: APLICA StatePatch ← AQUI O STATE MUDA
```go
if result.StatePatch != nil {
    for k, v := range result.StatePatch {
        instance.State[k] = v
    }
}
// instance.State["discount_rate"] = float64(10)
// State AGORA: {"discount_rate": 10}
```

### Linha 418: OutputHandles nao vazio → pula

### Linha 431: ResolveNextNodes
```
nextNodes = graph.ResolveNextNodes("node-2", ["out"])
  → adjacency["node-2"]["out"] = "node-3"
  → nextNodes = ["node-3"]
```

### Linha 448-455: Avanca + KV checkpoint
```
currentNodeID = "node-3"
instance.ActiveNodeIDs = ["node-3"]
KV Put("inst:inst-001", instance)    ← State JA TEM discount_rate=10
continue
```

### Estado apos iteracao 1

```
Status:        "running"
ActiveNodeIDs: ["node-3"]
State:         {"discount_rate": 10}     ← MUDOU! era 0, agora 10
NodeOutputs:   {}
ExecutionPath: [
  {nodeId: "node-1", status: "completed", outputHandle: "out"},
  {nodeId: "node-2", status: "completed", outputHandle: "out"},
]
```

---

## ITERACAO 2: condition (node-3)

### Entrada do loop

```
step = 2
currentNodeID = "node-3"
```

### Linhas 285-302: Busca node + executor

```
node = {ID: "node-3", Type: "core/condition"}
executor = &ConditionExecutor{evaluator: conditionEvaluator}
```

### Linhas 318-330: Monta execCtx

```go
execCtx = &NodeExecutionContext{
    State:        {"discount_rate": 10},     ← valor ATUALIZADO pelo step anterior
    EventPayload: {"orderId": "ORD-001", "total": 200},
    NodeOutputs:  {},
    NodeID:       "node-3",
    NodeType:     "core/condition",
    ParsedConfig: &ConditionNodeConfig{
        Condition: ConditionGroup{
            Logic: "AND",
            Items: [{
                Type: "condition",
                Data: ConditionItem{
                    Field:    FieldValue{Type: "state", Value: "discount_rate"},
                    Operator: "greater_than",
                    Value:    FieldValue{Type: "literal", Value: "5"},
                },
            }],
        },
    },
}
```

### Linhas 332-334: Executa

**Dentro do ConditionExecutor.Execute():**
```go
func (e *ConditionExecutor) Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
    cfg := execCtx.ParsedConfig.(*ConditionNodeConfig)
    if cfg == nil {
        return &NodeExecutionResult{OutputHandles: []string{"false"}}, nil
    }

    // 1. Avalia grupo de condicoes
    result, err := e.evaluator.EvaluateGroup(
        &cfg.Condition,              // ConditionGroup{Logic: "AND", Items: [...]}
        execCtx.Timezone,            // ""
        execCtx.EventPayload,        // {"orderId": "ORD-001", "total": 200}
        execCtx.State,               // {"discount_rate": 10}
        execCtx.NodeOutputs,         // {}
    )

    // Dentro do EvaluateGroup:
    //   Item 0: field = state["discount_rate"] = 10
    //           operator = "greater_than"
    //           value = literal "5" = 5
    //           10 > 5 → TRUE
    //   Logic AND, 1 item, todos TRUE → grupo TRUE

    // result = true

    // 2. Retorna handle baseado no resultado
    if result {
        return &NodeExecutionResult{OutputHandles: []string{"true"}}, nil
    }
    return &NodeExecutionResult{OutputHandles: []string{"false"}}, nil
}
```

**result:**
```
OutputHandles: ["true"]         ← 10 > 5 = verdadeiro
StatePatch:    nil
tudo mais nil
```

### Linha 367: pathEntry.OutputHandle = "true"

### Linha 408: StatePatch nil → pula (State NAO muda)

### Linha 431: ResolveNextNodes ← AQUI DECIDE O CAMINHO
```
nextNodes = graph.ResolveNextNodes("node-3", ["true"])
  → adjacency["node-3"]["true"] = "node-4"    ✓ MATCH
  → adjacency["node-3"]["false"] = "node-5"   ✗ nao usado
  → nextNodes = ["node-4"]
```

**Se discount_rate fosse 3 (menor que 5):**
```
OutputHandles seria ["false"]
nextNodes = graph.ResolveNextNodes("node-3", ["false"])
  → adjacency["node-3"]["false"] = "node-5"
  → Pularia direto pro end, nodes 4 nunca executaria
```

### Linha 448-455: Avanca + KV checkpoint
```
currentNodeID = "node-4"
instance.ActiveNodeIDs = ["node-4"]
KV Put
continue
```

### Estado apos iteracao 2

```
Status:        "running"
ActiveNodeIDs: ["node-4"]
State:         {"discount_rate": 10}     ← sem mudanca
NodeOutputs:   {}
ExecutionPath: [
  {nodeId: "node-1", outputHandle: "out"},
  {nodeId: "node-2", outputHandle: "out"},
  {nodeId: "node-3", outputHandle: "true"},   ← registra QUAL handle saiu
]
```

---

## ITERACAO 3: log (node-4)

### Entrada do loop

```
step = 3
currentNodeID = "node-4"
```

### Linhas 285-302: Busca node + executor

```
node = {ID: "node-4", Type: "core/log"}
executor = &LogExecutor{}
```

### Linhas 318-330: Monta execCtx

```go
execCtx = &NodeExecutionContext{
    State:        {"discount_rate": 10},
    EventPayload: {"orderId": "ORD-001", "total": 200},
    NodeID:       "node-4",
    NodeType:     "core/log",
    ParsedConfig: &LogNodeConfig{
        Message: "Discount: ${state.discount_rate}%",
        Level:   "info",
    },
}
```

### Linhas 332-334: Executa

**Dentro do LogExecutor.Execute():**
```go
func (e *LogExecutor) Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
    cfg := execCtx.ParsedConfig.(*LogNodeConfig)
    if cfg == nil {
        return &NodeExecutionResult{OutputHandles: []string{"out"}}, nil
    }

    // 1. Interpola a mensagem
    message := interpolateMessage(cfg.Message, execCtx.State, execCtx.EventPayload)

    // interpolateMessage:
    //   "Discount: ${state.discount_rate}%"
    //   → encontra ${state.discount_rate}
    //   → busca execCtx.State["discount_rate"] = 10
    //   → substitui: "Discount: 10%"

    // 2. Resolve level
    level := parseLogLevel(cfg.Level)   // "info" → LogInfo

    // 3. Cria LogEntry
    return &NodeExecutionResult{
        OutputHandles: []string{"out"},
        LogEntries: []LogEntry{{
            Level:     LogInfo,
            Message:   "Discount: 10%",
            Timestamp: time.Now(),
            NodeID:    "node-4",
            NodeType:  "core/log",
        }},
    }, nil
}
```

**result:**
```
OutputHandles: ["out"]
StatePatch:    nil
NodeOutput:    nil
LogEntries:    [{level: "info", message: "Discount: 10%", nodeId: "node-4"}]
Error:         nil
NodeState:     nil
```

### Linhas 406-416: Runtime aplica resultado

```
Linha 406: ExecutionPath append

Linha 408: StatePatch nil → pula (State NAO muda)

Linha 414: NodeOutput nil → pula
```

**O que acontece com LogEntries?**
Na versao ATUAL do codigo, LogEntries sao publicadas DEPOIS do loop (ou inline se houver hook).
Na versao NOVA (doc 09-executors.md), runtime publica no WORKFLOW-LOGS stream:
```
s.publishLogEntries(instance.ID, result.LogEntries)
  → subject: workflow.logs.inst-001
  → payload: {instanceId: "inst-001", entries: [{level: "info", message: "Discount: 10%", ...}]}
  → Consumer ClickHouse persiste para consulta futura
  → Fire-and-forget (nao bloqueia execucao)
```

### Linha 431: ResolveNextNodes
```
nextNodes = graph.ResolveNextNodes("node-4", ["out"])
  → adjacency["node-4"]["out"] = "node-5"
  → nextNodes = ["node-5"]
```

### Linha 448-455: Avanca + KV checkpoint
```
currentNodeID = "node-5"
instance.ActiveNodeIDs = ["node-5"]
KV Put
continue
```

### Estado apos iteracao 3

```
Status:        "running"
ActiveNodeIDs: ["node-5"]
State:         {"discount_rate": 10}     ← sem mudanca (log nao modifica state)
NodeOutputs:   {}
ExecutionPath: [
  {nodeId: "node-1", outputHandle: "out"},
  {nodeId: "node-2", outputHandle: "out"},
  {nodeId: "node-3", outputHandle: "true"},
  {nodeId: "node-4", outputHandle: "out"},
]
```

---

## ITERACAO 4: end (node-5) — TERMINAL

### Entrada do loop

```
step = 4
currentNodeID = "node-5"
```

### Linhas 285-302: Busca node + executor

```
node = {ID: "node-5", Type: "core/end"}
executor = &EndExecutor{resolver: valueResolver}
```

### Linhas 318-330: Monta execCtx

```go
execCtx = &NodeExecutionContext{
    State:        {"discount_rate": 10},
    EventPayload: {"orderId": "ORD-001", "total": 200},
    NodeOutputs:  {},
    NodeID:       "node-5",
    NodeType:     "core/end",
    ParsedConfig: &EndNodeConfig{
        TerminateWithError: false,
    },
}
```

### Linhas 332-334: Executa

**Dentro do EndExecutor.Execute():**
```go
func (e *EndExecutor) Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
    cfg := execCtx.ParsedConfig.(*EndNodeConfig)
    if cfg == nil {
        return &NodeExecutionResult{OutputHandles: []string{}}, nil
    }

    // terminateWithError == false → termina normalmente
    if !cfg.TerminateWithError {
        return &NodeExecutionResult{
            OutputHandles: []string{},     // ← VAZIO! Sinaliza terminal
        }, nil
    }

    // Se terminateWithError == true (NAO e o caso aqui):
    // errorMsg = resolver.Resolve(&cfg.ErrorMessage, ...)
    // return &NodeExecutionResult{
    //     OutputHandles: []string{},
    //     Error: &ExecutionError{Code: cfg.ErrorCode, Message: errorMsg},
    // }, nil
}
```

**result:**
```
OutputHandles: []          ← VAZIO!
StatePatch:    nil
NodeOutput:    nil
LogEntries:    nil
Error:         nil         ← sem erro
NodeState:     nil
```

### Linhas 367-405: Checks intermediarios

```
Linha 367: len(OutputHandles) > 0? NAO (vazio) → pathEntry.OutputHandle fica ""
Linha 371: result.Error == nil → pula
Linha 385: result.NodeState nil (sem waitType) → pula
```

### Linha 406: ExecutionPath append

### Linha 408: StatePatch nil → pula

### Linha 418-429: DETECCAO TERMINAL ← AQUI O WORKFLOW ACABA

```go
if len(result.OutputHandles) == 0 {          // [] vazio → TRUE!
    now := time.Now()
    instance.Status = entities.StatusCompleted    // ← "completed"
    instance.ActiveNodeIDs = []string{currentNodeID}  // ["node-5"]
    instance.CompletedAt = &now                   // timestamp
    instance.Updated = now

    // KV CHECKPOINT FINAL
    if err := s.kvCheckpoint(kvKey, instance); err != nil {
        return err
    }

    // NOTIFICA ARCHIVER
    s.publishStateEvent(instance, "workflow.state.completed")
    // → Archiver consumer recebe
    // → KV Get("inst:inst-001") busca estado completo
    // → MongoDB Upsert FULL (~2-5KB) — documento permanente
    // → KV Delete("inst:inst-001") — cleanup (workflow encerrou)

    return nil    // ← SAI DO LOOP. executeInline retorna.
}
```

### Estado FINAL

```
Status:        "completed"
ActiveNodeIDs: ["node-5"]
State:         {"discount_rate": 10}
NodeOutputs:   {}
CompletedAt:   2026-03-10T14:00:00.005Z
ExecutionPath: [
  {nodeId: "node-1", nodeType: "core/start",     status: "completed", outputHandle: "out",  durationMs: 0},
  {nodeId: "node-2", nodeType: "core/set_state",  status: "completed", outputHandle: "out",  durationMs: 0},
  {nodeId: "node-3", nodeType: "core/condition",   status: "completed", outputHandle: "true", durationMs: 0},
  {nodeId: "node-4", nodeType: "core/log",         status: "completed", outputHandle: "out",  durationMs: 0},
  {nodeId: "node-5", nodeType: "core/end",         status: "completed", outputHandle: "",     durationMs: 0},
]
```

---

## De volta ao HandleTrigger

```
Linha 148: s.executeInline retorna nil (sem erro)

Linha 154: logger.Info("Instance inst-001 status=completed currentNode=node-5")

Linha 157: msg.Ack()
  → Consumer NATS confirma processamento da message
  → Message removida do WORKFLOW-TRIGGER stream
  → Worker livre para processar PROXIMO workflow
```

---

## Resumo: O que aconteceu de VERDADE

```
HandleTrigger:
  1. Unmarshal trigger message
  2. FindById → WorkflowDefinition
  3. BuildGraph → ExecutionGraph (adjacency + ParsedConfigs)
  4. Cria WorkflowInstance (State com defaults, EventPayload copiado)
  5. KV Create (checkpoint inicial)
  6. Publish "workflow.state.created" → Archiver → MongoDB InsertOne leve

executeInline (loop com 5 iteracoes):

  Iteracao 0 — start:
    execCtx com State atual → executor.Execute → result{OutputHandles: ["out"]}
    PathEntry append → StatePatch nil (nada muda) → ResolveNextNodes("out") = node-2
    KV Put checkpoint → continue

  Iteracao 1 — set_state:
    execCtx com State atual → executor.Execute → result{OutputHandles: ["out"], StatePatch: {"discount_rate": 10}}
    PathEntry append → APLICA StatePatch: State["discount_rate"] = 10 → ResolveNextNodes("out") = node-3
    KV Put checkpoint → continue

  Iteracao 2 — condition:
    execCtx com State ATUALIZADO (discount_rate=10) → executor.Execute
    EvaluateGroup: 10 > 5 = true → result{OutputHandles: ["true"]}
    PathEntry append → ResolveNextNodes("true") = node-4     ← AQUI DECIDE O CAMINHO
    KV Put checkpoint → continue

  Iteracao 3 — log:
    execCtx com State atual → executor.Execute
    Interpola: "Discount: ${state.discount_rate}%" → "Discount: 10%"
    result{OutputHandles: ["out"], LogEntries: [{message: "Discount: 10%"}]}
    PathEntry append → Publica logs → ResolveNextNodes("out") = node-5
    KV Put checkpoint → continue

  Iteracao 4 — end:
    execCtx com State atual → executor.Execute
    terminateWithError=false → result{OutputHandles: []}     ← VAZIO
    len(OutputHandles)==0 → TERMINAL!
    Status = "completed", CompletedAt = now
    KV Put checkpoint FINAL
    Publish "workflow.state.completed" → Archiver → MongoDB Upsert FULL → KV Delete
    return nil

HandleTrigger:
  msg.Ack() → worker livre
```

```
Tempo total: ~5ms (tudo inline, sem async)
KV checkpoints: 6 (1 criacao + 5 iteracoes)
MongoDB writes: 2 (Archiver: 1 created + 1 completed)
NATS messages: 3 (1 trigger recebida + 1 state.created + 1 state.completed)
```

---

## Decisao Chave: OutputHandles controla TUDO

```
OutputHandles = ["out"]           → continua para proximo node (linear)
OutputHandles = ["true"]          → segue pelo caminho "true" (branch)
OutputHandles = ["false"]         → segue pelo caminho "false" (branch)
OutputHandles = ["case_a"]        → segue pelo caso que matchou (switch)
OutputHandles = ["out_1","out_2"] → spawna goroutines paralelas (fanout)
OutputHandles = []                → TERMINAL — workflow acabou
NodeState com waitType             → SUSPENDE — espera callback/timer/signal
OutputHandles + Error             → FALHA — workflow falhou
```

O executor e uma funcao pura que retorna COMANDOS.
O runtime executa os comandos: aplica patches, faz checkpoint, resolve edges, avanca.
O loop `for` do executeInline e o coracao de tudo.
