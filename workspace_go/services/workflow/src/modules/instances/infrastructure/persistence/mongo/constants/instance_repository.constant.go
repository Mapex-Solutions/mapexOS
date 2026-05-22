package constants

import (
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// CollectionName is the MongoDB collection name for workflow instance configs.
const CollectionName = "instances"

// Indexes defines the indexes for the instances collection.
var Indexes = []model.IndexDefinition{
	// Primary lookup: orgId + created for paginated listing.
	{
		Name: "idx_org_created",
		Keys: map[string]int{
			"orgId":   1,
			"created": -1,
		},
	},
	// Lookup by definition: which instances use this definition.
	{
		Name: "idx_definition",
		Keys: map[string]int{"definitionId": 1},
	},
}
