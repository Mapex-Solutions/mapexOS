package configMap

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// DefaultConfiguration defines all configuration keys for the Vault service.
var DefaultConfiguration = []config.ConfigDefinition{

	/* HTTP Server Configuration */
	{Key: "http_port", Env: "HTTP_PORT", Type: "int", Default: 5010},
	{Key: "http_address", Env: "HTTP_ADDRESS", Type: "string", Default: "0.0.0.0"},

	{Key: "service_name", Env: "SERVICE_NAME", Type: "string", Default: "mapexVault"},
	{Key: "service_version", Env: "SERVICE_VERSION", Type: "string", Default: "1.0.0"},

	/* MongoDB Configuration */
	{Key: "mongo_uri", Env: "MONGO_URI", Type: "string", Default: "mongodb://localhost:27017/?replicaSet=rs0"},
	{Key: "mongo_database", Env: "MONGO_DATABASE", Type: "string", Default: "mapex_vault"},
	{Key: "mongo_pool_size", Env: "MONGO_POOL_SIZE", Type: "int", Default: 10},
	{Key: "mongo_monitor_interval", Env: "MONGO_MONITOR_INTERVAL", Type: "int", Default: 10},

	/* Redis Configuration (Shared DB5 for authorization middleware) */
	{Key: "redis_host", Env: "REDIS_HOST", Type: "string", Default: "localhost"},
	{Key: "redis_port", Env: "REDIS_PORT", Type: "int", Default: 6379},
	{Key: "redis_username", Env: "REDIS_USERNAME", Type: "string", Default: ""},
	{Key: "redis_password", Env: "REDIS_PASSWORD", Type: "string", Default: ""},
	{Key: "redis_shared_db", Env: "REDIS_SHARED_DB", Type: "int", Default: 5},

	/* NATS Configuration */
	{Key: "nats_url", Env: "NATS_URL", Type: "string", Default: "nats://localhost:4222"},
	{Key: "nats_username", Env: "NATS_USERNAME", Type: "string", Default: "service"},
	{Key: "nats_password", Env: "NATS_PASSWORD", Type: "string", Default: "service_secret", Sensitive: true},
	{Key: "nats_client_name", Env: "NATS_CLIENT_NAME", Type: "string", Default: "mapexVault-service"},

	/* Authentication (JWT for external API) */
	{Key: "auth_strategy", Env: "AUTH_STRATEGY", Type: "string", Default: "jwt"},
	{Key: "auth_secret", Env: "AUTH_SECRET", Type: "string", Default: "a-string-secret-at-least-256-bits-long", Sensitive: true},
	{Key: "auth_jwks_url", Env: "AUTH_JWKS_URL", Type: "string", Default: ""},
	{Key: "auth_algorithm", Env: "AUTH_ALGORITHM", Type: "string", Default: "HS256"},

	{Key: "auth_roles_source", Env: "AUTH_ROLES_SOURCE", Type: "string", Default: "token"},
	{Key: "auth_roles_path", Env: "AUTH_ROLES_PATH", Type: "string", Default: "roles"},
	{Key: "auth_roles_api_url", Env: "AUTH_ROLES_API_URL", Type: "string", Default: ""},

	/* Internal API key for service-to-service communication */
	{Key: "internal_api_key", Env: "INTERNAL_API_KEY", Type: "string", Default: "5230c2e2-e245-468d-89e8-94154cf520d0", Sensitive: true},

	/* Environment */
	{Key: "go_env", Env: "GO_ENV", Type: "string", Default: "dev"},
	{Key: "log_level", Env: "LOG_LEVEL", Type: "string", Default: ""},
	{Key: "ctx_timeout", Env: "CTX_TIMEOUT", Type: "int", Default: 4},

	/* Metrics */
	{Key: "metrics_go_collector", Env: "METRICS_GO_COLLECTOR", Type: "bool", Default: true},
	{Key: "metrics_process_collector", Env: "METRICS_PROCESS_COLLECTOR", Type: "bool", Default: true},

	/* Permission Middleware */
	{Key: "mapexos_url", Env: "MAPEXOS_URL", Type: "string", Default: "http://localhost:5000"},

	/* Credential Encryption (Envelope: Master Key → DEK → Data) */
	// Must match the value used by scripts/prebuild/pki/seed-encryptor at
	// PKI generation time — the service uses this key to decrypt the CA
	// private keys seeded into pkiCertificateAuthorities.
	{Key: "credential_master_key", Env: "CREDENTIAL_MASTER_KEY", Type: "string", Default: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", Sensitive: true},
}
