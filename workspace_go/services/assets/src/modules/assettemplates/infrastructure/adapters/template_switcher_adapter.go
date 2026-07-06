package adapters

import (
	"context"

	"assets/src/modules/assettemplates/application/ports"

	assetsDtos "assets/src/modules/assets/application/dtos"
	assetsPorts "assets/src/modules/assets/application/ports"
)

// templateSwitcherAdapter implements ports.TemplateSwitcherPort by updating the
// asset's template via the assets service (never a direct collection write).
type templateSwitcherAdapter struct {
	assets assetsPorts.AssetServicePort
}

// NewTemplateSwitcherAdapter returns a TemplateSwitcherPort over the assets
// service.
func NewTemplateSwitcherAdapter(assets assetsPorts.AssetServicePort) ports.TemplateSwitcherPort {
	return &templateSwitcherAdapter{assets: assets}
}

func (a *templateSwitcherAdapter) SwitchTemplate(ctx context.Context, assetID, targetTemplateID string) error {
	dto := &assetsDtos.AssetUpdateDTO{AssetTemplateID: &targetTemplateID}
	_, err := a.assets.UpdateAssetById(ctx, &assetID, dto)
	return err
}

var _ ports.TemplateSwitcherPort = (*templateSwitcherAdapter)(nil)
