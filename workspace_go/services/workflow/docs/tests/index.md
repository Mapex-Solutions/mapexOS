# Tests

## Prerequisites

- Go 1.25+ installed

## Run

```bash
go test ./... -count=1
```

## Test Coverage

### Total: 277 tests across 18 test files

**Engine — Operators (123 tests):**
- `engine/domain/operators/datetime/datetime_operators_test.go` — 46 tests
- `engine/domain/operators/comparison/comparison_operators_test.go` — 33 tests
- `engine/domain/operators/stringops/string_operators_test.go` — 23 tests
- `engine/domain/operators/group/group_operators_test.go` — 21 tests

**Engine — Evaluators (26 tests):**
- `engine/domain/evaluators/condition_evaluator_test.go` — 13 tests
- `engine/domain/evaluators/value_resolver_test.go` — 13 tests

**Definitions (68 tests):**
- `definitions/domain/services/node_validator_test.go` — 39 tests
- `definitions/application/services/definition_service_helpers_test.go` — 19 tests
- `definitions/domain/services/cycle_detector_test.go` — 10 tests

**Runtime — Executors (13 tests):**
- `runtime/domain/executors/inline/inline_executors_test.go` — 7 tests
- `runtime/domain/executors/async/async_executors_test.go` — 5 tests
- `runtime/domain/executors/control/control_executors_test.go` — 5 tests (some share subtests)
- `runtime/domain/executors/executor_registry_test.go` — 1 test

**Runtime — Domain Services (14 tests):**
- `runtime/domain/services/graph_builder_test.go` — 12 tests
- `runtime/domain/services/config_parsing_test.go` — 2 tests

**Archiver (15 tests):**
- `archiver/application/services/archiver_service_test.go` — 15 tests

**Reconciler (8 tests):**
- `reconciler/application/services/reconciler_service_test.go` — 8 tests

**Integration (5 tests):**
- `main_test.go` — 5 tests

### Test Strategy

All tests are **unit tests** with no external dependencies (no NATS, no MongoDB, no I/O).
Mocks are used for ports: `ConditionEvaluatorPort`, `ValueResolverPort`, `KVStore`, `Publisher`,
`ArchiverRepository`.

Each executor test constructs a `NodeExecutionContext` with the appropriate parsed config,
executes the node, and asserts `OutputHandles`, `StatePatch`, `WaitRequest`, and `LogEntries`.

Engine operator tests cover type coercion, edge cases (nil values, empty strings), and
cross-type comparisons (numeric, string, datetime, boolean).

Definition tests cover cycle detection (topological sort), node validation (per-type config
requirements), and service helper functions (diffCodeNodes, contractNodes, extractOrgId).

### Coverage by Module

| Module | Tests | Coverage Focus |
|--------|-------|---------------|
| Engine operators | 123 | Comparison, string, datetime, group operators + type coercion |
| Definitions | 68 | Cycle detection, node validation, service helpers |
| Engine evaluators | 26 | Condition evaluation, value resolution |
| Archiver | 15 | Batch processing, state event routing, KV cleanup |
| Runtime domain | 14 | Graph building, config parsing |
| Runtime executors | 13 | Inline/async/control node execution |
| Reconciler | 8 | Timer sweep, expiration handling |
| Integration | 5 | Main bootstrap validation |
