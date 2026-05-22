# API Specification

This document reflects the canonical DTOs from `packages/contracts/services/triggers/triggers`.

## TriggerCreate
Required fields:
- `name`
- `triggerType`
- `category`
- `enabled`
- `config`

Template and tenancy fields:
- `isSystem`, `isTemplate` control template visibility.
- `orgId` and `pathKey` are **overwritten** by the service based on RequestContext.

## TriggerUpdate
All fields are optional and replace values when provided. `config` replaces the entire config block.

## TriggerResponse
Includes persisted fields and `config` (union type).

## TriggerQuery
Supports filtering and pagination:
- Filters: `id`, `name`, `triggerType`, `category`, `enabled`, `orgId`, `pathKey`, `isSystem`, `isTemplate`
- Pagination: `page`, `perPage`, `sort`
- Hierarchy: `includeChildren` from `BaseQueryDTO`

## TriggerConfig (Union Type)
Exactly one config block must be populated and **must match** `triggerType`.

### Technical
- `http`: `endpoint`, `method`, optional `headers`, `body`, `timeout`
- `mqtt`: `broker`, `port`, `topic`, `qos`, optional `username`, `password`, `clientId`, `message`, `useTLS`
- `rabbitmq`: `host`, `port`, `username`, `password`, `publishMode`, optional `exchange`, `exchangeType`, `routingKey`, `queue`, `message`, `useTLS`
- `nats`: `server`, `subject`, optional `username`, `password`, `token`, `message`, `useTLS`
- `websocket`: `url`, optional `message`, `headers`

### Communication
- `email`: `to`, `subject`, optional `cc`, `bcc`, `body`, `htmlBody`
- `teams`: `webhookUrl`, `title`, `text`, optional `themeColor`
- `slack`: `webhookUrl`, `message`, optional `channel`, `username`, `iconEmoji`

## TriggerExecuteEvent (NATS)
Subject pattern: `trigger.{triggerId}.execute`

Fields:
- `triggerId`: target trigger
- `executionId`: unique execution id
- `eventTrackerId`: end‑to‑end tracing id
- `source`: origin service (`router` or `ruleengine`)
- `payload`: free‑form data used for placeholder resolution
- `orgId`: tenant id
- `pathKey`: hierarchy path
- `created`: timestamp string
