package steps

// Bag keys this package writes. Other packages reading these keys
// import the constants from here.
const (
	// BagKeyDefinitionID is the Mongo ObjectID hex of the workflow
	// definition created by CreateDefinition. The CreateInstance step
	// reads it to populate the instance's definitionId field.
	BagKeyDefinitionID = "workflow.definitionID"

	// BagKeyDefinitionVersion is the version returned by the workflow
	// service on definition create. Required as a separate field on
	// the InstanceCreate body (definitionVersion).
	BagKeyDefinitionVersion = "workflow.definitionVersion"
)
