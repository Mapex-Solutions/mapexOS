# Phase 2: Event pipeline (RabbitMQ trigger)

## What this test proves

RabbitMQ trigger fires from the full telemetry pipeline: a POST to
the gateway is transformed by the js-executor against the asset's
template script, routed by the trigger router, and the RabbitMQ
trigger executor publishes to an ephemeral RabbitMQ container
started via testcontainers-go.

## Outcome on PASS

- Telemetry POSTed to the gateway lands as an event on the asset's
  stream.
- The js-executor runs the template script and emits a standardized
  payload.
- The route group's trigger router selects the RabbitMQ trigger.
- `events_trigger` records at least one successful execution.
- The trigger's last resolved request data contains the saga-scoped
  queue prefix `saga-mq-`.
- Asset teardown cascades cleanly through Compensate.

## Outcome on FAIL

- Gateway rejects the POST (data source API key wrong / not
  provisioned).
- Template script missing or the js-executor cannot run it (no
  standardized payload, router sees nothing).
- Route group not wired to the trigger, or the asset's
  RouteGroupIds drift from the configured trigger RG.
- testcontainers-go cannot pull / start the RabbitMQ image, or
  trigger RabbitMQ config (URI / queue) drifts from the saga-managed
  container.

## How to run

```bash
cd e2e_tests
go test -tags=saga -v ./journey/automations/trigger_rabbitmq/phase2_event_pipeline
```

The entry point is `TestJourney` in `journey_test.go`.

## Composition

`Run` composes phase0 (IAM bootstrap from
`journey/iot/mqtt_broker_auth/phase0_iam_bootstrap`) in front of this
phase's `Items()`. Phase1 is not included — phase2 brings up its own
RabbitMQ container, RabbitMQ trigger, route groups, data source,
template, and connectivity asset, then POSTs the telemetry that
fires the trigger.
