# Go Structs — Entities, Configs e Interfaces

Decisões: B1-B7

---

## Problema

O frontend gera um JSON DSL (nodes, edges, variables) no TypeScript. O backend precisa de structs Go que espelhem esse JSON sem conversão intermediária. O engine (conditionEngine, valueResolver) trabalha com structs nativas — conversão JSON → struct nativa acontece APENAS nas bordas (HTTP handler, NATS handler).

---

## Como Estamos Resolvendo

Structs divididas por domínio:

- **Definitions** (`modules/definitions/domain/entities/`): WorkflowNode, WorkflowEdge, FieldValue, Variables, Conditions, WorkflowDefinition
- **Runtime** (`modules/runtime/domain/entities/`): WorkflowInstance, NodeExecutionContext, NodeExecutionResult, PathEntry, ExecutionError, LogEntry, ExecutionGraph
- **Configs** (`modules/runtime/domain/entities/`): Typed configs para cada node type (17 structs)

---

## Structs de Definition

### WorkflowNode + WorkflowEdge + FieldValue

```go
// ═══════════════════════════════════════════════════
// Pacote: modules/definitions/domain/entities/
// ═══════════════════════════════════════════════════

// ── Node ──────────────────────────────────────────

type WorkflowNode struct {
    ID           string                 `json:"id"           bson:"id"`
    Type         string                 `json:"type"         bson:"type"`         // "core/condition", "core/delay", etc.
    Label        string                 `json:"label"        bson:"label"`
    Position     Position               `json:"position"     bson:"position"`
    Config       map[string]interface{} `json:"config"       bson:"config"`       // Varia por tipo — tipado no executor
    ParentNodeID string                 `json:"parentNodeId" bson:"parentNodeId"` // Para group_frame
}

type Position struct {
    X float64 `json:"x" bson:"x"`
    Y float64 `json:"y" bson:"y"`
}

// ── Edge ──────────────────────────────────────────

type WorkflowEdge struct {
    ID           string  `json:"id"           bson:"id"`
    Source       string  `json:"source"       bson:"source"`
    SourceHandle string  `json:"sourceHandle" bson:"sourceHandle"` // "out", "true", "false", "step_1", "case_xxx", "default", "body", "done", "matched"
    Target       string  `json:"target"       bson:"target"`
    TargetHandle string  `json:"targetHandle" bson:"targetHandle"`
    Label        string  `json:"label"        bson:"label"`
    PathOffsetX  float64 `json:"pathOffsetX"  bson:"pathOffsetX"`  // Visual offset
    PathOffsetY  float64 `json:"pathOffsetY"  bson:"pathOffsetY"`
}

// ── FieldValue (usado em conditions, set_state, mappings) ──

type FieldValueType string

const (
    FieldValueEvent      FieldValueType = "event"       // Acessa EventPayload
    FieldValueState      FieldValueType = "state"        // Acessa State (variáveis do usuário)
    FieldValueVariable   FieldValueType = "variable"     // Acessa Variables (mesma resolução que state)
    FieldValueLiteral    FieldValueType = "literal"      // Valor estático
    FieldValueNodeOutput FieldValueType = "node_output"  // Acessa NodeOutputs[nodeId]
)
// Loop item/index são injetados no State do usuário via StatePatch.

type FieldValue struct {
    Type   FieldValueType `json:"type"             bson:"type"`
    Value  string         `json:"value"            bson:"value"`            // Literal ou path (state.myVar)
    Mode   string         `json:"mode,omitempty"   bson:"mode,omitempty"`   // "dynamic"|"manual" (só type=event)
    NodeID string         `json:"nodeId,omitempty" bson:"nodeId,omitempty"` // Só type=node_output
}
```

### Variables + CaptureFields

```go
type VariableType string

const (
    VarTypeString  VariableType = "string"
    VarTypeNumber  VariableType = "number"
    VarTypeBoolean VariableType = "boolean"
    VarTypeJSON    VariableType = "json"
)

type WorkflowVariable struct {
    Field        string       `json:"field"        bson:"field"`
    Type         VariableType `json:"type"         bson:"type"`
    DefaultValue interface{}  `json:"defaultValue" bson:"defaultValue"`
    Description  string       `json:"description"  bson:"description"`
    Durable      bool         `json:"durable"      bson:"durable"`      // Persiste entre runs
}

type CaptureField struct {
    Field       string       `json:"field"       bson:"field"`
    Type        VariableType `json:"type"        bson:"type"`
    Description string       `json:"description" bson:"description"`
}
```

