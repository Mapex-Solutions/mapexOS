package listsclient

import "github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"

// ListsClient is the outbound HTTP adapter to mapexIam's internal lists resolve
// endpoint.
type ListsClient struct {
	client *httpclient.HTTPClient
}
