package ports

import (
	"context"
	"errors"

	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets_templates"
)

// ErrTemplateNotFound is returned by the marketplace client when the catalog has
// no bundle for the requested (vendor, slug). The install maps it to an HTTP 404.
// It lives on the port so the application layer can detect it without importing
// the infrastructure adapter.
var ErrTemplateNotFound = errors.New("marketplace template not found")

// MarketplaceBundleFetch is the result of fetching one template bundle from the
// marketplace: the decoded fields for building the local record, the EXACT raw
// bytes for sha256 verification (never re-encode these — the published hash is
// computed over these exact bytes the endpoint serves verbatim), and the
// identity/integrity metadata read from the response headers.
type MarketplaceBundleFetch struct {
	Bundle          v1.MarketplaceBundle
	RawBytes        []byte
	MarketplaceGuid string
	DeclaredSha256  string
}

// MarketplaceClientPort fetches asset template bundles from the mapexMarketplace
// catalog service.
type MarketplaceClientPort interface {
	// FetchBundle retrieves the bundle for (vendor, slug). Returns a distinct
	// not-found error (the caller maps it to 404) when no such template exists.
	FetchBundle(ctx context.Context, vendor, slug string) (*MarketplaceBundleFetch, error)
}
