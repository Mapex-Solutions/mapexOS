package instances

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/*
 * DTOs (API Layer) for Workflow Instance configs
 */

// InstanceId represents the params DTO for /:instanceId
type InstanceId struct {
	InstanceId string `params:"instanceId" validate:"required"`
}

// InstanceCreate represents the body DTO for creating an instance config.
type InstanceCreate struct {
	DefinitionID      model.ObjectId         `json:"definitionId" validate:"required"`
	DefinitionVersion int                    `json:"definitionVersion" validate:"required,min=1"`
	DefinitionName    string                 `json:"definitionName,omitempty"`
	Name              string                 `json:"name" validate:"required,min=1"`
	Description       string                 `json:"description,omitempty"`
	OrgID             *model.ObjectId        `json:"orgId,omitempty"`
	PathKey           string                 `json:"pathKey,omitempty"`
	ExternalInputs    map[string]interface{} `json:"externalInputs,omitempty"`
	IsTemplate        bool                   `json:"isTemplate"`
	UniqueExecution   bool                   `json:"uniqueExecution"`
	WorkflowUUID      string                 `json:"workflowUUID,omitempty" validate:"omitempty"`
}

// InstanceUpdate represents the body DTO for updating an instance config.
type InstanceUpdate struct {
	Name            *string                `json:"name,omitempty" validate:"omitempty,min=1"`
	Description     *string                `json:"description,omitempty"`
	ExternalInputs  map[string]interface{} `json:"externalInputs,omitempty"`
	IsTemplate      *bool                  `json:"isTemplate,omitempty"`
	UniqueExecution *bool                  `json:"uniqueExecution,omitempty"`
	WorkflowUUID    *string                `json:"workflowUUID,omitempty"`
	Enabled         *bool                  `json:"enabled,omitempty"`
}

// InstanceQuery represents the query parameters for listing instance configs.
type InstanceQuery struct {
	query.BaseQueryDTO

	DefinitionID    *string `query:"definitionId" validate:"omitempty"`
	Name            *string `query:"name" validate:"omitempty,max=100"`
	Enabled         *bool   `query:"enabled" validate:"omitempty"`
	UniqueExecution *bool   `query:"uniqueExecution" validate:"omitempty"`
}

// InstanceResponse represents the API response for a workflow instance config.
type InstanceResponse struct {
	ID                *common.ObjectID       `json:"_id,omitempty"`
	DefinitionID      *common.ObjectID       `json:"definitionId,omitempty"`
	DefinitionVersion *int                   `json:"definitionVersion,omitempty"`
	DefinitionName    *string                `json:"definitionName,omitempty"`
	Name              *string                `json:"name,omitempty"`
	Description       *string                `json:"description,omitempty"`
	OrgID             *common.ObjectID       `json:"orgId,omitempty"`
	PathKey           *string                `json:"pathKey,omitempty"`
	ExternalInputs    map[string]interface{} `json:"externalInputs,omitempty"`
	IsSystem          *bool                  `json:"isSystem,omitempty"`
	IsTemplate        *bool                  `json:"isTemplate,omitempty"`
	UniqueExecution   *bool                  `json:"uniqueExecution,omitempty"`
	WorkflowUUID      *string                `json:"workflowUUID,omitempty"`
	Enabled           *bool                  `json:"enabled,omitempty"`
	Created           *common.NullTime       `json:"created,omitempty"`
	Updated           *common.NullTime       `json:"updated,omitempty"`
}

func (r *InstanceResponse) SetCreated(t *common.NullTime) { r.Created = t }
func (r *InstanceResponse) SetUpdated(t *common.NullTime) { r.Updated = t }

// ExecuteRequest represents the optional body for executing a workflow instance.
// eventPayload is optional test/debug data; workflowUUID is optional for idempotent retry.
type ExecuteRequest struct {
	EventPayload map[string]interface{} `json:"eventPayload,omitempty"`
	WorkflowUUID string                 `json:"workflowUUID,omitempty"`
}

// ExecuteResponse represents the response from executing a workflow instance.
type ExecuteResponse struct {
	WorkflowUUID string      `json:"workflowUUID"`
	Status       string      `json:"status"`
	ErrorInfo    interface{} `json:"errorInfo,omitempty"`
}
