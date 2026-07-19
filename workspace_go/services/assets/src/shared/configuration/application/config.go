package configMap

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// DefaultConfiguration defines all configuration keys for the Assets service.
//
// This configuration follows the centralized pattern from MapexOS service,
// providing a single source of truth for all service configuration.
//
// Configuration sources (in priority order):
//  1. Environment variables
//  2. Default values defined here
var DefaultConfiguration = []config.ConfigDefinition{

	/** HTTP Server Configuration */
	{Key: "http_port", Env: "HTTP_PORT", Type: "int", Default: 5002},
	{Key: "http_address", Env: "HTTP_ADDRESS", Type: "string", Default: "0.0.0.0"},

	{Key: "service_name", Env: "SERVICE_NAME", Type: "string", Default: "assets"},
	{Key: "service_version", Env: "SERVICE_VERSION", Type: "string", Default: "1.1.0"},

	/* MongoDB Configuration */
	{Key: "mongo_uri", Env: "MONGO_URI", Type: "string", Default: "mongodb://localhost:27017/?replicaSet=rs0"},
	{Key: "mongo_database", Env: "MONGO_DATABASE", Type: "string", Default: "assets"},
	{Key: "mongo_pool_size", Env: "MONGO_POOL_SIZE", Type: "int", Default: 10},
	{Key: "mongo_monitor_interval", Env: "MONGO_MONITOR_INTERVAL", Type: "int", Default: 10},

	/* Redis Configuration */
	{Key: "redis_host", Env: "REDIS_HOST", Type: "string", Default: "localhost"},
	{Key: "redis_port", Env: "REDIS_PORT", Type: "int", Default: 6379},
	{Key: "redis_username", Env: "REDIS_USERNAME", Type: "string", Default: ""},
	{Key: "redis_password", Env: "REDIS_PASSWORD", Type: "string", Default: ""},
	{Key: "redis_db", Env: "REDIS_DB", Type: "int", Default: 0},
	{Key: "redis_shared_db", Env: "REDIS_SHARED_DB", Type: "int", Default: 5},

	/** NATS Configuration (Core only — JetStream + presence consumer) */
	// Single NATS connection to CORE. The assets MS authenticates as the
	// shared 'service' user. Device-side MQTT auth runs INSIDE the
	// mapex-mqtt-broker plugin off the AssetReadModel (no HTTP callout
	// to this service); presence advisories arrive via the broker plugin
	// on mapexos.presence.advisory. Neither path uses the
	// service-account NATS connection.
	{Key: "nats_url", Env: "NATS_URL", Type: "string", Default: "nats://localhost:4222"},
	{Key: "nats_username", Env: "NATS_USERNAME", Type: "string", Default: "service"},
	{Key: "nats_password", Env: "NATS_PASSWORD", Type: "string", Default: "service_secret", Sensitive: true},
	{Key: "nats_client_name", Env: "NATS_CLIENT_NAME", Type: "string", Default: "assets-service"},
	{Key: "nats_kv_replicas", Env: "NATS_KV_REPLICAS", Type: "int", Default: 1}, // 1 = single-node safe; a NATS cluster can raise it for KV/lease HA

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
	{Key: "log_level", Env: "LOG_LEVEL", Type: "string", Default: ""},

	// ctx_timeout is used by the global middleware to establish a request timeout in seconds
	// value on the context, which is then passed down to all handlers.
	{Key: "ctx_timeout", Env: "CTX_TIMEOUT", Type: "int", Default: 4},

	/** Metrics Configuration */
	{Key: "metrics_go_collector", Env: "METRICS_GO_COLLECTOR", Type: "bool", Default: true},
	{Key: "metrics_process_collector", Env: "METRICS_PROCESS_COLLECTOR", Type: "bool", Default: true},

	/** Permission Middleware Configuration */
	{Key: "mapexos_url", Env: "MAPEXOS_URL", Type: "string", Default: "http://localhost:5000"},

	/** Internal Service URLs */
	{Key: "router_service_url", Env: "ROUTER_SERVICE_URL", Type: "string", Default: "http://localhost:5003"},

	/** mapexVault Configuration */
	// mqttcerts.OnMount calls GET /internal/pki/intermediate_ca_bundle
	// against this URL with the vault API key on every boot. The vault
	// API key is DISTINCT from this service's internal_api_key so each
	// side can be rotated independently.
	{Key: "mapex_vault_url", Env: "MAPEX_VAULT_URL", Type: "string", Default: "http://localhost:5010"},
	{Key: "mapex_vault_api_key", Env: "MAPEX_VAULT_API_KEY", Type: "string", Default: "5230c2e2-e245-468d-89e8-94154cf520d0", Sensitive: true},

	/** Marketplace Configuration */
	// The asset_templates marketplace catalog is open and unauthenticated (no API
	// key). The install client fetches a bundle by (vendor, slug) from here. Defaults
	// to the deployed marketplace so it matches the UI, which points at the same host;
	// override per environment to target a local catalog. Base host only — the client
	// appends /api/v1/asset_templates/{vendor}/{slug}.
	{Key: "asset_marketplace_url", Env: "ASSET_MARKETPLACE_URL", Type: "string", Default: "https://marketplace.mapexos.io"},

	/** mapexIam Configuration */
	// The install resolves org-scoped classification through mapexIam's internal
	// /internal/lists/resolve endpoint, authenticated with internal_api_key above.
	{Key: "mapexiam_url", Env: "MAPEXIAM_URL", Type: "string", Default: "http://localhost:5000"},

	/** MinIO/S3 Configuration */
	{Key: "object_store_endpoint", Env: "OBJECT_STORE_ENDPOINT", Type: "string", Default: "localhost:9000"},
	{Key: "object_store_access_key", Env: "OBJECT_STORE_ACCESS_KEY", Type: "string", Default: "svc-assets", Sensitive: true},
	{Key: "object_store_secret_key", Env: "OBJECT_STORE_SECRET_KEY", Type: "string", Default: "svc-assets-secret-change-me", Sensitive: true},
	{Key: "object_store_use_ssl", Env: "OBJECT_STORE_USE_SSL", Type: "bool", Default: false},
	{Key: "object_store_region", Env: "OBJECT_STORE_REGION", Type: "string", Default: "us-east-1"},
	{Key: "object_store_auth_is_needed", Env: "OBJECT_STORE_AUTH_IS_NEEDED", Type: "bool", Default: true},
	{Key: "minio_templates_bucket", Env: "MINIO_TEMPLATES_BUCKET", Type: "string", Default: "mapex-templates"},
	{Key: "minio_assets_bucket", Env: "MINIO_ASSETS_BUCKET", Type: "string", Default: "mapex-assets"},
	{Key: "minio_asset_auth_bucket", Env: "MINIO_ASSET_AUTH_BUCKET", Type: "string", Default: "mapex-asset-auth"},

	/**
	 * OTA Firmware Store (S3-compatible).
	 *
	 * Independent from the TieredCache MinIO above so firmware artifacts can
	 * live in a different bucket — or a different provider (e.g. AWS S3) —
	 * without touching the read-model/auth caches. The client speaks the S3
	 * protocol via minio-go, so pointing at real S3 is config-only: set the
	 * endpoint/region/use_ssl/credentials here. Defaults mirror the local
	 * MinIO so dev works out of the box.
	 */
	// Firmware store shares the OBJECT_STORE_* credential (svc-assets) but keeps
	// its own endpoint — the presigned URLs it mints must be browser-reachable.
	{Key: "firmware_store_endpoint", Env: "FIRMWARE_STORE_ENDPOINT", Type: "string", Default: "localhost:9000"},
	{Key: "firmware_store_use_ssl", Env: "FIRMWARE_STORE_USE_SSL", Type: "bool", Default: false},
	{Key: "firmware_store_region", Env: "FIRMWARE_STORE_REGION", Type: "string", Default: "us-east-1"},
	{Key: "firmware_store_bucket", Env: "FIRMWARE_STORE_BUCKET", Type: "string", Default: "mapex-firmware"},

	/**
	 * OTA behavior (remote firmware update). TTLs and intervals are in SECONDS;
	 * pacing/retry are counts. All overridable per environment.
	 */
	{Key: "ota_presigned_url_ttl", Env: "OTA_PRESIGNED_URL_TTL", Type: "int", Default: 1800}, // presigned PUT/GET TTL (30m)
	{Key: "ota_orphan_gc_ttl", Env: "OTA_ORPHAN_GC_TTL", Type: "int", Default: 86400},        // abandon-check for never-finalized uploads (24h)
	{Key: "ota_scan_interval", Env: "OTA_SCAN_INTERVAL", Type: "int", Default: 60},           // global pacing scan tick (1m)
	{Key: "ota_rate_per_minute", Env: "OTA_RATE_PER_MINUTE", Type: "int", Default: 60},       // default dispatch rate when a plan omits it
	{Key: "ota_max_attempts", Env: "OTA_MAX_ATTEMPTS", Type: "int", Default: 2},              // retry cap per execution
	{Key: "ota_live_state_ttl", Env: "OTA_LIVE_STATE_TTL", Type: "int", Default: 86400},      // Redis live execution state TTL (24h)

	/** TieredCache Configuration (L0=RAM, L1=Disk, L2=MinIO) */
	{Key: "cache_l0_max_size", Env: "CACHE_L0_MAX_SIZE", Type: "int", Default: 268435456}, // 256MB
	{Key: "cache_l0_max_items", Env: "CACHE_L0_MAX_ITEMS", Type: "int", Default: 100000},
	{Key: "cache_l0_ttl_seconds", Env: "CACHE_L0_TTL_SECONDS", Type: "int", Default: 300}, // 5min
	{Key: "cache_l1_enabled", Env: "CACHE_L1_ENABLED", Type: "bool", Default: true},
	{Key: "cache_l1_dir", Env: "CACHE_L1_DIR", Type: "string", Default: "/tmp/mapexos/cache"},
	{Key: "cache_l1_max_size", Env: "CACHE_L1_MAX_SIZE", Type: "int", Default: 10737418240}, // 10GB
	{Key: "cache_l1_ttl_seconds", Env: "CACHE_L1_TTL_SECONDS", Type: "int", Default: 3600},  // 1h

	/** Health Monitor Configuration */
	{Key: "health_monitor_scan_interval", Env: "HEALTH_MONITOR_SCAN_INTERVAL", Type: "int", Default: 60}, // value in SECONDS (600 = 10 minutes)
	{Key: "health_monitor_batch_size", Env: "HEALTH_MONITOR_BATCH_SIZE", Type: "int", Default: 500},
}
