# Walkthrough: Workflow de 8 Nodes — Passo a Passo

Este documento rastreia a execucao COMPLETA de um workflow real, mostrando exatamente o que acontece em cada step: o que o executor recebe, o que retorna, o que o runtime faz com o resultado, e como o estado da instancia muda.

---

## O Workflow: "Order Processing"

Trigger: evento `order.created`

```
                                            ┌── true ──→ [5:log] → [6:delay] → [7:trigger_event] → [8:end]
[1:start] → [2:set_state] → [3:code] → [4:condition]
                                            └── false ──→ [8:end]
```

### Nodes

```
ID       Type                Config
───────  ──────────────────  ──────────────────────────────────────
node-1   core/start          (nenhum)
node-2   core/set_state      operation: "set", targetField: "discount_rate", valueSource: {type: "literal", value: "10"}
node-3   core/code           script: "const total = event.total; const rate = state.discount_rate / 100; return { output: { finalTotal: total * (1 - rate) }, statePatch: { final_total: total * (1 - rate) } };", timeout: 30
node-4   core/condition      condition: {logic: "AND", items: [{field: {type: "state", value: "final_total"}, operator: "greater_than", value: {type: "literal", value: "150"}}]}
node-5   core/log            message: "Processing order ${event.orderId} total=${state.final_total} BRL", level: "info"
node-6   core/delay          duration: 10, unit: "s"
node-7   core/trigger_event  eventType: "order.processed", payloadMapping: [{key: "orderId", value: {type: "event", value: "orderId"}}, {key: "finalTotal", value: {type: "state", value: "final_total"}}]
node-8   core/end            terminateWithError: false
```

### Edges

```
Source   SourceHandle   Target
───────  ────────────   ───────
node-1   out            node-2
node-2   out            node-3
node-3   out            node-4
node-4   true           node-5
node-4   false          node-8
node-5   out            node-6
node-6   out            node-7
node-7   out            node-8
```

### Variables (definidas na WorkflowDefinition)

```
Field           Type     Default  Description
──────────────  ───────  ───────  ──────────────────────
discount_rate   number   0        Taxa de desconto (%)
final_total     number   0        Total apos desconto
```

### Event Payload (enviado pelo Trigger Service)

```json
{
  "orderId": "ORD-2024-001",
  "customerName": "Maria Silva",
  "total": 200.00,
  "currency": "BRL"
}
```

---

## Fase 0: Trigger — Criacao da Instancia

O Trigger Service detecta evento `order.created` e publica no `WORKFLOW-TRIGGER`:

```
NATS message → subject: workflow.trigger.{workflowId}
payload: {
  workflowId: "wf-abc",
  orgId: "org-123",
  eventPayload: {"orderId": "ORD-2024-001", "customerName": "Maria Silva", "total": 200.00, "currency": "BRL"}
}
```

O Runtime consumer recebe e chama `HandleTrigger()`:

```
1. Busca WorkflowDefinition via TieredCache (L0 RAM → L1 Disk → L2 MinIO)
2. Cria WorkflowInstance:
```

```
instance = {
  _id:                ObjectID("inst-001"),
  workflowId:         ObjectID("wf-abc"),
  workflowName:       "Order Processing",
  orgId:              ObjectID("org-123"),
  definitionVersion:  3,
  version:            1,
  status:             "created",

  activeNodeIds:      ["node-1"],            ← start node
  nodeStates:         {},                    ← vazio — nenhum node executou
  state:              {"discount_rate": 0, "final_total": 0},   ← defaults das variables
  eventPayload:       {"orderId": "ORD-2024-001", "customerName": "Maria Silva", "total": 200.00, "currency": "BRL"},
  nodeOutputs:        {},
  executionPath:      [],

  variables:          [{field: "discount_rate", type: "number", defaultValue: 0}, {field: "final_total", type: "number", defaultValue: 0}],
  externalInputs:     {},
  captureFields:      [],

  errorInfo:          nil,
  parentInstanceId:   nil,                   ← nao e subworkflow
  parentNodeId:       "",

  startedAt:          2026-03-10T14:00:00.000Z,
  completedAt:        nil,
  created:            2026-03-10T14:00:00.000Z,
  updated:            2026-03-10T14:00:00.000Z,
}
```

