# Triggers Service Documentation

## Overview
The Triggers service is the outbound execution layer of MapexOS. It takes **trigger execution events** produced by RuleEngine and Router, resolves dynamic placeholders, and runs the configured action (HTTP, MQTT, Email, Slack, etc.). This isolates integrations from business logic so rules stay portable while delivery mechanisms evolve independently.

## Responsibilities
- Execute trigger actions consumed from NATS (`trigger.*.execute`).
- Resolve placeholders in trigger configs using event payload data.
- Provide CRUD APIs for trigger configuration and templates.
- Cache trigger definitions for low-latency execution.

## Non‑Responsibilities
- Rule evaluation (RuleEngine).
- Event ingestion (HTTP Gateway).
- Event persistence and analytics (Events service).

## Primary Data Flow
1. RuleEngine/Router publishes a `TriggerExecuteEvent` to NATS.
2. Triggers fetches the trigger config (Redis cache → Mongo fallback).
3. Placeholders are resolved using the event payload.
4. Executor runs the action (HTTP, Email, MQTT, etc.).
5. Execution result is published to `events.trigger`.

## Docs Map
- [Architecture](architecture/index.md)
- [API Specification](architecture/api-specification/index.md)
- [Endpoints](endpoints/index.md)
- [Configuration](configuration/index.md)
- [Operations](operations/index.md)
- [Observability](observability/index.md)
- [Tests](tests/index.md)
- [Benchmarks](benchmarks/index.md)
