package marketplaceclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"assets/src/modules/assettemplates/application/ports"

	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets_templates"
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

var _ ports.MarketplaceClientPort = (*MarketplaceClient)(nil)

// NewMarketplaceClient builds the outbound HTTP adapter to the mapexMarketplace
// catalog. Unauthenticated — the catalog is open and read-only, so no API key is
// attached.
func NewMarketplaceClient() ports.MarketplaceClientPort {
	base, _ := config.GetStringValue("asset_marketplace_url")
	return &MarketplaceClient{
		client: httpclient.New(httpclient.Config{
			BaseURL: base,
			Timeout: 10 * time.Second,
		}),
	}
}

// FetchBundle fetches the bundle for (vendor, slug). The endpoint serves the raw
// on-disk file verbatim (not an envelope), so the body read here is exactly the
// bytes the published sha256 is computed over — kept untouched in RawBytes for
// the hard-verify, then decoded separately into Bundle. Uses Raw so the 404 and
// the identity/integrity headers are visible (the wrapped Get would hide them).
func (c *MarketplaceClient) FetchBundle(ctx context.Context, vendor, slug string) (*ports.MarketplaceBundleFetch, error) {
	resp, err := c.client.Raw(ctx, http.MethodGet, "/api/v1/asset_templates/"+vendor+"/"+slug, nil)
	if err != nil {
		return nil, fmt.Errorf("transport: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrTemplateNotFound
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("marketplace status=%d body=%s", resp.StatusCode, string(body))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var bundle v1.MarketplaceBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return nil, fmt.Errorf("decode bundle: %w", err)
	}

	return &ports.MarketplaceBundleFetch{
		Bundle:          bundle,
		RawBytes:        raw,
		MarketplaceGuid: resp.Header.Get("X-Marketplace-Guid"),
		DeclaredSha256:  resp.Header.Get("X-Marketplace-Sha256"),
	}, nil
}
