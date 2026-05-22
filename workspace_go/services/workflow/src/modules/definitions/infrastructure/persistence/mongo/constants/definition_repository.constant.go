package constants

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// CollectionName is the MongoDB collection name for workflow definitions.
const CollectionName = "definitions"

// Indexes defines the indexes for the definitions collection.
var Indexes = []model.IndexDefinition{
	// Listing: org + date sort
	{
		Name: "idx_org_created",
		Keys: map[string]int{"orgId": 1, "created": -1},
	},
	// Listing: org + enabled filter (active workflows dropdown)
	{
		Name: "idx_org_enabled",
		Keys: map[string]int{"orgId": 1, "enabled": 1, "created": -1},
	},
	// Trigger path lookup (HTTP/webhook triggers resolve by pathKey + orgId)
	{
		Name: "idx_pathkey_org",
		Keys: map[string]int{"pathKey": 1, "orgId": 1},
	},
	// Name search
	{
		Name: "idx_name",
		Keys: map[string]int{"name": 1},
	},
}
