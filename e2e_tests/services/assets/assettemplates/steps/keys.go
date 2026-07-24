package steps

// Bag keys this package writes. Other packages reading these keys import
// the constants from here.
const (
	// BagKeyTemplateID is the asset template id returned by
	// CreateTemplate. Asset creation reads it to bind the new asset
	// to the saga template.
	BagKeyTemplateID = "assets.assetTemplateID"

	// BagKeyMigrationPlanID is the template-migration plan id returned by
	// CreateMigrationPlan. The migration asserts read it to poll the plan and
	// its per-asset executions.
	BagKeyMigrationPlanID = "assets.migrationPlanID"
)

// TemplateIDKey is the label-scoped bag key for CreateTemplateWithLabel, so a
// journey that provisions several templates (e.g. an OTA source + target) keeps
// each id under its own key. Readers consume the id with the same label.
func TemplateIDKey(label string) string {
	return BagKeyTemplateID + ":" + label
}
