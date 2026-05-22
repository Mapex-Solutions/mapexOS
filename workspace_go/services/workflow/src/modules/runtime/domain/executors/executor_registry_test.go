package executors

import (
	"context"
	"errors"
	"testing"

	"workflow/src/modules/runtime/domain/entities"
)

type stubExecutor struct {
	nodeType string
}

func (s *stubExecutor) NodeType() string { return s.nodeType }
func (s *stubExecutor) Execute(_ context.Context, _ *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	return &entities.NodeExecutionResult{OutputHandles: []string{"out"}}, nil
}

func TestExecutorRegistry(t *testing.T) {
	t.Run("register and get", func(t *testing.T) {
		registry := NewExecutorRegistry()
		registry.Register(&stubExecutor{nodeType: "core/test"})

		executor, err := registry.Get("core/test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if executor.NodeType() != "core/test" {
			t.Fatalf("expected core/test, got %s", executor.NodeType())
		}
	})

	t.Run("get unknown returns ErrExecutorNotFound", func(t *testing.T) {
		registry := NewExecutorRegistry()

		_, err := registry.Get("unknown/type")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, entities.ErrExecutorNotFound) {
			t.Fatalf("expected ErrExecutorNotFound, got %v", err)
		}
	})

	t.Run("register overwrites", func(t *testing.T) {
		registry := NewExecutorRegistry()
		registry.Register(&stubExecutor{nodeType: "core/test"})
		registry.Register(&stubExecutor{nodeType: "core/test"})

		executor, err := registry.Get("core/test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if executor.NodeType() != "core/test" {
			t.Fatalf("expected core/test, got %s", executor.NodeType())
		}
	})
}
