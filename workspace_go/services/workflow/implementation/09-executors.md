# Node Executors — 17 tipos

Decisões: C1-C17

---

## Execution Model (Temporal-inspired)

O executor é uma **função pura**: recebe contexto de leitura, retorna comandos para o runtime.

```
Executor = Função PURA
  Input:  NodeExecutionContext (readonly: state, nodeStates, config, eventPayload, nodeOutputs)
  Output: NodeExecutionResult  (comandos: outputHandles, statePatch, nodeState, nodeOutput, logs, error)

  NÃO modifica state
  NÃO publica em NATS
  NÃO acessa KV/MongoDB

Runtime = Orquestrador de I/O
  Recebe os comandos do executor e EXECUTA:
  1. Aplica StatePatch     → instance.State[key] = value
  2. Aplica NodeState      → instance.NodeStates[nodeId] = {...}
  3. Aplica NodeOutput     → instance.NodeOutputs[nodeId] = value
  4. Publica LogEntries    → WORKFLOW-LOGS stream
  5. Registra PathEntry    → instance.ExecutionPath
  6. Se NodeState tem waitType → suspende, despacha, KV checkpoint
  7. Se sem waitType         → KV checkpoint, avança pelo OutputHandles
```

Inspirado no Temporal.io: workflow code produz Commands, Server executa.
A diferença: Temporal faz replay do histórico completo. MapexOS continua do último KV checkpoint.

---

## Classificação por comportamento

- **INLINE** (7): executa e segue imediatamente. Sem suspensão.
- **ASYNC** (5): retorna NodeState com `waitType` → runtime suspende e despacha para serviço externo.
- **CONTROL** (5): controla fluxo de execução (fanout, merge, sequence, loop, wait_for).

---

## Handle Reference (source of truth = frontend DSL)

```
Node Type        Output Handles                           Significado
─────────────    ──────────────────────────────────────   ─────────────────────────
start            "out"                                    Saída única — inicia workflow
end              [] (vazio)                               Terminal — workflow acabou
condition        "true", "false"                          Resultado da condição
switch           "case_{id}", "default"                   Caso(s) que matcharam ou fallback
set_state        "out"                                    Saída única — state modificado
log              "out"                                    Saída única — log emitido
goto             "out"                                    Saída lógica (resolved no BuildGraph)
trigger_event    "out"                                    Saída após callback do Trigger Service
code             "out"                                    Saída após callback do JS executor
subworkflow      "out"                                    Saída após workflow filho completar
delay            "out"                                    Saída após timer expirar
wait_signal      "out"                                    Saída após signal recebido
wait_for         "matched"                                Condição satisfeita
fanout           "out_1", "out_2", ..., "out_N"           N branches paralelas
merge            "out" ou [] (esperando)                  Branches convergem
sequence         "step_1", ..., "step_N", "done"          Steps em ordem, "done" ao final
loop             "body", "done"                           "body" = próxima iteração, "done" = acabou
```

---

## Padrão ASYNC — Fluxo Unificado

Todos os 5 async nodes seguem o MESMO padrão. A diferença é quem responde.

```
1. EXECUTOR retorna NodeState com waitType
   → Executor NÃO publica nada. Apenas retorna dados.

2. RUNTIME detecta waitType no NodeState
   → instance.NodeStates[nodeId] = resultado do executor
   → instance.ActiveNodeIDs = [nodeId]
   → instance.Status = waiting
   → KV checkpoint

3. RUNTIME despacha baseado no node.Type:
   → core/delay         → publica TimerRegistration no WORKFLOW-RECONCILER
   → core/wait_signal   → (nada — espera chamada HTTP externa)
   → core/wait_for      → (nada — reavalia em cada interação)
   → core/code          → publica CodeExecutionRequest no WORKFLOW-JS-CODE
   → core/subworkflow   → publica SubworkflowTrigger no WORKFLOW-TRIGGER
   → core/trigger_event → publica TriggerEventRequest no stream do Trigger Service

4. SERVIÇO EXTERNO processa e responde
   → Publica no WORKFLOW-RESUME (subject único: workflow.resume):
     {instanceId, nodeId, status: "ok"|"error", output, statePatch, error}

5. RUNTIME recebe resume
   → Carrega instância do KV
   → Aplica output/statePatch do resume
   → Limpa NodeStates[nodeId]
   → Segue pela edge do OutputHandle
   → Continua execução inline
```

