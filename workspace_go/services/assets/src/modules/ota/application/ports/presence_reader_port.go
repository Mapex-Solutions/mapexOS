package ports

import "context"

// PresenceReaderPort reads device online/offline state from the State Engine
// (the healthmonitor presence read model). Implemented by an adapter that only
// READS healthmonitor's read model — it does NOT modify healthmonitor.
type PresenceReaderPort interface {
	IsOnline(ctx context.Context, orgID, assetUUID string) (bool, error)
}

// OTAAssetInfo is the minimal asset-registry data the reconciler needs to
// dispatch: the device identity (assetUUID) and its transport (protocol).
type OTAAssetInfo struct {
	AssetUUID string
	Protocol  string
}

// AssetReaderPort resolves an asset's dispatch info by its id. Executions store
// the asset's Mongo id (assetId); dispatch needs the assetUUID + protocol.
type AssetReaderPort interface {
	GetAssetInfo(ctx context.Context, assetID string) (*OTAAssetInfo, error)
}
