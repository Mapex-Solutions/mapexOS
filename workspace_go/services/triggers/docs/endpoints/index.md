# Endpoints

## HTTP API
Base path: `/api/v1/triggers`

All endpoints require auth and use RequestContext to populate `orgId` and `pathKey`.

### Summary
| Method | Path | Permission | Description |
|---|---|---|---|
| GET | `/` | `triggers.list` | List triggers (paginated, filtered) |
| GET | `/counter` | `triggers.list` | Count triggers (cached) |
| POST | `/` | `triggers.create` | Create trigger |
| GET | `/:id` | `triggers.read` | Get trigger by ID |
| PATCH | `/:id` | `triggers.update` | Update trigger |
| DELETE | `/:id` | `triggers.delete` | Delete trigger |

### Request Processing
All requests go through validation, permission check, and organization context injection before being processed.

### List + Filter
`GET /api/v1/triggers`

Query parameters:
- `projection` (comma‑separated fields)
- `page` (default: 1)
- `perPage` (default: 10)
- `sort` (default: `created:desc`)
- `includeChildren` (default: `false`)
- `id`, `name`, `triggerType`, `category`, `enabled`, `orgId`, `pathKey`, `isSystem`, `isTemplate`

### Counter
`GET /api/v1/triggers/counter`

Returns total count for the current org context (cached).

### Create
`POST /api/v1/triggers`

Uses `TriggerCreate`. The service validates `config` against `triggerType`.

### Read
`GET /api/v1/triggers/:id`

### Update
`PATCH /api/v1/triggers/:id`

Uses `TriggerUpdate` (partial update).

### Delete
`DELETE /api/v1/triggers/:id`

## NATS
### Trigger Execution (Consumed)
Subject pattern: `trigger.{triggerId}.execute`

Payload (TriggerExecuteEvent):
- `triggerId` (string)
- `executionId` (string)
- `eventTrackerId` (string)
- `source` (string: `router` | `ruleengine`)
- `payload` (map[string]any) used by placeholder resolver
- `orgId` (ObjectId)
- `pathKey` (string)
- `created` (string)

### Trigger Execution Result (Published)
Subject: `events.trigger`

Result payload includes:
- `triggerId`, `triggerName`, `triggerType`, `category`
- `eventTrackerId`, `orgId`, `pathKey`, `source`
- `success`, `durationMs`, `error`
- `requestData` (serialized config after placeholder resolution)

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`