### Condition System

```go
type GroupLogicOperator string

const (
    LogicAND  GroupLogicOperator = "AND"
    LogicOR   GroupLogicOperator = "OR"
    LogicNAND GroupLogicOperator = "NAND"
    LogicNOR  GroupLogicOperator = "NOR"
)

type ConditionItem struct {
    ID       string     `json:"id"       bson:"id"`
    Name     string     `json:"name"     bson:"name"`
    Field    FieldValue `json:"field"    bson:"field"`    // Lado esquerdo
    Operator string     `json:"operator" bson:"operator"` // "equals", "contains", "greater_than", etc.
    Value    FieldValue `json:"value"    bson:"value"`    // Lado direito
}

type ConditionGroupItem struct {
    Type string      `json:"type" bson:"type"` // "condition" ou "group"
    Data interface{} `json:"data" bson:"data"` // ConditionItem ou ConditionGroup
}

type ConditionGroup struct {
    ID    string               `json:"id"    bson:"id"`
    Name  string               `json:"name"  bson:"name"`
    Logic GroupLogicOperator   `json:"logic" bson:"logic"`
    Items []ConditionGroupItem `json:"items" bson:"items"`
}

type SwitchCase struct {
    ID        string         `json:"id"        bson:"id"`        // Também é o handle ID
    Name      string         `json:"name"      bson:"name"`
    Condition ConditionGroup `json:"condition" bson:"condition"`
}

type RetryPolicy struct {
    Enabled           bool     `json:"enabled"           bson:"enabled"`
    MaxAttempts       int      `json:"maxAttempts"       bson:"maxAttempts"`
    InitialInterval   string   `json:"initialInterval"   bson:"initialInterval"`   // "1s", "5m"
    BackoffMultiplier float64  `json:"backoffMultiplier" bson:"backoffMultiplier"`
    MaxInterval       string   `json:"maxInterval"       bson:"maxInterval"`       // "5m", "1h"
    NonRetryableErrors []string `json:"nonRetryableErrors" bson:"nonRetryableErrors"`
}
```

### WorkflowDefinition

```go
type WorkflowDefinition struct {
    ID                model.ObjectId          `json:"_id"                bson:"_id"`
    OrgID             *model.ObjectId         `json:"orgId"              bson:"orgId"`
    Name              string                  `json:"name"               bson:"name"`
    Description       string                  `json:"description"        bson:"description"`
    Enabled           bool                    `json:"enabled"            bson:"enabled"`
    IsTemplate        bool                    `json:"isTemplate"         bson:"isTemplate"`
    DefinitionVersion int                     `json:"definitionVersion"  bson:"definitionVersion"`
    Timezone          FieldValue              `json:"timezone"           bson:"timezone"`
    RetryPolicy       RetryPolicy             `json:"retryPolicy"       bson:"retryPolicy"`
    Variables         []WorkflowVariable      `json:"variables"          bson:"variables"`
    CaptureFields     []CaptureField          `json:"captureFields"      bson:"captureFields"`
    Nodes             []WorkflowNode          `json:"nodes"              bson:"nodes"`
    Edges             []WorkflowEdge          `json:"edges"              bson:"edges"`
    Metadata          DefinitionMetadata      `json:"metadata"           bson:"metadata"`
    PathKey           string                  `json:"pathKey"            bson:"pathKey"`
    Scope             string                  `json:"scope"              bson:"scope"`
    Created           time.Time               `json:"created"            bson:"created"`
    Updated           time.Time               `json:"updated"            bson:"updated"`
}

type DefinitionMetadata struct {
    CanvasViewport CanvasViewport `json:"canvasViewport" bson:"canvasViewport"`
}

type CanvasViewport struct {
    X    float64 `json:"x"    bson:"x"`
    Y    float64 `json:"y"    bson:"y"`
    Zoom float64 `json:"zoom" bson:"zoom"`
}
```

---

## Structs de Runtime

### WorkflowInstance

