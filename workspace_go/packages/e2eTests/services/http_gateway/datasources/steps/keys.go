package steps

// Bag keys this package writes. Other packages reading these keys import
// the constants from here.
const (
	// BagKeyDataSourceID is the Mongo ObjectID hex of the datasource
	// created by CreateDataSource. The HTTP heartbeat / event steps
	// reference it as the `ds` query parameter on /api/v1/heartbeat
	// and /api/v1/events.
	BagKeyDataSourceID = "httpGateway.dataSourceID"

	// BagKeyDataSourceApiKey is the plaintext apiKey embedded into the
	// datasource on create. The heartbeat / event steps present it on
	// the X-API-Key request header so the http_gateway auth middleware
	// accepts the request.
	BagKeyDataSourceApiKey = "httpGateway.dataSourceApiKey"
)
