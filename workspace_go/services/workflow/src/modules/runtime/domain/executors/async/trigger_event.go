package async

import (
	"context"
	"fmt"
	"time"

	enginePorts "workflow/src/modules/engine/application/ports"
	"workflow/src/modules/runtime/domain/entities"
)

/*
 * TRIGGER EVENT EXECUTOR
 * Publishes an event to the Trigger Service and waits for a callback.
 * Resolves payload field values via ValueResolver before publishing.
 * Returns NodeState with waitType "callback" — the RuntimeService publishes
 * the actual event to NATS and sets up the callback subject.
 */

// TriggerEventExecutor publishes an event to the Trigger Service and waits for a callback.
// Resolves payload field values via ValueResolver before publishing.
type TriggerEventExecutor struct {
	resolver enginePorts.ValueResolverPort
}

// NewTriggerEventExecutor creates a new TriggerEventExecutor with the given value resolver.
func NewTriggerEventExecutor(resolver enginePorts.ValueResolverPort) entities.NodeExecutor {
	return &TriggerEventExecutor{resolver: resolver}
}

// NodeType returns "core/trigger_event".
func (e *TriggerEventExecutor) NodeType() string {
	return "core/trigger_event"
}

// Execute resolves the payload mapping fields and returns a NodeState with waitType "callback"
// for the RuntimeService to publish the event to NATS.
func (e *TriggerEventExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.TriggerEventNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("trigger_event: missing or invalid config")
	}

	if cfg.EventType == "" {
		return nil, fmt.Errorf("trigger_event: eventType is required")
	}

	// Resolve payload field values
	payload := make(map[string]interface{})
	for _, field := range cfg.PayloadMapping {
		val := field.Value
		resolved, err := e.resolver.Resolve(&val, execCtx.EventPayload, execCtx.State, execCtx.NodeOutputs, execCtx.ExternalInputs)
		if err != nil {
			return nil, fmt.Errorf("trigger_event: failed to resolve field %s: %w", field.Key, err)
		}
		payload[field.Key] = resolved
	}

	expiresAt := CalculateExpiresAt(execCtx.Timeout, 30*time.Second)

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
		NodeState: map[string]interface{}{
			"waitType":     "callback",
			"eventType":    cfg.EventType,
			"payload":      payload,
			"expiresAt":    expiresAt,
			"enableOutput": IsEnableOutput(execCtx.Timeout),
		},
	}, nil
}
