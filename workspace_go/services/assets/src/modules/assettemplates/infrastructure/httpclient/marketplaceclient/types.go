package marketplaceclient

import (
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// MarketplaceClient is the outbound HTTP adapter to the mapexMarketplace catalog.
type MarketplaceClient struct {
	client *httpclient.HTTPClient
}