```
3. BuildGraph(definition) → ExecutionGraph com adjacency, ParsedConfigs, gotoPairs
4. KV Put("inst:inst-001", instance)          ← checkpoint inicial
5. Publica no WORKFLOW-STATE: {instanceId: "inst-001", status: "created"}   ← Archiver cria doc leve no MongoDB
6. Inicia executeInline(ctx, instance, graph, "node-1")
```

---

## Step 1: start (node-1) — INLINE

### Executor recebe

```go
NodeExecutionContext{
  InstanceID:   ObjectID("inst-001"),
  State:        {"discount_rate": 0, "final_total": 0},
  EventPayload: {"orderId": "ORD-2024-001", "customerName": "Maria Silva", "total": 200.00, "currency": "BRL"},
  NodeOutputs:  {},
  NodeStates:   {},
  NodeID:       "node-1",
  NodeType:     "core/start",
  ParsedConfig: nil,
  Timezone:     "America/Sao_Paulo",
}
```

### Executor retorna

```go
NodeExecutionResult{
  OutputHandles: ["out"],
  StatePatch:    nil,
  NodeState:     nil,
  NodeOutput:    nil,
  LogEntries:    nil,
  Error:         nil,
}
```

### Runtime aplica

```
1. StatePatch nil     → nada muda no instance.State
2. NodeState nil      → nada muda no instance.NodeStates
3. NodeOutput nil     → nada muda no instance.NodeOutputs
4. LogEntries nil     → nada a publicar
5. PathEntry registrada:
   {nodeId: "node-1", nodeType: "core/start", status: "completed", enteredAt: T+0ms, exitedAt: T+0ms, durationMs: 0, outputHandle: "out"}
6. KV Put("inst:inst-001", instance)                    ← checkpoint
7. Resolve OutputHandle "out":
   adjacency["node-1"]["out"] = "node-2"
   → proximo node: node-2
```

### Estado da instancia apos step 1

```
activeNodeIds: ["node-2"]
status:        "running"
state:         {"discount_rate": 0, "final_total": 0}      ← sem mudanca
nodeStates:    {}                                            ← sem mudanca
nodeOutputs:   {}                                            ← sem mudanca
executionPath: [
  {nodeId: "node-1", nodeType: "core/start", status: "completed", outputHandle: "out", durationMs: 0}
]
```

---

## Step 2: set_state (node-2) — INLINE

### Executor recebe

```go
NodeExecutionContext{
  // ... (mesmo de antes, exceto:)
  State:        {"discount_rate": 0, "final_total": 0},
  NodeID:       "node-2",
  NodeType:     "core/set_state",
  ParsedConfig: &SetStateNodeConfig{
    Operation:   "set",
    TargetField: "discount_rate",
    ValueSource: FieldValue{Type: "literal", Value: "10"},
  },
}
```

### Executor executa internamente

```
1. Chama ValueResolver.Resolve(&FieldValue{Type: "literal", Value: "10"}, state, event, nodeOutputs)
   → retorna: float64(10)    (literal convertido para numero)

2. Operation "set":
   → StatePatch = {"discount_rate": float64(10)}
```

### Executor retorna

```go
NodeExecutionResult{
  OutputHandles: ["out"],
  StatePatch:    {"discount_rate": float64(10)},
  NodeState:     nil,
  NodeOutput:    nil,
  LogEntries:    nil,
  Error:         nil,
}
```

### Runtime aplica

```
1. StatePatch {"discount_rate": 10}
   → instance.State["discount_rate"] = 10
   → instance.State agora: {"discount_rate": 10, "final_total": 0}

2. NodeState nil      → nada
3. NodeOutput nil     → nada
4. PathEntry: {nodeId: "node-2", status: "completed", outputHandle: "out", durationMs: 0}
5. KV Put("inst:inst-001", instance)                    ← checkpoint
6. Resolve "out" → adjacency["node-2"]["out"] = "node-3"
```

### Estado da instancia apos step 2

```
activeNodeIds: ["node-3"]
state:         {"discount_rate": 10, "final_total": 0}     ← discount_rate atualizado!
nodeStates:    {}
nodeOutputs:   {}
executionPath: [
  {nodeId: "node-1", ...},
  {nodeId: "node-2", nodeType: "core/set_state", status: "completed", outputHandle: "out", durationMs: 0}
]
```

