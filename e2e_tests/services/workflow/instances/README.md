# Module e2e: workflow / instances

## Scope

Saga building blocks for the `instances` module of the workflow
service (`workspace_go/services/workflow/src/modules/instances/`).
This package covers the instance half of the workflow lifecycle —
the runnable binding that pairs a previously created definition
(by id + version) with concrete external inputs and an operator-
chosen execution policy. It pairs the `CreateInstance` step with
the canonical `SagaSimpleInstance` payload (a literal UI capture)
and injects `definitionId`, `definitionVersion`, and
`definitionName` from bag keys published by the sibling
`workflow/definitions` package, then publishes the resulting
instance id so route groups of `kind=workflow` can reference it.
There are no `Test*` functions here: the suite lives in
`journey/iot/connectivity_actions_*/phase1_workflow/` and loads
these steps to drive the workflow service end-to-end.

## Endpoints exercised

- `POST /api/v1/workflow_instances` — register a workflow instance
  bound to an existing definition. The step posts the canonical
  literal payload with `definitionId`, `definitionVersion`,
  `definitionName`, and `name` overridden per run.
- `DELETE /api/v1/workflow_instances/{id}` — invoked from the
  Compensate path of `CreateInstance`; idempotent against `404`.

## Fixtures

The package has no JSON fixtures; the payload body is a Go literal
captured from the platform UI and parameterized per run by the
`payloads/` builder.

| File                                  | Purpose                                                                                  |
|---------------------------------------|------------------------------------------------------------------------------------------|
| `payloads/saga_simple_instance.go`    | Canonical `Device Status` instance — empty `externalInputs`, `isTemplate=false`, `uniqueExecution=false`, no `pathKey` and no `workflowUUID`; `name` is rewritten to `saga-workflow-inst-<runID>` and `definitionId` / `definitionVersion` / `definitionName` are overridden by `CreateInstance` from the bag. |

The `steps/` folder carries `create_instance.go` (the saga step with
its Compensate twin) and `keys.go` (the bag key constant
`BagKeyInstanceID` that the route group `kind=workflow` payload
reads).

## How to run

This package contains no `*_test.go` files; running `go test
./services/workflow/instances` will print `[no test files]`. The
step is consumed by the IoT connectivity-action journeys:

```bash
cd e2e_tests

# Phase 1 of the HTTP and MQTT connectivity-action journeys
go test -tags=saga ./journey/iot/connectivity_actions_http/phase1_workflow
go test -tags=saga ./journey/iot/connectivity_actions_mqtt/phase1_workflow
```

## Outcome on pass

When the consuming journey passes, this package has proven that the
workflow instances module honours `POST /api/v1/workflow_instances`
end-to-end: an instance bound to a real definition + version is
persisted with a `_id` returned in the create response, the saga
can publish that id for a downstream route group of kind=workflow to
reference, and the Compensate `DELETE` removes the record cleanly.

## Requirements

- `workflow` reachable on port `5007` (override via `WORKFLOW_URL`).
- `mapexos` reachable on port `5000` for the admin token bootstrap.
- A workflow definition must already exist on the bag — the
  `CreateInstance` step fails fast if `BagKeyDefinitionID` or
  `BagKeyDefinitionVersion` is missing, so journeys must run
  `definitions.CreateDefinition` first.
- Stack started from `mapexOSDeploy`; the seed admin user
  `admin@mapex.local` and the root organization
  `0000000000000000000aa001` must be present (provisioned by
  `mongodb-init`).

## Notes

- The workflow service distinguishes definition (template) from
  instance (runnable binding). This package owns the instance half
  only; definition lifecycle is the sibling `workflow/definitions/`
  package.
- `definitionName` is kept aligned with the runID-derived name the
  definitions package writes so the operator listing instances for a
  saga run sees a consistent label across the definition + instance
  pair.
- The canonical payload is a verbatim UI capture; a parse error at
  runtime is treated as a developer error (the step panics) rather
  than a saga-recoverable failure.
- The placeholder values in `sagaSimpleInstanceJSON`
  (`definitionId`, `definitionVersion`, `definitionName`) are
  always overridden at runtime — they exist only to keep the
  literal a valid POST body in isolation.
