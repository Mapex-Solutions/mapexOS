package ports

import (
	"workflow/src/modules/definitions/domain/entities"

	definitionsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/definitions"
)

// Port-level type aliases — expose domain entities through the port boundary.
// Other modules import these types from ports, NEVER from domain/entities directly.

// DefinitionInvalidatePayload re-exports the cross-service contract for the
// FANOUT cache-invalidation payload. The application layer imports this alias
// instead of reaching into interfaces/message (which would cross the Hexagonal layering).
type DefinitionInvalidatePayload = definitionsContract.DefinitionInvalidatePayload

type WorkflowDefinition = entities.WorkflowDefinition
type WorkflowNode = entities.WorkflowNode
type WorkflowVariable = entities.WorkflowVariable
type ExternalInput = entities.ExternalInput
type FieldValue = entities.FieldValue
type FieldValueType = entities.FieldValueType
type ConditionGroup = entities.ConditionGroup
type ConditionGroupItem = entities.ConditionGroupItem
type ConditionItem = entities.ConditionItem
type GroupLogicOperator = entities.GroupLogicOperator
type SwitchCase = entities.SwitchCase
type DefinitionStatus = entities.DefinitionStatus
type ErrorHandlerConfig = entities.ErrorHandlerConfig

// Re-export constants
const (
	FieldValueState      = entities.FieldValueState
	FieldValueLiteral    = entities.FieldValueLiteral
	FieldValueEvent      = entities.FieldValueEvent
	FieldValueVariable   = entities.FieldValueVariable
	FieldValueInput      = entities.FieldValueInput
	FieldValueNodeOutput = entities.FieldValueNodeOutput
	LogicAND             = entities.LogicAND
	LogicOR              = entities.LogicOR
	LogicNAND            = entities.LogicNAND
	LogicNOR             = entities.LogicNOR
	StatusValid          = entities.StatusValid
)
