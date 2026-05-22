package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// WorkflowInstance is the root aggregate for an instance configuration.
// It holds the immutable inputs needed to run a workflow execution:
// which definition to use, what external inputs were provided, and ownership info.
// One instance can trigger N executions (1:N relationship).
type WorkflowInstance struct {
	ID                model.ObjectId         `bson:"_id,omitempty"`
	DefinitionID      model.ObjectId         `bson:"definitionId"`
	DefinitionVersion int                    `bson:"definitionVersion"`
	DefinitionName    string                 `bson:"definitionName"`
	Name              string                 `bson:"name"`
	Description       string                 `bson:"description"`
	OrgID             *model.ObjectId        `bson:"orgId"`
	PathKey           string                 `bson:"pathKey"`
	ExternalInputs    map[string]interface{} `bson:"externalInputs"`
	IsSystem          bool                   `bson:"isSystem"`
	IsTemplate        bool                   `bson:"isTemplate"`
	UniqueExecution   bool                   `bson:"uniqueExecution"`
	WorkflowUUID      string                 `bson:"workflowUUID,omitempty"`
	Enabled           bool                   `bson:"enabled"`
	RetentionDays     uint16                 `bson:"retentionDays"` // ClickHouse retention override (0 = use org policy)
	Created           time.Time              `bson:"created"`
	Updated           time.Time              `bson:"updated"`
}

// GetCreated returns the creation timestamp.
func (i *WorkflowInstance) GetCreated() time.Time { return i.Created }

// GetUpdated returns the last update timestamp.
func (i *WorkflowInstance) GetUpdated() time.Time { return i.Updated }
