package executors

import (
	"fmt"
	"strings"

	"workflow/src/modules/runtime/domain/entities"
)

// corePrefix identifies core node types (e.g., "core/condition", "core/delay").
// Nodes with this prefix are looked up in the executor map.
// All other nodes are routed to the plugin executor.
const corePrefix = "core/"

// ExecutorRegistry maps node types to their executor implementations.
// Core nodes (core/*) are looked up by exact type in the map.
// Plugin nodes (anything else) are routed directly to the plugin executor without map lookup.
type ExecutorRegistry struct {
	executors      map[string]entities.NodeExecutor
	pluginExecutor entities.NodeExecutor
}

// NewExecutorRegistry creates a new empty ExecutorRegistry.
func NewExecutorRegistry() *ExecutorRegistry {
	return &ExecutorRegistry{
		executors: make(map[string]entities.NodeExecutor),
	}
}

// Register adds a core node executor to the registry, keyed by its NodeType().
func (r *ExecutorRegistry) Register(executor entities.NodeExecutor) {
	r.executors[executor.NodeType()] = executor
}

// SetPluginExecutor sets the executor used for all non-core node types.
func (r *ExecutorRegistry) SetPluginExecutor(executor entities.NodeExecutor) {
	r.pluginExecutor = executor
}

// Get retrieves the executor for the given node type.
// Core nodes (core/*) are looked up in the map.
// Plugin nodes are routed directly to the plugin executor — no map lookup, zero CPU wasted.
func (r *ExecutorRegistry) Get(nodeType string) (entities.NodeExecutor, error) {
	if strings.HasPrefix(nodeType, corePrefix) {
		executor, ok := r.executors[nodeType]
		if !ok {
			return nil, fmt.Errorf("%w: %s", entities.ErrExecutorNotFound, nodeType)
		}
		return executor, nil
	}

	if r.pluginExecutor != nil {
		return r.pluginExecutor, nil
	}

	return nil, fmt.Errorf("%w: %s (no plugin executor registered)", entities.ErrExecutorNotFound, nodeType)
}
