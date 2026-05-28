# Configuration

All configuration is defined in `src/shared/configuration/application/configMap.ts`.

## Service
| ENV | Default | Description |
|---|---|---|
| `HTTP_PORT` | `8000` | HTTP server port |
| `HTTP_ADDRESS` | `0.0.0.0` | HTTP bind address |
| `SERVICE_NAME` | `js-executor` | Service name |
| `SERVICE_VERSION` | `1.0.0` | Service version |
| `NODE_ENV` | `dev` | Environment |
| `LOG_LEVEL` | `""` | Overrides default log level |

## Redis
| ENV | Default | Description |
|---|---|---|
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_USERNAME` | `` | Redis username |
| `REDIS_PASSWORD` | `` | Redis password |
| `REDIS_DB` | `1` | Redis DB |

## Redis Lock
| ENV | Default | Description |
|---|---|---|
| `REDIS_LOCK_DRIFT_FACTOR` | `0.01` | Lock drift factor |
| `REDIS_LOCK_RETRY_COUNT` | `10` | Retry count |
| `REDIS_LOCK_RETRY_DELAY` | `200` | Retry delay (ms) |
| `REDIS_LOCK_RETRY_JITTER` | `200` | Retry jitter (ms) |

## NATS
| ENV | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `NATS_USERNAME` | `service` | NATS username |
| `NATS_PASSWORD` | `service_secret` | NATS password |
| `NATS_CLIENT_NAME` | `js-executor-service` | NATS client name |

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
| `ASSETS_SERVICE_URL` | `http://localhost:5002` | Assets service URL |

## MinIO / S3
| ENV | Default | Description |
|---|---|---|
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO endpoint |
| `MINIO_ACCESS_KEY` | `mapex_admin` | Access key |
| `MINIO_SECRET_KEY` | `mapex_admin_secret_change_me` | Secret key |
| `MINIO_USE_SSL` | `false` | Use SSL |
| `MINIO_REGION` | `us-east-1` | Region |
| `MINIO_ASSETS_BUCKET` | `mapex-assets` | Assets bucket |
| `MINIO_TEMPLATES_BUCKET` | `mapex-templates` | Templates bucket |
| `MINIO_BYTECODE_BUCKET` | `mapex-bytecode` | Bytecode bucket |

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
| `CACHE_FALLBACK_TIMEOUT` | `5000` | Fallback timeout (ms) |
| `CACHE_BYTECODE_SKIP_L0` | `true` | Skip bytecode L0 cache |

## Execution Tuning
| ENV | Default | Description |
|---|---|---|
| `CPU_LIMIT` | `4` | CPU limit for auto‑tuning |
| `PISCINA_WORKERS` | `0` | Worker threads (0 = auto) |
| `ISOLATE_MEMORY_LIMIT_MB` | `32` | V8 isolate memory limit per worker |
| `WORKER_SCRIPT_TIMEOUT_MS` | `10000` | Script timeout in ms |
| `CONTEXT_RECYCLE_INTERVAL` | `10000` | Context recycle interval |
| `NATS_CONSUMER_BATCH_SIZE` | `0` | NATS batch size (0 = auto) |
| `NATS_CONSUMER_FETCH_TIMEOUT` | `0` | NATS fetch timeout (0 = default) |
| `NATS_CONSUMER_MAX_ACK_PENDING` | `0` | Max ack pending (0 = auto) |
| `CONCURRENCY_CHUNK_SIZE` | `0` | Concurrency chunk size (0 = auto) |
| `EVENTS_PER_WORKER` | `0` | Events per worker batch (0 = auto) |
