package entities

import (
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// FieldVocabulary is a single canonical dynamic-field name in the curated,
// multi-tenant authoring vocabulary served to the asset-template UI.
//
// Scope rules:
//   - IsSystem=true, OrgID=nil  -> central platform standard, shared to every org.
//   - IsSystem=false, OrgID=set -> an org's own additions (read-only here).
//
// Value is the English canonical field name. Hint carries the localized
// human-readable description keyed by locale ("en-US", "pt-BR").
type FieldVocabulary struct {
	ID       model.ObjectId    `bson:"_id,omitempty"`
	Value    string            `bson:"value"`
	Hint     map[string]string `bson:"hint"`
	Type     string            `bson:"type"`
	Unit     string            `bson:"unit"`
	Category string            `bson:"category"`
	IsSystem bool              `bson:"isSystem"`
	OrgID    *model.ObjectId   `bson:"orgId,omitempty"`
	Enabled  bool              `bson:"enabled"`
}