---

## Step 3: code (node-3) — ASYNC (suspende)

### Executor recebe

```go
NodeExecutionContext{
  State:        {"discount_rate": 10, "final_total": 0},
  EventPayload: {"orderId": "ORD-2024-001", "customerName": "Maria Silva", "total": 200.00, "currency": "BRL"},
  NodeOutputs:  {},
  NodeStates:   {},
  NodeID:       "node-3",
  NodeType:     "core/code",
  ParsedConfig: &CodeNodeConfig{
    Script:  "const total = event.total; const rate = state.discount_rate / 100; return { output: { finalTotal: total * (1 - rate) }, statePatch: { final_total: total * (1 - rate) } };",
    Timeout: 30,
  },
}
```

### Executor retorna

```go
NodeExecutionResult{
  OutputHandles: ["out"],
  StatePatch:    nil,
  NodeState:     {
    "waitType": "callback",
    "script":   "const total = event.total; ...",
    "timeout":  30,
  },
  NodeOutput:    nil,
  LogEntries:    nil,
  Error:         nil,
}
```

### Runtime detecta ASYNC

```
1. NodeState["waitType"] = "callback" → SUSPENSAO!

2. instance.NodeStates["node-3"] = {"waitType": "callback", "script": "...", "timeout": 30}
3. instance.ActiveNodeIDs = ["node-3"]
4. instance.Status = "waiting"
5. PathEntry: {nodeId: "node-3", status: "waiting", enteredAt: T+1ms}    ← sem exitedAt!
6. KV Put("inst:inst-001", instance)                    ← checkpoint com status waiting

7. DESPACHO baseado em nodeType "core/code":
   Publica no WORKFLOW-JS-CODE stream:
   {
     instanceId: "inst-001",
     nodeId:     "node-3",
     orgId:      "org-123",
     workflowId: "wf-abc",
     script:     "const total = event.total; ...",
     timeout:    30,
     state:      {"discount_rate": 10, "final_total": 0},
     eventPayload: {"orderId": "ORD-2024-001", ...},
     variables:  [{"field": "discount_rate", ...}, {"field": "final_total", ...}],
     nodeOutputs: {}
   }

8. executeInline PARA. Consumer retorna (ACK da message original).
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   INSTANCIA SUSPENSA — esperando callback do JS executor
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Estado da instancia SUSPENSA

```
status:        "waiting"
activeNodeIds: ["node-3"]
state:         {"discount_rate": 10, "final_total": 0}
nodeStates:    {
  "node-3": {"waitType": "callback", "script": "...", "timeout": 30}
}
nodeOutputs:   {}
executionPath: [
  {nodeId: "node-1", status: "completed", ...},
  {nodeId: "node-2", status: "completed", ...},
  {nodeId: "node-3", status: "waiting", enteredAt: T+1ms, exitedAt: nil}    ← ABERTO
]
```

---

## Intervalo: js-workflow-executor processa o script

```
js-workflow-executor (servico separado) recebe a mensagem:

1. Busca script compilado via TieredCache
2. Cria V8 Isolate (32MB heap, timeout 30s)
3. Injeta contexto:
   const event = {"orderId": "ORD-2024-001", "customerName": "Maria Silva", "total": 200.00, "currency": "BRL"};
   const state = {"discount_rate": 10, "final_total": 0};
   const variables = {"discount_rate": {"field": "discount_rate", "type": "number"}, ...};
   const nodes = {};

4. Executa script:
   const total = event.total;                    // 200.00
   const rate = state.discount_rate / 100;       // 0.10
   const finalTotal = total * (1 - rate);        // 200.00 * 0.90 = 180.00
   return {
     output: { finalTotal: 180.00 },
     statePatch: { final_total: 180.00 }
   };

5. Publica no WORKFLOW-RESUME:
   {
     instanceId: "inst-001",
     nodeId:     "node-3",
     status:     "ok",
     output:     {"finalTotal": 180.00},
     statePatch: {"final_total": 180.00}
   }
