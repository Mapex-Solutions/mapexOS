package repositories

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/*
 * REPOSITORY TYPES
 * Structs used as parameters/results in repository port contracts.
 */

// WaitingUpdate contains the fields to update when an execution suspends.
// Written by the Archiver on "waiting" events (~150 bytes per update).
type WaitingUpdate struct {
	WorkflowUUID  string    `bson:"workflowUUID"`
	Status        string    `bson:"status"`
	ActiveNodeIDs []string  `bson:"activeNodeIds"`
	Updated       time.Time `bson:"updated"`
}

// LightweightExecution is the minimal data inserted on "created" events.
// Provides listing visibility in the frontend while the workflow executes.
type LightweightExecution struct {
	ID            model.ObjectId  `bson:"_id,omitempty"`
	WorkflowUUID  string          `bson:"workflowUUID"`
	InstanceID    model.ObjectId  `bson:"instanceId"`
	DefinitionID  model.ObjectId  `bson:"definitionId"`
	WorkflowName  string          `bson:"workflowName"`
	InstanceName  string          `bson:"instanceName"`
	DefinitionName string         `bson:"definitionName"`
	OrgID         *model.ObjectId `bson:"orgId"`
	PathKey       string          `bson:"pathKey"`
	Version       int             `bson:"version"`
	Status        string          `bson:"status"`
	ActiveNodeIDs []string        `bson:"activeNodeIds"`
	TriggerSource  string          `bson:"triggerSource,omitempty"`
	StartedAt      time.Time       `bson:"startedAt"`
	Created        time.Time       `bson:"created"`
	Updated        time.Time       `bson:"updated"`
}
