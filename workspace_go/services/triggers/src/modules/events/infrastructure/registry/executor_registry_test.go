package registry

import (
	"testing"
)

/**
 * NewExecutorRegistry Tests
 */

func TestNewExecutorRegistry_RegistersAllExecutors(t *testing.T) {
	registry := NewExecutorRegistry()

	expectedTypes := []string{
		"http",
		"mqtt",
		"rabbitmq",
		"nats",
		"websocket",
		"email",
		"teams",
		"slack",
	}

	for _, triggerType := range expectedTypes {
		executor, exists := registry.GetExecutor(triggerType)
		if !exists {
			t.Errorf("NewExecutorRegistry() should register executor for type %q", triggerType)
			continue
		}
		if executor == nil {
			t.Errorf("NewExecutorRegistry() executor for type %q should not be nil", triggerType)
		}
	}
}

func TestNewExecutorRegistry_ExactCount(t *testing.T) {
	registry := NewExecutorRegistry()

	// Cast to access internal state (same package)
	r, ok := registry.(*executorRegistry)
	if !ok {
		t.Fatal("NewExecutorRegistry() should return *executorRegistry")
	}

	expectedCount := 8
	actualCount := len(r.executors)

	if actualCount != expectedCount {
		t.Errorf("NewExecutorRegistry() registered %d executors, want %d", actualCount, expectedCount)
	}
}

/**
 * GetExecutor Tests
 */

func TestGetExecutor_ExistingType(t *testing.T) {
	registry := NewExecutorRegistry()

	executor, exists := registry.GetExecutor("http")

	if !exists {
		t.Fatal("GetExecutor(\"http\") should return true")
	}
	if executor == nil {
		t.Fatal("GetExecutor(\"http\") should return non-nil executor")
	}
}

func TestGetExecutor_UnknownType(t *testing.T) {
	registry := NewExecutorRegistry()

	executor, exists := registry.GetExecutor("unknown_type")

	if exists {
		t.Error("GetExecutor(\"unknown_type\") should return false")
	}
	if executor != nil {
		t.Error("GetExecutor(\"unknown_type\") should return nil executor")
	}
}

func TestGetExecutor_EmptyString(t *testing.T) {
	registry := NewExecutorRegistry()

	executor, exists := registry.GetExecutor("")

	if exists {
		t.Error("GetExecutor(\"\") should return false")
	}
	if executor != nil {
		t.Error("GetExecutor(\"\") should return nil executor")
	}
}

/**
 * GetType Consistency Tests
 */

func TestGetExecutor_TypeConsistency(t *testing.T) {
	registry := NewExecutorRegistry()

	expectedTypes := []string{
		"http",
		"mqtt",
		"rabbitmq",
		"nats",
		"websocket",
		"email",
		"teams",
		"slack",
	}

	for _, triggerType := range expectedTypes {
		executor, exists := registry.GetExecutor(triggerType)
		if !exists {
			t.Errorf("GetExecutor(%q) should return true", triggerType)
			continue
		}

		// Verify that executor.GetType() matches the registry key
		if executor.GetType() != triggerType {
			t.Errorf("GetExecutor(%q).GetType() = %q, want %q", triggerType, executor.GetType(), triggerType)
		}
	}
}

func TestGetExecutor_CaseSensitive(t *testing.T) {
	registry := NewExecutorRegistry()

	// Registry keys are lowercase — uppercase should not match
	casesToTest := []string{"HTTP", "Http", "MQTT", "Mqtt", "NATS", "Nats"}

	for _, triggerType := range casesToTest {
		_, exists := registry.GetExecutor(triggerType)
		if exists {
			t.Errorf("GetExecutor(%q) should return false (case-sensitive lookup)", triggerType)
		}
	}
}