### Quem responde cada tipo:

```
waitType      Quem despacha         Quem responde              Quando
────────────  ────────────────────  ─────────────────────────  ──────────────────────────
timer         WORKFLOW-RECONCILER   Reconciler (sweep 1s)      Timer expira
signal        (nada)                HTTP API externa            Sistema externo decide
condition     (nada)                Qualquer resume passivo     Quando state muda
callback      WORKFLOW-JS-CODE      js-workflow-executor        Script termina (~ms a ~30s)
callback      WORKFLOW-TRIGGER      Workflow filho              Filho completa
callback      TRIGGER-SERVICE       Trigger Service             Evento processado
```

### Todos voltam pelo mesmo caminho:

```
WORKFLOW-RESUME ← {instanceId, nodeId, status, output?, statePatch?, error?}
```

Subject único. instanceId + nodeId no payload. Sem subjects dinâmicos.

---

## Executors INLINE

### C1. start

```
Complexidade: Trivial
Config: nenhum

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    —
  NodeState:     —
  NodeOutput:    —
  LogEntries:    —
  Error:         —

Notas: Passthrough. O state já foi inicializado pelo runtime (Variables com defaults).
```

### C2. end

```
Complexidade: Simple
Config: { terminateWithError, errorCode, errorMessage }

Executor retorna:
  OutputHandles: [] (vazio — sinaliza workflow completo)
  StatePatch:    —
  NodeState:     —
  NodeOutput:    —
  LogEntries:    —
  Error:         se terminateWithError=true → {code: errorCode, message: resolvedMessage}

Runtime detecta:
  OutputHandles vazio + sem Error → status=completed
  OutputHandles vazio + com Error → status=failed com ErrorInfo

Dependência: ValueResolver (resolver errorMessage de FieldValue → string)
```

### C3. condition

```
Complexidade: Medium
Config: { condition: ConditionGroup, selectedTemplateIds }

Executor retorna:
  OutputHandles: ["true"] se condição verdadeira
                 ["false"] se condição falsa
  StatePatch:    —
  NodeState:     —
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se avaliação falhar

Lógica:
  1. Chama ConditionEvaluator.EvaluateGroup(condition, timezone, event, state, nodeOutputs)
  2. true → ["true"], false → ["false"]

Dependência: ConditionEvaluator (engine module)
```

### C4. switch

```
Complexidade: Medium
Config: { cases: SwitchCase[], matchMode: "first"|"all" }

Executor retorna:
  OutputHandles: ["case_{id}"]      (1 handle, mode first — primeiro que matcha)
                 ["case_{id}", ...]  (N handles, mode all — todos que matcham)
                 ["default"]         (nenhum matchou)
  StatePatch:    —
  NodeState:     —
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se avaliação falhar

Lógica:
  1. Para cada case, avalia case.Condition via ConditionEvaluator
  2. mode "first": retorna handle do PRIMEIRO match, para
  3. mode "all": retorna handles de TODOS que matcharam
  4. Nenhum match → ["default"]

Dependência: ConditionEvaluator (engine module)
```

### C5. set_state

```
Complexidade: Simple
Config: { operation, targetField, valueSource }

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    {targetField: newValue}
  NodeState:     —
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se resolução ou operação falhar

Operações:
  "set":       StatePatch = {targetField: resolvedValue}
  "increment": StatePatch = {targetField: currentValue + resolvedValue}
  "decrement": StatePatch = {targetField: currentValue - resolvedValue}
  "append":    StatePatch = {targetField: append(currentArray, resolvedValue)}
  "remove":    StatePatch = {targetField: removeFromArray(currentArray, resolvedValue)}

Dependência: ValueResolver (resolver valueSource de FieldValue → valor concreto)
```

### C6. log

