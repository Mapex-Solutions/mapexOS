package constants

// FieldVocabularyCategoryOrder is the canonical enumeration of field
// vocabulary categories, in their fixed order. This is domain vocabulary:
// the set of categories a canonical dynamic-field name can belong to. The
// presentation layer renders groups in this order and omits empty ones.
var FieldVocabularyCategoryOrder = []string{
	"climate",
	"air-quality",
	"presence-counting",
	"agriculture-liquids",
	"electrical",
	"mechanical-position",
	"weather",
	"device-health",
}
