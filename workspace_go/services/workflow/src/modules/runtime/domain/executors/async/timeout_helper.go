package async

import (
	"fmt"
	"time"

	"workflow/src/modules/runtime/domain/entities"
)

/**
 * Timeout helper for async executors.
 * Calculates expiresAt from node-level TimeoutConfig or falls back to a default duration.
 * Used by all async executors to include expiresAt in NodeState for NATS Schedule pickup.
 */

// CalculateExpiresAt resolves the absolute expiration time for an async node.
// Priority: node timeout config → executor default duration.
func CalculateExpiresAt(timeout *entities.TimeoutConfig, defaultDuration time.Duration) time.Time {
	if timeout == nil || timeout.Duration <= 0 {
		return time.Now().Add(defaultDuration)
	}
	d, err := ParseDuration(timeout.Duration, timeout.Unit)
	if err != nil {
		return time.Now().Add(defaultDuration)
	}
	return time.Now().Add(d)
}

// IsEnableOutput returns whether the timeout should route to an output handle instead of failing.
func IsEnableOutput(timeout *entities.TimeoutConfig) bool {
	return timeout != nil && timeout.EnableOutput
}

// ParseDuration converts a duration value + unit string to time.Duration.
func ParseDuration(value int, unit string) (time.Duration, error) {
	switch unit {
	case "seconds", "s":
		return time.Duration(value) * time.Second, nil
	case "minutes", "m":
		return time.Duration(value) * time.Minute, nil
	case "hours", "h":
		return time.Duration(value) * time.Hour, nil
	case "days", "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "months":
		return time.Duration(value) * 30 * 24 * time.Hour, nil
	case "years":
		return time.Duration(value) * 365 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown duration unit: %s", unit)
	}
}