```

---

## Resume 1: Runtime recebe callback do code node

O consumer WORKFLOW-RESUME recebe a mensagem e chama `HandleResume()`:

```
1. Carrega instancia do KV: KV Get("inst:inst-001") → instance (status=waiting)
2. Valida: nodeId "node-3" esta em instance.ActiveNodeIDs? SIM
3. Valida: instance.NodeStates["node-3"]["waitType"] = "callback"? SIM

4. Aplica resume data:
   a. output → instance.NodeOutputs["node-3"] = {"finalTotal": 180.00}
   b. statePatch → instance.State["final_total"] = 180.00
      → instance.State agora: {"discount_rate": 10, "final_total": 180.00}

5. Limpa wait: delete(instance.NodeStates, "node-3")
   → instance.NodeStates agora: {}

6. Atualiza PathEntry do node-3:
   {nodeId: "node-3", status: "completed", exitedAt: T+50ms, durationMs: 49, outputHandle: "out"}

7. instance.Status = "running"
8. KV Put("inst:inst-001", instance)    ← checkpoint

9. Resolve OutputHandle "out":
   adjacency["node-3"]["out"] = "node-4"
   → Continua executeInline a partir de node-4
```

---

## Step 4: condition (node-4) — INLINE

### Executor recebe

```go
NodeExecutionContext{
  State:        {"discount_rate": 10, "final_total": 180.00},    ← atualizado pelo code!
  EventPayload: {"orderId": "ORD-2024-001", ...},
  NodeOutputs:  {"node-3": {"finalTotal": 180.00}},              ← output do code disponivel!
  NodeStates:   {},
  NodeID:       "node-4",
  NodeType:     "core/condition",
  ParsedConfig: &ConditionNodeConfig{
    Condition: ConditionGroup{
      Logic: "AND",
      Items: [{
        Type: "condition",
        Data: ConditionItem{
          Field:    FieldValue{Type: "state", Value: "final_total"},
          Operator: "greater_than",
          Value:    FieldValue{Type: "literal", Value: "150"},
        },
      }],
    },
  },
}
```

### Executor executa internamente

```
1. Chama ConditionEvaluator.EvaluateGroup(&condition, timezone, event, state, nodeOutputs)
2. Resolve field: state["final_total"] = 180.00
3. Resolve value: literal "150" = 150.00
4. Avalia: 180.00 > 150.00 → TRUE
5. Logica AND com 1 item TRUE → grupo TRUE
```

### Executor retorna

```go
NodeExecutionResult{
  OutputHandles: ["true"],        ← condicao verdadeira!
  StatePatch:    nil,
  NodeState:     nil,
  NodeOutput:    nil,
  LogEntries:    nil,
  Error:         nil,
}
```

### Runtime aplica

```
1. Nada a aplicar (tudo nil)
2. PathEntry: {nodeId: "node-4", status: "completed", outputHandle: "true", durationMs: 0}
3. KV Put checkpoint
4. Resolve "true" → adjacency["node-4"]["true"] = "node-5"
```

### Estado apos step 4

```
activeNodeIds: ["node-5"]
state:         {"discount_rate": 10, "final_total": 180.00}
nodeStates:    {}
nodeOutputs:   {"node-3": {"finalTotal": 180.00}}
executionPath: [
  {nodeId: "node-1", status: "completed", outputHandle: "out"},
  {nodeId: "node-2", status: "completed", outputHandle: "out"},
  {nodeId: "node-3", status: "completed", outputHandle: "out", durationMs: 49},     ← durou 49ms (async)
  {nodeId: "node-4", status: "completed", outputHandle: "true", durationMs: 0},
]
```

**Nota**: Se `final_total` fosse <= 150, o OutputHandle seria `"false"` e o runtime seguiria para `node-8` (end) diretamente, pulando nodes 5-7.

---

## Step 5: log (node-5) — INLINE

### Executor recebe

```go
NodeExecutionContext{
  State:        {"discount_rate": 10, "final_total": 180.00},
  EventPayload: {"orderId": "ORD-2024-001", "customerName": "Maria Silva", "total": 200.00, "currency": "BRL"},
  NodeID:       "node-5",
  NodeType:     "core/log",
  ParsedConfig: &LogNodeConfig{
    Message: "Processing order ${event.orderId} total=${state.final_total} BRL",
    Level:   "info",
  },
}
```

### Executor executa internamente

```
1. Interpola message:
   "Processing order ${event.orderId} total=${state.final_total} BRL"
   → ${event.orderId} resolve para "ORD-2024-001"
   → ${state.final_total} resolve para "180"
   → "Processing order ORD-2024-001 total=180 BRL"
