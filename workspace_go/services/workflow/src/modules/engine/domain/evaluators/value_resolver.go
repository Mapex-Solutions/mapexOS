package evaluators

import (
	"fmt"
	"strings"

	defPorts "workflow/src/modules/definitions/application/ports"
)

/*
 * VALUE RESOLVER
 * Resolves FieldValue to actual values from 6 sources:
 * literal, event, state, variable, input, node_output.
 */

type ValueResolver struct{}

/*
 * NewValueResolver creates a new ValueResolver instance.
 */
func NewValueResolver() *ValueResolver {
	return &ValueResolver{}
}

/*
 * Resolve extracts the actual value from a FieldValue based on its source type.
 *
 * Sources:
 * - literal: Returns the Value string. If the value contains '{{...}}' placeholders,
 *   they are interpolated against (event, state, input, output) using renderTemplate
 *   (see value_resolver_template.go).
 * - event: Extracts from eventPayload using Value as dot-path
 * - state: Extracts from state using Value as dot-path
 * - variable: Extracts from state using Value as dot-path (variables merged into state)
 * - node_output: Extracts from nodeOutputs[NodeID] using Value as dot-path
 */
func (r *ValueResolver) Resolve(
	field *defPorts.FieldValue,
	eventPayload map[string]interface{},
	state map[string]interface{},
	nodeOutputs map[string]interface{},
	externalInputs map[string]interface{},
) (interface{}, error) {
	if field == nil {
		return nil, ErrInvalidFieldValue
	}

	switch field.Type {
	case defPorts.FieldValueLiteral:
		v := field.Value
		if !strings.Contains(v, "{{") {
			return v, nil
		}
		return renderTemplate(v, eventPayload, state, externalInputs, nodeOutputs), nil

	case defPorts.FieldValueEvent:
		return r.resolveFromMap(field.Value, eventPayload)

	case defPorts.FieldValueState:
		return r.resolveFromMap(field.Value, state)

	case defPorts.FieldValueVariable:
		return r.resolveFromMap(field.Value, state)

	case defPorts.FieldValueInput:
		// Frontend stores value as "input.{field}" — strip prefix
		path := field.Value
		if strings.HasPrefix(path, "input.") {
			path = strings.TrimPrefix(path, "input.")
		}
		return r.resolveFromMap(path, externalInputs)

	case defPorts.FieldValueNodeOutput:
		return r.resolveNodeOutput(field, nodeOutputs)

	default:
		return nil, fmt.Errorf("%w: unknown source type '%s'", ErrInvalidSource, field.Type)
	}
}

/*
 * resolveNodeOutput extracts a value from a specific node's output map.
 * Requires FieldValue.NodeID to identify the target node.
 */
func (r *ValueResolver) resolveNodeOutput(
	field *defPorts.FieldValue,
	nodeOutputs map[string]interface{},
) (interface{}, error) {
	if field.NodeID == "" {
		return nil, ErrInvalidNodeID
	}

	if nodeOutputs == nil {
		return nil, fmt.Errorf("%w: nodeOutputs is nil", ErrInvalidSource)
	}

	nodeData, exists := nodeOutputs[field.NodeID]
	if !exists {
		return nil, fmt.Errorf("%w: node '%s' not found in outputs", ErrFieldNotFound, field.NodeID)
	}

	// If Value is empty, return the entire node output
	if field.Value == "" {
		return nodeData, nil
	}

	nodeMap, ok := nodeData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: node '%s' output is not a map", ErrFieldNotFound, field.NodeID)
	}

	return r.resolveFromMap(field.Value, nodeMap)
}

/*
 * resolveFromMap extracts a value from a map using a dot-notation path.
 * Supports nested paths like "device.sensor.temperature".
 */
func (r *ValueResolver) resolveFromMap(
	path string,
	source map[string]interface{},
) (interface{}, error) {
	if source == nil {
		return nil, ErrInvalidSource
	}

	if path == "" {
		return nil, ErrInvalidPath
	}

	parts := strings.Split(path, ".")

	var current interface{} = source
	for _, part := range parts {
		if part == "" {
			continue
		}

		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: cannot traverse '%s' - not a map", ErrFieldNotFound, part)
		}

		value, exists := currentMap[part]
		if !exists {
			return nil, fmt.Errorf("%w: '%s' in path '%s'", ErrFieldNotFound, part, path)
		}

		current = value
	}

	return current, nil
}

/*
 * BuildDescription creates a human-readable description of the field source.
 */
func (r *ValueResolver) BuildDescription(field *defPorts.FieldValue) string {
	if field == nil {
		return "nil"
	}

	switch field.Type {
	case defPorts.FieldValueLiteral:
		return field.Value
	case defPorts.FieldValueEvent:
		return fmt.Sprintf("event.%s", field.Value)
	case defPorts.FieldValueState:
		return fmt.Sprintf("state.%s", field.Value)
	case defPorts.FieldValueVariable:
		return fmt.Sprintf("variable.%s", field.Value)
	case defPorts.FieldValueInput:
		return fmt.Sprintf("input.%s", field.Value)
	case defPorts.FieldValueNodeOutput:
		if field.NodeID != "" {
			return fmt.Sprintf("node[%s].%s", field.NodeID, field.Value)
		}
		return fmt.Sprintf("nodeOutput.%s", field.Value)
	default:
		return fmt.Sprintf("%s.%s", field.Type, field.Value)
	}
}
