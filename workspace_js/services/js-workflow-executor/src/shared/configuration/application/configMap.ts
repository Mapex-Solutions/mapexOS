import { ConfigDefinition } from '@mapexos/microservices';

import {
	HTTP_PORT_DEFAULT,
	HTTP_ADDRESS_DEFAULT,
	SERVICE_NAME_DEFAULT,
	SERVICE_VERSION_DEFAULT,
	NATS_URL_DEFAULT,
	NATS_USERNAME_DEFAULT,
	NATS_PASSWORD_DEFAULT,
	NATS_CLIENT_NAME_DEFAULT,
	AUTH_STRATEGY_DEFAULT,
	AUTH_SECRET_DEFAULT,
	AUTH_JWKS_URL_DEFAULT,
	AUTH_ALGORITHM_DEFAULT,
	AUTH_ROLES_SOURCE_DEFAULT,
	AUTH_ROLES_PATH_DEFAULT,
	AUTH_ROLES_API_URL_DEFAULT,
	INTERNAL_API_KEY_DEFAULT,
	WORKFLOW_SERVICE_URL_DEFAULT,
	NODE_ENV_DEFAULT,
	LOG_LEVEL_DEFAULT,
	OBJECT_STORE_ENDPOINT_DEFAULT,
	OBJECT_STORE_ACCESS_KEY_DEFAULT,
	OBJECT_STORE_SECRET_KEY_DEFAULT,
	OBJECT_STORE_USE_SSL_DEFAULT,
	OBJECT_STORE_REGION_DEFAULT,
	OBJECT_STORE_AUTH_IS_NEEDED_DEFAULT,
	MINIO_WORKFLOWS_BUCKET_DEFAULT,
	CACHE_L0_MAX_SIZE_DEFAULT,
	CACHE_L0_MAX_ITEMS_DEFAULT,
	CACHE_L0_TTL_SECONDS_DEFAULT,
	CACHE_L1_ENABLED_DEFAULT,
	CACHE_L1_DIR_DEFAULT,
	CACHE_L1_MAX_SIZE_DEFAULT,
	CACHE_L1_TTL_SECONDS_DEFAULT,
	CACHE_FALLBACK_TIMEOUT_DEFAULT,
	CACHE_BYTECODE_SKIP_L0_DEFAULT,
	CPU_LIMIT_DEFAULT,
	PISCINA_WORKERS_DEFAULT,
	ISOLATE_MEMORY_LIMIT_MB_DEFAULT,
	WORKER_SCRIPT_TIMEOUT_MS_DEFAULT,
	CONTEXT_RECYCLE_INTERVAL_DEFAULT,
	NATS_CONSUMER_BATCH_SIZE_DEFAULT,
	NATS_CONSUMER_FETCH_TIMEOUT_DEFAULT,
	NATS_CONSUMER_MAX_ACK_PENDING_DEFAULT,
} from './constants';

