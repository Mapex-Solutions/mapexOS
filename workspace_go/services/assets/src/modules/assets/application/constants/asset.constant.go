package constants

import "time"

// HealthStatusAllowedRouterKinds enumerates the router kinds permitted to
// appear inside route groups referenced by HealthMonitor config. Must stay
// aligned with router-side constants.HealthStatusAllowedRouterKinds.
// Used by validateHealthMonitorConfig at Create/Update time.
//
// Cross-service FANOUT subjects/streams previously declared here now live in
// packages/contracts/services/assets/assets/constants.go — NATS subjects
// consumed by Router, JS-Executor, and Events are cross-service and cannot
// live in a service-local application/constants.
var HealthStatusAllowedRouterKinds = map[string]bool{
	"trigger":  true,
	"workflow": true,
}

/**
 * CACHE LIFETIME (application-level behavior)
 *
 * These TTLs express "how long the application allows cached values to
 * remain authoritative" — they are application behavior, not a property
 * of any specific cache technology (any KV cache with TTL semantics can
 * enforce them). Redis key construction (prefixes, key format) is an
 * infrastructure concern and stays in infrastructure/cache/redis.
 */
const (
	// AuthCacheTTL is the lifetime for Auth Callout cache entries that
	// the broker callout consults to validate device CONNECT requests.
	AuthCacheTTL = 24 * time.Hour

	// CounterCacheTTL is the lifetime for per-org asset counter cache
	// entries. 6 hours.
	CounterCacheTTL = 6 * time.Hour
)

/**
 * MQTT CREDENTIAL DEFAULTS (application-level behavior)
 */
const (
	// MqttPasswordLength is the length (in characters) of the random
	// alphanumeric password the platform generates for the operator on
	// the create / change-password flow. Long enough to make brute
	// force on the broker auth callout impractical, short enough to
	// type once when copy-paste is unavailable (factory provisioning).
	MqttPasswordLength = 24

	// MqttPasswordAlphabet is the character set used by the platform
	// password generator. Lower + upper + digits — symbols are excluded
	// because some broker firmwares strip them on CONNECT.
	MqttPasswordAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// MqttPasswordBcryptCost is the bcrypt cost applied when hashing the
	// stored MQTT password. Picked so a single hash on commodity hardware
	// stays under ~250ms — high enough to slow offline cracking, low
	// enough to keep the auth callout p99 acceptable on cold reads.
	MqttPasswordBcryptCost = 10
)
