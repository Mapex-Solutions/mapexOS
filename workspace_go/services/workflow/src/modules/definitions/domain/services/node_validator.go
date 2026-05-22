package services

import (
	"workflow/src/modules/definitions/domain/constants"
	"workflow/src/modules/definitions/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/*
 * NODE VALIDATOR
 * Validates each node's Config map based on its Type.
 * Called on CREATE/UPDATE only — rejects invalid definitions with 400.
 *
 * NodeValidationError lives in types.go (same package).
 */

// validDelayUnits are accepted values for delay node unit field.
var validDelayUnits = map[string]bool{
	"s": true, "seconds": true,
	"m": true, "minutes": true,
	"h": true, "hours": true,
	"d": true, "days": true,
}

// validSetStateOps are accepted values for set_state operation field.
var validSetStateOps = map[string]bool{
	"set":       true,
	"increment": true,
	"decrement": true,
	"append":    true,
	"remove":    true,
}

// validGotoRoles are accepted values for goto role field.
var validGotoRoles = map[string]bool{
	"sender":   true,
	"receiver": true,
}

// ValidateNodes checks all node configs and returns errors per node.
// Returns nil if all nodes are valid.
func ValidateNodes(nodes []entities.WorkflowNode) []NodeValidationError {
	var result []NodeValidationError

	for _, n := range nodes {
		errs := validateNode(n)
		if len(errs) > 0 {
			result = append(result, NodeValidationError{
				NodeID:   n.ID,
				NodeType: n.Type,
				Errors:   errs,
			})
		}
	}

	return result
}

// validateNode validates a single node's config based on its type.
func validateNode(n entities.WorkflowNode) []string {
	cfg := n.Config

	switch n.Type {
	// Visual-only nodes — skip
	case constants.NodeTypeTextNote, constants.NodeTypeGroupFrame:
		return nil

	// No validation needed
	case constants.NodeTypeStart, constants.NodeTypeEnd, constants.NodeTypeLog:
		return nil

	case constants.NodeTypeCondition:
		return validateCondition(cfg)
	case constants.NodeTypeCode:
		return validateCode(cfg)
	case constants.NodeTypeSetState:
		return validateSetState(cfg)
	case constants.NodeTypeSwitch:
		return validateSwitch(cfg)
	case constants.NodeTypeSubworkflow:
		return validateSubworkflow(cfg)
	case constants.NodeTypeDelay:
		return validateDelay(cfg)
	case constants.NodeTypeWaitSignal:
		return validateWaitSignal(cfg)
	case constants.NodeTypeLoop:
		return validateLoop(cfg)
	case constants.NodeTypeFanout:
		return validateFanout(cfg)
	case constants.NodeTypeMerge:
		return validateMerge(cfg)
	case constants.NodeTypeSequence:
		return validateSequence(cfg)
	case constants.NodeTypeTriggerEvent:
		return validateTriggerEvent(cfg)
	case constants.NodeTypeWaitFor:
		return validateWaitFor(cfg)
	case constants.NodeTypeGoto:
		return validateGoto(cfg)

	default:
		// Unknown node type — no validation
		return nil
	}
}

func validateCondition(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetSlice(cfg, "items") == nil {
		errs = append(errs, "items is required")
	}
	return errs
}

func validateCode(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetString(cfg, "script") == "" {
		errs = append(errs, "script is required")
	}
	return errs
}

func validateSetState(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetString(cfg, "targetField") == "" {
		errs = append(errs, "targetField is required")
	}
	op := model.MapGetString(cfg, "operation")
	if op == "" {
		errs = append(errs, "operation is required")
	} else if !validSetStateOps[op] {
		errs = append(errs, "operation must be one of: set, increment, decrement, append, remove")
	}
	return errs
}

func validateSwitch(cfg map[string]interface{}) []string {
	var errs []string
	cases := model.MapGetSlice(cfg, "cases")
	if len(cases) == 0 {
		errs = append(errs, "cases is required and must not be empty")
	}
	return errs
}

func validateSubworkflow(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetString(cfg, "workflowId") == "" {
		errs = append(errs, "workflowId is required")
	}
	return errs
}

func validateDelay(cfg map[string]interface{}) []string {
	var errs []string
	dur := model.MapGetInt(cfg, "duration")
	if dur <= 0 {
		errs = append(errs, "duration must be greater than 0")
	}
	unit := model.MapGetString(cfg, "unit")
	if unit == "" {
		errs = append(errs, "unit is required")
	} else if !validDelayUnits[unit] {
		errs = append(errs, "unit must be one of: s, seconds, m, minutes, h, hours, d, days")
	}
	return errs
}

func validateWaitSignal(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetString(cfg, "signalName") == "" {
		errs = append(errs, "signalName is required")
	}
	return errs
}

func validateLoop(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetMap(cfg, "source") == nil {
		errs = append(errs, "source is required")
	}
	return errs
}

func validateFanout(cfg map[string]interface{}) []string {
	var errs []string
	b := model.MapGetInt(cfg, "branches")
	if b <= 0 {
		errs = append(errs, "branches must be greater than 0")
	} else if b > 20 {
		errs = append(errs, "branches must not exceed 20")
	}
	mode := model.MapGetString(cfg, "mode")
	if mode != "" && mode != "waitAll" && mode != "firstCompleted" {
		errs = append(errs, "mode must be 'waitAll' or 'firstCompleted'")
	}
	return errs
}

func validateMerge(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetInt(cfg, "branches") <= 0 {
		errs = append(errs, "branches must be greater than 0")
	}
	return errs
}

func validateSequence(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetInt(cfg, "steps") <= 0 {
		errs = append(errs, "steps must be greater than 0")
	}
	return errs
}

func validateTriggerEvent(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetString(cfg, "eventType") == "" {
		errs = append(errs, "eventType is required")
	}
	return errs
}

func validateWaitFor(cfg map[string]interface{}) []string {
	var errs []string
	if model.MapGetString(cfg, "field") == "" {
		errs = append(errs, "field is required")
	}
	if model.MapGetString(cfg, "operator") == "" {
		errs = append(errs, "operator is required")
	}
	return errs
}

func validateGoto(cfg map[string]interface{}) []string {
	var errs []string
	role := model.MapGetString(cfg, "role")
	if role == "" {
		errs = append(errs, "role is required")
	} else if !validGotoRoles[role] {
		errs = append(errs, "role must be one of: sender, receiver")
	}
	if model.MapGetString(cfg, "pairLabel") == "" {
		errs = append(errs, "pairLabel is required")
	}
	return errs
}
