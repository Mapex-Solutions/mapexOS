# Module e2e: workflow / definitions

## Scope

Saga building blocks for the `definitions` module of the workflow
service (`workspace_go/services/workflow/src/modules/definitions/`).
This package covers the definition half of the workflow lifecycle —
the static DAG description an operator authors in the canvas (nodes,
edges, states, retry policy, external inputs and signals). It pairs
the `CreateDefinition` step with the canonical `SagaSimpleDefinition`
payload (a literal UI capture of a Start → Set State → Code → End
flow) and publishes the returned id + version on the bag so the
sibling `instances` package can bind an instance to the freshly
created definition. There are no `Test*` functions here: the suite
lives in `journey/iot/connectivity_actions_*/phase1_workflow/` and
loads these steps to drive the workflow service end-to-end.

## Endpoints exercised

- `POST /api/v1/workflow_definitions` — register a workflow
  definition. The step posts the canonical literal payload and only
  parameterizes `name` on `runID` so concurrent saga runs do not
  collide on Mongo's name-unique index.
- `DELETE /api/v1/workflow_definitions/{id}` — invoked from the
  Compensate path of `CreateDefinition`; idempotent against `404`.

## Fixtures

The package has no JSON fixtures; the payload body is a Go literal
captured from the platform UI (DevTools → Network) and parameterized
per run by the `payloads/` builder.

| File                                    | Purpose                                                                                  |
|-----------------------------------------|------------------------------------------------------------------------------------------|
| `payloads/saga_simple_definition.go`    | Canonical `Device Status` definition — Start → Set State → Code → End, one state field, one external signal, and one installed plugin (`telegram`); name is rewritten to `saga-workflow-def-<runID>`. |

The `steps/` folder carries `create_definition.go` (the saga step
with its Compensate twin) and `keys.go` (the bag key constants
`BagKeyDefinitionID` and `BagKeyDefinitionVersion` the instances
package reads to bind its create payload).

## How to run

This package contains no `*_test.go` files; running `go test
./services/workflow/definitions` will print `[no test files]`. The
step is consumed by the IoT connectivity-action journeys:

```bash
cd e2e_tests

# Phase 1 of the HTTP and MQTT connectivity-action journeys
go test -tags=saga ./journey/iot/connectivity_actions_http/phase1_workflow
go test -tags=saga ./journey/iot/connectivity_actions_mqtt/phase1_workflow
```

## Outcome on pass

When the consuming journey passes, this package has proven that the
workflow definitions module honours `POST /api/v1/workflow_definitions`
end-to-end for the smallest realistic DAG the runtime accepts: the
definition is persisted with a `_id` and a `definitionVersion`, both
returned in the create response, and the Compensate `DELETE` removes
the record cleanly.

## Requirements

- `workflow` reachable on port `5007` (override via `WORKFLOW_URL`).
- `mapexos` reachable on port `5000` for the admin token bootstrap.
- Stack started from `mapexOSDeploy`; the seed admin user
  `admin@mapex.local` and the root organization
  `0000000000000000000aa001` must be present (provisioned by
  `mongodb-init`).

## Notes

- The workflow service distinguishes definition (template) from
  instance (runnable binding). This package owns the definition half
  only; instance lifecycle is the sibling `workflow/instances/`
  package.
- `definitionVersion` is published on the bag as an `int` and the
  instances `CreateInstance` step type-asserts it back — keep the
  type stable when changing the response decoder.
- The canonical payload is a verbatim UI capture so the request shape
  stays in lockstep with what the operator sends from the canvas; a
  parse error at runtime is treated as a developer error (the step
  panics) rather than a saga-recoverable failure.
- The payload references one installed plugin (`telegram`); the
  workflow service must have the plugin registered or the create call
  will fail validation.