```
Complexidade: Trivial
Config: { message, level }

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    —
  NodeState:     —
  NodeOutput:    —
  LogEntries:    [{level: level, message: interpolatedMessage, timestamp, nodeId, nodeType}]
  Error:         —

Lógica:
  1. Interpola message: ${state.field} → valor do state, ${event.field} → valor do event
  2. Level: "debug"|"info"|"warn"|"error" (default: "info")
  3. Runtime publica LogEntry no WORKFLOW-LOGS stream (fire-and-forget)
```

### C17. goto (portal virtual)

```
Complexidade: Simple
Config: { role: "sender"|"receiver", pairLabel, pairColor }

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    —
  NodeState:     —
  NodeOutput:    —
  LogEntries:    —
  Error:         Se sender sem receiver matching → GOTO_NO_RECEIVER

Lógica:
  1. Assert ParsedConfig como *GotoNodeConfig
  2. Se role == "sender":
     → Valida que GraphBuilder injetou edge via graph.HasEdge(nodeID, "out")
     → Se NÃO tem edge → retorna ExecutionError{Code: "GOTO_NO_RECEIVER"}
     → Se tem edge → passthrough, retorna ["out"]
  3. Se role == "receiver":
     → Passthrough, retorna ["out"]

  O routing é resolvido no BuildGraph:
    - Sender: adjacency[senderId]["out"] = gotoPairs[pairLabel] (aponta pro receiver)
    - Receiver: adjacency[receiverId]["out"] = edges normais desenhadas pelo usuário
  Sem edge visual — conexão lógica via pairLabel.

Prevenção de loops infinitos: MaxInlineSteps = 500 (E2)
```

---

## Executors ASYNC

### C7. trigger_event

```
Complexidade: Medium
Config: { eventType, payloadMapping: [{key, value: FieldValue}] }

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    —
  NodeState:     {
    "waitType":  "callback",
    "eventType": config.eventType,
    "payload":   resolvedPayloadMap
  }
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se eventType vazio ou resolução falhar

Lógica:
  1. Valida eventType não vazio
  2. Resolve cada payloadMapping via ValueResolver → monta payload map
  3. Retorna NodeState com dados do evento

Runtime despacha:
  Publica no stream do Trigger Service:
  {instanceId, nodeId, orgId, workflowId, eventType, payload}

Trigger Service responde:
  Publica no WORKFLOW-RESUME: {instanceId, nodeId, status: "ok"}

Dependência: ValueResolver
```

### C8. code (JavaScript isolado)

```
Complexidade: Medium
Config: { script, timeout }

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    —
  NodeState:     {
    "waitType": "callback",
    "script":   config.script,
    "timeout":  config.timeout
  }
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se config inválido

Runtime despacha:
  Publica no WORKFLOW-JS-CODE stream:
  {instanceId, nodeId, orgId, workflowId,
   script, timeout,
   eventPayload, state, variables, nodeOutputs}

js-workflow-executor processa:
  1. Busca script compilado via TieredCache (L0 RAM → L1 Disk → L2 MinIO → fallback HTTP)
  2. Cria V8 Isolate (32MB heap, timeout configurado)
  3. Injeta contexto: event, state, variables, nodes
  4. Executa script → captura result.output e result.statePatch
  5. Publica no WORKFLOW-RESUME:
     Sucesso: {instanceId, nodeId, status: "ok", output: result.output, statePatch: result.statePatch}
     Erro:    {instanceId, nodeId, status: "error", error: {code: "SCRIPT_ERROR", message: "..."}}

O que o script do usuário vê:
  const event = { temperature: 25.5 };
  const state = { counter: 3 };
  const variables = { threshold: 30 };
  const nodes = { condition_1: { matched: true } };

  const result = {
    output: { converted: event.temperature * 1.8 + 32 },
    statePatch: { counter: state.counter + 1 }
  };

Segurança (isolamento):
  V8 Isolate: 32MB heap, zero acesso a Node.js/FS/network
  Timeout: hard limit (config.timeout)
  OOM recovery: isolate.isDisposed → NACK → recria
  Context recycling: a cada 10K execuções

Serviço separado: js-workflow-executor (NÃO compartilha com js-executor de ingestão)
```