2. Level "info" → LogInfo
```

### Executor retorna

```go
NodeExecutionResult{
  OutputHandles: ["out"],
  StatePatch:    nil,
  NodeState:     nil,
  NodeOutput:    nil,
  LogEntries:    [{
    Level:     "info",
    Message:   "Processing order ORD-2024-001 total=180 BRL",
    Timestamp: 2026-03-10T14:00:00.052Z,
    NodeID:    "node-5",
    NodeType:  "core/log",
  }],
  Error:         nil,
}
```

### Runtime aplica

```
1. LogEntries presente → publica no WORKFLOW-LOGS stream (fire-and-forget):
   subject: workflow.logs.inst-001
   payload: {instanceId: "inst-001", entry: {level: "info", message: "Processing order ORD-2024-001 total=180 BRL", ...}}

   Este log vai pro ClickHouse via consumer dedicado (outro modulo).

2. PathEntry: {nodeId: "node-5", status: "completed", outputHandle: "out", durationMs: 0}
3. KV Put checkpoint
4. Resolve "out" → adjacency["node-5"]["out"] = "node-6"
```

---

## Step 6: delay (node-6) — ASYNC (suspende)

### Executor recebe

```go
NodeExecutionContext{
  State:        {"discount_rate": 10, "final_total": 180.00},
  NodeID:       "node-6",
  NodeType:     "core/delay",
  ParsedConfig: &DelayNodeConfig{
    Duration: 10,
    Unit:     "s",
  },
}
```

### Executor executa internamente

```
1. Unit "s" → multiplier = time.Second
2. expiresAt = time.Now() + 10*time.Second = 2026-03-10T14:00:10.053Z
```

### Executor retorna

```go
NodeExecutionResult{
  OutputHandles: ["out"],
  StatePatch:    nil,
  NodeState:     {
    "waitType":  "timer",
    "expiresAt": "2026-03-10T14:00:10.053Z",
  },
  NodeOutput:    nil,
  LogEntries:    nil,
  Error:         nil,
}
```

### Runtime detecta ASYNC

```
1. NodeState["waitType"] = "timer" → SUSPENSAO!

2. instance.NodeStates["node-6"] = {"waitType": "timer", "expiresAt": "2026-03-10T14:00:10.053Z"}
3. instance.ActiveNodeIDs = ["node-6"]
4. instance.Status = "waiting"
5. PathEntry: {nodeId: "node-6", status: "waiting", enteredAt: T+52ms}
6. KV Put("inst:inst-001", instance)

7. DESPACHO baseado em nodeType "core/delay":
   Publica no WORKFLOW-RECONCILER:
   {
     instanceId: "inst-001",
     nodeId:     "node-6",
     waitType:   "timer",
     expiresAt:  "2026-03-10T14:00:10.053Z"
   }

8. executeInline PARA.
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   INSTANCIA SUSPENSA — esperando timer de 10 segundos
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Estado SUSPENSA (segunda vez)

```
status:        "waiting"
activeNodeIds: ["node-6"]
state:         {"discount_rate": 10, "final_total": 180.00}
nodeStates:    {
  "node-6": {"waitType": "timer", "expiresAt": "2026-03-10T14:00:10.053Z"}
}
nodeOutputs:   {"node-3": {"finalTotal": 180.00}}
executionPath: [
  {nodeId: "node-1", status: "completed"},
  {nodeId: "node-2", status: "completed"},
  {nodeId: "node-3", status: "completed", durationMs: 49},
  {nodeId: "node-4", status: "completed"},
  {nodeId: "node-5", status: "completed"},
  {nodeId: "node-6", status: "waiting", exitedAt: nil},         ← ABERTO
]
```

---

## Intervalo: Reconciler processa o timer

