package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	"assets/src/modules/assettemplates/domain/entities"

	authContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// syncTemplateL2 writes the template scripts to MinIO (L2). On failure
// publishes a retry hint to the MAPEXOS-L2-WRITES stream so the in-
// module fallback consumer can reconcile against current Mongo state
// once MinIO recovers. Templates carry no auth credentials — only the
// single mapex-templates bucket is touched here.
func (s *AssetTemplateService) syncTemplateL2(c ctx.Context, t *entities.Assettemplate) {
	if err := s.deps.TemplateStoragePort.WriteScripts(c, t); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:AssetTemplate] L2 write failed for %s: %v", t.ID.Hex(), err))
		s.publishL2RetryTemplate(c, t)
	}
}

// deleteTemplateL2 removes the template scripts from L2. Best-effort
// — a failure is logged but doesn't propagate to the caller. Mongo is
// the source of truth and the broker doesn't consume templates.
func (s *AssetTemplateService) deleteTemplateL2(c ctx.Context, t *entities.Assettemplate) {
	orgId := s.getTemplateOrgId(t)
	if err := s.deps.TemplateStoragePort.DeleteScripts(c, orgId, t.ID.Hex()); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:AssetTemplate] L2 delete failed for %s: %v", t.ID.Hex(), err))
	}
}

// publishL2RetryTemplate sends a retry hint to the L2 writes stream
// with `Nats-Msg-Id: template:{id}` so NATS-native dedup (5s window)
// coalesces rapid successive failures on the same template.
func (s *AssetTemplateService) publishL2RetryTemplate(c ctx.Context, t *entities.Assettemplate) {
	payload, err := json.Marshal(map[string]string{"templateId": t.ID.Hex()})
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:AssetTemplate] L2 retry marshal failed for %s", t.ID.Hex()))
		return
	}
	msgId := "template:" + t.ID.Hex()
	if pubErr := s.deps.L2WritesPublisher.PublishRetry(c, authContract.L2WritesTemplateSubject, msgId, payload); pubErr != nil {
		logger.Error(pubErr, fmt.Sprintf("[SERVICE:AssetTemplate] L2 retry publish failed for %s", t.ID.Hex()))
	}
}

// fetchTemplateByID looks up the current template state from Mongo.
// Used by ProcessL2WriteRetry — the retry stream message carries only
// the id; the actual scripts are rebuilt from current Mongo state.
func (s *AssetTemplateService) fetchTemplateByID(c ctx.Context, templateId string) (*entities.Assettemplate, error) {
	t, err := s.deps.AssetTemplateRepo.FindById(c, &templateId)
	if err != nil {
		return nil, fmt.Errorf("l2 retry: find template %s: %w", templateId, err)
	}
	return t, nil
}
