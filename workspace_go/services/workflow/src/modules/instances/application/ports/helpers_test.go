package ports

import (
	"testing"

	defPorts "workflow/src/modules/definitions/application/ports"
)

func TestInitializeState_WithDefaults(t *testing.T) {
	defs := []defPorts.WorkflowVariable{
		{Field: "counter", DefaultValue: "0"},
		{Field: "status", DefaultValue: "idle"},
		{Field: "noDefault", DefaultValue: nil},
	}

	state := InitializeState(defs)

	if state["counter"] != "0" {
		t.Fatalf("expected counter='0', got %v", state["counter"])
	}
	if state["status"] != "idle" {
		t.Fatalf("expected status='idle', got %v", state["status"])
	}
	if _, exists := state["noDefault"]; exists {
		t.Fatal("expected noDefault to not be in state")
	}
}

func TestInitializeState_Empty(t *testing.T) {
	state := InitializeState(nil)
	if len(state) != 0 {
		t.Fatalf("expected empty state, got %d entries", len(state))
	}
}

func TestInitializeExternalInputs_MergesWithDefaults(t *testing.T) {
	defs := []defPorts.ExternalInput{
		{Field: "sensor1", DefaultValue: "default-uuid"},
		{Field: "userName", DefaultValue: "Anonymous"},
	}

	provided := map[string]interface{}{
		"sensor1": "actual-uuid-123",
	}

	inputs := InitializeExternalInputs(defs, provided)

	if inputs["sensor1"] != "actual-uuid-123" {
		t.Fatalf("expected sensor1='actual-uuid-123', got %v", inputs["sensor1"])
	}
	if inputs["userName"] != "Anonymous" {
		t.Fatalf("expected userName='Anonymous', got %v", inputs["userName"])
	}
}

func TestInitializeExternalInputs_ProvidedTakesPrecedence(t *testing.T) {
	defs := []defPorts.ExternalInput{
		{Field: "key1", DefaultValue: "default"},
	}

	provided := map[string]interface{}{
		"key1": "override",
		"key2": "extra",
	}

	inputs := InitializeExternalInputs(defs, provided)

	if inputs["key1"] != "override" {
		t.Fatalf("expected key1='override', got %v", inputs["key1"])
	}
	if inputs["key2"] != "extra" {
		t.Fatalf("expected key2='extra', got %v", inputs["key2"])
	}
}

func TestInitializeExternalInputs_NilProvided(t *testing.T) {
	defs := []defPorts.ExternalInput{
		{Field: "sensor1", DefaultValue: "default-uuid"},
	}

	inputs := InitializeExternalInputs(defs, nil)

	if inputs["sensor1"] != "default-uuid" {
		t.Fatalf("expected sensor1='default-uuid', got %v", inputs["sensor1"])
	}
}
