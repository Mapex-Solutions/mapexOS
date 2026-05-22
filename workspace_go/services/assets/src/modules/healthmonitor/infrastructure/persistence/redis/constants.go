package redis

// Redis key prefixes for health monitor state.
// Appended to the key prefix from the client config.
const (
	// RedisKeyLastSeen is the prefix for the per-org last-seen ZSET.
	// Full key: {clientPrefix}hm:{orgId}
	RedisKeyLastSeen = "hm:"

	// RedisKeyMiss is the prefix for the per-org miss-counter HASH.
	// Full key: {clientPrefix}hm:miss:{orgId}
	RedisKeyMiss = "hm:miss:"

	// RedisKeyAlerted is the prefix for the per-org alerted (offline) SET.
	// Full key: {clientPrefix}hm:alerted:{orgId}
	RedisKeyAlerted = "hm:alerted:"

	// RedisKeyOrgs is the global SET of orgs with active monitoring.
	RedisKeyOrgs = "hm:orgs"

	// RedisKeyKnown is the prefix for the per-org known-online SET.
	// Tracks assets confirmed online at least once.
	// Full key: {clientPrefix}hm:known:{orgId}
	RedisKeyKnown = "hm:known:"

	// RedisKeyLastConnect is the prefix for the per-org last-connect HASH
	// (assetUUID → unixSeconds). Distinct from RedisKeyLastSeen because
	// semantics differ: lastSeenAt moves on every heartbeat, lastConnectAt
	// only on a $SYS.ACCOUNT.*.CONNECT advisory. The anti-race invariant
	// in the presence consumer (disconnect.timestamp > lastConnectAt)
	// requires a stable CONNECT timestamp that heartbeats do NOT push
	// forward.
	// Full key: {clientPrefix}hm:lc:{orgId}
	RedisKeyLastConnect = "hm:lc:"
)
