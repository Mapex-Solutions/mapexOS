package ports

import (
	defPorts "workflow/src/modules/definitions/application/ports"
)

/*
 * ENGINE PORTS
 * Contracts for the engine module's condition evaluation and value resolution.
 */

// ConditionEvaluatorPort defines the contract for evaluating condition groups.
type ConditionEvaluatorPort interface {
	// EvaluateGroup recursively evaluates a ConditionGroup with short-circuit logic.
	EvaluateGroup(
		group *defPorts.ConditionGroup,
		timezone string,
		eventPayload map[string]interface{},
		state map[string]interface{},
		nodeOutputs map[string]interface{},
		externalInputs map[string]interface{},
	) (bool, error)
}

// ValueResolverPort defines the contract for resolving FieldValue to actual values.
type ValueResolverPort interface {
	// Resolve extracts the actual value from a FieldValue based on its source type.
	Resolve(
		field *defPorts.FieldValue,
		eventPayload map[string]interface{},
		state map[string]interface{},
		nodeOutputs map[string]interface{},
		externalInputs map[string]interface{},
	) (interface{}, error)

	// BuildDescription creates a human-readable description of the field source.
	BuildDescription(field *defPorts.FieldValue) string
}