### C9. subworkflow

```
Complexidade: Complex
Config: { workflowId, workflowName, executionMode, terminationPolicy, timeout, inputMappings, outputMappings }

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    —
  NodeState:     {
    "waitType":           "callback",
    "workflowId":         config.workflowId,
    "workflowName":       config.workflowName,
    "executionMode":      config.executionMode,
    "terminationPolicy":  config.terminationPolicy,
    "inputData":          resolvedInputMappings,
    "outputMappings":     config.outputMappings,
    "timeout":            config.timeout
  }
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se config inválido ou resolução de input falhar

Lógica:
  1. Resolve inputMappings via ValueResolver → monta inputData map
  2. Retorna NodeState com dados do subworkflow

Runtime despacha:
  Publica no WORKFLOW-TRIGGER stream:
  {workflowId, eventPayload: inputData,
   parentInstanceId: instance.ID, parentNodeId: nodeId}

Workflow filho executa normalmente. Quando completa:
  Publica no WORKFLOW-RESUME do PAI:
  {instanceId: parentInstanceId, nodeId: parentNodeId, status: "ok",
   output: {outputMappings do filho}}

TerminationPolicy (Temporal.io inspired):
  "terminate" (default): pai completa/falha → filho cancelado forçosamente
  "abandon": pai completa/falha → filho continua independente
  "request_cancel": pai completa/falha → signal de cancel ao filho

Dependência: ValueResolver
```

### C14. delay (timer)

```
Complexidade: Simple
Config: { duration, unit }

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    —
  NodeState:     {
    "waitType":  "timer",
    "expiresAt": now() + duration (calculado pela unit)
  }
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se unit inválida

Units: "s"|"seconds", "m"|"minutes", "h"|"hours", "d"|"days"

Runtime despacha:
  Publica no WORKFLOW-RECONCILER:
  {instanceId, nodeId, waitType: "timer", expiresAt}

Reconciler:
  Priority queue in-memory → tick 1s → timer expirou?
  Sim → publica no WORKFLOW-RESUME: {instanceId, nodeId, status: "ok"}
```

### C15. wait_signal (espera sinal externo)

```
Complexidade: Medium
Config: { signalName, timeout, maxTimeoutCycles, mappings }

Executor retorna:
  OutputHandles: ["out"]
  StatePatch:    —
  NodeState:     {
    "waitType":         "signal",
    "signalName":       config.signalName,
    "maxTimeoutCycles": config.maxTimeoutCycles,
    "timeoutCount":     0
  }
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se config inválido

Runtime:
  Se timeout configurado → também registra timer no Reconciler (expiresAt = now + timeout)

Resume via HTTP API:
  POST /api/v1/workflows/instances/{instanceId}/signal
  Body: {signalName: "approval", data: {approved: true}}

  Handler valida signalName → publica no WORKFLOW-RESUME:
  {instanceId, nodeId, status: "ok", output: signalData}

Timeout via Reconciler:
  Timer expira → resume com status "timeout"
  Runtime incrementa timeoutCount
  Se timeoutCount >= maxTimeoutCycles → workflow falha
  Se não → renova timeout, permanece waiting
```

---

## Executors CONTROL

### C10. fanout (inline goroutines)

```
Complexidade: Complex
Config: { branches: N }

Executor retorna:
  OutputHandles: ["out_1", "out_2", ..., "out_N"]
  StatePatch:    —
  NodeState:     —
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se branches <= 0 ou branches > MaxFanoutBranches

Hard limit: MaxFanoutBranches = 20

Lógica:
  1. Valida N <= 20
  2. Retorna N handles dinâmicos: out_1, out_2, ..., out_N

Runtime resolve fanout:
  1. ResolveNextNodes retorna N target nodes (um por handle)
  2. Spawna N goroutines (sync.WaitGroup)
  3. Cada goroutine executa inline com cópia isolada do state
  4. Cada branch para quando: encontra merge, encontra async, ou termina
  5. Se todas branches sync → merge state patches, continua
  6. Se alguma branch async → checkpoint com NodeStates de cada async node

Fanout com múltiplos async simultâneos:
  Cada branch async gera uma entrada em instance.NodeStates.
  Cada resume resolve uma branch. Quando merge condition é met → continua.
```

