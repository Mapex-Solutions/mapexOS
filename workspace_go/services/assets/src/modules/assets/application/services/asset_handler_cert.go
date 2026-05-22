package services

import (
	ctx "context"
	"fmt"
	"time"

	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// loadAssetForCertSync resolves the asset by its UUID prior to a cert
// reflection (issue / replace). The cert lifecycle is the source of
// truth for `currentCert`, but Mongo is the source of truth for the
// rest of the entity — this loads "before" so fanoutUpdateSideEffects
// can compute the right diff (health-state cleanup is a no-op for the
// cert path because the patch never flips HealthMonitor.Enabled).
func (s *AssetService) loadAssetForCertSync(c ctx.Context, assetUUID string) (*entities.Asset, error) {
	asset, err := s.deps.AssetRepo.FindByAssetUUID(c, &assetUUID)
	if err != nil || asset == nil || asset.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Asset not found for UUID: " + assetUUID},
		}
	}
	return asset, nil
}

// applyCurrentCertPatch persists the new cert subdoc onto the asset
// via the standard $set update path. The patch is intentionally
// scoped to `currentCert` + `updated` so unrelated fields (auth mode,
// route groups, health monitor) are not touched — only the cert
// metadata changes during an IssueCert reflection.
func (s *AssetService) applyCurrentCertPatch(
	c ctx.Context,
	before *entities.Asset,
	in ports.AssetCertificateInput,
) (*entities.Asset, error) {
	idHex := before.ID.Hex()
	payload := map[string]any{
		"currentCert": entities.AssetCertificate{
			Serial:      in.Serial,
			Fingerprint: in.Fingerprint,
			SubjectCN:   in.SubjectCN,
			IssuedAt:    in.IssuedAt,
			ExpiresAt:   in.ExpiresAt,
		},
		"updated": time.Now().UTC(),
	}
	updated, err := s.deps.AssetRepo.FindByIdAndUpdate(c, &idHex, payload)
	if err != nil {
		return nil, fmt.Errorf("set currentCert: %w", err)
	}
	if updated == nil || updated.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Asset not found after cert reflection"},
		}
	}
	return updated, nil
}

// findAssetByCertSerial locates the asset whose `currentCert.serial`
// matches the supplied serial. Used by the revoke reflection where
// the only handle the mqttcerts module has is the serial pulled from
// the URL path. A typical platform carries one active cert per asset
// so the filter is naturally unique; defensive in the rare case of
// stale Mongo state we still take the first match.
func (s *AssetService) findAssetByCertSerial(c ctx.Context, serial string) (*entities.Asset, error) {
	filters := model.Map{"currentCert.serial": serial}
	pagination := &model.PaginationOpts{Page: 1, PerPage: 1}
	result, err := s.deps.AssetRepo.FindWithFilters(c, filters, pagination, nil)
	if err != nil {
		return nil, fmt.Errorf("find asset by cert serial: %w", err)
	}
	if result == nil || len(result.Items) == 0 {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Asset not found for cert serial: " + serial},
		}
	}
	asset := result.Items[0]
	return &asset, nil
}

// applyClearCurrentCert drops the `currentCert` subdoc from the
// asset. The repository's FindByIdAndUpdate goes through `$set`, so
// the cleared field is written as JSON null rather than fully
// removed — functionally equivalent for downstream readers (the BSON
// decoder yields a nil *AssetCertificate either way) and keeps the
// path free of a bespoke $unset method on the repo.
func (s *AssetService) applyClearCurrentCert(c ctx.Context, before *entities.Asset) (*entities.Asset, error) {
	idHex := before.ID.Hex()
	payload := map[string]any{
		"currentCert": nil,
		"updated":     time.Now().UTC(),
	}
	updated, err := s.deps.AssetRepo.FindByIdAndUpdate(c, &idHex, payload)
	if err != nil {
		return nil, fmt.Errorf("clear currentCert: %w", err)
	}
	if updated == nil || updated.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Asset not found after cert clear"},
		}
	}
	return updated, nil
}