```go
// ═══════════════════════════════════════════════════
// Pacote: modules/runtime/domain/entities/
// ═══════════════════════════════════════════════════

type InstanceStatus string

const (
    StatusCreated   InstanceStatus = "created"
    StatusRunning   InstanceStatus = "running"
    StatusWaiting   InstanceStatus = "waiting"
    StatusCompleted InstanceStatus = "completed"
    StatusFailed    InstanceStatus = "failed"
    StatusCancelled InstanceStatus = "cancelled"
)

type WorkflowInstance struct {
    ID                model.ObjectId                       `json:"_id"               bson:"_id"`
    WorkflowID        model.ObjectId                       `json:"workflowId"        bson:"workflowId"`
    WorkflowName      string                               `json:"workflowName"      bson:"workflowName"`
    OrgID             *model.ObjectId                      `json:"orgId"             bson:"orgId"`
    DefinitionVersion int                                  `json:"definitionVersion" bson:"definitionVersion"`
    Version           int                                  `json:"version"           bson:"version"`            // CAS (NATS KV revision)
    Status            InstanceStatus                       `json:"status"            bson:"status"`

    // ── Execution tracking ──────────────────────
    ActiveNodeIDs     []string                             `json:"activeNodeIds"     bson:"activeNodeIds"`      // Nodes ativos (1 em linear, N em fanout)
    NodeStates        map[string]map[string]interface{}    `json:"nodeStates"        bson:"nodeStates"`         // Estado interno por node (loop counter, wait info, etc.)
    State             map[string]interface{}               `json:"state"             bson:"state"`              // Variáveis do USUÁRIO (set_state, loop_item, etc.)
    EventPayload      map[string]interface{}               `json:"eventPayload"      bson:"eventPayload"`       // Payload do trigger original
    NodeOutputs       map[string]interface{}               `json:"nodeOutputs"       bson:"nodeOutputs"`        // Outputs por nodeId (code result, subworkflow result)
    ExecutionPath     []PathEntry                          `json:"executionPath"     bson:"executionPath"`      // Histórico de execução (DAG visualization)

    // ── Definition data (copiados na criação) ───
    Variables         []WorkflowVariable                   `json:"variables"         bson:"variables"`          // Definição das variáveis (types, defaults, descriptions)
    ExternalInputs    map[string]interface{}               `json:"externalInputs"    bson:"externalInputs"`     // Inputs externos configurados na definition
    CaptureFields     []CaptureField                       `json:"captureFields"     bson:"captureFields"`      // Campos a capturar do evento

    // ── State defaults (variable defaults for code execution) ──
    StateDefaults     map[string]interface{}               `json:"stateDefaults"     bson:"stateDefaults"`

    // ── Error ───────────────────────────────────
    ErrorInfo         *ExecutionError                      `json:"errorInfo"         bson:"errorInfo"`

    // ── Subworkflow context ─────────────────────
    ParentInstanceID  *model.ObjectId                      `json:"parentInstanceId"  bson:"parentInstanceId"`   // ID da instância PAI (só em filhos)
    ParentNodeID      string                               `json:"parentNodeId"      bson:"parentNodeId"`       // NodeID no pai que é o subworkflow node
    Depth             int                                  `json:"depth"             bson:"depth"`              // Subworkflow recursion depth

    // ── Timestamps ──────────────────────────────
    StartedAt         *time.Time                           `json:"startedAt"         bson:"startedAt"`
    CompletedAt       *time.Time                           `json:"completedAt"       bson:"completedAt"`
    Created           time.Time                            `json:"created"           bson:"created"`
    Updated           time.Time                            `json:"updated"           bson:"updated"`
}
```

### NodeExecutionContext + NodeExecutionResult

```go
// Input para TODOS os 17 executors (readonly — executor NÃO modifica)
type NodeExecutionContext struct {
    // Instance state (readonly)
    InstanceID     model.ObjectId
    State          map[string]interface{}                // Variáveis do usuário
    EventPayload   map[string]interface{}                // Payload do trigger
    NodeOutputs    map[string]interface{}                // Outputs de nodes anteriores
    NodeStates     map[string]map[string]interface{}     // Estado interno de TODOS os nodes
    ExternalInputs map[string]interface{}                // Inputs externos configurados na definition
    Depth          int                                   // Subworkflow recursion depth

    // Node info
    NodeID       string                                // ID do node sendo executado
    NodeType     string                                // "core/condition", "core/delay", etc.
    ParsedConfig interface{}                           // Typed config struct (parsed pelo GraphBuilder)
    Label        string                                // Label do node

    // Graph
    Graph *ExecutionGraph                              // Para resolver output handles e adjacência

    // Timezone
    Timezone string                                    // Resolved from definition
}

// Output de TODOS os 17 executors (comandos para o runtime)
type NodeExecutionResult struct {
    // Control flow
    OutputHandles []string                             // Handles para seguir no grafo (ver 09-executors.md)

    // State mutations (runtime aplica)
    StatePatch    map[string]interface{}                // Delta para merge no instance.State (variáveis do USUÁRIO)
    NodeState     map[string]interface{}                // Estado deste node → instance.NodeStates[nodeId]

    // Output data
    NodeOutput    interface{}                           // Output do node (code result, subworkflow result)

    // Logs
    LogEntries    []LogEntry                           // Step logs → WORKFLOW-LOGS stream

    // Error
    Error         *ExecutionError                      // Erro estruturado (end node com terminateWithError)
}
```

