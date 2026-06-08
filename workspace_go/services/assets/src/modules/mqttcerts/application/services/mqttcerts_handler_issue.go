package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	appConsts "assets/src/modules/mqttcerts/application/constants"
	dtos "assets/src/modules/mqttcerts/application/dtos"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

func (s *MqttCertsService) validateIssueRequest(req *dtos.IssueCertRequest) error {
	if req == nil {
		return errors.New("nil request")
	}
	if req.AssetUUID == "" {
		return errors.New("assetUUID required")
	}
	return nil
}

// assertGatewayCertEligible rejects issuing a gateway cert unless the asset is a
// LoRaWAN gateway in cert auth mode, so the gateway_certs endpoint can't mint a
// cert for a device or an eui-mode gateway.
func (s *MqttCertsService) assertGatewayCertEligible(ctx context.Context, assetUUID string) error {
	asset, err := s.deps.AssetService.GetByUUID(ctx, assetUUID)
	if err != nil {
		return fmt.Errorf("load asset: %w", err)
	}
	lw := asset.Protocol.Lorawan
	if lw == nil || lw.Kind != assetsContract.LorawanKindGateway || lw.Gateway == nil {
		return errors.New("asset is not a lorawan gateway")
	}
	if lw.Gateway.AuthMode != assetsContract.LorawanGatewayAuthModeCert {
		return errors.New("gateway authMode is not cert")
	}
	return nil
}

// signNewDeviceCert pulls the CA from RAM and emits a fresh cert + key.
// SubjectCN is the bare assetUUID — the broker plugin uses it as the
// MQTT username and is globally unique (Mongo idx_asset_uuid_unique).
// Tenant scoping flows from the auth projection's orgId, not the wire.
// TTL is resolved from the asset's `protocol.mqtt.certTTL` so the
// operator's per-asset preference (configured in the wizard's Step4)
// drives validity rather than a platform-wide constant.
func (s *MqttCertsService) signNewDeviceCert(ctx context.Context, assetUUID string, _ *reqCtx.RequestContext) (*signedCertBundle, error) {
	ca := s.deps.CAStore.Get()
	if ca == nil {
		return nil, ErrCANotReady
	}
	ttlDays, err := s.resolveCertTTLDays(ctx, assetUUID)
	if err != nil {
		return nil, err
	}
	subjectCN := assetUUID
	certPEM, keyPEM, serial, fingerprint, err := s.deps.Signer.SignDeviceCert(ca, subjectCN, ttlDays)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &signedCertBundle{
		certPEM:     certPEM,
		keyPEM:      keyPEM,
		serialHex:   fmt.Sprintf("%X", serial),
		fingerprint: fingerprint,
		subjectCN:   subjectCN,
		issuedAt:    now,
		expiresAt:   now.AddDate(0, 0, ttlDays),
	}, nil
}

// resolveCertTTLDays loads the asset and projects its certTTL into a
// day count clamped to [MinDeviceCertTTLDays, MaxDeviceCertTTLDays].
// Falls back to the platform default when the asset has no override
// (legacy assets created before the field landed). An invalid unit
// or out-of-range product yields a validation error so the operator
// hears about it instead of silently getting a default cert.
func (s *MqttCertsService) resolveCertTTLDays(ctx context.Context, assetUUID string) (int, error) {
	asset, err := s.deps.AssetService.GetByUUID(ctx, assetUUID)
	if err != nil {
		return 0, fmt.Errorf("load asset for cert TTL: %w", err)
	}
	// The certTTL override lives on the mqtt block (cert-mode device) or the
	// lorawan gateway block. Either absent falls back to the platform default.
	var value int
	var unit string
	switch {
	case asset.Protocol.Mqtt != nil && asset.Protocol.Mqtt.CertTTL != nil:
		value, unit = asset.Protocol.Mqtt.CertTTL.Value, asset.Protocol.Mqtt.CertTTL.Unit
	case asset.Protocol.Lorawan != nil && asset.Protocol.Lorawan.Gateway != nil && asset.Protocol.Lorawan.Gateway.CertTTL != nil:
		value, unit = asset.Protocol.Lorawan.Gateway.CertTTL.Value, asset.Protocol.Lorawan.Gateway.CertTTL.Unit
	default:
		return appConsts.DefaultDeviceCertTTLDays, nil
	}
	if value <= 0 {
		return 0, errors.New("certTTL.value must be > 0")
	}
	mult := appConsts.CertTTLUnitToDays(unit)
	if mult == 0 {
		return 0, fmt.Errorf("certTTL.unit invalid: %q", unit)
	}
	days := value * mult
	if days < appConsts.MinDeviceCertTTLDays || days > appConsts.MaxDeviceCertTTLDays {
		return 0, fmt.Errorf("certTTL out of range: %d days (allowed %d..%d)",
			days, appConsts.MinDeviceCertTTLDays, appConsts.MaxDeviceCertTTLDays)
	}
	return days, nil
}

// orgIDFromContext extracts the org id from the request context.
// Prefers OrgContext (X-Org-Context header) over scoped fallbacks.
// Returns "unknown" when nothing is available — broker plugin will
// deny via the cross-tenant guard before the cert lands anyway.
func (s *MqttCertsService) orgIDFromContext(rc *reqCtx.RequestContext) string {
	if rc == nil {
		return "unknown"
	}
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		return *rc.OrgContext
	}
	if len(rc.ScopedOrgIds) > 0 {
		return rc.ScopedOrgIds[0]
	}
	return "unknown"
}

// buildIssueResponse maps the signed bundle + CA chain to the wire DTO.
func (s *MqttCertsService) buildIssueResponse(bundle *signedCertBundle) *dtos.IssueCertResponse {
	ca := s.deps.CAStore.Get()
	chain := append([]byte{}, ca.CertPEM...)
	// NOTE: root cert isn't held in this RAM store — only the intermediate.
	// The chain returned here is intermediate-only; the device needs the
	// full chain for verification, but most TLS stacks accept the leaf
	// against the intermediate as long as the root is in the device's
	// trust store. Operators wanting the full chain can fetch from
	// mapexVault's /internal/pki/ca_chain endpoint and embed during
	// device provisioning.
	return &dtos.IssueCertResponse{
		Serial:      bundle.serialHex,
		Fingerprint: bundle.fingerprint,
		SubjectCN:   bundle.subjectCN,
		IssuedAt:    bundle.issuedAt,
		ExpiresAt:   bundle.expiresAt,
		CertPEM:     bundle.certPEM,
		KeyPEM:      bundle.keyPEM,
		CAChainPEM:  chain,
	}
}
