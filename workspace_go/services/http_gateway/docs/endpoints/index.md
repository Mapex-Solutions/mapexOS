# Endpoints

## HTTP API

### Events (external — no platform auth, auth per Data Source)
| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/events` | Ingest webhook event (auth per Data Source) |

**Flow:** Request validation -> Data Source authentication -> event processing

**Query parameters (required):**
- `ds` (string, MongoID): Data Source ID

**Auth:** Depends on Data Source config (`oauth2`, `jwt`, `apiKey`, `ip_whitelist`, `none`)

### Data Sources (internal — platform auth required)
| Method | Path | Permission | Description |
|---|---|---|---|
| GET | `/api/v1/data_sources` | `datasources.list` | List data sources (paginated, filtered) |
| POST | `/api/v1/data_sources` | `datasources.create` | Create data source |
| GET | `/api/v1/data_sources/:dataSourceId` | `datasources.read` | Get data source by ID |
| PATCH | `/api/v1/data_sources/:dataSourceId` | `datasources.update` | Update data source |
| DELETE | `/api/v1/data_sources/:dataSourceId` | `datasources.delete` | Delete data source |

**Flow:** User authentication -> request validation -> permission check -> operation execution

### Observability
| Method | Path | Description |
|---|---|---|
| GET | `/metrics` | Prometheus metrics |
| GET | `/health` | Health check (MongoDB, Redis app, Redis shared, NATS core) |

## NATS (outbound)
| Subject | Direction | Description |
|---|---|---|
| `processor.js.execute` | publish | Forward authenticated event to JS-Executor pipeline |
| `events.raw` | publish (fire-and-forget) | Auth failure security event with `success=false` |
