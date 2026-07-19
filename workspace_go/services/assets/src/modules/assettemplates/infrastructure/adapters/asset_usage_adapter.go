package adapters

import (
	"context"

	"assets/src/modules/assettemplates/application/ports"

	assetsPorts "assets/src/modules/assets/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// assetUsageAdapter implements ports.AssetUsagePort by counting assets that
// reference a template through the assets service (never a direct collection read).
type assetUsageAdapter struct {
	assets assetsPorts.AssetServicePort
}

// NewAssetUsageAdapter returns an AssetUsagePort over the assets service.
func NewAssetUsageAdapter(assets assetsPorts.AssetServicePort) ports.AssetUsagePort {
	return &assetUsageAdapter{assets: assets}
}

func (a *assetUsageAdapter) CountAssetsUsingTemplate(ctx context.Context, requestContext *reqCtx.RequestContext, templateID string) (int64, error) {
	return a.assets.CountAssetsByTemplate(ctx, requestContext, templateID)
}

var _ ports.AssetUsagePort = (*assetUsageAdapter)(nil)
