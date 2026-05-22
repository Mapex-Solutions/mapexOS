# Endpoints

## HTTP API

### Events Queries (base: `/api/v1/events`)
| Method | Path | Description |
|---|---|---|
| GET | `/raw` | Query raw events (filters, pagination) |
| GET | `/jsexec` | Query JS‑Executor debug events |
| GET | `/router` | Query router execution history |
| GET | `/businessrule` | Query business rule execution history |
| GET | `/trigger` | Query trigger execution history |
| POST | `/store/query` | Query processed events with EVA filters |
| GET | `/store/:eventTrackerId` | Get processed event detail with EVA field names |

**Query/Body contracts**
- DTOs are defined under `workspace_go/packages/contracts/services/events/**`.
- Each endpoint validates its DTO via middleware before execution.

### Retention Policies (base: `/api/v1/retention`)
| Method | Path | Description |
|---|---|---|
| GET | `/` | List retention policies (paginated, filtered) |
| PUT | `/` | Upsert retention policy by org + type |
| GET | `/:retentionPolicyId` | Get retention policy by ID |
| DELETE | `/:retentionPolicyId` | Delete retention policy |

## NATS Consumers

### Events Module
| Stream | Subject | Description |
|---|---|---|
| `EVENTS` | `events.save` | Store processed events into main table |
| `EVENTS-RAW` | `events.raw` | Store raw events for debugging |
| `EVENTS-JSEXEC` | `events.logs.jsexecutor` | Store JS-Executor debug logs |
| `MAPEXOS-DLQ` | `dlq.mapexos` | Store failed events (DLQ) |
| `EVENTS-ROUTER` | `events.router` | Store router execution events |
| `EVENTS-BUSINESSRULE` | `events.businessrule` | Store business rule execution events |
| `EVENTS-TRIGGER` | `events.trigger` | Store trigger execution events |

### Retention Module
| Stream | Subject | Description |
|---|---|---|
| `MAPEXOS` | `mapexos.events.organization.created` | Create default retention policies for new organizations |

## Observability
- Health: `GET /health`
- Metrics: `GET /metrics`
