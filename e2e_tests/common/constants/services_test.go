package constants

import (
	"testing"
	"time"
)

// TestGetTimeoutMultiplier_ParseAndClamp proves the env is parsed and clamped so it
// can only extend a wait, never shorten it.
func TestGetTimeoutMultiplier_ParseAndClamp(t *testing.T) {
	t.Setenv("SAGA_TIMEOUT_MULTIPLIER", "3")
	if got := getTimeoutMultiplier(); got != 3.0 {
		t.Errorf("multiplier %q = %v, want 3.0", "3", got)
	}

	t.Setenv("SAGA_TIMEOUT_MULTIPLIER", "0.5")
	if got := getTimeoutMultiplier(); got != 1.0 {
		t.Errorf("multiplier %q = %v, want clamped to 1.0", "0.5", got)
	}

	t.Setenv("SAGA_TIMEOUT_MULTIPLIER", "not-a-number")
	if got := getTimeoutMultiplier(); got != 1.0 {
		t.Errorf("multiplier %q = %v, want 1.0 fallback", "not-a-number", got)
	}
}

// TestScaleTimeout proves ScaleTimeout applies the default 1.0 multiplier (env unset)
// as a pass-through.
func TestScaleTimeout(t *testing.T) {
	if got := ScaleTimeout(10 * time.Second); got != 10*time.Second {
		t.Errorf("ScaleTimeout(10s) = %v, want 10s with the default multiplier", got)
	}
}