```
Reconciler (modulo separado) — goroutine sweep a cada 1 segundo:

1. Recebeu TimerRegistration via WORKFLOW-RECONCILER consumer
2. Adicionou na priority queue: {instanceId: "inst-001", nodeId: "node-6", expiresAt: T+10s}

... 10 segundos passam ...

3. Tick: now() >= expiresAt → timer expirou!
4. Confirma via KV: KV Get("inst:inst-001") → status=waiting, NodeStates["node-6"]["waitType"]="timer" ✓
5. Publica no WORKFLOW-RESUME:
   {
     instanceId: "inst-001",
     nodeId:     "node-6",
     status:     "ok"
   }
```

---

## Resume 2: Runtime recebe timer do Reconciler

```
1. KV Get("inst:inst-001") → instance (status=waiting)
2. Valida: nodeId "node-6" em ActiveNodeIDs? SIM
3. Valida: NodeStates["node-6"]["waitType"] = "timer"? SIM

4. Resume sem data (timer nao tem output/statePatch)
5. Limpa wait: delete(instance.NodeStates, "node-6")
6. Atualiza PathEntry node-6: status: "completed", exitedAt: T+10.053s, durationMs: 10001
7. instance.Status = "running"
8. KV Put checkpoint

9. Resolve "out" → adjacency["node-6"]["out"] = "node-7"
   → Continua executeInline a partir de node-7
```

---

## Step 7: trigger_event (node-7) — ASYNC (suspende)

### Executor recebe

```go
NodeExecutionContext{
  State:        {"discount_rate": 10, "final_total": 180.00},
  EventPayload: {"orderId": "ORD-2024-001", "customerName": "Maria Silva", "total": 200.00, "currency": "BRL"},
  NodeOutputs:  {"node-3": {"finalTotal": 180.00}},
  NodeID:       "node-7",
  NodeType:     "core/trigger_event",
  ParsedConfig: &TriggerEventNodeConfig{
    EventType: "order.processed",
    PayloadMapping: [
      {Key: "orderId", Value: FieldValue{Type: "event", Value: "orderId"}},
      {Key: "finalTotal", Value: FieldValue{Type: "state", Value: "final_total"}},
    ],
  },
}
```

### Executor executa internamente

```
1. Valida eventType "order.processed" nao vazio ✓
2. Resolve payloadMapping[0]: event["orderId"] = "ORD-2024-001"
3. Resolve payloadMapping[1]: state["final_total"] = 180.00
4. payload = {"orderId": "ORD-2024-001", "finalTotal": 180.00}
```

### Executor retorna

```go
NodeExecutionResult{
  OutputHandles: ["out"],
  StatePatch:    nil,
  NodeState:     {
    "waitType":  "callback",
    "eventType": "order.processed",
    "payload":   {"orderId": "ORD-2024-001", "finalTotal": 180.00},
  },
  NodeOutput:    nil,
  LogEntries:    nil,
  Error:         nil,
}
```

### Runtime detecta ASYNC

```
1. NodeState["waitType"] = "callback" → SUSPENSAO!

2. instance.NodeStates["node-7"] = {"waitType": "callback", "eventType": "order.processed", "payload": {...}}
3. instance.ActiveNodeIDs = ["node-7"]
4. instance.Status = "waiting"
5. PathEntry: {nodeId: "node-7", status: "waiting"}
6. KV Put checkpoint

7. DESPACHO baseado em nodeType "core/trigger_event":
   Publica no stream do Trigger Service:
   {
     instanceId: "inst-001",
     nodeId:     "node-7",
     orgId:      "org-123",
     workflowId: "wf-abc",
     eventType:  "order.processed",
     payload:    {"orderId": "ORD-2024-001", "finalTotal": 180.00}
   }

8. executeInline PARA.
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   INSTANCIA SUSPENSA — esperando Trigger Service processar
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## Intervalo: Trigger Service processa o evento

```
Trigger Service (servico SEPARADO):

1. Recebe TriggerEventRequest via NATS consumer
2. Processa evento "order.processed" — pode:
   - Disparar webhooks configurados
   - Triggerar OUTROS workflows que ouvem "order.processed"
   - Enviar notificacoes
   - Publicar em RabbitMQ, Kafka, etc.

3. Apos processar, responde:
   Publica no WORKFLOW-RESUME:
   {
     instanceId: "inst-001",
     nodeId:     "node-7",
     status:     "ok"
   }
