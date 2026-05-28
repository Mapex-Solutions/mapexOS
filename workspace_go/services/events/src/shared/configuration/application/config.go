package configMap

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// DefaultConfiguration defines all configuration keys for the Events service.
//
// This configuration follows the centralized pattern from MapexOS service,
// providing a single source of truth for all service configuration.
//
// Configuration sources (in priority order):
//  1. Environment variables
//  2. Default values defined here
var DefaultConfiguration = []config.ConfigDefinition{

	/** HTTP Server Configuration */
	{Key: "http_port", Env: "HTTP_PORT", Type: "int", Default: 5004},
	{Key: "http_address", Env: "HTTP_ADDRESS", Type: "string", Default: "0.0.0.0"},

	{Key: "service_name", Env: "SERVICE_NAME", Type: "string", Default: "events"},
	{Key: "service_version", Env: "SERVICE_VERSION", Type: "string", Default: "1.0.0"},

	/* ClickHouse Configuration */
	{Key: "clickhouse_host", Env: "CLICKHOUSE_HOST", Type: "string", Default: "localhost"},
	{Key: "clickhouse_port", Env: "CLICKHOUSE_PORT", Type: "int", Default: 9440},
	{Key: "clickhouse_database", Env: "CLICKHOUSE_DATABASE", Type: "string", Default: "mapexos"},
	{Key: "clickhouse_username", Env: "CLICKHOUSE_USERNAME", Type: "string", Default: "mapexos_user"},
	{Key: "clickhouse_password", Env: "CLICKHOUSE_PASSWORD", Type: "string", Default: "mapexos_password", Sensitive: true},

	/** NATS Configuration */
	{Key: "nats_url", Env: "NATS_URL", Type: "string", Default: "nats://localhost:4222"},
	{Key: "nats_username", Env: "NATS_USERNAME", Type: "string", Default: "service"},
	{Key: "nats_password", Env: "NATS_PASSWORD", Type: "string", Default: "service_secret", Sensitive: true},
	{Key: "nats_client_name", Env: "NATS_CLIENT_NAME", Type: "string", Default: "events-service"},

	/**
	* Defines the authentication strategy used by the platform (jwt or oauth2) and
	* how user roles are retrieved (token, db, or api).
	* Supports static JWT secrets (HS256) and external identity providers using JWKS (RS256).
	* Allows dynamic role extraction from tokens, local databases, or external APIs.
	 */
	{Key: "auth_strategy", Env: "AUTH_STRATEGY", Type: "string", Default: "jwt"},
	{Key: "auth_secret", Env: "AUTH_SECRET", Type: "string", Default: "a-string-secret-at-least-256-bits-long", Sensitive: true},
	{Key: "auth_jwks_url", Env: "AUTH_JWKS_URL", Type: "string", Default: ""},
	{Key: "auth_algorithm", Env: "AUTH_ALGORITHM", Type: "string", Default: "HS256"},

	{Key: "auth_roles_source", Env: "AUTH_ROLES_SOURCE", Type: "string", Default: "token"}, // token | db | api
	{Key: "auth_roles_path", Env: "AUTH_ROLES_PATH", Type: "string", Default: "roles"},
	{Key: "auth_roles_api_url", Env: "AUTH_ROLES_API_URL", Type: "string", Default: ""},

	/**
	*	My apiKey for authentication
	* Please replace this with your own apiKey for authentication
	* This api Is used for internal communication between services
	 */
	{Key: "internal_api_key", Env: "INTERNAL_API_KEY", Type: "string", Default: "5230c2e2-e245-468d-89e8-94154cf520d0", Sensitive: true},

	/** Environment administration */
	{Key: "go_env", Env: "GO_ENV", Type: "string", Default: "dev"},

	// log_level overrides the default log level (debug, info, warn, error).
	// When empty, defaults to debug (dev) or info (production).
	{Key: "log_level", Env: "LOG_LEVEL", Type: "string", Default: ""},

	// ctx_timeout is used by the global middleware to establish a request timeout in seconds
	// value on the context, which is then passed down to all handlers.
	{Key: "ctx_timeout", Env: "CTX_TIMEOUT", Type: "int", Default: 4},

	/** Permission Middleware Configuration */
	{Key: "mapexos_url", Env: "MAPEXOS_URL", Type: "string", Default: "http://localhost:5000"},

	/** Assets Service Configuration for internal API calls */
	{Key: "assets_url", Env: "ASSETS_URL", Type: "string", Default: "http://localhost:5002"},

	/* MongoDB Configuration */
	{Key: "mongo_uri", Env: "MONGO_URI", Type: "string", Default: "mongodb://localhost:27017/?replicaSet=rs0"},
	{Key: "mongo_database", Env: "MONGO_DATABASE", Type: "string", Default: "events"},
	{Key: "mongo_pool_size", Env: "MONGO_POOL_SIZE", Type: "int", Default: 10},
	{Key: "mongo_monitor_interval", Env: "MONGO_MONITOR_INTERVAL", Type: "int", Default: 10},

	/* Redis Configuration */
	{Key: "redis_host", Env: "REDIS_HOST", Type: "string", Default: "localhost"},
	{Key: "redis_port", Env: "REDIS_PORT", Type: "int", Default: 6379},
	{Key: "redis_username", Env: "REDIS_USERNAME", Type: "string", Default: ""},
	{Key: "redis_password", Env: "REDIS_PASSWORD", Type: "string", Default: ""},
	{Key: "redis_db", Env: "REDIS_DB", Type: "int", Default: 0},
	{Key: "redis_shared_db", Env: "REDIS_SHARED_DB", Type: "int", Default: 5},

	/** NATS Batch Processing Configuration */
	// Same value used for NATS fetch batch and ClickHouse bulk insert.
	// ClickHouse recommends 10K-100K rows per insert for optimal throughput.
	{Key: "nats_batch_size", Env: "NATS_BATCH_SIZE", Type: "int", Default: 10000},
	{Key: "nats_fetch_timeout", Env: "NATS_FETCH_TIMEOUT", Type: "int", Default: 1}, // in seconds

	/** Metrics Configuration */
	{Key: "metrics_go_collector", Env: "METRICS_GO_COLLECTOR", Type: "bool", Default: true},
	{Key: "metrics_process_collector", Env: "METRICS_PROCESS_COLLECTOR", Type: "bool", Default: true},

	/** MinIO/S3 Configuration */
	{Key: "minio_endpoint", Env: "MINIO_ENDPOINT", Type: "string", Default: "localhost:9000"},
	{Key: "minio_access_key", Env: "MINIO_ACCESS_KEY", Type: "string", Default: "mapex_admin", Sensitive: true},
	{Key: "minio_secret_key", Env: "MINIO_SECRET_KEY", Type: "string", Default: "mapex_admin_secret_change_me", Sensitive: true},
	{Key: "minio_use_ssl", Env: "MINIO_USE_SSL", Type: "bool", Default: false},
	{Key: "minio_region", Env: "MINIO_REGION", Type: "string", Default: "us-east-1"},
	{Key: "minio_templates_bucket", Env: "MINIO_TEMPLATES_BUCKET", Type: "string", Default: "mapex-templates"},

	/** TieredCache Configuration (L0=RAM, L1=Disk, L2=S3) */
	{Key: "cache_l0_max_size", Env: "CACHE_L0_MAX_SIZE", Type: "int", Default: 268435456},  // 256MB
	{Key: "cache_l0_max_items", Env: "CACHE_L0_MAX_ITEMS", Type: "int", Default: 100000},
	{Key: "cache_l0_ttl_seconds", Env: "CACHE_L0_TTL_SECONDS", Type: "int", Default: 300},  // 5min
	{Key: "cache_l1_enabled", Env: "CACHE_L1_ENABLED", Type: "bool", Default: true},
	{Key: "cache_l1_dir", Env: "CACHE_L1_DIR", Type: "string", Default: "/tmp/mapexos/cache"},
	{Key: "cache_l1_max_size", Env: "CACHE_L1_MAX_SIZE", Type: "int", Default: 10737418240}, // 10GB
	{Key: "cache_l1_ttl_seconds", Env: "CACHE_L1_TTL_SECONDS", Type: "int", Default: 3600}, // 1h
}
