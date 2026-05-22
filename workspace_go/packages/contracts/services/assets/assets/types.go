package assets

// AssetInvalidatePayload is the cross-service wire contract for the FANOUT subject mapexos.fanout.asset.invalidate. Published by the assets service; consumed by router and js-executor.
type AssetInvalidatePayload struct {
	OrgId     string `json:"orgId"`
	AssetUUID string `json:"assetUUID"`
}
