package collection

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// CollectionName is the MongoDB collection name for memberships.
const CollectionName = "memberships"

// Indexes defines the indexes for the memberships collection.
var Indexes = []model.IndexDefinition{
	{
		Name: "idx_membership_unique",
		Keys: map[string]int{
			"assigneeType": 1,
			"assigneeId":   1,
			"orgId":        1,
		},
		Unique: true,
	},
	{
		Name: "idx_assignee",
		Keys: map[string]int{
			"assigneeType": 1,
			"assigneeId":   1,
		},
	},
	{
		Name: "idx_org",
		Keys: map[string]int{"orgId": 1},
	},
}
