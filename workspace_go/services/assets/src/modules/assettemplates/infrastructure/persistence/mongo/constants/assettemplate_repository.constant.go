package constants

import (
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
	// Per-org install key: a given org installs a given marketplace template at
	// most once. Partial so the unique constraint applies only to installed
	// records — hand-created templates (no marketplaceGuid) are never indexed and
	// so never collide on a shared null. The same index also guarantees a single
	// SHARED CONTENT document per guid: the shared doc has no orgId, so its index
	// key is {orgId: null, marketplaceGuid: guid} — one unique slot per guid,
	// distinct from each per-org link's {orgId: <org>, marketplaceGuid: guid}.
	{
		Name:                    "idx_org_marketplace_guid_unique",
		Keys:                    map[string]int{"orgId": 1, "marketplaceGuid": 1},
		Unique:                  true,
		PartialFilterExpression: bson.M{"marketplaceGuid": bson.M{"$exists": true}},
	},
	// Hydration lookup: read paths fetch the shared content by marketplaceGuid.
	{
		Name: "idx_marketplace_guid",
		Keys: map[string]int{"marketplaceGuid": 1},
	},
}
