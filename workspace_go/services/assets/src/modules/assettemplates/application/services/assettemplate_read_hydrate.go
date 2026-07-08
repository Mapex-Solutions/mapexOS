package services

import (
	ctx "context"

	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/domain/entities"
)

// applyMarketplaceSource stamps the response's origin, derived from whether the
// record is a marketplace-installed link (has a marketplaceGuid) or hand-created.
func applyMarketplaceSource(dto *dtos.AssetTemplateResponse, entity *entities.Assettemplate) {
	source := "local"
	if entity.MarketplaceGuid != nil {
		source = "marketplace"
	}
	dto.Source = &source
}

// hydrateSharedContent fills a marketplace record's heavy body (scripts, dynamic
// and available fields) from the shared content cache — those live once per
// marketplaceGuid, not on the per-org link record. A cache miss leaves them
// empty; a no-op for hand-created records.
func (s *AssetTemplateService) hydrateSharedContent(c ctx.Context, dto *dtos.AssetTemplateResponse, entity *entities.Assettemplate) {
	if entity.MarketplaceGuid == nil {
		return
	}
	bundle, err := s.getCachedContent(c, *entity.MarketplaceGuid)
	if err != nil || bundle == nil {
		return
	}
	dto.ScriptTest = bundle.ScriptTest
	dto.ScriptProcessor = bundle.ScriptProcessor
	dto.ScriptValidator = stringPtrOrNil(bundle.ScriptValidator)
	dto.ScriptConversion = stringPtrOrNil(bundle.ScriptConversion)
	dto.AvailableFields = bundle.AvailableFields
	dto.DynamicFields = mapMarketplaceDynamicFields(bundle.DynamicFields)
}

// mapMarketplaceDynamicFields converts the wire dynamic fields into the response
// dynamic-field shape.
func mapMarketplaceDynamicFields(in []dtos.MarketplaceDynamicField) []dtos.DynamicField {
	if len(in) == 0 {
		return nil
	}
	out := make([]dtos.DynamicField, len(in))
	for i, f := range in {
		out[i] = dtos.DynamicField{
			FieldId: f.FieldId,
			Field:   f.Field,
			Value:   f.Value,
			Type:    f.Type,
			Status:  f.Status,
		}
	}
	return out
}

// stringPtrOrNil returns a pointer to s, or nil when s is empty.
func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
