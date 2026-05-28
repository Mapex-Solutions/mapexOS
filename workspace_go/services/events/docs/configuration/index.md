# Configuration

Configuration is provided via environment variables.

## Service
| ENV | Default | Description |
|---|---|---|
| `HTTP_PORT` | `5004` | HTTP server port |
| `HTTP_ADDRESS` | `0.0.0.0` | HTTP bind address |
| `SERVICE_NAME` | `events` | Service name |
| `SERVICE_VERSION` | `1.0.0` | Service version |
| `GO_ENV` | `dev` | Environment (`dev`/`production`) |
| `LOG_LEVEL` | `""` | Overrides default log level (debug for dev, info for prod) |
| `CTX_TIMEOUT` | `4` | Request timeout in seconds |

## ClickHouse
| ENV | Default | Description |
|---|---|---|
| `CLICKHOUSE_HOST` | `localhost` | ClickHouse host |
| `CLICKHOUSE_PORT` | `9440` | ClickHouse port |
| `CLICKHOUSE_DATABASE` | `mapexos` | Database name |
| `CLICKHOUSE_USERNAME` | `mapexos_user` | Username |
| `CLICKHOUSE_PASSWORD` | `mapexos_password` | Password |

## MongoDB
| ENV | Default | Description |
|---|---|---|
| `MONGO_URI` | `mongodb://localhost:27017/?replicaSet=rs0` | Mongo connection URI |
| `MONGO_DATABASE` | `events` | Database name |
| `MONGO_POOL_SIZE` | `10` | MongoDB pool size |
| `MONGO_MONITOR_INTERVAL` | `10` | Monitor interval (seconds) |

## Redis
| ENV | Default | Description |
|---|---|---|
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_USERNAME` | `` | Redis username |
| `REDIS_PASSWORD` | `` | Redis password |
| `REDIS_DB` | `0` | Redis DB for app cache |
| `REDIS_SHARED_DB` | `5` | Redis DB for shared cache |

## NATS
| ENV | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `NATS_USERNAME` | `service` | NATS username |
| `NATS_PASSWORD` | `service_secret` | NATS password |
| `NATS_CLIENT_NAME` | `events-service` | NATS client name |
| `NATS_BATCH_SIZE` | `5000` | Batch size (NATS fetch + ClickHouse insert) |
| `NATS_FETCH_TIMEOUT` | `1` | Fetch timeout in seconds |

## Auth / Permissions
| ENV | Default | Description |
|---|---|---|
| `AUTH_STRATEGY` | `jwt` | Auth strategy (`jwt` or `oauth2`) |
| `AUTH_SECRET` | `a-string-secret-at-least-256-bits-long` | JWT secret (HS256) |
| `AUTH_JWKS_URL` | `` | JWKS URL for OAuth2/JWT RS256 |
| `AUTH_ALGORITHM` | `HS256` | JWT algorithm |
| `AUTH_ROLES_SOURCE` | `token` | Roles source (`token`, `db`, `api`) |
| `AUTH_ROLES_PATH` | `roles` | Path in token payload |
| `AUTH_ROLES_API_URL` | `` | Roles API URL |
| `INTERNAL_API_KEY` | `5230c2e2-e245-468d-89e8-94154cf520d0` | Internal service auth key |

## External Services
| ENV | Default | Description |
|---|---|---|
| `MAPEXOS_URL` | `http://localhost:5000` | MapexOS API URL |
| `ASSETS_URL` | `http://localhost:5002` | Assets service URL |

## Metrics
| ENV | Default | Description |
|---|---|---|
| `METRICS_GO_COLLECTOR` | `true` | Enable Go runtime metrics |
| `METRICS_PROCESS_COLLECTOR` | `true` | Enable process metrics |

## MinIO / S3 (Templates Cache)
| ENV | Default | Description |
|---|---|---|
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO endpoint |
| `MINIO_ACCESS_KEY` | `mapex_admin` | Access key |
| `MINIO_SECRET_KEY` | `mapex_admin_secret_change_me` | Secret key |
| `MINIO_USE_SSL` | `false` | Use SSL |
| `MINIO_REGION` | `us-east-1` | Region |
| `MINIO_TEMPLATES_BUCKET` | `mapex-templates` | Templates bucket |

## Tiered Cache
| ENV | Default | Description |
|---|---|---|
| `CACHE_L0_MAX_SIZE` | `268435456` | L0 max size (bytes) |
| `CACHE_L0_MAX_ITEMS` | `100000` | L0 max items |
| `CACHE_L0_TTL_SECONDS` | `300` | L0 TTL (seconds) |
| `CACHE_L1_ENABLED` | `true` | Enable L1 disk cache |
| `CACHE_L1_DIR` | `/tmp/mapexos/cache` | L1 cache directory |
| `CACHE_L1_MAX_SIZE` | `10737418240` | L1 max size (bytes) |
| `CACHE_L1_TTL_SECONDS` | `3600` | L1 TTL (seconds) |
