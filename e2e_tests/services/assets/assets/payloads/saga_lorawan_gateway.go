package payloads

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
)

// SagaGatewayFrequencyPlan is the plan the saga gateway registers with. EU is a
// safe default present in the vendored TTS plan set.
const SagaGatewayFrequencyPlan = "EU_863_870"

// SagaLorawanGateway returns a LoRaWAN gateway asset payload: protocol=lorawan,
// kind=gateway, authMode=eui (UDP, registered-EUI only). A gateway carries no
// device keys. The assetUUID is a 16-hex EUI derived from the runID so reruns do
// not collide and the value is a valid gateway EUI.
func SagaLorawanGateway(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	eui := euiFromRunID(runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-lorawan-gateway-%s", runID),
			Enabled:         true,
			DebugEnabled:    true,
			AssetUUID:       eui,
			AssetTemplateID: templateID,
			RouteGroupIds:   []string{routeGroupID},
			Protocol: contracts.ProtocolType{
				Type: "lorawan",
				Lorawan: &contracts.LorawanConfig{
					Kind: contracts.LorawanKindGateway,
					Gateway: &contracts.LorawanGatewayConfig{
						AuthMode:        contracts.LorawanGatewayAuthModeEUI,
						FrequencyPlanID: SagaGatewayFrequencyPlan,
					},
				},
			},
		},
	}
}

// SagaLorawanGatewayCert returns a LoRaWAN gateway asset payload in cert mode:
// protocol=lorawan, kind=gateway, authMode=cert (Basics Station, mTLS via the
// platform PKI). The mTLS cert is issued in a later step (POST
// /api/v1/gateway_certs); CertTTL is short (1 day) so the journey does not lean
// on the platform default and the cert metadata round-trip is exercised. Same EUI
// derivation as the eui variant.
func SagaLorawanGatewayCert(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	eui := euiFromRunID(runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-lorawan-gateway-cert-%s", runID),
			Enabled:         true,
			DebugEnabled:    true,
			AssetUUID:       eui,
			AssetTemplateID: templateID,
			RouteGroupIds:   []string{routeGroupID},
			Protocol: contracts.ProtocolType{
				Type: "lorawan",
				Lorawan: &contracts.LorawanConfig{
					Kind: contracts.LorawanKindGateway,
					Gateway: &contracts.LorawanGatewayConfig{
						AuthMode:        contracts.LorawanGatewayAuthModeCert,
						FrequencyPlanID: SagaGatewayFrequencyPlan,
						CertTTL:         &contracts.CertTTLConfig{Value: 1, Unit: "day"},
					},
				},
			},
		},
	}
}

// SagaGatewayAPIKey is the known plaintext Basics Station token the key-mode
// saga gateway registers with: 64 hex chars (32 bytes), within bcrypt's input
// limit. The assert bcrypt-compares the projection's apiKeyHash against it.
const SagaGatewayAPIKey = "00112233445566778899AABBCCDDEEFF00112233445566778899AABBCCDDEEFF"

// SagaLorawanGatewayKey returns a LoRaWAN gateway asset payload in key mode:
// protocol=lorawan, kind=gateway, authMode=key (Basics Station bearer token). The
// token (SagaGatewayAPIKey) is request-only; the assets service hashes it and the
// projection carries only the bcrypt hash. Same EUI derivation as the other
// variants.
func SagaLorawanGatewayKey(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	eui := euiFromRunID(runID)
	return &AssetCreateBuilder{
		spec: contracts.AssetCreate{
			Name:            fmt.Sprintf("saga-lorawan-gateway-key-%s", runID),
			Enabled:         true,
			DebugEnabled:    true,
			AssetUUID:       eui,
			AssetTemplateID: templateID,
			RouteGroupIds:   []string{routeGroupID},
			Protocol: contracts.ProtocolType{
				Type: "lorawan",
				Lorawan: &contracts.LorawanConfig{
					Kind: contracts.LorawanKindGateway,
					Gateway: &contracts.LorawanGatewayConfig{
						AuthMode:        contracts.LorawanGatewayAuthModeKey,
						FrequencyPlanID: SagaGatewayFrequencyPlan,
						APIKey:          SagaGatewayAPIKey,
					},
				},
			},
		},
	}
}

// GatewayEUI derives a stable uppercase 16-hex gateway EUI unique to a (runID,
// label) pair, so a journey can provision more than one gateway — and reruns can
// coexist — without an EUI collision on the LNS.
func GatewayEUI(runID, label string) string {
	return euiFromSeed(runID + "-gw-" + label)
}

// SagaLorawanGatewayFor returns an eui-mode gateway builder whose EUI (and name)
// are label-scoped, so several gateways can coexist in one journey. Otherwise
// identical to SagaLorawanGateway.
func SagaLorawanGatewayFor(label string) func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	return func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
		b := SagaLorawanGateway(runID, templateID, routeGroupID)
		b.spec.Name = fmt.Sprintf("saga-lorawan-gateway-%s-%s", label, runID)
		b.spec.AssetUUID = GatewayEUI(runID, label)
		b.spec.HealthMonitor = lorawanHealthMonitor()
		gatewayClearOptional(b, templateID, routeGroupID)
		return b
	}
}

// SagaLorawanGatewayKeyFor is the key-mode (Basics Station bearer token) analogue
// of SagaLorawanGatewayFor. The token stays the shared SagaGatewayAPIKey.
func SagaLorawanGatewayKeyFor(label string) func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	return func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
		b := SagaLorawanGatewayKey(runID, templateID, routeGroupID)
		b.spec.Name = fmt.Sprintf("saga-lorawan-gateway-key-%s-%s", label, runID)
		b.spec.AssetUUID = GatewayEUI(runID, label)
		b.spec.HealthMonitor = lorawanHealthMonitor()
		gatewayClearOptional(b, templateID, routeGroupID)
		return b
	}
}

// gatewayClearOptional drops the template / route-group fields when they were not
// supplied, so a gateway-only journey (which omits both — a LoRaWAN gateway asset
// requires neither) does not send an empty template id or a [""] route-group list.
func gatewayClearOptional(b *AssetCreateBuilder, templateID, routeGroupID string) {
	if templateID == "" {
		b.spec.AssetTemplateID = ""
	}
	if routeGroupID == "" {
		b.spec.RouteGroupIds = nil
	}
}

// euiFromSeed derives a stable, collision-resistant uppercase 16-hex EUI (8 bytes)
// from an arbitrary seed by hashing it — unlike euiFromRunID (which hex-encodes the
// raw string and truncates, so it only reflects the seed's first 8 characters), all
// 16 hex chars here carry entropy, so distinct seeds (labels, reruns) yield distinct
// EUIs. Deterministic: the same seed always maps to the same EUI, so a sim-driving
// step can recompute the exact identity an asset was provisioned with.
func euiFromSeed(seed string) string {
	sum := sha256.Sum256([]byte(seed))
	return strings.ToUpper(hex.EncodeToString(sum[:8]))
}

// euiFromRunID maps the runID to a stable uppercase 16-hex EUI (pad/truncate).
func euiFromRunID(runID string) string {
	h := strings.ToUpper(hex.EncodeToString([]byte(runID)))
	if len(h) >= 16 {
		return h[:16]
	}
	return (h + "0000000000000000")[:16]
}
