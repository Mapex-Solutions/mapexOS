package constants

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// CollectionName is the MongoDB collection name for assets.
const CollectionName = "assets"

// Indexes defines the indexes for the assets collection.
var Indexes = []model.IndexDefinition{
	{
		Name: "idx_org_created",
		Keys: map[string]int{"orgId": 1, "created": -1},
	},
	{
		Name:   "idx_asset_uuid_unique",
		Keys:   map[string]int{"assetUUID": 1},
		Unique: true,
	},
	{
		Name:   "idx_mqtt_username",
		Keys:   map[string]int{"protocol.mqtt.username": 1},
		Unique: true,
		Sparse: true,
	},
	{
		Name: "idx_template_id",
		Keys: map[string]int{"assetTemplateId": 1},
	},
	{
		Name: "idx_pathkey",
		Keys: map[string]int{"pathKey": 1},
	},
}
