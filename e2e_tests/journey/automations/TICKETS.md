# Internal tickets — automations journeys

Tracking file for follow-up work not yet executed. Not part of the
journey READMEs (those describe only the working tests).

## Open

_(empty — all scoped items shipped)_

## Closed

- ~~MQTT/NATS broker creds in trigger config~~ — replaced with in-process broker (mochi-mqtt) and embedded server (nats-server/v2). No external broker needed.
- ~~RabbitMQ test broker~~ — replaced with testcontainers-go ephemeral container per test run.
- ~~Phase 2 (event pipeline) for every trigger journey~~ — `PostRawEvent` step in http_gateway/datasources + `phase2_event_pipeline/` package under every `trigger_<type>/`. Exercises `gateway → js-executor → router → trigger` end-to-end.
- ~~Retrofit existing IoT journeys with the new docs standard~~ — `journey/iot/**` now carries per-item `// comment` annotations and README.md + README_pt.md pairs at every phase and journey level.
- ~~`run-tests.sh saga trigger-<type>`~~ — added `saga` subcommand with shortcuts for each trigger type, help/list updated, bash completion extended.
