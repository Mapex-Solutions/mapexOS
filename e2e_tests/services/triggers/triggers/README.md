# Module e2e: triggers / triggers

## Scope

Saga building blocks for the `triggers` module of the triggers service
(`workspace_go/services/triggers/src/modules/triggers/`). This package
exercises the public CRUD contract of triggers across the eight kinds
supported by the platform — HTTP, Email (SMTP), WebSocket, Slack,
Teams, MQTT, NATS, and RabbitMQ — by composing one create step per
kind together with an in-process sink, broker, or container the saga
journeys stand up to receive the trigger output. There are no
`Test*` functions here: the suite lives in `journey/automations/` and
loads these steps, payload builders, and asserts to drive the trigger
service end-to-end.

## Endpoints exercised

- `POST /api/v1/triggers` — register a trigger of any of the eight
  supported kinds. Each `Create*Trigger` step posts a verbatim UI
  capture and parameterizes only the name and the destination
  (endpoint URL, broker host/port, server URL, etc.) so the trigger
  fires against the saga-managed sink.
- `DELETE /api/v1/triggers/{id}` — invoked from the Compensate path of
  every create step; idempotent against `404`.

## Fixtures

The package has no JSON fixtures; payload bodies are Go literals
captured from the platform UI and parameterized per run by the
`payloads/` builders.

| File                              | Purpose                                                                                              |
|-----------------------------------|-------------------------------------------------------------------------------------------------------|
| `payloads/saga_simple_trigger.go` | Generic HTTP trigger; endpoint rewritten to `constants.TriggerSinkURL`.                               |
| `payloads/saga_email_trigger.go`  | Email trigger pointing at the in-process SMTP sink.                                                   |
| `payloads/saga_websocket_trigger.go` | WebSocket trigger targeting the in-process WS upgrade server.                                      |
| `payloads/saga_slack_trigger.go`  | Slack trigger; webhook URL rewritten to the HTTP sink so the post is captured locally.                |
| `payloads/saga_teams_trigger.go`  | Microsoft Teams trigger; webhook URL rewritten to the HTTP sink.                                      |
| `payloads/saga_mqtt_trigger.go`   | MQTT trigger pointing at the embedded mochi-mqtt broker the saga starts.                              |
| `payloads/saga_nats_trigger.go`   | NATS trigger pointing at the embedded NATS server the saga starts.                                    |
| `payloads/saga_rabbitmq_trigger.go` | RabbitMQ trigger pointing at the ephemeral RabbitMQ container the saga starts (testcontainers).     |

The `steps/` folder mirrors one `Create*Trigger` per payload plus the
sink/broker/container lifecycle steps (`start_test_sink.go`,
`start_smtp_sink.go`, `start_websocket_sink.go`, `start_mqtt_broker.go`,
`start_nats_server.go`, `start_rabbitmq_container.go`) and
`compensate_helpers.go` for the shared teardown logic.

The `asserts/` folder carries `assert_sink_hit.go` (HTTP, WebSocket,
generic counters) and `assert_smtp_received.go` (subject/to/from/body
validation without an external MX parser).

## How to run

This package contains no `*_test.go` files; running `go test
./services/triggers/triggers` will print `[no test files]`. The steps
are consumed by the automation journeys under
`journey/automations/trigger_*/`:

```bash
cd e2e_tests

# All eight trigger journeys
go test -tags=saga ./journey/automations/...

# A single trigger kind end-to-end
go test -tags=saga ./journey/automations/trigger_http/...
go test -tags=saga ./journey/automations/trigger_email/...
```

## Outcome on pass

When the consuming journeys pass, this package has proven that the
triggers module honours `POST /api/v1/triggers` (and the Compensate
`DELETE`) for every supported kind, and that the resolved trigger
configuration actually fires against a live destination — observed
through an in-process sink, broker, or container rather than the
events service round-trip.

## Requirements

- `triggers` reachable on port `5006` (override via `TRIGGERS_URL`).
- `mapexos` reachable on port `5000` for the admin token bootstrap.
- Ports `11010` (HTTP sink), `11025` (SMTP sink), `11026` (WebSocket
  sink) free on the host; the broker/server/container steps bind on
  OS-assigned ports.
- Docker available for the RabbitMQ testcontainers step.

## Notes

- Each trigger kind has its own bag key namespace (see `steps/keys.go`)
  so a journey may compose several triggers in the same run without
  collisions.
- The HTTP, Slack, and Teams triggers share the same generic HTTP sink
  — the assert just compares the captured body against the kind's
  expected webhook envelope.
- The embedded NATS / mochi-mqtt servers and the RabbitMQ container
  are torn down by the saga Compensate path; failed runs may leave
  containers behind that `docker ps` will reveal.
- Eight trigger kinds are exercised inline by the create steps; new
  kinds added to the module require a new `create_<kind>_trigger.go` +
  `saga_<kind>_trigger.go` pair, not a change to existing files.
