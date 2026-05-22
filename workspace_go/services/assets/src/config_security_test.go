package main

import (
	"testing"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"

	configApp "assets/src/shared/configuration/application"
)

// sensitiveKeysForService enumerates every configuration key that MUST
// carry Sensitive: true in DefaultConfiguration. Adding a new credential
// without updating this list fails the regression test below.
var sensitiveKeysForService = []string{
	"nats_password",
	"auth_secret",
	"internal_api_key",
	"mapex_vault_api_key",
	"minio_access_key",
	"minio_secret_key",
}

// TestSecurity_SensitiveKeysAreMarked catches regressions where a new
// credential key joins DefaultConfiguration without Sensitive: true —
// its dev default would silently leak when GO_ENV is not dev.
func TestSecurity_SensitiveKeysAreMarked(t *testing.T) {
	marked := map[string]bool{}
	for _, d := range configApp.DefaultConfiguration {
		if d.Sensitive {
			marked[d.Key] = true
		}
	}
	for _, key := range sensitiveKeysForService {
		if !marked[key] {
			t.Errorf("config key %q must have Sensitive: true — dev default would leak when GO_ENV is not dev", key)
		}
	}
}

// TestSecurity_NonDevWithDefaultsHasViolations simulates startup with no
// env-var overrides (every key resolves to its hardcoded Default). The
// validator must surface every entry in the curated list.
func TestSecurity_NonDevWithDefaultsHasViolations(t *testing.T) {
	resolved := map[string]interface{}{}
	for _, d := range configApp.DefaultConfiguration {
		resolved[d.Key] = d.Default
	}

	violations := config.FindSensitiveDefaultsInUse(configApp.DefaultConfiguration, resolved)
	if len(violations) != len(sensitiveKeysForService) {
		t.Errorf("expected %d violations matching the curated list, got %d: %v",
			len(sensitiveKeysForService), len(violations), violations)
	}
}

// TestSecurity_OverridesPassValidation verifies that when every sensitive
// env var is set to a value different from its Default, no violations
// remain — InitConfig would not fatal in this state.
func TestSecurity_OverridesPassValidation(t *testing.T) {
	resolved := map[string]interface{}{}
	for _, d := range configApp.DefaultConfiguration {
		if d.Sensitive {
			resolved[d.Key] = "PROD_OVERRIDE_" + d.Key
		} else {
			resolved[d.Key] = d.Default
		}
	}

	if v := config.FindSensitiveDefaultsInUse(configApp.DefaultConfiguration, resolved); len(v) != 0 {
		t.Errorf("prod with overrides must produce no violations, got %v", v)
	}
}
