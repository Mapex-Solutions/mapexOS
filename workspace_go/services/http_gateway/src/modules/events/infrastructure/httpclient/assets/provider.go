package assets

import (
	"time"

	"http_gateway/src/modules/events/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	configuration "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// NewOTAJobsPort creates the OTAJobsPort over an HTTP client configured for the
// Asset MS internal API.
//
// Configuration values:
//   - assets_service_url: Base URL of the Asset MS
//   - internal_api_key: API key for authenticating internal requests
func NewOTAJobsPort() ports.OTAJobsPort {
	assetsServiceURL, _ := configuration.GetStringValue("assets_service_url")
	apiKey, _ := configuration.GetStringValue("internal_api_key")

	client := httpclient.New(httpclient.Config{
		BaseURL: assetsServiceURL,
		APIKey:  apiKey,
		Timeout: 5 * time.Second,
	})

	return NewOTAJobsAdapter(client)
}
