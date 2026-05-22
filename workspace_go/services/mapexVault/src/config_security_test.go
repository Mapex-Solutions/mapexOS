package main

import (
	"testing"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"

	configApp "mapexVault/src/shared/configuration/application"
)

var sensitiveKeysForService = []string{
	"nats_password",
	"auth_secret",
	"internal_api_key",
	"credential_master_key",
}

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