```

---

## Resume 3: Runtime recebe callback do Trigger Service

```
1. KV Get("inst:inst-001") → instance (status=waiting)
2. Valida: nodeId "node-7" em ActiveNodeIDs? SIM
3. Valida: NodeStates["node-7"]["waitType"] = "callback"? SIM

4. Resume sem data (trigger_event nao retorna output)
5. Limpa wait: delete(instance.NodeStates, "node-7")
6. Atualiza PathEntry node-7: status: "completed"
7. instance.Status = "running"
8. KV Put checkpoint

9. Resolve "out" → adjacency["node-7"]["out"] = "node-8"
   → Continua executeInline a partir de node-8
```

---

## Step 8: end (node-8) — INLINE (terminal)

### Executor recebe

```go
NodeExecutionContext{
  State:        {"discount_rate": 10, "final_total": 180.00},
  EventPayload: {"orderId": "ORD-2024-001", ...},
  NodeOutputs:  {"node-3": {"finalTotal": 180.00}},
  NodeStates:   {},
  NodeID:       "node-8",
  NodeType:     "core/end",
  ParsedConfig: &EndNodeConfig{
    TerminateWithError: false,
  },
}
```

### Executor retorna

```go
NodeExecutionResult{
  OutputHandles: [],         ← VAZIO — sinaliza workflow completo
  StatePatch:    nil,
  NodeState:     nil,
  NodeOutput:    nil,
  LogEntries:    nil,
  Error:         nil,        ← nil — sem erro
}
```

### Runtime detecta TERMINAL

```
1. OutputHandles VAZIO + Error nil → WORKFLOW COMPLETOU COM SUCESSO

2. instance.Status = "completed"
3. instance.CompletedAt = 2026-03-10T14:00:10.058Z
4. instance.ActiveNodeIDs = []
5. PathEntry: {nodeId: "node-8", status: "completed", durationMs: 0}
6. KV Put("inst:inst-001", instance)    ← checkpoint final

7. Publica no WORKFLOW-STATE:
   {instanceId: "inst-001", status: "completed"}
   → Archiver consumer recebe → KV Get (full state) → MongoDB Upsert FULL → KV Delete (cleanup)

8. executeInline retorna. Consumer ACK.
```

---

## Estado FINAL da instancia

```json
{
  "_id":                "inst-001",
  "workflowId":         "wf-abc",
  "workflowName":       "Order Processing",
  "orgId":              "org-123",
  "definitionVersion":  3,
  "version":            12,
  "status":             "completed",

  "activeNodeIds":      [],
  "nodeStates":         {},
  "state":              {"discount_rate": 10, "final_total": 180.00},
  "eventPayload":       {"orderId": "ORD-2024-001", "customerName": "Maria Silva", "total": 200.00, "currency": "BRL"},
  "nodeOutputs":        {"node-3": {"finalTotal": 180.00}},

  "executionPath": [
    {"nodeId": "node-1", "nodeType": "core/start",         "status": "completed", "outputHandle": "out",  "durationMs": 0},
    {"nodeId": "node-2", "nodeType": "core/set_state",     "status": "completed", "outputHandle": "out",  "durationMs": 0},
    {"nodeId": "node-3", "nodeType": "core/code",          "status": "completed", "outputHandle": "out",  "durationMs": 49},
    {"nodeId": "node-4", "nodeType": "core/condition",     "status": "completed", "outputHandle": "true", "durationMs": 0},
    {"nodeId": "node-5", "nodeType": "core/log",           "status": "completed", "outputHandle": "out",  "durationMs": 0},
    {"nodeId": "node-6", "nodeType": "core/delay",         "status": "completed", "outputHandle": "out",  "durationMs": 10001},
    {"nodeId": "node-7", "nodeType": "core/trigger_event", "status": "completed", "outputHandle": "out",  "durationMs": 5},
    {"nodeId": "node-8", "nodeType": "core/end",           "status": "completed", "outputHandle": "",     "durationMs": 0}
  ],

  "variables":          [{"field": "discount_rate", ...}, {"field": "final_total", ...}],
  "externalInputs":     {},
  "captureFields":      [],

  "errorInfo":          null,
  "parentInstanceId":   null,
  "parentNodeId":       "",

  "startedAt":          "2026-03-10T14:00:00.000Z",
  "completedAt":        "2026-03-10T14:00:10.058Z",
  "created":            "2026-03-10T14:00:00.000Z",
  "updated":            "2026-03-10T14:00:10.058Z"
}
```

---

## Resumo da Execucao

```
Timeline:

