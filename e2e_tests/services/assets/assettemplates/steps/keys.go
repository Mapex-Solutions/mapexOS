package steps

// Bag keys this package writes. Other packages reading these keys import
// the constants from here.
const (
	// BagKeyTemplateID is the asset template id returned by
	// CreateTemplate. Asset creation reads it to bind the new asset
	// to the saga template.
	BagKeyTemplateID = "assets.assetTemplateID"
)
