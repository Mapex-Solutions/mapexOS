package services

import (
	"context"
	"time"

	"assets/src/modules/mqttcerts/domain/constants"
	"assets/src/modules/mqttcerts/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// persistRevokedRow inserts a row in mqttRevokedCertificates. The
// caller resolves the serial -> assetUUID link via
// AssetService.ClearCurrentCertBySerial before invoking this helper
// so the audit row carries the asset link, which lets
// ListRevokedByAsset surface revoked rows in the asset details
// drawer. Detailed metadata (fingerprint, SubjectCN, issuedAt) is
// not denormalized here — operators that need it can join against
// the asset's prior CurrentCert history once history tracking lands.
func (s *MqttCertsService) persistRevokedRow(ctx context.Context, rc *reqCtx.RequestContext, serial string, assetUUID string, reason constants.RevocationReason) error {
	now := time.Now().UTC()
	row := &entities.RevokedCertificate{
		Serial:    serial,
		AssetUUID: assetUUID,
		OrgID:     s.orgIDFromContext(rc),
		RevokedAt: now,
		Reason:    reason,
		Created:   now,
		Updated:   now,
	}
	return s.deps.RevokedRepo.Create(ctx, row)
}
