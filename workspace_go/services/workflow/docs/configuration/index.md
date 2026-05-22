# Configuration

Configuration is provided via environment variables.

## Service

| ENV | Default | Description |
|-----|---------|-------------|
| `HTTP_PORT` | `5010` | HTTP server port |
| `HTTP_ADDRESS` | `0.0.0.0` | HTTP bind address |
| `SERVICE_NAME` | `workflow` | Service name (used in consumer names, logs) |
| `SERVICE_VERSION` | `1.0.0` | Service version |
| `GO_ENV` | `dev` | Environment (`dev` or `prod`) |
| `LOG_LEVEL` | — | Override log level (`debug`, `info`, `warn`, `error`, `silent`) |
| `CTX_TIMEOUT` | `4` | Global request timeout in seconds |

## MongoDB

| ENV | Default | Description |
|-----|---------|-------------|
| `MONGO_URI` | `mongodb://localhost:27017/?replicaSet=rs0` | MongoDB connection URI |
| `MONGO_DATABASE` | `dev-workflow` | Database name |
| `MONGO_POOL_SIZE` | `10` | Connection pool size |
| `MONGO_MONITOR_INTERVAL` | `10` | Health monitor interval in seconds |

## Redis (Shared DB only)

| ENV | Default | Description |
|-----|---------|-------------|
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_USERNAME` | — | Redis username |
| `REDIS_PASSWORD` | — | Redis password |
| `REDIS_SHARED_DB` | `5` | Shared DB number (auth middleware cache) |

> **Note:** No App Redis is used. TieredCache + MinIO replaces application-level caching.

## NATS

| ENV | Default | Description |
|-----|---------|-------------|
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `NATS_USERNAME` | `service` | NATS auth username |
| `NATS_PASSWORD` | `service_secret` | NATS auth password |
| `NATS_CLIENT_NAME` | `workflow-service` | NATS client identifier |

> **Note:** `nats_fetch_timeout` is used by consumers via `config.GetIntValue` but is not declared
> in the service configuration. It relies on the framework default (typically 30s). If you need to
> override it, add a `NATS_FETCH_TIMEOUT` entry to `DefaultConfiguration` in `config.go`.

## Authentication

| ENV | Default | Description |
|-----|---------|-------------|
| `AUTH_STRATEGY` | `jwt` | Auth strategy (`jwt` or `oauth2`) |
| `AUTH_SECRET` | `a-string-secret-at-least-256-bits-long` | JWT secret for HS256 |
| `AUTH_JWKS_URL` | — | JWKS URL for RS256 |
| `AUTH_ALGORITHM` | `HS256` | JWT algorithm (`HS256` or `RS256`) |
| `AUTH_ROLES_SOURCE` | `token` | Roles source (`token`, `db`, `api`) |
| `AUTH_ROLES_PATH` | `roles` | JSON path in token for roles |
| `AUTH_ROLES_API_URL` | — | External API URL for role lookup |
| `INTERNAL_API_KEY` | `5230c2e2-e245-468d-89e8-94154cf520d0` | Internal service API key |

## MinIO / S3

| ENV | Default | Description |
|-----|---------|-------------|
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO endpoint |
| `MINIO_ACCESS_KEY` | `mapexos_admin` | Access key |
| `MINIO_SECRET_KEY` | `mapexos_admin_secret_change_me` | Secret key |
| `MINIO_USE_SSL` | `false` | Enable SSL |
| `MINIO_REGION` | `us-east-1` | S3 region |
| `MINIO_DEFINITIONS_BUCKET` | `mapex-workflows` | Bucket for workflow definitions |

## TieredCache

| ENV | Default | Description |
|-----|---------|-------------|
| `CACHE_L0_MAX_SIZE` | `268435456` | L0 (RAM) max size in bytes (256MB) |
| `CACHE_L0_MAX_ITEMS` | `100000` | L0 max item count |
| `CACHE_L0_TTL_SECONDS` | `300` | L0 TTL (5 minutes) |
| `CACHE_L1_DIR` | `/tmp/mapexos/cache` | L1 (disk) cache directory |
| `CACHE_L1_MAX_SIZE` | `10737418240` | L1 max size in bytes (10GB) |
| `CACHE_L1_TTL_SECONDS` | `3600` | L1 TTL (1 hour) |

## Metrics

| ENV | Default | Description |
|-----|---------|-------------|
| `METRICS_GO_COLLECTOR` | `true` | Enable Go runtime metrics |
| `METRICS_PROCESS_COLLECTOR` | `true` | Enable process metrics |

## Integration

| ENV | Default | Description |
|-----|---------|-------------|
| `MAPEXOS_URL` | `http://localhost:5000` | MapexOS internal API base URL |