### C11. merge

```
Complexidade: Complex
Config: { branches: N, strategy: "all"|"any"|"first" }

Executor retorna:
  OutputHandles: ["out"] se condition met
                 [] (vazio) se ainda esperando branches
  StatePatch:    —
  NodeState:     {
    "expectedBranches":  config.branches,
    "completedBranches": completedCount
  }
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se config inválido

Strategies:
  "all" (default): espera TODAS as branches completarem → ["out"]
  "any":           primeira branch que completar → ["out"], cancela restantes
  "first":         primeira branch → ["out"], cancela restantes

Merge de state patches (all strategy):
  Branch 0 patch: { counter: 5, name: "alice" }
  Branch 1 patch: { counter: 10, items: [1,2,3] }
  Resultado (last-write-wins, branch index maior ganha): { counter: 10, name: "alice", items: [1,2,3] }

Strategy "any" — cancelamento:
  Quando primeira branch completa → marca outras branches como cancelled
  Se branches canceladas têm async nodes ativos → runtime publica cancel
```

### C12. sequence (inline sequencial)

```
Complexidade: Medium
Config: { steps: N }

Executor retorna:
  OutputHandles: ["step_{currentStep}"] se ainda tem steps
                 ["done"] se todos completaram
  StatePatch:    —
  NodeState:     {
    "counter": currentStep,
    "total":   config.steps
  }
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se config inválido

Lógica:
  1. Lê counter do NodeStates[nodeId] (0 se primeiro acesso)
  2. Se counter < total → retorna ["step_{counter+1}"], incrementa counter
  3. Se counter >= total → retorna ["done"]

Cada step é uma saída diferente no grafo. Step 1 vai pro sub-grafo 1, step 2 pro sub-grafo 2, etc.
Quando o sub-grafo volta pro sequence node, o counter avança.
```

### C13. loop (inline iteração)

```
Complexidade: Complex
Config: { source: FieldValue → array }

Executor retorna:
  OutputHandles: ["body"] se ainda tem items
                 ["done"] se todos processados
  StatePatch:    {
    "loop_item":  items[currentIndex],    ← item atual no state do USUÁRIO
    "loop_index": currentIndex            ← index atual no state do USUÁRIO
  }
  NodeState:     {
    "counter":     currentIndex + 1,      ← estado INTERNO (não mistura com state do usuário)
    "total":       len(items),
    "currentItem": items[currentIndex]
  }
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se source não é array ou excede MaxLoopIterations

Lógica:
  1. Resolve source via ValueResolver → array
  2. Lê counter do NodeStates[nodeId] (0 se primeiro acesso)
  3. Se counter < total:
     → StatePatch injeta loop_item e loop_index no state do usuário
     → NodeState incrementa counter interno
     → OutputHandles: ["body"]
  4. Se counter >= total:
     → OutputHandles: ["done"]
     → NodeState mantém counter final (debug)

Hard limit: MaxLoopIterations = 10000

StatePatch vs NodeState:
  StatePatch → vai pro instance.State (loop_item, loop_index — o USUÁRIO usa nos nodes do body)
  NodeState  → vai pro instance.NodeStates[nodeId] (counter — estado INTERNO do loop)
  São mapas DIFERENTES. Sem ambiguidade.

Dependência: ValueResolver (resolver source de FieldValue → array)
```

### C16. wait_for (EVENT-DRIVEN — zero polling)

