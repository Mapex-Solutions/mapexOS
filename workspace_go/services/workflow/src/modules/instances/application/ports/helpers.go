package ports

import (
	defPorts "workflow/src/modules/definitions/application/ports"
)

// InitializeState builds the initial state map from workflow variable definitions.
// Each variable with a non-nil DefaultValue is included in the resulting map.
//
// Defined here (application/ports) rather than in domain/entities to keep the
// instance domain free of cross-module imports: the function accepts types
// exposed via the definitions module's ports.
func InitializeState(stateDefs []defPorts.WorkflowVariable) map[string]interface{} {
	state := make(map[string]interface{})
	for _, def := range stateDefs {
		if def.DefaultValue != nil {
			state[def.Field] = def.DefaultValue
		}
	}
	return state
}

// InitializeExternalInputs merges trigger-provided external inputs with definition defaults.
// Trigger values take precedence over definition defaults.
func InitializeExternalInputs(defs []defPorts.ExternalInput, provided map[string]interface{}) map[string]interface{} {
	inputs := make(map[string]interface{})
	for _, def := range defs {
		if def.DefaultValue != nil {
			inputs[def.Field] = def.DefaultValue
		}
	}
	for k, v := range provided {
		inputs[k] = v
	}
	return inputs
}
