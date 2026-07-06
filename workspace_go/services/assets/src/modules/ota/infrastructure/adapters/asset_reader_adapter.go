package adapters

import (
	"context"

	"assets/src/modules/ota/application/ports"

	assetsPorts "assets/src/modules/assets/application/ports"
)

// assetReaderAdapter implements ports.AssetReaderPort by reading the asset via
// the assets service (assetId → assetUUID + protocol).
type assetReaderAdapter struct {
	assets assetsPorts.AssetServicePort
}

// NewAssetReaderAdapter returns an AssetReaderPort over the assets service.
func NewAssetReaderAdapter(assets assetsPorts.AssetServicePort) ports.AssetReaderPort {
	return &assetReaderAdapter{assets: assets}
}

func (a *assetReaderAdapter) GetAssetInfo(ctx context.Context, assetID string) (*ports.OTAAssetInfo, error) {
	asset, err := a.assets.GetAssetById(ctx, &assetID)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, nil
	}
	info := &ports.OTAAssetInfo{}
	if asset.AssetUUID != nil {
		info.AssetUUID = *asset.AssetUUID
	}
	if asset.Protocol != nil {
		info.Protocol = asset.Protocol.Type
	}
	return info, nil
}

var _ ports.AssetReaderPort = (*assetReaderAdapter)(nil)
