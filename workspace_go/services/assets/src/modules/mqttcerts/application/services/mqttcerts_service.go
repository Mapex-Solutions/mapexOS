package services

import (
	"context"
	"errors"
	"fmt"

	assetPorts "assets/src/modules/assets/application/ports"
	mqttDi "assets/src/modules/mqttcerts/application/di"
	dtos "assets/src/modules/mqttcerts/application/dtos"
	mqttPorts "assets/src/modules/mqttcerts/application/ports"
	domConsts "assets/src/modules/mqttcerts/domain/constants"

	common "github.com/Mapex-Solutions/mapexGoKit/microservices/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time port checks.
var (
	_ mqttPorts.MqttCertsServicePort = (*MqttCertsService)(nil)
	_ common.Mountable               = (*MqttCertsService)(nil)
)

// Domain-level errors mapped to HTTP at the boundary.
var (
	ErrCANotReady          = errors.New("ca not ready")
	ErrAssetAlreadyHasCert = errors.New("asset already has currentCert; pass force=true to replace")
)

// New constructs the service. Returns the port interface.
func New(deps mqttDi.MqttCertsServiceDI) mqttPorts.MqttCertsServicePort {
	return &MqttCertsService{deps: deps}
}

// OnMount runs synchronously during DI lifecycle. ONE attempt with a
// short timeout to fetch the intermediate CA bundle from mapexVault.
// On success: cache in RAM + flip caReady. On failure: spawn a goroutine
// that retries with exponential backoff (no max attempts).
func (s *MqttCertsService) OnMount() {
	if ok := s.tryBootstrapSync(); ok {
		s.markCAReady()
		return
	}
	logger.Warn("[SERVICE:MqttCerts] OnMount: initial fetch failed, spawning retry goroutine")
	go s.spawnBootstrapGoroutine(context.Background())
}

// IsCAReady is the cheap atomic load used by the require_ca_ready middleware.
func (s *MqttCertsService) IsCAReady() bool {
	return s.deps.CAStore.IsReady()
}

// IssueCert signs a new device cert for the asset, reflects the
// cert metadata back onto the asset entity (via the assets module's
// SetCurrentCert port), and returns the PEM bundle once. PEM key
// bytes are never persisted server-side — only the serial,
// fingerprint, SubjectCN and validity window land on Mongo.
//
// The cert sign and the asset-side reflection are NOT atomic. On a
// reflection failure the operator gets a 5xx and the signed cert is
// effectively dropped on the floor (no broker entry, no zip
// download); from the device's perspective this is identical to a
// pre-sign failure. A successful reflection is the only path that
// returns the bundle to the operator.
func (s *MqttCertsService) IssueCert(ctx context.Context, rc *reqCtx.RequestContext, req *dtos.IssueCertRequest) (*dtos.IssueCertResponse, error) {
	if !s.IsCAReady() {
		return nil, ErrCANotReady
	}
	if err := s.validateIssueRequest(req); err != nil {
		return nil, err
	}
	bundle, err := s.signNewDeviceCert(ctx, req.AssetUUID, rc)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	if err := s.reflectIssuedCertOnAsset(ctx, req.AssetUUID, bundle); err != nil {
		return nil, fmt.Errorf("reflect cert on asset: %w", err)
	}
	return s.buildIssueResponse(bundle), nil
}

// RevokeCert clears the asset's `currentCert` (so the broker plugin
// drops the L1 entry on the next FANOUT-triggered invalidation), then
// persists an audit row in mqttRevokedCertificates. The audit row
// carries the resolved assetUUID so ListRevokedByAsset queries hit it.
func (s *MqttCertsService) RevokeCert(ctx context.Context, rc *reqCtx.RequestContext, serial string, reason string) error {
	r := domConsts.RevocationReason(reason)
	if r == "" {
		r = domConsts.ReasonUserAction
	}
	assetUUID, err := s.deps.AssetService.ClearCurrentCertBySerial(ctx, serial)
	if err != nil {
		return fmt.Errorf("clear currentCert on asset: %w", err)
	}
	return s.persistRevokedRow(ctx, rc, serial, assetUUID, r)
}

// reflectIssuedCertOnAsset mirrors the freshly-signed cert metadata
// onto `asset.currentCert` via the assets module's port and runs the
// standard L2 + FANOUT side effects so the broker plugin sees the new
// serial on the next CONNECT. Pulled into a helper to keep IssueCert
// readable as orchestration.
func (s *MqttCertsService) reflectIssuedCertOnAsset(ctx context.Context, assetUUID string, bundle *signedCertBundle) error {
	return s.deps.AssetService.SetCurrentCert(ctx, assetUUID, assetPorts.AssetCertificateInput{
		Serial:      bundle.serialHex,
		Fingerprint: bundle.fingerprint,
		SubjectCN:   bundle.subjectCN,
		IssuedAt:    bundle.issuedAt,
		ExpiresAt:   bundle.expiresAt,
	})
}

// ListRevokedByAsset reads the revoked rows for an asset.
func (s *MqttCertsService) ListRevokedByAsset(ctx context.Context, rc *reqCtx.RequestContext, assetUUID string) ([]*dtos.RevokedCertResponse, error) {
	rows, err := s.deps.RevokedRepo.FindByAssetUUID(ctx, assetUUID)
	if err != nil {
		return nil, fmt.Errorf("repo find: %w", err)
	}
	return s.mapToResponseDTOs(rows), nil
}

// HardDeleteByAssetUUID drops all revoked rows for the asset. Called by
// the assets module's asset-deletion side-effect path. NO L2 / fanout
// here — the caller (asset service) handles those.
func (s *MqttCertsService) HardDeleteByAssetUUID(ctx context.Context, assetUUID string) error {
	n, err := s.deps.RevokedRepo.DeleteByAssetUUID(ctx, assetUUID)
	if err != nil {
		return fmt.Errorf("repo delete: %w", err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:MqttCerts] hard-delete revoked rows count=%d assetUUID=%s", n, assetUUID))
	return nil
}
