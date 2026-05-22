# Configuration

Configuration is provided via environment variables.

## Service
| ENV | Default | Description |
|---|---|---|
| `HTTP_PORT` | `5000` | HTTP server port |
| `HTTP_ADDRESS` | `0.0.0.0` | HTTP bind address |
| `SERVICE_NAME` | `mapexos` | Service name |
| `SERVICE_VERSION` | `1.0.0` | Service version |
| `GO_ENV` | `dev` | Environment (`dev`/`production`) |
| `LOG_LEVEL` | `""` | Overrides default log level |
| `CTX_TIMEOUT` | `4` | Request timeout in seconds |

## MongoDB
| ENV | Default | Description |
|---|---|---|
| `MONGO_URI` | `mongodb://localhost:27017/?replicaSet=rs0` | Mongo connection URI |
| `MONGO_DATABASE` | `mapexos` | Database name |
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
| `NATS_CLIENT_NAME` | `mapexos-service` | NATS client name |

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
