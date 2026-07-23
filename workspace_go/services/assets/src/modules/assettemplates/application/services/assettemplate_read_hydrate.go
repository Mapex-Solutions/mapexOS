package services

import (
	ctx "context"
	"fmt"

	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/domain/entities"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// applyMarketplaceSource stamps the response's origin: the derived source string
// ("marketplace" when a marketplaceGuid is present, else "local") and the
// isMarketplace flag the UI reads to render the template read-only.
func applyMarketplaceSource(dto *dtos.AssetTemplateResponse, entity *entities.Assettemplate) {
	source := "local"
	if entity.MarketplaceGuid != nil {
		source = "marketplace"
	}
	dto.Source = &source

	isMarketplace := entity.IsMarketplace
	dto.IsMarketplace = &isMarketplace
}

// hydrateSharedContent fills a per-org marketplace LINK's heavy body (scripts,
// dynamic and available fields) from the durable shared content document, which
// holds the body once per marketplaceGuid. It is a no-op for hand-created
// templates and for the shared content document itself (org-less; already carries
// the body). A missing shared document is logged and leaves the body empty.
func (s *AssetTemplateService) hydrateSharedContent(c ctx.Context, dto *dtos.AssetTemplateResponse, entity *entities.Assettemplate) {
	if !entity.IsMarketplace || entity.OrgID == nil || entity.MarketplaceGuid == nil {
		return
	}

	shared, err := s.deps.AssetTemplateRepo.FindMarketplaceContentByGuid(c, *entity.MarketplaceGuid)
	if err != nil || shared == nil {
		logger.Warn(fmt.Sprintf("[SERVICE:AssetTemplate] Missing shared content for marketplace link %s (guid=%s)", entity.ID.Hex(), *entity.MarketplaceGuid))
		return
	}

	body, _ := mapper.EntityToDto[entities.Assettemplate, dtos.AssetTemplateResponse](shared)
	dto.ScriptTest = body.ScriptTest
	dto.ScriptProcessor = body.ScriptProcessor
	dto.ScriptValidator = body.ScriptValidator
	dto.ScriptConversion = body.ScriptConversion
	dto.AvailableFields = body.AvailableFields
	dto.DynamicFields = body.DynamicFields
}

// resolveContentTemplate returns the durable shared content document when the
// given record is a per-org marketplace link (which carries no body), so the L2
// write and the internal fallback response never serialize empty scripts. Any
// other record (local, system, or the shared content doc itself) passes through.
func (s *AssetTemplateService) resolveContentTemplate(c ctx.Context, template *entities.Assettemplate) *entities.Assettemplate {
	if !template.IsMarketplace || template.OrgID == nil || template.MarketplaceGuid == nil {
		return template
	}
	shared, err := s.deps.AssetTemplateRepo.FindMarketplaceContentByGuid(c, *template.MarketplaceGuid)
	if err != nil || shared == nil {
		return template
	}
	return shared
}
