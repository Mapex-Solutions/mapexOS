package redis

import (
	"triggers/src/modules/triggers/application/ports"
)

// Compile-time check to ensure KeyBuilder implements TriggerCacheKeyBuilderPort.
var _ ports.TriggerCacheKeyBuilderPort = (*KeyBuilder)(nil)

// NewKeyBuilder returns a new Redis key builder wired as the
// TriggerCacheKeyBuilderPort port implementation.
func NewKeyBuilder() ports.TriggerCacheKeyBuilderPort {
	return &KeyBuilder{}
}

// TriggerKey returns the Redis key for a single trigger cache entry.
// Key format: TRIGGER:{triggerId}
func (k *KeyBuilder) TriggerKey(triggerId string) string {
	return BuildTriggerCacheKey(triggerId)
}

// CounterKey returns the Redis key for the per-org trigger counter cache.
// Key format: counter:triggers:{orgId}
func (k *KeyBuilder) CounterKey(orgId string) string {
	return BuildCounterCacheKey(orgId)
}