**Como o runtime detecta suspensão:**
```
Se result.NodeState != nil && result.NodeState["waitType"] != nil:
    → Suspende instância
    → Despacha baseado no node.Type (ver 09-executors.md)
Se result.NodeState == nil ou sem waitType:
    → Execução continua normalmente pelo OutputHandles
```

### Supporting Types

```go
// ── PathEntry (visualização do DAG — "Event History" do MapexOS) ──

type PathEntry struct {
    NodeID       string     `json:"nodeId"                 bson:"nodeId"`
    NodeType     string     `json:"nodeType"               bson:"nodeType"`
    Status       string     `json:"status"                 bson:"status"`       // completed|error|waiting
    EnteredAt    time.Time  `json:"enteredAt"              bson:"enteredAt"`
    ExitedAt     *time.Time `json:"exitedAt,omitempty"     bson:"exitedAt,omitempty"`
    DurationMs   int64      `json:"durationMs"             bson:"durationMs"`
    OutputHandle string     `json:"outputHandle,omitempty" bson:"outputHandle,omitempty"`
    Error        *string    `json:"error,omitempty"        bson:"error,omitempty"`
}

// ── ExecutionError ───────────────────────────────

type ExecutionError struct {
    Code       string    `json:"code"       bson:"code"`       // "TIMEOUT", "SCRIPT_ERROR", "CONDITION_ERROR"
    Message    string    `json:"message"    bson:"message"`
    NodeID     string    `json:"nodeId"     bson:"nodeId"`
    NodeType   string    `json:"nodeType"   bson:"nodeType"`
    Timestamp  time.Time `json:"timestamp"  bson:"timestamp"`
    StackTrace string    `json:"stackTrace" bson:"stackTrace"` // Opcional (code node errors)
}

// ── LogEntry (step log → ClickHouse via WORKFLOW-LOGS) ──

type LogLevel string

const (
    LogDebug LogLevel = "debug"
    LogInfo  LogLevel = "info"
    LogWarn  LogLevel = "warn"
    LogError LogLevel = "error"
)

type LogEntry struct {
    Level      LogLevel               `json:"level"`
    Message    string                 `json:"message"`
    Timestamp  time.Time              `json:"timestamp"`
    NodeID     string                 `json:"nodeId"`
    NodeType   string                 `json:"nodeType"`
    Data       map[string]interface{} `json:"data,omitempty"`
}

// ── ExecutionGraph (construído da definition pelo BuildGraph) ──

type ExecutionGraph struct {
    Adjacency     map[string]map[string]string    // nodeId → handleId → targetNodeId
    Nodes         map[string]*WorkflowNode        // lookup rápido por ID
    GotoPairs     map[string]string               // pairLabel → receiverNodeId
    ParsedConfigs map[string]interface{}           // nodeId → typed config struct
    Timezone      string                          // Resolved from WorkflowDefinition.Timezone (literal value)
}
```

### Node Executor Interface

```go
// ═══════════════════════════════════════════════════
// Pacote: modules/runtime/domain/entities/
// ═══════════════════════════════════════════════════

type NodeExecutor interface {
    Execute(ctx context.Context, execCtx *NodeExecutionContext) (*NodeExecutionResult, error)
    NodeType() string // "core/condition", "core/delay", etc.
}
```

---

## Typed Configs (por node type)

Cada executor recebe `ParsedConfig interface{}` que é a struct tipada do seu tipo.
Conversão `Config map[string]interface{}` → struct tipada acontece no `BuildGraph`.

