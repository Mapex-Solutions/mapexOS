package configMap

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// DefaultConfiguration defines all configuration keys for the Workflow service.
//
// Configuration sources (in priority order):
//  1. Environment variables
//  2. Default values defined here
var DefaultConfiguration = []config.ConfigDefinition{

	/* HTTP Server Configuration */
	{Key: "http_port", Env: "HTTP_PORT", Type: "int", Default: 5007},
	{Key: "http_address", Env: "HTTP_ADDRESS", Type: "string", Default: "0.0.0.0"},

	{Key: "service_name", Env: "SERVICE_NAME", Type: "string", Default: "workflow"},
	{Key: "service_version", Env: "SERVICE_VERSION", Type: "string", Default: "1.0.0"},

	/* MongoDB Configuration */
	{Key: "mongo_uri", Env: "MONGO_URI", Type: "string", Default: "mongodb://localhost:27017/?replicaSet=rs0"},
	{Key: "mongo_database", Env: "MONGO_DATABASE", Type: "string", Default: "workflow"},
	{Key: "mongo_pool_size", Env: "MONGO_POOL_SIZE", Type: "int", Default: 10},
	{Key: "mongo_monitor_interval", Env: "MONGO_MONITOR_INTERVAL", Type: "int", Default: 10},

	/* Redis Configuration */
	// Workflow service only uses Shared Redis (DB 5) for authorization middleware.
	// NO App Redis — TieredCache + MinIO handles all caching needs.
	{Key: "redis_host", Env: "REDIS_HOST", Type: "string", Default: "localhost"},
	{Key: "redis_port", Env: "REDIS_PORT", Type: "int", Default: 6379},
	{Key: "redis_username", Env: "REDIS_USERNAME", Type: "string", Default: ""},
	{Key: "redis_password", Env: "REDIS_PASSWORD", Type: "string", Default: ""},
	{Key: "redis_shared_db", Env: "REDIS_SHARED_DB", Type: "int", Default: 5},

	/* NATS Configuration (Core - JetStream, Domain Events) */
	{Key: "nats_url", Env: "NATS_URL", Type: "string", Default: "nats://localhost:4222"},
	{Key: "nats_username", Env: "NATS_USERNAME", Type: "string", Default: "service"},
	{Key: "nats_password", Env: "NATS_PASSWORD", Type: "string", Default: "service_secret", Sensitive: true},
	{Key: "nats_client_name", Env: "NATS_CLIENT_NAME", Type: "string", Default: "workflow-service"},

	/*
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

	/*
	*	Internal API key for service-to-service communication
	 */
	{Key: "internal_api_key", Env: "INTERNAL_API_KEY", Type: "string", Default: "5230c2e2-e245-468d-89e8-94154cf520d0", Sensitive: true},

	/* Vault MS Configuration */
	{Key: "vault_url", Env: "VAULT_URL", Type: "string", Default: "http://localhost:5010"},

	/* Environment administration */
	{Key: "go_env", Env: "GO_ENV", Type: "string", Default: "dev"},
	{Key: "log_level", Env: "LOG_LEVEL", Type: "string", Default: ""},

	// ctx_timeout is used by the global middleware to establish a request timeout in seconds
	{Key: "ctx_timeout", Env: "CTX_TIMEOUT", Type: "int", Default: 4},

	/* Metrics Configuration */
	{Key: "metrics_go_collector", Env: "METRICS_GO_COLLECTOR", Type: "bool", Default: true},
	{Key: "metrics_process_collector", Env: "METRICS_PROCESS_COLLECTOR", Type: "bool", Default: true},

	/* Permission Middleware Configuration */
	{Key: "mapexos_url", Env: "MAPEXOS_URL", Type: "string", Default: "http://localhost:5000"},

	/* MinIO/S3 Configuration */
	{Key: "minio_endpoint", Env: "MINIO_ENDPOINT", Type: "string", Default: "localhost:9000"},
	{Key: "minio_access_key", Env: "MINIO_ACCESS_KEY", Type: "string", Default: "mapex_admin", Sensitive: true},
	{Key: "minio_secret_key", Env: "MINIO_SECRET_KEY", Type: "string", Default: "mapex_admin_secret_change_me", Sensitive: true},
	{Key: "minio_use_ssl", Env: "MINIO_USE_SSL", Type: "bool", Default: false},
	{Key: "minio_region", Env: "MINIO_REGION", Type: "string", Default: "us-east-1"},
	{Key: "minio_definitions_bucket", Env: "MINIO_DEFINITIONS_BUCKET", Type: "string", Default: "mapex-workflows"},

	/* TieredCache Configuration — Definitions (L0=RAM, L1=Disk, L2=MinIO) */
	{Key: "cache_l0_max_size", Env: "CACHE_L0_MAX_SIZE", Type: "int", Default: 268435456}, // 256MB
	{Key: "cache_l0_max_items", Env: "CACHE_L0_MAX_ITEMS", Type: "int", Default: 100000},
	{Key: "cache_l0_ttl_seconds", Env: "CACHE_L0_TTL_SECONDS", Type: "int", Default: 300}, // 5min
	{Key: "cache_l1_dir", Env: "CACHE_L1_DIR", Type: "string", Default: "/tmp/mapexos/cache"},
	{Key: "cache_l1_max_size", Env: "CACHE_L1_MAX_SIZE", Type: "int", Default: 10737418240}, // 10GB
	{Key: "cache_l1_ttl_seconds", Env: "CACHE_L1_TTL_SECONDS", Type: "int", Default: 3600},  // 1h

	/* TieredCache Configuration — Plugins (L0=RAM, L1=Disk, NO L2) */
	{Key: "plugins_cache_l0_max_size", Env: "PLUGINS_CACHE_L0_MAX_SIZE", Type: "int", Default: 5242880}, // 5MB
	{Key: "plugins_cache_l0_max_items", Env: "PLUGINS_CACHE_L0_MAX_ITEMS", Type: "int", Default: 1000},
	{Key: "plugins_cache_l0_ttl_seconds", Env: "PLUGINS_CACHE_L0_TTL_SECONDS", Type: "int", Default: 300},  // 5min
	{Key: "plugins_cache_l1_max_size", Env: "PLUGINS_CACHE_L1_MAX_SIZE", Type: "int", Default: 5368709120}, // 5GB
	{Key: "plugins_cache_l1_ttl_seconds", Env: "PLUGINS_CACHE_L1_TTL_SECONDS", Type: "int", Default: 1800}, // 30min

	/* TieredCache Configuration — Instances (L0=RAM, L1=Disk, NO L2) */
	{Key: "instances_cache_l0_max_size", Env: "INSTANCES_CACHE_L0_MAX_SIZE", Type: "int", Default: 10485760}, // 10MB
	{Key: "instances_cache_l0_max_items", Env: "INSTANCES_CACHE_L0_MAX_ITEMS", Type: "int", Default: 50000},
	{Key: "instances_cache_l0_ttl_seconds", Env: "INSTANCES_CACHE_L0_TTL_SECONDS", Type: "int", Default: 300}, // 5min
	{Key: "instances_cache_l1_max_size", Env: "INSTANCES_CACHE_L1_MAX_SIZE", Type: "int", Default: 5368709120}, // 5GB
	{Key: "instances_cache_l1_ttl_seconds", Env: "INSTANCES_CACHE_L1_TTL_SECONDS", Type: "int", Default: 1800}, // 30min

	/* Credential Encryption (Envelope: Master Key → DEK → Data) */
	{Key: "credential_master_key", Env: "CREDENTIAL_MASTER_KEY", Type: "string", Default: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", Sensitive: true},
}