```
Complexidade: Medium
Config: { field, operator, compareTo, timeout, maxTimeoutCycles }

Executor retorna:
  SE condição JÁ é verdadeira:
    OutputHandles: ["matched"]
    NodeState:     {} (vazio — sem wait)

  SE condição NÃO é verdadeira:
    OutputHandles: ["matched"]
    NodeState:     {
      "waitType":         "condition",
      "field":            config.field,
      "operator":         config.operator,
      "compareTo":        config.compareTo,
      "maxTimeoutCycles": config.maxTimeoutCycles,
      "timeoutCount":     0
    }

  StatePatch:    —
  NodeOutput:    —
  LogEntries:    —
  Error:         Go error se avaliação falhar

Lógica:
  1. Avalia condição: state[field] {operator} compareTo
  2. Se TRUE → retorna NodeState vazio (sem wait, continua inline)
  3. Se FALSE → retorna NodeState com waitType "condition" → runtime suspende

Re-avaliação (ZERO POLLING):
  Qualquer interação com a instância (signal, callback, resume) passa pelo Runtime.
  Runtime carrega instância → ANTES de continuar, verifica TODOS os NodeStates com waitType "condition":
    1. Aplica dados recebidos ao state (se houver StatePatch no resume)
    2. Re-avalia: state[field] {operator} compareTo
    3. Se MET → limpa NodeState, segue pela edge "matched"
    4. Se NOT MET → permanece waiting

  Custo: ~1µs por re-avaliação (condição simples).
  Zero goroutine. Zero timer periódico. Apenas eventos reais triggam re-avaliação.

Timeout via Reconciler:
  Se timeout configurado → timer registrado no Reconciler
  Timer expira → resume com timeout
  timeoutCount++ → se >= maxTimeoutCycles → workflow falha

Dependência: ConditionEvaluator (engine module)
```

---

## Async — Fluxos Detalhados

### Delay: Timer

```
EXECUTOR retorna NodeState: {waitType: "timer", expiresAt: "2026-03-10T15:30:00Z"}

RUNTIME:
  → instance.NodeStates["delay_1"] = {waitType: "timer", expiresAt: ...}
  → instance.Status = waiting
  → KV checkpoint
  → Publica no WORKFLOW-RECONCILER: {instanceId, nodeId, expiresAt}

RECONCILER:
  → Priority queue → tick 1s → timer expirou
  → Publica no WORKFLOW-RESUME: {instanceId: "abc", nodeId: "delay_1", status: "ok"}

RUNTIME recebe resume:
  → Carrega instância do KV
  → Limpa NodeStates["delay_1"]
  → Segue pela edge "out"
```

### Wait Signal: Signal Externo

```
EXECUTOR retorna NodeState: {waitType: "signal", signalName: "approval", maxTimeoutCycles: 5}

RUNTIME:
  → instance.NodeStates["wait_signal_1"] = {waitType: "signal", ...}
  → instance.Status = waiting
  → KV checkpoint
  → (nada — espera chamada HTTP)

SISTEMA EXTERNO:
  POST /api/workflows/instances/{id}/signal
  Body: {signalName: "approval", data: {approved: true, approver: "John"}}

HTTP HANDLER:
  → Valida signalName
  → Publica no WORKFLOW-RESUME: {instanceId, nodeId, status: "ok", output: {approved: true, ...}}

RUNTIME recebe resume:
  → instance.NodeOutputs["wait_signal_1"] = {approved: true, approver: "John"}
  → Limpa NodeStates["wait_signal_1"]
  → Segue pela edge "out"
```

### Wait For: Condição no State (Event-Driven)

```
EXECUTOR retorna NodeState: {waitType: "condition", field: "status", operator: "==", compareTo: "done"}

RUNTIME:
  → instance.NodeStates["wait_for_1"] = {waitType: "condition", ...}
  → instance.Status = waiting
  → KV checkpoint
  → (nada — fica esperando passivamente)

QUALQUER RESUME que chega para esta instância:
  → Runtime carrega instância do KV
  → Aplica StatePatch do resume (se houver)
  → Verifica TODOS os NodeStates com waitType "condition"
  → Avalia: state["status"] == "done"
  → Se TRUE → limpa NodeStates["wait_for_1"], segue pela edge "matched"
  → Se FALSE → permanece esperando
```

### Code: JavaScript Isolado

