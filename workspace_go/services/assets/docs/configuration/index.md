# Configuration

Configuration is provided via environment variables.

## Service
| ENV | Default | Description |
|---|---|---|
| `HTTP_PORT` | `5002` | HTTP server port |
| `HTTP_ADDRESS` | `0.0.0.0` | HTTP bind address |
| `SERVICE_NAME` | `assets` | Service name |
| `SERVICE_VERSION` | `1.0.0` | Service version |
| `GO_ENV` | `dev` | Environment (`dev`/`production`). Drives env-scoped NATS stream/subject names. |
| `LOG_LEVEL` | `""` | Overrides default log level |
| `CTX_TIMEOUT` | `4` | Request timeout in seconds |

## MongoDB
| ENV | Default | Description |
|---|---|---|
| `MONGO_URI` | `mongodb://localhost:27017/?replicaSet=rs0` | Mongo connection URI |
| `MONGO_DATABASE` | `assets` | Database name |
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

## NATS (Core + JetStream)
| ENV | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `NATS_USERNAME` | `service` | NATS username |
| `NATS_PASSWORD` | `service_secret` | NATS password |
| `NATS_CLIENT_NAME` | `assets-service` | NATS client name |
| `STREAM_REPLICAS` | `1` | Replica count used when this service creates its owned streams (`MAPEXOS-ASSETS-HEARTBEAT`, `MAPEXOS-ASSETS-MQTT-PRESENCE`, `MAPEXOS-ASSETS-HEALTH-MONITOR`). |

## mapexVault (MQTT PKI source of truth)
| ENV | Default | Description |
|---|---|---|
| `MAPEX_VAULT_URL` | `http://localhost:5010` | mapexVault HTTP base URL. Used by `mqttcerts` `OnMount` to fetch the intermediate CA at startup. |
| `MAPEX_VAULT_API_KEY` | `` (required) | API key for mapexVault `pki/` endpoints. |

## MQTT PKI (signing parameters used by `mqttcerts`)
| ENV | Default | Description |
|---|---|---|
| `PKI_DEVICE_CERT_TTL_DAYS` | `365` | TTL for newly issued device leaf certs (days). |
| `PKI_KEY_ALGORITHM` | `ECDSA-P256` | Signing algorithm. Only `ECDSA-P256` is supported today; the env var is kept to flag future rotation work. |
| `PKI_ON_MOUNT_RETRY_BASE_SEC` | `2` | Base seconds for exponential backoff when `OnMount` fails to pull the intermediate CA. |
| `PKI_ON_MOUNT_RETRY_MAX_SEC` | `60` | Cap for the exponential backoff between `OnMount` retries. |

The intermediate CA itself is fetched from mapexVault, not configured here. See `documentation/architecture/mqtt-pki.md` for the bootstrap flow (operator runs `scripts/prebuild/pki/generate-pki.sh`; mapexVault seeds Mongo via `VAULT_PKI_BOOTSTRAP_DIR` on first start; this service pulls from mapexVault at runtime).

## Auth / Permissions (platform JWT layer)
| ENV | Default | Description |
|---|---|---|
| `AUTH_STRATEGY` | `jwt` | Auth strategy (`jwt` or `oauth2`) |
| `AUTH_SECRET` | `a-string-secret-at-least-256-bits-long` | JWT secret (HS256) |
| `AUTH_JWKS_URL` | `` | JWKS URL for OAuth2/JWT RS256 |
| `AUTH_ALGORITHM` | `HS256` | JWT algorithm |
| `AUTH_ROLES_SOURCE` | `token` | Roles source (`token`, `db`, `api`) |
| `AUTH_ROLES_PATH` | `roles` | Path in token payload |
| `AUTH_ROLES_API_URL` | `` | Roles API URL |
| `INTERNAL_API_KEY` | `5230c2e2-e245-468d-89e8-94154cf520d0` | Internal service auth key. Gates every `/internal/*` route via the standard `apikeymw.ApiKeyAuthMiddleware` (`X-API-Key` header) — including the read-model L3 fallback `GET /internal/assets/:assetUUID` that the mapex-mqtt-broker plugin reads on TieredCache miss. The broker MUST be configured to forward this value in `X-API-Key`. |

## External Services
| ENV | Default | Description |
|---|---|---|
| `MAPEXOS_URL` | `http://localhost:5000` | MapexOS API URL |
| `ROUTER_SERVICE_URL` | `http://localhost:5003` | Router service URL |

## MinIO / S3
| ENV | Default | Description |
|---|---|---|
| `MINIO_ENDPOINT` | `localhost:9000` | MinIO endpoint |
| `MINIO_ACCESS_KEY` | `mapex_admin` | Access key |
| `MINIO_SECRET_KEY` | `mapex_admin_secret_change_me` | Secret key |
| `MINIO_USE_SSL` | `false` | Use SSL |
| `MINIO_REGION` | `us-east-1` | Region |
| `MINIO_TEMPLATES_BUCKET` | `mapex-templates` | Templates bucket |
| `MINIO_ASSETS_BUCKET` | `mapex-assets` | Assets bucket |

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

## Health monitoring
| ENV | Default | Description |
|---|---|---|
| `HEALTHMONITOR_SCAN_INTERVAL_SEC` | `60` | Interval between scheduled offline scans. |
| `HEALTHMONITOR_OFFLINE_GRACE_SEC` | `120` | Default grace period before flipping an asset to `offline` (per-asset `thresholdMinutes` overrides this). |
