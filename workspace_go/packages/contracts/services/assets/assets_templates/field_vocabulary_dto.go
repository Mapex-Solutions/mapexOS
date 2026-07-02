package assetstemplate

// FieldVocabularyQuery represents the query parameters for the field
// vocabulary endpoint. The requested language selects which localized
// hint/label string is returned (en-US fallback when missing).
type FieldVocabularyQuery struct {
	Lang string `query:"lang" validate:"omitempty,max=10"`
}

// FieldVocabularyField is a single canonical dynamic-field entry already
// resolved to the requested language. The value is the English canonical
// field name; hint is the localized human-readable description.
type FieldVocabularyField struct {
	Value string `json:"value"`
	Hint  string `json:"hint"`
	Type  string `json:"type"`
	Unit  string `json:"unit"`
}

// FieldVocabularyGroup buckets fields under a category. Label is the
// localized category label already resolved to the requested language.
type FieldVocabularyGroup struct {
	Category string                 `json:"category"`
	Label    string                 `json:"label"`
	Fields   []FieldVocabularyField `json:"fields"`
}

// FieldVocabularyResponse is the read-only payload served to the asset
// template authoring UI: the curated, multi-tenant field vocabulary
// grouped by category in a fixed display order.
type FieldVocabularyResponse struct {
	Groups []FieldVocabularyGroup `json:"groups"`
}
