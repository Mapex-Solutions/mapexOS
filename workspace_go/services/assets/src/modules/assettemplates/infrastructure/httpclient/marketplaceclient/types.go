package marketplaceclient

import (
	"errors"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// ErrTemplateNotFound is returned when the marketplace has no bundle for the
// requested (vendor, slug); the install maps it to an HTTP 404.
var ErrTemplateNotFound = errors.New("marketplace template not found")

// MarketplaceClient is the outbound HTTP adapter to the mapexMarketplace catalog.
type MarketplaceClient struct {
	client *httpclient.HTTPClient
}
