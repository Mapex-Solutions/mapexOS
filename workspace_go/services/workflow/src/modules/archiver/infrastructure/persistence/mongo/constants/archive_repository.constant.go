package constants

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// CollectionName is the MongoDB collection for workflow executions.
// The Archiver writes execution state to this collection with BulkWrite operations.
const CollectionName = "executions"

// Indexes defines the indexes for the executions collection.
var Indexes = []model.IndexDefinition{
	// Unique workflowUUID — prevents duplicate executions with same UUID
	{
		Name:   "idx_workflow_uuid_unique",
		Keys:   map[string]int{"workflowUUID": 1},
		Unique: true,
	},
	// Listing: org + status + date sort
	{
		Name: "idx_org_status_created",
		Keys: map[string]int{"orgId": 1, "status": 1, "created": -1},
	},
	// Listing: instance executions
	{
		Name: "idx_instance_created",
		Keys: map[string]int{"instanceId": 1, "created": -1},
	},
	// Listing: definition executions
	{
		Name: "idx_definition_created",
		Keys: map[string]int{"definitionId": 1, "created": -1},
	},
	// End-to-end tracking
	{
		Name:   "idx_event_tracker",
		Keys:   map[string]int{"eventTrackerId": 1},
		Sparse: true,
	},
	// TTL: auto-delete terminal executions after expireAt date
	{
		Name:               "idx_expire_at_ttl",
		Keys:               map[string]int{"expireAt": 1},
		Sparse:             true,
		ExpireAfterSeconds: &ttlZero,
	},
}

// ttlZero is used for MongoDB TTL index where the field itself contains the exact expiry time.
var ttlZero int32 = 0