// InitConfig initializes the configuration module with the provided definitions.
export const defaultConfiguration: ConfigDefinition[] = [

	/** HTTP Server Configuration */
	{ key: 'http_port', env: 'HTTP_PORT', type: 'int', default: HTTP_PORT_DEFAULT },
	{ key: 'http_address', env: 'HTTP_ADDRESS', type: 'string', default: HTTP_ADDRESS_DEFAULT },

	{ key: 'service_name', env: 'SERVICE_NAME', type: 'string', default: SERVICE_NAME_DEFAULT },
	{ key: 'service_version', env: 'SERVICE_VERSION', type: 'string', default: SERVICE_VERSION_DEFAULT },

	/** NATS Configuration */
	{ key: 'nats_url', env: 'NATS_URL', type: 'string', default: NATS_URL_DEFAULT },
	{ key: 'nats_username', env: 'NATS_USERNAME', type: 'string', default: NATS_USERNAME_DEFAULT },
	{ key: 'nats_password', env: 'NATS_PASSWORD', type: 'string', default: NATS_PASSWORD_DEFAULT, sensitive: true },
	{ key: 'nats_client_name', env: 'NATS_CLIENT_NAME', type: 'string', default: NATS_CLIENT_NAME_DEFAULT },

	/**
	 * Defines the authentication strategy used by the platform (jwt or oauth2) and
	 * how user roles are retrieved (token, db, or api).
	 */
	{ key: 'auth_strategy', env: 'AUTH_STRATEGY', type: 'string', default: AUTH_STRATEGY_DEFAULT },
	{ key: 'auth_secret', env: 'AUTH_SECRET', type: 'string', default: AUTH_SECRET_DEFAULT, sensitive: true },
	{ key: 'auth_jwks_url', env: 'AUTH_JWKS_URL', type: 'string', default: AUTH_JWKS_URL_DEFAULT },
	{ key: 'auth_algorithm', env: 'AUTH_ALGORITHM', type: 'string', default: AUTH_ALGORITHM_DEFAULT },

	{ key: 'auth_roles_source', env: 'AUTH_ROLES_SOURCE', type: 'string', default: AUTH_ROLES_SOURCE_DEFAULT },
	{ key: 'auth_roles_path', env: 'AUTH_ROLES_PATH', type: 'string', default: AUTH_ROLES_PATH_DEFAULT },
	{ key: 'auth_roles_api_url', env: 'AUTH_ROLES_API_URL', type: 'string', default: AUTH_ROLES_API_URL_DEFAULT },

	/** Internal API key for service-to-service communication */
	{ key: 'internal_api_key', env: 'INTERNAL_API_KEY', type: 'string', default: INTERNAL_API_KEY_DEFAULT, sensitive: true },

	/** Workflow Service URL (fallback HTTP to fetch script source when L2 misses) */
	{ key: 'workflow_service_url', env: 'WORKFLOW_SERVICE_URL', type: 'string', default: WORKFLOW_SERVICE_URL_DEFAULT },

	/** Environment administration */
	{ key: 'node_env', env: 'NODE_ENV', type: 'string', default: NODE_ENV_DEFAULT },

	/** Log level override (silent, fatal, error, warn, info, debug, trace). */
	{ key: 'log_level', env: 'LOG_LEVEL', type: 'string', default: LOG_LEVEL_DEFAULT },

	/** MinIO/S3 Configuration */
	{ key: 'object_store_endpoint', env: 'OBJECT_STORE_ENDPOINT', type: 'string', default: OBJECT_STORE_ENDPOINT_DEFAULT },
	{ key: 'object_store_access_key', env: 'OBJECT_STORE_ACCESS_KEY', type: 'string', default: OBJECT_STORE_ACCESS_KEY_DEFAULT, sensitive: true },
	{ key: 'object_store_secret_key', env: 'OBJECT_STORE_SECRET_KEY', type: 'string', default: OBJECT_STORE_SECRET_KEY_DEFAULT, sensitive: true },
	{ key: 'object_store_use_ssl', env: 'OBJECT_STORE_USE_SSL', type: 'bool', default: OBJECT_STORE_USE_SSL_DEFAULT },
	{ key: 'object_store_region', env: 'OBJECT_STORE_REGION', type: 'string', default: OBJECT_STORE_REGION_DEFAULT },
	{ key: 'object_store_auth_is_needed', env: 'OBJECT_STORE_AUTH_IS_NEEDED', type: 'bool', default: OBJECT_STORE_AUTH_IS_NEEDED_DEFAULT },
	{ key: 'minio_workflows_bucket', env: 'MINIO_WORKFLOWS_BUCKET', type: 'string', default: MINIO_WORKFLOWS_BUCKET_DEFAULT },

	/** TieredCache Configuration (L0=RAM, L1=Disk, L2=MinIO) */
	{ key: 'cache_l0_max_size', env: 'CACHE_L0_MAX_SIZE', type: 'int', default: CACHE_L0_MAX_SIZE_DEFAULT },  // 256MB
	{ key: 'cache_l0_max_items', env: 'CACHE_L0_MAX_ITEMS', type: 'int', default: CACHE_L0_MAX_ITEMS_DEFAULT },
	{ key: 'cache_l0_ttl_seconds', env: 'CACHE_L0_TTL_SECONDS', type: 'int', default: CACHE_L0_TTL_SECONDS_DEFAULT },  // 5min
	{ key: 'cache_l1_enabled', env: 'CACHE_L1_ENABLED', type: 'bool', default: CACHE_L1_ENABLED_DEFAULT },
	{ key: 'cache_l1_dir', env: 'CACHE_L1_DIR', type: 'string', default: CACHE_L1_DIR_DEFAULT },
	{ key: 'cache_l1_max_size', env: 'CACHE_L1_MAX_SIZE', type: 'int', default: CACHE_L1_MAX_SIZE_DEFAULT }, // 10GB
	{ key: 'cache_l1_ttl_seconds', env: 'CACHE_L1_TTL_SECONDS', type: 'int', default: CACHE_L1_TTL_SECONDS_DEFAULT }, // 1h

	/** TieredCache Fallback Configuration (HTTP API when L2 misses) */
	{ key: 'cache_fallback_timeout', env: 'CACHE_FALLBACK_TIMEOUT', type: 'int', default: CACHE_FALLBACK_TIMEOUT_DEFAULT }, // ms

	/** Bytecode Cache - Skip L0 (RAM) since Script Registry already caches in RAM */
	{ key: 'cache_bytecode_skip_l0', env: 'CACHE_BYTECODE_SKIP_L0', type: 'bool', default: CACHE_BYTECODE_SKIP_L0_DEFAULT },

	/** CPU Limit — the SINGLE knob for auto-tuning.
	 *  Set to match your pod/container CPU limit. */
	{ key: 'cpu_limit', env: 'CPU_LIMIT', type: 'int', default: CPU_LIMIT_DEFAULT },

	/** Piscina Worker Threads (0 = auto: CPU_LIMIT - 1, minimum 1) */
	{ key: 'piscina_workers', env: 'PISCINA_WORKERS', type: 'int', default: PISCINA_WORKERS_DEFAULT },

	/** V8 Isolate Configuration (per worker thread) */
	{ key: 'isolate_memory_limit_mb', env: 'ISOLATE_MEMORY_LIMIT_MB', type: 'int', default: ISOLATE_MEMORY_LIMIT_MB_DEFAULT },

	/** Worker Script Execution Timeout in ms */
	{ key: 'worker_script_timeout_ms', env: 'WORKER_SCRIPT_TIMEOUT_MS', type: 'int', default: WORKER_SCRIPT_TIMEOUT_MS_DEFAULT },

	/** Recycle V8 context every N events to prevent memory leaks */
	{ key: 'context_recycle_interval', env: 'CONTEXT_RECYCLE_INTERVAL', type: 'int', default: CONTEXT_RECYCLE_INTERVAL_DEFAULT },

	/** NATS Consumer Tuning (0 = auto from CPU_LIMIT) */
	{ key: 'nats_consumer_batch_size', env: 'NATS_CONSUMER_BATCH_SIZE', type: 'int', default: NATS_CONSUMER_BATCH_SIZE_DEFAULT },
	{ key: 'nats_consumer_fetch_timeout', env: 'NATS_CONSUMER_FETCH_TIMEOUT', type: 'int', default: NATS_CONSUMER_FETCH_TIMEOUT_DEFAULT }, // 1000ms = 1s fetch timeout
	{ key: 'nats_consumer_max_ack_pending', env: 'NATS_CONSUMER_MAX_ACK_PENDING', type: 'int', default: NATS_CONSUMER_MAX_ACK_PENDING_DEFAULT },
];