```go
// ═══════════════════════════════════════════════════
// Pacote: modules/runtime/domain/entities/
// ═══════════════════════════════════════════════════

type ConditionNodeConfig struct {
    Condition           ConditionGroup  `json:"condition"`
    SelectedTemplateIds []string        `json:"selectedTemplateIds"`
}

type SwitchNodeConfig struct {
    Cases               []SwitchCase `json:"cases"`
    MatchMode           string       `json:"matchMode"` // "first" | "all"
    SelectedTemplateIds []string     `json:"selectedTemplateIds"`
}

type SetStateNodeConfig struct {
    Operation           string     `json:"operation"`   // "set"|"increment"|"decrement"|"remove"|"append"
    TargetField         string     `json:"targetField"`
    ValueSource         FieldValue `json:"valueSource"`
    SelectedTemplateIds []string   `json:"selectedTemplateIds"`
}

type LogNodeConfig struct {
    Message string `json:"message"`
    Level   string `json:"level"` // "info"|"warn"|"error"|"debug"
}

type CodeNodeConfig struct {
    Script  string `json:"script"`
    Timeout int    `json:"timeout"` // seconds
}

type DelayNodeConfig struct {
    Duration int    `json:"duration"`
    Unit     string `json:"unit"` // "s"|"m"|"h"|"d" ou "seconds"|"minutes"|"hours"|"days"
}

type WaitSignalNodeConfig struct {
    SignalName       string            `json:"signalName"`
    Timeout          string            `json:"timeout"`          // "10m", "30s"
    MaxTimeoutCycles int               `json:"maxTimeoutCycles"`
    Mappings         []SignalMapping   `json:"mappings"`
}

type SignalMapping struct {
    ParamName string     `json:"paramName"`
    Value     FieldValue `json:"value"`
}

type WaitForNodeConfig struct {
    Field            string     `json:"field"`
    Operator         string     `json:"operator"`
    CompareTo        FieldValue `json:"compareTo"`
    Timeout          string     `json:"timeout"`          // "5m", "1h" — Reconciler timeout
    MaxTimeoutCycles int        `json:"maxTimeoutCycles"`
}
// Nota: wait_for é EVENT-DRIVEN (re-avalia em cada interação), NÃO polling.

type FanoutNodeConfig struct {
    Branches int `json:"branches"` // Gera out_1, out_2, ... (max 20)
}

type MergeNodeConfig struct {
    Branches int    `json:"branches"` // Espera N branches
    Strategy string `json:"strategy"` // "all"|"any"|"first"
}

type SequenceNodeConfig struct {
    Steps int `json:"steps"` // Gera step_1, step_2, ...
}

type LoopNodeConfig struct {
    Source FieldValue `json:"source"` // Array para iterar (resolved via ValueResolver)
}

type SubworkflowNodeConfig struct {
    WorkflowID     string            `json:"workflowId"`
    WorkflowName   string            `json:"workflowName"`
    ExecutionMode  string            `json:"executionMode"` // "sync" (espera) | "fire_and_forget"
    Timeout        TimeoutConfig     `json:"timeout"`
    InputMappings  []InputMapping    `json:"inputMappings"`
    OutputMappings []OutputMapping   `json:"outputMappings"`
}
// FUTURO: TerminationPolicy (Temporal.io inspired) — "terminate"|"abandon"|"request_cancel"

type TimeoutConfig struct {
    Duration int    `json:"duration"`
    Unit     string `json:"unit"` // "s"|"m"|"h"|"d"
}

type InputMapping struct {
    ChildParamName string     `json:"childParamName"`
    Value          FieldValue `json:"value"`
}

type OutputMapping struct {
    OutputName string `json:"outputName"`
    StateField string `json:"stateField"`
}

type TriggerEventNodeConfig struct {
    EventType      string                `json:"eventType"`
    PayloadMapping []TriggerPayloadField `json:"payloadMapping"`
}

type TriggerPayloadField struct {
    Key   string     `json:"key"`
    Value FieldValue `json:"value"`
}

type EndNodeConfig struct {
    TerminateWithError bool       `json:"terminateWithError"`
    ErrorCode          string     `json:"errorCode"`
    ErrorMessage       FieldValue `json:"errorMessage"`
}

type GotoNodeConfig struct {
    Role      string `json:"role"`      // "sender"|"receiver"
    PairLabel string `json:"pairLabel"` // Match sender↔receiver
    PairColor string `json:"pairColor"` // Visual only
}
```
