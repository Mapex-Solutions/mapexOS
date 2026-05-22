import { ConfigDefinition } from '@mapexos/microservices';

// InitConfig initializes the configuration module with the provided definitions.
export const defaultConfiguration: ConfigDefinition[] = [

	/** HTTP Server Configuration */
	{ key: 'http_port', env: 'HTTP_PORT', type: 'int', default: 8000 },
	{ key: 'http_address', env: 'HTTP_ADDRESS', type: 'string', default: '0.0.0.0' },

	{ key: 'service_name', env: 'SERVICE_NAME', type: 'string', default: 'js-executor' },
	{ key: 'service_version', env: 'SERVICE_VERSION', type: 'string', default: '1.0.0' },

	/* Redis Configuration */
	{ key: 'redis_host', env: 'REDIS_HOST', type: 'string', default: 'localhost' },
	{ key: 'redis_port', env: 'REDIS_PORT', type: 'int', default: 6379 },
	{ key: 'redis_username', env: 'REDIS_USERNAME', type: 'string', default: '' },
	{ key: 'redis_password', env: 'REDIS_PASSWORD', type: 'string', default: '' },
	{ key: 'redis_db', env: 'REDIS_DB', type: 'int', default: 1 },

	/* Redis Lock Configuration */
	{ key: 'redis_lock_drift_factor', env: 'REDIS_LOCK_DRIFT_FACTOR', type: 'int', default: 0.01 },
	{ key: 'redis_lock_retry_count', env: 'REDIS_LOCK_RETRY_COUNT', type: 'int', default: 10 },
	{ key: 'redis_lock_retry_delay', env: 'REDIS_LOCK_RETRY_DELAY', type: 'int', default: 200 },
	{ key: 'redis_lock_retry_jitter', env: 'REDIS_LOCK_RETRY_JITTER', type: 'int', default: 200 },

	/** NATS Configuration */
	{ key: 'nats_url', env: 'NATS_URL', type: 'string', default: 'nats://localhost:4222' },
	{ key: 'nats_username', env: 'NATS_USERNAME', type: 'string', default: 'service' },
	{ key: 'nats_password', env: 'NATS_PASSWORD', type: 'string', default: 'service_secret', sensitive: true },
	{ key: 'nats_client_name', env: 'NATS_CLIENT_NAME', type: 'string', default: 'js-executor-service' },

	/**
	 * Defines the authentication strategy used by the platform (jwt or oauth2) and
	 * how user roles are retrieved (token, db, or api).
	 * Supports static JWT secrets (HS256) and external identity providers using JWKS (RS256).
	 * Allows dynamic role extraction from tokens, local databases, or external APIs.
	 */
	{ key: 'auth_strategy', env: 'AUTH_STRATEGY', type: 'string', default: 'jwt' },
	{ key: 'auth_secret', env: 'AUTH_SECRET', type: 'string', default: 'a-string-secret-at-least-256-bits-long', sensitive: true },
	{ key: 'auth_jwks_url', env: 'AUTH_JWKS_URL', type: 'string', default: '' },
	{ key: 'auth_algorithm', env: 'AUTH_ALGORITHM', type: 'string', default: 'HS256' },

	{ key: 'auth_roles_source', env: 'AUTH_ROLES_SOURCE', type: 'string', default: 'token' }, // token | db | api
	{ key: 'auth_roles_path', env: 'AUTH_ROLES_PATH', type: 'string', default: 'roles' },
	{ key: 'auth_roles_api_url', env: 'AUTH_ROLES_API_URL', type: 'string', default: '' },

	/**
	 *	My apiKey for authentication
	 * Please replace this with your own apiKey for authentication
	 * This api Is used for internal communication between services
	 */
	{ key: 'internal_api_key', env: 'INTERNAL_API_KEY', type: 'string', default: '5230c2e2-e245-468d-89e8-94154cf520d0', sensitive: true },

	/** Assets Service Configuration */
	{ key: 'assets_service_url', env: 'ASSETS_SERVICE_URL', type: 'string', default: 'http://localhost:5002' },

	/** Environment administration */
	{ key: 'node_env', env: 'NODE_ENV', type: 'string', default: 'dev' },

	/** Log level override (silent, fatal, error, warn, info, debug, trace).
	 *  When set, overrides the auto-derived level from NODE_ENV.
	 *  Use LOG_LEVEL=silent for performance testing. */
	{ key: 'log_level', env: 'LOG_LEVEL', type: 'string', default: '' },

	/** MinIO/S3 Configuration */
	{ key: 'minio_endpoint', env: 'MINIO_ENDPOINT', type: 'string', default: 'localhost:9000' },
	{ key: 'minio_access_key', env: 'MINIO_ACCESS_KEY', type: 'string', default: 'mapexos_admin', sensitive: true },
	{ key: 'minio_secret_key', env: 'MINIO_SECRET_KEY', type: 'string', default: 'mapexos_admin_secret_change_me', sensitive: true },
	{ key: 'minio_use_ssl', env: 'MINIO_USE_SSL', type: 'bool', default: false },
	{ key: 'minio_region', env: 'MINIO_REGION', type: 'string', default: 'us-east-1' },
	{ key: 'minio_assets_bucket', env: 'MINIO_ASSETS_BUCKET', type: 'string', default: 'mapex-assets' },
	{ key: 'minio_templates_bucket', env: 'MINIO_TEMPLATES_BUCKET', type: 'string', default: 'mapex-templates' },
	{ key: 'minio_bytecode_bucket', env: 'MINIO_BYTECODE_BUCKET', type: 'string', default: 'mapex-bytecode' },

	/** TieredCache Configuration (L0=RAM, L1=Disk, L2=MinIO) */
	{ key: 'cache_l0_max_size', env: 'CACHE_L0_MAX_SIZE', type: 'int', default: 268435456 },  // 256MB
	{ key: 'cache_l0_max_items', env: 'CACHE_L0_MAX_ITEMS', type: 'int', default: 100000 },
	{ key: 'cache_l0_ttl_seconds', env: 'CACHE_L0_TTL_SECONDS', type: 'int', default: 300 },  // 5min
	{ key: 'cache_l1_enabled', env: 'CACHE_L1_ENABLED', type: 'bool', default: true },
	{ key: 'cache_l1_dir', env: 'CACHE_L1_DIR', type: 'string', default: '/tmp/mapexos/cache' },
	{ key: 'cache_l1_max_size', env: 'CACHE_L1_MAX_SIZE', type: 'int', default: 10737418240 }, // 10GB
	{ key: 'cache_l1_ttl_seconds', env: 'CACHE_L1_TTL_SECONDS', type: 'int', default: 3600 }, // 1h

	/** TieredCache Fallback Configuration (HTTP API when L2 misses) */
	{ key: 'cache_fallback_timeout', env: 'CACHE_FALLBACK_TIMEOUT', type: 'int', default: 5000 }, // ms

	/** Bytecode Cache - Skip L0 (RAM) since Script Registry already caches in RAM */
	{ key: 'cache_bytecode_skip_l0', env: 'CACHE_BYTECODE_SKIP_L0', type: 'bool', default: true },

	/** CPU Limit — the SINGLE knob for auto-tuning.
	 *  All worker/chunk/batch values are derived from this.
	 *  Set to match your pod/container CPU limit. */
	{ key: 'cpu_limit', env: 'CPU_LIMIT', type: 'int', default: 4 },

	/** Piscina Worker Threads (0 = auto: CPU_LIMIT - 1, minimum 1) */
	{ key: 'piscina_workers', env: 'PISCINA_WORKERS', type: 'int', default: 0 },

	/** V8 Isolate Configuration (per worker thread) */
	{ key: 'isolate_memory_limit_mb', env: 'ISOLATE_MEMORY_LIMIT_MB', type: 'int', default: 32 },

	/** Worker Script Execution Timeout in ms */
	{ key: 'worker_script_timeout_ms', env: 'WORKER_SCRIPT_TIMEOUT_MS', type: 'int', default: 10000 },

	/** Recycle V8 context every N events to prevent memory leaks */
	{ key: 'context_recycle_interval', env: 'CONTEXT_RECYCLE_INTERVAL', type: 'int', default: 10000 },

	/** NATS Consumer Tuning (0 = auto from CPU_LIMIT) */
	{ key: 'nats_consumer_batch_size', env: 'NATS_CONSUMER_BATCH_SIZE', type: 'int', default: 0 },          // 0 = auto (CPU_LIMIT × 500)
	{ key: 'nats_consumer_fetch_timeout', env: 'NATS_CONSUMER_FETCH_TIMEOUT', type: 'int', default: 1000 }, // 1000ms = 1s fetch timeout
	{ key: 'nats_consumer_max_ack_pending', env: 'NATS_CONSUMER_MAX_ACK_PENDING', type: 'int', default: 0 }, // 0 = auto (batchSize × 2)

	/** Concurrency chunk size (0 = auto from PISCINA_WORKERS × 8) */
	{ key: 'concurrency_chunk_size', env: 'CONCURRENCY_CHUNK_SIZE', type: 'int', default: 0 },

	/** Events per piscina.run() call for batch workers (0 = auto: 500) */
	{ key: 'events_per_worker', env: 'EVENTS_PER_WORKER', type: 'int', default: 0 },
];