```
EXECUTOR retorna NodeState: {waitType: "callback", script: "return state.price * 1.1;", timeout: 30}

RUNTIME:
  → instance.NodeStates["code_1"] = {waitType: "callback", ...}
  → instance.Status = waiting
  → KV checkpoint
  → Publica no WORKFLOW-JS-CODE: {instanceId, nodeId, script, timeout, state, eventPayload, variables, nodeOutputs}

JS-WORKFLOW-EXECUTOR:
  → Busca script via TieredCache
  → Cria V8 sandbox isolada
  → Injeta event, state, variables, nodes
  → Executa script
  → Publica no WORKFLOW-RESUME: {instanceId, nodeId, status: "ok", output: 110.0, statePatch: {calculated_price: 110.0}}

RUNTIME recebe resume:
  → instance.State["calculated_price"] = 110.0 (aplica statePatch)
  → instance.NodeOutputs["code_1"] = 110.0
  → Limpa NodeStates["code_1"]
  → Segue pela edge "out"
```

### Subworkflow: Workflow Filho

```
EXECUTOR retorna NodeState: {waitType: "callback", workflowId: "wf-456", terminationPolicy: "terminate", inputData: {address: "Rua X"}}

RUNTIME:
  → instance.NodeStates["subwf_1"] = {waitType: "callback", ...}
  → instance.Status = waiting
  → KV checkpoint
  → Publica no WORKFLOW-TRIGGER: {workflowId: "wf-456", eventPayload: {address: "Rua X"}, parentInstanceId: "abc", parentNodeId: "subwf_1"}

RUNTIME (outro consumer, mesmo serviço):
  → Cria instância FILHA (parentInstanceId: "abc", parentNodeId: "subwf_1")
  → Executa workflow filho normalmente

WORKFLOW FILHO COMPLETA:
  → Publica no WORKFLOW-RESUME: {instanceId: "abc", nodeId: "subwf_1", status: "ok", output: {valid: true}}

RUNTIME recebe resume do PAI:
  → Aplica outputMappings (filho → pai)
  → instance.NodeOutputs["subwf_1"] = {valid: true}
  → Salva childInstanceId em NodeStates (histórico)
  → Limpa wait do NodeStates["subwf_1"]
  → Segue pela edge "out"
```

### Trigger Event: Publica Evento

```
EXECUTOR retorna NodeState: {waitType: "callback", eventType: "user.created", payload: {name: "John"}}

RUNTIME:
  → instance.NodeStates["trigger_1"] = {waitType: "callback", ...}
  → instance.Status = waiting
  → KV checkpoint
  → Publica no stream do Trigger Service: {instanceId, nodeId, eventType, payload}

TRIGGER SERVICE:
  → Processa evento (pode disparar outros workflows, webhooks, etc.)
  → Publica no WORKFLOW-RESUME: {instanceId, nodeId, status: "ok"}

RUNTIME recebe resume:
  → Limpa NodeStates["trigger_1"]
  → Segue pela edge "out"
```

---

## Checklist de implementação

```
Executors inline (modules/runtime/domain/executors/inline/):
  ✅ start.go
  ✅ end.go
  ✅ condition.go
  ✅ switch_node.go
  ✅ set_state.go
  ✅ log.go
  ✅ goto_node.go

Executors async (modules/runtime/domain/executors/async/):
  ✅ trigger_event.go
  ✅ code.go
  ✅ subworkflow.go
  ✅ delay.go
  ✅ wait_signal.go

Executors control (modules/runtime/domain/executors/control/):
  ✅ fanout.go
  ✅ merge.go
  ✅ sequence.go
  ✅ loop.go
  ✅ wait_for.go

Registry:
  ✅ executor_registry.go — map[string]NodeExecutor (nodeType → executor)

DDD Compliance (concluído):
  ✅ InstanceStateRepository — abstrai NATS KV (domain/repositories + infra/persistence/nats)
  ✅ RuntimePublisherPort — abstrai NATS Publisher (application/ports + infra/messaging/nats)
  ✅ Message types em interfaces/message/types.go (não duplicados no service)
  ✅ GoDoc em todos os tipos/funções/métodos exportados do domain layer
  ✅ GoTo sender validation — error GOTO_NO_RECEIVER para sender sem receiver matching
  ✅ Timezone resolution — WorkflowDefinition → ExecutionGraph → NodeExecutionContext

Serviço novo:
  □ js-workflow-executor (separado do js-executor — ver C8 acima)
```
