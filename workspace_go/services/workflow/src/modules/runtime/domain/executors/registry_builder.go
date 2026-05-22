package executors

import (
	runtimePorts "workflow/src/modules/runtime/application/ports"
	enginePorts "workflow/src/modules/engine/application/ports"
	pluginPorts "workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/runtime/domain/executors/async"
	"workflow/src/modules/runtime/domain/executors/control"
	"workflow/src/modules/runtime/domain/executors/inline"
)

// BuildRegistry constructs the ExecutorRegistry with all 17 core node executors
// and the generic plugin executor for marketplace nodes.
// Called once at service startup; the registry is shared across all executions.
func BuildRegistry(
	evaluator enginePorts.ConditionEvaluatorPort,
	resolver enginePorts.ValueResolverPort,
	vaultService runtimePorts.VaultPort,
	pluginRepo pluginPorts.PluginManifestRepository,
) *ExecutorRegistry {
	registry := NewExecutorRegistry()

	// Inline executors (7)
	registry.Register(inline.NewStartExecutor())
	registry.Register(inline.NewEndExecutor(resolver))
	registry.Register(inline.NewConditionExecutor(evaluator))
	registry.Register(inline.NewSwitchExecutor(evaluator))
	registry.Register(inline.NewSetStateExecutor(resolver))
	registry.Register(inline.NewLogExecutor())
	registry.Register(inline.NewGotoExecutor())

	// Async executors (5)
	registry.Register(async.NewDelayExecutor())
	registry.Register(async.NewWaitSignalExecutor())
	registry.Register(async.NewCodeExecutor())
	registry.Register(async.NewSubworkflowExecutor(resolver))
	registry.Register(async.NewTriggerEventExecutor(resolver))

	// Control executors (5)
	registry.Register(control.NewFanoutExecutor())
	registry.Register(control.NewMergeExecutor())
	registry.Register(control.NewSequenceExecutor())
	registry.Register(control.NewLoopExecutor(resolver))
	registry.Register(control.NewWaitForExecutor(evaluator))

	// Plugin executor (generic — handles all non-core node types)
	registry.SetPluginExecutor(async.NewPluginExecutor(vaultService, pluginRepo, resolver))

	return registry
}
