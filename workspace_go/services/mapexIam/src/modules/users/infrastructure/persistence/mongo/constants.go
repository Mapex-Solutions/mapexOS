package collection

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// CollectionName is the MongoDB collection name for users.
const CollectionName = "users"

// Indexes defines the indexes for the users collection.
var Indexes = []model.IndexDefinition{
	{
		Name:   "idx_email_unique",
		Keys:   map[string]int{"email": 1},
		Unique: true,
	},
}
