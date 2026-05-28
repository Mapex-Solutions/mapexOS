# Configuration

Configuration is provided via environment variables.

Configuration sources (in priority order):
1. Environment variables
2. Default values defined in code

## Service
| ENV | Default | Description |
|---|---|---|
| `HTTP_PORT` | `5003` | HTTP server port |
| `HTTP_ADDRESS` | `0.0.0.0` | HTTP bind address |
| `SERVICE_NAME` | `router` | Service name used in logs, NATS consumer names, health |
| `SERVICE_VERSION` | `1.0.0` | Service version reported in health endpoint |
| `GO_ENV` | `dev` | Environment (`dev` / `production`) |
| `LOG_LEVEL` | `""` | Overrides default log level (`debug` for dev, `info` for production) |
| `CTX_TIMEOUT` | `4` | Request timeout in seconds (applied by ContextInjector middleware) |

## MongoDB
| ENV | Default | Description |
|---|---|---|
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection URI |
| `MONGO_DATABASE` | `router` | Database name |
| `MONGO_POOL_SIZE` | `10` | MongoDB connection pool size |
| `MONGO_MONITOR_INTERVAL` | `10` | Server monitor interval in seconds |

## Redis
| ENV | Default | Description |
|---|---|---|
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_USERNAME` | `""` | Redis username |
| `REDIS_PASSWORD` | `""` | Redis password |
| `REDIS_DB` | `0` | Redis DB index for app cache (RouteGroup cache) |
| `REDIS_SHARED_DB` | `5` | Redis DB index for shared cache (coverage, permissions) |

## NATS
| ENV | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `NATS_USERNAME` | `service` | NATS username |
| `NATS_PASSWORD` | `service_secret` | NATS password |
| `NATS_CLIENT_NAME` | `router-service` | NATS client name |
| `NATS_BATCH_SIZE` | `8000` | Messages per batch fetch |
| `NATS_FETCH_TIMEOUT` | `5` | Fetch timeout in seconds |

## Auth / Permissions
| ENV | Default | Description |
|---|---|---|
| `AUTH_STRATEGY` | `jwt` | Auth strategy (`jwt` or `oauth2`) |
| `AUTH_SECRET` | `a-string-secret-at-least-256-bits-long` | JWT secret (HS256) |
| `AUTH_JWKS_URL` | `""` | JWKS URL for OAuth2/JWT RS256 |
| `AUTH_ALGORITHM` | `HS256` | JWT algorithm |
| `AUTH_ROLES_SOURCE` | `token` | Roles source (`token`, `db`, `api`) |
| `AUTH_ROLES_PATH` | `roles` | JSON path in token payload for roles |
| `AUTH_ROLES_API_URL` | `""` | External roles API URL |
| `INTERNAL_API_KEY` | `5230c2e2-e245-468d-89e8-94154cf520d0` | API key for internal MS-to-MS communication |

## External Services
| ENV | Default | Description |
|---|---|---|
| `MAPEXOS_URL` | `http://localhost:5000` | MapexOS API URL (permission middleware) |
| `ASSETS_URL` | `http://localhost:5002` | Assets service URL (TieredCache fallback) |

## MinIO / S3 (TieredCache L2)
| ENV | Default | Description |
|---|---|---|
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO endpoint |
| `MINIO_ACCESS_KEY` | `mapex_admin` | Access key |
| `MINIO_SECRET_KEY` | `mapex_admin_secret_change_me` | Secret key |
| `MINIO_USE_SSL` | `false` | Use SSL |
| `MINIO_REGION` | `us-east-1` | Region |
| `MINIO_ASSETS_BUCKET` | `mapex-assets` | Assets bucket name |
| `MINIO_TEMPLATES_BUCKET` | `mapex-templates` | Templates bucket name |

## TieredCache
| ENV | Default | Description |
|---|---|---|
| `CACHE_L0_MAX_SIZE` | `268435456` | L0 (RAM) max size in bytes (256 MB) |
| `CACHE_L0_MAX_ITEMS` | `100000` | L0 (RAM) max number of items |
| `CACHE_L0_TTL_SECONDS` | `300` | L0 (RAM) TTL in seconds (5 min) |
| `CACHE_L1_ENABLED` | `true` | Enable L1 (Disk) cache |
| `CACHE_L1_DIR` | `/tmp/mapexos/cache` | L1 (Disk) cache directory |
| `CACHE_L1_MAX_SIZE` | `10737418240` | L1 (Disk) max size in bytes (10 GB) |
| `CACHE_L1_TTL_SECONDS` | `3600` | L1 (Disk) TTL in seconds (1 hour) |
| `CACHE_FALLBACK_TIMEOUT` | `5` | Fallback HTTP timeout in seconds |

## Metrics
| ENV | Default | Description |
|---|---|---|
| `METRICS_GO_COLLECTOR` | `true` | Enable Go runtime metrics collector |
| `METRICS_PROCESS_COLLECTOR` | `true` | Enable process metrics collector |
