# Endpoints

## HTTP API

### External API (JWT Authentication)

Base path: `/api/v1/route_groups`

**Flow:** Authentication (JWT) -> request validation -> permission check -> operation execution.

| Method | Path | Permission | Description |
|---|---|---|---|
| `GET` | `/api/v1/route_groups` | `RouteGroupList` | List route groups (paginated, filtered) |
| `GET` | `/api/v1/route_groups/counter` | `RouteGroupList` | Count route groups (cached, 6h TTL) |
| `POST` | `/api/v1/route_groups` | `RouteGroupCreate` | Create a new route group |
| `GET` | `/api/v1/route_groups/:routeGroupId` | `RouteGroupRead` | Get route group by ID |
| `PATCH` | `/api/v1/route_groups/:routeGroupId` | `RouteGroupUpdate` | Update route group by ID |
| `DELETE` | `/api/v1/route_groups/:routeGroupId` | `RouteGroupDelete` | Delete route group by ID |

Each request is validated, then permission-checked, then scoped to the user's organization.

#### GET /api/v1/route_groups

Query parameters:

| Parameter | Type | Default | Description |
|---|---|---|---|
| `projection` | string | _(all fields)_ | Comma-separated fields to return |
| `page` | int | `1` | Page number |
| `perPage` | int | `20` | Items per page |
| `sort` | string | `created:desc` | Sort field and direction |
| `includeChildren` | bool | `false` | Include child org route groups (hierarchical) |
| `name` | string | | Filter by name |
| `enabled` | bool | | Filter by enabled status |
| `version` | string | | Filter by version |
| `isTemplate` | bool | | Filter by template flag |
| `isSystem` | bool | | Filter by system flag |

Headers:
- `Authorization: Bearer <JWT>` (required)
- `X-Org-Context: <orgId>` (optional, for org-scoped filtering)

#### GET /api/v1/route_groups/counter

Returns `{ "count": <number> }`. Uses cached count (6h TTL), invalidated on create/delete.

#### POST /api/v1/route_groups

Body: `RouteGroupCreateDTO` (validated via middleware).

#### GET /api/v1/route_groups/:routeGroupId

Path params:
- `routeGroupId` (string, required): MongoDB ObjectId of the route group.

#### PATCH /api/v1/route_groups/:routeGroupId

Path params:
- `routeGroupId` (string, required): MongoDB ObjectId of the route group.

Body: `RouteGroupUpdateDTO` (validated via middleware). All fields are optional (partial update).

#### DELETE /api/v1/route_groups/:routeGroupId

Path params:
- `routeGroupId` (string, required): MongoDB ObjectId of the route group.

### Internal API (API Key Authentication)

Base path: `/api/internal/v1/routegroups`

Authentication: `X-API-Key` header (configured via `INTERNAL_API_KEY`).

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/internal/v1/routegroups` | Get multiple route groups by IDs |

#### GET /api/internal/v1/routegroups

Query parameters:

| Parameter | Type | Required | Description |
|---|---|---|---|
| `ids` | string | yes | Comma-separated list of route group IDs |
| `projection` | string | no | Comma-separated fields to return |

Example: `GET /api/internal/v1/routegroups?ids=id1,id2,id3&projection=name,enabled`

## NATS

### Consumed

| Subject | Stream | Pattern | Description |
|---|---|---|---|
| `route.execute` | `ROUTE-GROUPS` | WorkQueue (durable, load-balanced) | Route execution events |
| `fanout.asset.invalidate` | `FANOUT` | FANOUT (ephemeral, broadcast) | Asset cache invalidation |

#### route.execute

Required payload fields:
- `orgId` (string)
- `assetUUID` (string)
- `event` (map)

Optional fields:
- `pathKey` (string)
- `eventTrackerId` (string)

#### fanout.asset.invalidate

Required payload fields:
- `orgId` (string)
- `assetUUID` (string)

### Published

| Subject | Description |
|---|---|
| `events.save` | Event persistence (`EventStoreDTO` with asset context and optional metadata) |
| `events.lake_house` | Lakehouse analytics (enriched event with asset identifiers) |
| `events.notification` | Notification dispatch (enriched event with `notificationId`) |
| `ruleengine.{businessRuleId}.execute` | Rule engine evaluation (enriched event with `businessRuleId`) |
| `trigger.{triggerId}.execute` | Trigger execution (`TriggerExecuteEvent` with original event as payload) |
| `events.router` | Routing history (`RouterHistoryEvent` with match results per router) |

## Observability Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/health` | none | Health check |
| `GET` | `/metrics` | none | Prometheus metrics |
