package configMap

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// DefaultConfiguration defines all configuration keys for the Triggers service.
//
// This configuration follows the centralized pattern from MapexOS service,
// providing a single source of truth for all service configuration.
//
// Configuration sources (in priority order):
//  1. Environment variables
//  2. Default values defined here
var DefaultConfiguration = []config.ConfigDefinition{

	/** HTTP Server Configuration */
	{Key: "http_port", Env: "HTTP_PORT", Type: "int", Default: 5006},
	{Key: "http_address", Env: "HTTP_ADDRESS", Type: "string", Default: "0.0.0.0"},

	{Key: "service_name", Env: "SERVICE_NAME", Type: "string", Default: "triggers"},
	{Key: "service_version", Env: "SERVICE_VERSION", Type: "string", Default: "1.0.0"},

	/* MongoDB Configuration */
	{Key: "mongo_uri", Env: "MONGO_URI", Type: "string", Default: "mongodb://localhost:27017"},
	{Key: "mongo_database", Env: "MONGO_DATABASE", Type: "string", Default: "triggers"},
	{Key: "mongo_pool_size", Env: "MONGO_POOL_SIZE", Type: "int", Default: 10},
	{Key: "mongo_monitor_interval", Env: "MONGO_MONITOR_INTERVAL", Type: "int", Default: 10},

	/* Redis Configuration */
	{Key: "redis_host", Env: "REDIS_HOST", Type: "string", Default: "localhost"},
	{Key: "redis_port", Env: "REDIS_PORT", Type: "int", Default: 6379},
	{Key: "redis_username", Env: "REDIS_USERNAME", Type: "string", Default: ""},
	{Key: "redis_password", Env: "REDIS_PASSWORD", Type: "string", Default: ""},
	{Key: "redis_db", Env: "REDIS_DB", Type: "int", Default: 0},
	{Key: "redis_shared_db", Env: "REDIS_SHARED_DB", Type: "int", Default: 5},

	/** NATS Configuration */
	{Key: "nats_url", Env: "NATS_URL", Type: "string", Default: "nats://localhost:4222"},
	{Key: "nats_username", Env: "NATS_USERNAME", Type: "string", Default: "service"},
	{Key: "nats_password", Env: "NATS_PASSWORD", Type: "string", Default: "service_secret", Sensitive: true},
	{Key: "nats_client_name", Env: "NATS_CLIENT_NAME", Type: "string", Default: "triggers-service"},
	{Key: "nats_batch_size", Env: "NATS_BATCH_SIZE", Type: "int", Default: 500},
	{Key: "nats_fetch_timeout", Env: "NATS_FETCH_TIMEOUT", Type: "int", Default: 1},

	/** Executor Configuration */
	{Key: "trigger_executor_workers", Env: "TRIGGER_EXECUTOR_WORKERS", Type: "int", Default: 50},

	/** Metrics Configuration */
	{Key: "metrics_go_collector", Env: "METRICS_GO_COLLECTOR", Type: "bool", Default: true},
	{Key: "metrics_process_collector", Env: "METRICS_PROCESS_COLLECTOR", Type: "bool", Default: true},

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
}
