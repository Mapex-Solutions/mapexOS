package constants

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// CollectionName is the MongoDB collection name for asset templates.
const CollectionName = "assets_templates"

// Indexes defines the indexes for the assets_templates collection.
var Indexes = []model.IndexDefinition{
	{
		Name: "idx_org",
		Keys: map[string]int{"orgId": 1},
	},
	{
		Name: "idx_pathkey",
		Keys: map[string]int{"pathKey": 1},
	},
	{
		Name: "idx_manufacturer",
		Keys: map[string]int{"manufacturerId": 1},
	},
	{
		Name: "idx_model",
		Keys: map[string]int{"modelId": 1},
	},
	{
		Name: "idx_category",
		Keys: map[string]int{"categoryId": 1},
	},
}
