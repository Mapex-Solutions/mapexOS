package constants

import model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"

// FieldVocabularyCollectionName is the MongoDB collection name for the
// curated field vocabulary served to the asset-template authoring UI.
const FieldVocabularyCollectionName = "field_vocabulary"

// FieldVocabularyIndexes defines the indexes for the field_vocabulary
// collection. The read path filters on enabled plus the isSystem/orgId
// visibility union, so those fields are indexed.
var FieldVocabularyIndexes = []model.IndexDefinition{
	{
		Name: "idx_enabled",
		Keys: map[string]int{"enabled": 1},
	},
	{
		Name: "idx_system",
		Keys: map[string]int{"isSystem": 1},
	},
	{
		Name: "idx_org",
		Keys: map[string]int{"orgId": 1},
	},
}
