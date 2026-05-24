# Module e2e: events / events

## Scope

This folder is reserved for the module e2e suite of the `events` module
of the events service
(`workspace_go/services/events/src/modules/events/`), but the test
package itself has not been written yet. What lives here today is the
shared `asserts/` package: saga oracles that query the events service's
public HTTP API on behalf of cross-module journeys (IoT pipelines,
trigger smokes, workflow event verification). Module-scoped CRUD tests
for the three event listings will land alongside these asserts in a
future iteration.

## Endpoints exercised

Indirectly — only via saga journeys, through the oracles in `asserts/`:

- `GET /api/v1/events/raw` — list raw ingest records (filtered by
  `threadId` + `startTime`); used by `AssertRawEventReceivedAfter` to
  prove gateway ingestion reached the events store.
- `GET /api/v1/events/trigger` — list trigger execution records
  (filtered by `triggerId`); used by `AssertTriggerEventReceivedAfter`,
  `AssertTriggerExecutedSuccessfullyEventually`, and
  `AssertLastTriggerRequestDataContains` to verify trigger delivery
  and inspect the resolved `requestData` payload.
- `GET /api/v1/events/workflow` — list workflow execution records
  (filtered by `instanceId` + `startTime`); used by
  `AssertWorkflowEventReceivedAfter` to confirm a workflow ran for the
  asset whose telemetry kicked it off.

## Test functions

None at the module-e2e layer. The asserts in this folder are consumed
by saga journeys (run with the `saga` build tag), not by `go test
./services/events/events`.

## Fixtures

No external fixtures — the asserts build their query strings from bag
values produced by upstream saga steps.

## How to run

The folder has no test file yet, so `go test ./services/events/events`
is a no-op. The asserts are exercised by every saga journey that
verifies downstream event delivery, for example:

```bash
cd e2e_tests

# Exercises AssertTriggerExecutedSuccessfullyEventually +
# AssertLastTriggerRequestDataContains via the trigger smokes
go test -tags=saga ./journey/automations/... -v

# Exercises AssertRawEventReceivedAfter + AssertWorkflowEventReceivedAfter
go test -tags=saga ./journey/iot/... -v
```

## Outcome on pass

When the consuming saga journeys pass, these asserts collectively
prove the events service is reachable on its public read paths and
that downstream NATS consumers correctly populate the raw, trigger,
and workflow ClickHouse tables exposed by those endpoints.

## Requirements

- `events` reachable on port `5004` (override via `EVENTS_URL`).
- ClickHouse, NATS, and the upstream producers (`http_gateway`,
  `triggers`, `workflow`) reachable so the events service actually
  has records to return.
- Stack started from `mapexOSDeploy`; the seed admin user
  `admin@mapex.local` and the root organization
  `0000000000000000000aa001` must be present.

## Notes

- The asserts deliberately go through the public HTTP API only — they
  never read ClickHouse, Mongo, or NATS subjects directly, so the
  contract under test stays the user-facing one.
- Each assert subtracts a 2-second slack from the `startTime` it
  applies, to absorb clock skew between the test runner and the
  events service.
- The ClickHouse insert pipeline carries a batch buffer, so first
  observations of a record can take ten to thirty seconds in DEV;
  poll budgets in `assert_trigger_success.go` reflect that.
