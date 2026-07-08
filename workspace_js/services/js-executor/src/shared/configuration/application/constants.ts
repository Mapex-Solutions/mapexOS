// Default values for the service configuration, grouped by context. Imported
// into the ConfigDefinition list in configMap.ts so the definition entries carry
// named constants instead of inline literals. Sensitive credentials keep a
// dev-friendly default that the ConfigModule guard refuses to run with in a
// non-dev NODE_ENV.

// HTTP server — listener port, bind address, and service identity.
export const HTTP_PORT_DEFAULT = 8000;
export const HTTP_ADDRESS_DEFAULT = '0.0.0.0';
export const SERVICE_NAME_DEFAULT = 'js-executor';
export const SERVICE_VERSION_DEFAULT = '1.1.0';

// Redis — connection to the service's private database.
export const REDIS_HOST_DEFAULT = 'localhost';
export const REDIS_PORT_DEFAULT = 6379;
export const REDIS_USERNAME_DEFAULT = '';
export const REDIS_PASSWORD_DEFAULT = '';
export const REDIS_DB_DEFAULT = 1;

// Redis lock — Redlock tuning for distributed locks.
export const REDIS_LOCK_DRIFT_FACTOR_DEFAULT = 0.01;
export const REDIS_LOCK_RETRY_COUNT_DEFAULT = 10;
export const REDIS_LOCK_RETRY_DELAY_DEFAULT = 200;
export const REDIS_LOCK_RETRY_JITTER_DEFAULT = 200;

// NATS — connection and credentials (the password is sensitive).
export const NATS_URL_DEFAULT = 'nats://localhost:4222';
export const NATS_USERNAME_DEFAULT = 'service';
export const NATS_PASSWORD_DEFAULT = 'service_secret';
export const NATS_CLIENT_NAME_DEFAULT = 'js-executor-service';

// Auth — strategy, JWT secret/JWKS, algorithm, and role resolution.
export const AUTH_STRATEGY_DEFAULT = 'jwt';
export const AUTH_SECRET_DEFAULT = 'a-string-secret-at-least-256-bits-long';
export const AUTH_JWKS_URL_DEFAULT = '';
export const AUTH_ALGORITHM_DEFAULT = 'HS256';
export const AUTH_ROLES_SOURCE_DEFAULT = 'token'; // token | db | api
export const AUTH_ROLES_PATH_DEFAULT = 'roles';
export const AUTH_ROLES_API_URL_DEFAULT = '';

// Internal API key — shared secret for service-to-service calls (sensitive).
export const INTERNAL_API_KEY_DEFAULT = '5230c2e2-e245-468d-89e8-94154cf520d0';

// Assets service — internal HTTP base URL.
export const ASSETS_SERVICE_URL_DEFAULT = 'http://localhost:5002';

// Environment administration.
export const NODE_ENV_DEFAULT = 'dev';

// Log level override — empty means auto-derive from NODE_ENV.
export const LOG_LEVEL_DEFAULT = '';

// Object store (MinIO/S3) — endpoint, scoped credentials (sensitive), buckets.
export const OBJECT_STORE_ENDPOINT_DEFAULT = 'localhost:9000';
export const OBJECT_STORE_ACCESS_KEY_DEFAULT = 'svc-jsexec';
export const OBJECT_STORE_SECRET_KEY_DEFAULT = 'svc-jsexec-secret-change-me';
export const OBJECT_STORE_USE_SSL_DEFAULT = false;
export const OBJECT_STORE_REGION_DEFAULT = 'us-east-1';
export const OBJECT_STORE_AUTH_IS_NEEDED_DEFAULT = true;
export const MINIO_ASSETS_BUCKET_DEFAULT = 'mapex-assets';
export const MINIO_TEMPLATES_BUCKET_DEFAULT = 'mapex-templates';
export const MINIO_BYTECODE_BUCKET_DEFAULT = 'mapex-bytecode';

// TieredCache — L0 (RAM) and L1 (disk) sizing and TTLs.
export const CACHE_L0_MAX_SIZE_DEFAULT = 268435456; // 256MB
export const CACHE_L0_MAX_ITEMS_DEFAULT = 100000;
export const CACHE_L0_TTL_SECONDS_DEFAULT = 300; // 5min
export const CACHE_L1_ENABLED_DEFAULT = true;
export const CACHE_L1_DIR_DEFAULT = '/tmp/mapexos/cache';
export const CACHE_L1_MAX_SIZE_DEFAULT = 10737418240; // 10GB
export const CACHE_L1_TTL_SECONDS_DEFAULT = 3600; // 1h

// TieredCache fallback — HTTP API timeout when L2 misses.
export const CACHE_FALLBACK_TIMEOUT_DEFAULT = 5000; // ms

// Bytecode cache — skip L0 since the Script Registry already caches in RAM.
export const CACHE_BYTECODE_SKIP_L0_DEFAULT = true;

// CPU limit — the single knob all worker/chunk/batch values derive from.
export const CPU_LIMIT_DEFAULT = 4;

// Piscina worker threads (0 = auto: CPU_LIMIT - 1, minimum 1).
export const PISCINA_WORKERS_DEFAULT = 0;

// V8 isolate memory limit per worker thread (MB).
export const ISOLATE_MEMORY_LIMIT_MB_DEFAULT = 32;

// Worker script execution timeout (ms).
export const WORKER_SCRIPT_TIMEOUT_MS_DEFAULT = 10000;

// Recycle the V8 context every N events to prevent memory leaks.
export const CONTEXT_RECYCLE_INTERVAL_DEFAULT = 10000;

// NATS consumer tuning (0 = auto from CPU_LIMIT).
export const NATS_CONSUMER_BATCH_SIZE_DEFAULT = 0; // 0 = auto (CPU_LIMIT × 500)
export const NATS_CONSUMER_FETCH_TIMEOUT_DEFAULT = 1000; // 1s fetch timeout
export const NATS_CONSUMER_MAX_ACK_PENDING_DEFAULT = 0; // 0 = auto (batchSize × 2)

// Concurrency chunk size (0 = auto from PISCINA_WORKERS × 8).
export const CONCURRENCY_CHUNK_SIZE_DEFAULT = 0;

// Events per piscina.run() call for batch workers (0 = auto: 500).
export const EVENTS_PER_WORKER_DEFAULT = 0;
