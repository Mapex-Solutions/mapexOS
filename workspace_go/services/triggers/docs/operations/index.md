# Operations

## Start
```bash
./bin/triggers
```

## Smoke Checks
- `GET /health`
- `GET /metrics`
- `GET /api/v1/triggers` with valid auth

## Dependencies
- MongoDB (trigger persistence)
- Redis (trigger cache + counter cache)
- NATS JetStream (trigger execution stream)

## NATS Stream
- Stream: `TRIGGERS`
- Subject: `trigger.*.execute`
- Consumer: `{SERVICE_NAME}-trigger-executor`
- Queue group: `{SERVICE_NAME}-TRIGGER-EXEC-GROUP`

## Scaling Guidance
- Increase `TRIGGER_EXECUTOR_WORKERS` for I/O‑heavy trigger loads.
- Scale horizontally when executor latency or queue lag grows.
- Monitor `triggers_trigger_processing_duration_seconds` and `triggers_message_total` for backlog signals.

## Failure Handling
- Disabled triggers are ACKed and do not execute.
- Invalid payloads are REJECTed (no retry).
- Execution errors are NACKed and retried up to 5 times, then sent to DLQ.

## Benchmarks
See [Benchmarks](../benchmarks/index.md) for load-testing methodology, scripts, and results.
