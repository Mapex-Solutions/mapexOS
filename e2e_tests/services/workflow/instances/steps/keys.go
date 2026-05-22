package steps

// Bag keys this package writes. Other packages reading these keys
// import the constants from here.
const (
	// BagKeyInstanceID is the Mongo ObjectID hex of the workflow
	// instance created by CreateInstance. Route groups of kind=workflow
	// reference it on Router.Workflow.Data.instanceId.
	BagKeyInstanceID = "workflow.instanceID"
)