T+0ms      HandleTrigger → cria instancia → KV Put → publish WORKFLOW-STATE "created"
T+0ms      Step 1: start        → INLINE → OutputHandle "out"      → KV Put
T+0ms      Step 2: set_state    → INLINE → OutputHandle "out"      → KV Put → discount_rate=10
T+1ms      Step 3: code         → ASYNC  → SUSPENDE (callback)     → KV Put → dispatch JS executor
                    ╌╌╌╌╌╌ instancia SUSPENSA ╌╌╌╌╌╌
T+50ms     Resume 1: JS executor responde → final_total=180        → KV Put
T+50ms     Step 4: condition    → INLINE → OutputHandle "true"     → KV Put → 180 > 150 ✓
T+51ms     Step 5: log          → INLINE → OutputHandle "out"      → KV Put → log publicado
T+52ms     Step 6: delay        → ASYNC  → SUSPENDE (timer 10s)    → KV Put → dispatch Reconciler
                    ╌╌╌╌╌╌ instancia SUSPENSA ╌╌╌╌╌╌
T+10.053s  Resume 2: Reconciler timer expira                        → KV Put
T+10.053s  Step 7: trigger_event → ASYNC → SUSPENDE (callback)     → KV Put → dispatch Trigger Service
                    ╌╌╌╌╌╌ instancia SUSPENSA ╌╌╌╌╌╌
T+10.058s  Resume 3: Trigger Service responde                       → KV Put
T+10.058s  Step 8: end          → INLINE → OutputHandles []        → KV Put → publish WORKFLOW-STATE "completed"
                    ━━━━━━ WORKFLOW COMPLETO ━━━━━━

Total: 10.058s (10s foram o delay)
KV checkpoints: 12 (1 criacao + 8 steps + 3 resumes)
NATS messages publicadas: 6 (1 state.created + 1 JS dispatch + 1 reconciler + 1 trigger dispatch + 1 state.completed + 1 log)
MongoDB writes: 2 (Archiver: 1 created leve + 1 completed FULL)
```

---

## Pontos Criticos Demonstrados

### 1. KV Checkpoint a cada step
Se o servico crashar entre step 2 e step 3, o runtime reinicia do step 2 (ultimo KV checkpoint). O `discount_rate` ja esta no state — nenhum dado perdido.

### 2. Async suspende e libera worker
O consumer NATS faz ACK apos cada suspensao. O worker fica livre para processar OUTRO workflow enquanto este espera o JS executor (49ms), o timer (10s), ou o Trigger Service.

### 3. Resume via subject unico
Todos os 3 resumes chegam pelo mesmo subject `workflow.resume`. O payload contém `instanceId` + `nodeId` — o runtime sabe exatamente qual node retomar.

### 4. StatePatch vs NodeState (separacao clara)
- Step 2 (set_state): StatePatch `{discount_rate: 10}` → vai para `instance.State`
- Step 3 (code): NodeState `{waitType: "callback", ...}` → vai para `instance.NodeStates["node-3"]`
- Resume 1: statePatch `{final_total: 180}` → vai para `instance.State`
- Nunca se misturam. Sem ambiguidade.

### 5. NodeOutputs acumula
Apos step 3, `instance.NodeOutputs["node-3"] = {finalTotal: 180}`. Qualquer node posterior pode acessar via `FieldValue{Type: "node_output", NodeID: "node-3"}`.

### 6. ExecutionPath = Event History
Cada PathEntry registra entrada, saida, duracao, e handle usado. O frontend pode renderizar o DAG com cores: verde (completed), amarelo (waiting), vermelho (error).

### 7. Condition branching
Se o total fosse $100 (final_total = $90 apos desconto), o condition retornaria `["false"]`, e o runtime pularia direto para node-8 (end). Nodes 5-7 nunca executariam. O executionPath mostraria apenas 5 entries.

### 8. MongoDB minimal
Archiver fez APENAS 2 writes: um InsertOne leve (~200 bytes) no `created`, e um Upsert FULL (~5KB) no `completed`. Zero writes intermediarias. KV foi a verdade durante toda a execucao.
