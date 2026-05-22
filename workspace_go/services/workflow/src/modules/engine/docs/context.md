# Bounded Context: Engine

**Service:** workflow
**Module path:** `src/modules/engine/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-04-29

## Purpose

Pure-computation module for condition evaluation and field-value resolution. Provides the two ports that Runtime executors and other modules call to: (a) evaluate a `ConditionGroup` tree against the current execution context (event payload, state, node outputs, external inputs) with short-circuit logic; (b) resolve a `FieldValue` (whose type may be `event`, `state`, `variable`, `input`, `literal`, or `nodeOutput`) to an actual value. Stateless, no I/O, no repositories, no NATS, no HTTP — only an `OperatorRegistry` initialised once at startup with 22 built-in operators.

## Module Layout (intentional §3 exemption)

This module exposes pure-domain evaluators directly, not a wrapping service. `application/services/engine_service.go` is a factory: it builds the `OperatorRegistry`, registers the 22 built-in operators, and returns the two port interfaces backed by `evaluators.ConditionEvaluator` and `evaluators.ValueResolver`. The port methods (`EvaluateGroup`, `Resolve`, `BuildDescription`) live where they conceptually belong — in `domain/evaluators/` — because they are stateless pure-computation primitives with no orchestration to surface in service.go. Wrapping them in an `EngineService` struct that only forwards calls would produce §3 anti-pattern A (1-line delegations), which is explicitly forbidden. Compile-time port checks (`var _ ports.ConditionEvaluatorPort = (*evaluators.ConditionEvaluator)(nil)`) in service.go cover the safety concern that an aggregator struct would otherwise enforce.

## Ubiquitous Language

| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Operator | Pluggable comparison/string/datetime/group predicate registered in the `OperatorRegistry` | Node operation (runtime plugin action) |
| ConditionEvaluator | Recursive evaluator that walks `ConditionGroup` tree with short-circuit AND/OR/NAND/NOR | Condition node executor (the walker-level wrapper in runtime) |
| ValueResolver | Resolves a `FieldValue` against runtime context maps | Plugin-executor templating (different namespace shape) |
| Between operator | Operator category with dedicated registry slot (registered via `RegisterBetween` alongside `RegisterCondition`) | Regular binary comparison |

## Published Events (driven — outbound)

_None._

## Consumed Events (driving — inbound)

_None._

## Driving Ports (what can call this module)

- Cross-module Go API only:
  - `ports.ConditionEvaluatorPort.EvaluateGroup(...)`
  - `ports.ValueResolverPort.Resolve(...)` + `BuildDescription(...)`
- Primary caller: runtime executors (`condition`, `switch`, `wait_for`, fanout filters).

## Driven Ports (what this module requires)

- Reads `definitions/domain/entities.ConditionGroup` and `FieldValue` types (value types only — no behavior coupling).
- Internal only: `OperatorRegistry` (in-memory, immutable after startup).

## Invariants and Business Rules

- `OperatorRegistry` is populated once at DI construction (`New`) and MUST be treated as immutable at runtime.
- Evaluation is pure: given the same inputs, output is deterministic; no side effects.
- Short-circuit semantics: AND stops on first false, OR on first true; NAND/NOR are negations of AND/OR respectively.
- `Between` operators are registered in two slots (condition + between) so both call sites resolve correctly.
- 22 operators total: 7 comparison (eq/neq/gt/gte/lt/lte/between), 5 string (contains/notContains/startsWith/endsWith/regex), 6 datetime (beforeDate/afterDate/beforeTime/afterTime/betweenDate/betweenTime), 4 group (and/or/nand/nor).

## Known Cross-Context Interactions

- Consumed by **runtime** (condition/switch/wait_for executors + walker-level filters) via DI-injected ports.
- Reads value types from **definitions** module (no cross-module entity import — only shared value types).

## Literal-Source Template Interpolation

When a `FieldValue` of type `literal` has a `value` containing `{{namespace.path}}` placeholders, the resolver interpolates them against the four runtime context maps. Workflows whose literal values contain no `{{` substring are byte-identical to the previous behavior thanks to a `strings.Contains` short-circuit at the top of the literal branch in `value_resolver.go`. The interpolation implementation lives in `value_resolver_template.go`.

| Namespace | Maps to | Example |
|-----------|---------|---------|
| `event` | The triggering event payload (`eventPayload`) | `{{event.user.name}}` |
| `state` | Workflow state variables (`state`) | `{{state.counter}}` |
| `input` | Workflow external inputs (`externalInputs`) | `{{input.threshold}}` |
| `output` | Output of a previous node, keyed by node ID (`nodeOutputs`) | `{{output.<nodeId>.field}}` |

### Semantics

- **Missing path** → empty string at the placeholder position. (Differs from the per-source `Resolve` branches which return `ErrFieldNotFound`.)
- **Scalar (number, bool, etc.)** → `fmt.Sprintf("%v", v)` stringification.
- **Object / slice** → `encoding/json.Marshal`. A `logger.Warn` line is emitted; this is intentional so authors notice unintended object stringification. Marshal failure → empty string.
- **Malformed `{{` without `}}`** → verbatim return, no error.
- **Array indexing** → use `0`-based integer segments inside the dot-path: `{{input.recipients.0}}` or `{{output.<nodeId>.items.2.name}}`. Out-of-range indices and non-integer segments at array positions resolve to empty string.
- Templates are **best-effort**: `renderTemplate` never returns an error. Any failure inside the interpolation degrades to empty string or verbatim.

### v1 Limitations

- Escape syntax for literal `{{` characters is NOT supported. Authors cannot output the `{{` character pair.
- No filters, pipes, or function calls (e.g. `{{event.name | uppercase}}`).

### Distinct from plugin executor

The plugin executor at `runtime/domain/executors/async/plugin.go` uses different namespace keys (`wf.state`, `wf.input`, `event`, `manifest`, `credentials`, `config`) for back-compat with all existing plugin manifests. The two resolvers coexist; the literal-template namespaces (`event/state/input/output`) are deliberately UI-aligned with the FieldSourceSelector component so authors write what they see.

### Output Catalog — which node types populate `output.<nodeId>`

The `output` namespace is keyed by node ID. Whether `{{output.<nodeId>.<field>}}` resolves to anything depends on whether the node type that produced the output emits a `NodeOutput`. The truth table comes from the runtime module: `runtime_handler_walker.go:123-124` writes synchronous outputs into `execution.NodeOutputs[currentNodeID]`, and `runtime_handler_resume.go:271-275` writes async-resumed outputs (callback or signal payload).

Synchronous executors that emit `NodeOutput`:

| Node type | Output shape | Source |
|-----------|--------------|--------|
| `loop` | `{ item, index }` per iteration | `runtime/domain/executors/control/loop.go:93` |

Asynchronous executors — output captured when the node resumes:

| Node type | Output shape | Source |
|-----------|--------------|--------|
| `code` | Whatever the JS script returns | `runtime/domain/executors/async/code.go` (resumed via callback) |
| `subworkflow` | The End-node payload of the sub-workflow | `runtime/domain/executors/async/subworkflow.go` |
| `plugin` (any marketplace plugin node — Telegram, HTTP, etc.) | The plugin action's HTTP response or callback payload | `runtime/domain/executors/async/plugin.go` |
| `wait_signal` | The signal payload (`SignalData`) | resumed via `runtime_handler_resume.go:275` |
| `wait_for` | The condition-resolved payload | resumed via `runtime_handler_resume.go:271` |
| `trigger_event` | The downstream triggered execution's output (if any) | `runtime/domain/executors/async/trigger_event.go` |

Node types that do NOT emit `NodeOutput` (so `{{output.<nodeId>.x}}` against them resolves to empty):

`start`, `end`, `condition`, `switch`, `set_state`, `goto`, `log`, `fanout`, `merge`, `sequence`.

Authoring tips:

- For `loop`, the safe references are `{{output.<loopNodeId>.item}}` and `{{output.<loopNodeId>.index}}` — both are always present once the loop body executes.
- For `code`, the keys depend on what the script returned. There is no static schema; refer to the script body when authoring templates.
- For `plugin` nodes, the response shape comes from the plugin manifest's action definition. Inspect the plugin's manifest (or the marketplace docs) to know which keys exist.
- For `subworkflow`, the keys are whatever the sub-workflow's End node assigned. Confirm by opening the sub-workflow's End-node config.
