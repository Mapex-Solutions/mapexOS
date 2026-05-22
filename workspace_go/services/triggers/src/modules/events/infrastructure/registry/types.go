package registry

import (
	"triggers/src/modules/events/application/ports"
)

// executorRegistry is the infrastructure implementation of ports.ExecutorRegistry.
//
// Following Hexagonal Architecture, this adapter:
// - Lives in the infrastructure layer
// - Implements the application port interface (ports.ExecutorRegistry)
// - Knows about all concrete executor implementations
// - Provides factory method for dependency injection
type executorRegistry struct {
	executors map[string]ports.TriggerExecutor
}

// Compile-time check to ensure executorRegistry implements ports.ExecutorRegistry interface
var _ ports.ExecutorRegistry = (*executorRegistry)(nil)
