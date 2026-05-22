package services

import (
	"assets/src/modules/assettemplates/domain/entities"
	ctx "context"
	"encoding/json"
	"fmt"

	templateContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assettemplates"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// writeScripts writes template scripts to object storage (L2 cache).
//
// This delegates to the TemplateStoragePort which handles entity-to-payload
// conversion and storage operations. The adapter encapsulates all MinIO details.
//
// Parameters:
//   - context: Context for cancellation and timeouts
//   - template: The template entity containing scripts
func (s *AssetTemplateService) writeScripts(context ctx.Context, template *entities.Assettemplate) {
	s.syncTemplateL2(context, template)
}

// deleteScripts removes template scripts from object storage (L2 cache).
//
// This delegates to the TemplateStoragePort which handles storage operations.
// The adapter encapsulates all MinIO details.
//
// Key format: {orgId}/{templateId}.json
// If IsSystem=true, orgId is "mapexos_public"
//
// Parameters:
//   - context: Context for cancellation and timeouts
//   - template: The template entity to delete
func (s *AssetTemplateService) deleteScripts(context ctx.Context, template *entities.Assettemplate) {
	s.deleteTemplateL2(context, template)
}

// getTemplateOrgId returns the org ID for storage key based on template visibility.
//
// Logic:
//   - IsSystem=true → returns "mapexos_public"
//   - IsSystem=false → returns the template's OrgID
func (s *AssetTemplateService) getTemplateOrgId(template *entities.Assettemplate) string {
	if template.IsSystem {
		return "mapexos_public"
	}
	if template.OrgID != nil {
		return template.OrgID.Hex()
	}
	return "mapexos_public"
}

// publishTemplateInvalidate publishes a FANOUT message to invalidate template cache.
// Consuming services (JS-Executor, Events) will invalidate their TieredCache.
// Called on: Create, Update, Delete
//
// Payload format: { "orgId": "...", "templateId": "..." }
// Cache key format: {orgId}/{templateId}
func (s *AssetTemplateService) publishTemplateInvalidate(context ctx.Context, template *entities.Assettemplate) {
	if template == nil || template.ID.IsZero() {
		return
	}

	// Build payload with orgId + templateId for cache key construction
	invalidatePayload := templateContract.TemplateInvalidatePayload{
		OrgId:      s.getTemplateOrgId(template),
		TemplateId: template.ID.Hex(),
	}

	payload, err := json.Marshal(invalidatePayload)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:AssetTemplate] Failed to marshal FANOUT payload for template %s: %v", template.ID.Hex(), err))
		return
	}

	if err := s.deps.NatsBus.PublishFanout(context, templateContract.FanoutTemplateSubject, payload); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:AssetTemplate] Failed to publish FANOUT for template %s: %v", template.ID.Hex(), err))
		return
	}

	logger.Debug(fmt.Sprintf("[SERVICE:AssetTemplate] FANOUT published for template %s/%s", invalidatePayload.OrgId, invalidatePayload.TemplateId))
}
