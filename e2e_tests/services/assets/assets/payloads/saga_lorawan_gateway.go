package payloads

import (
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
// label) pair, so a journey can provision more than one gateway without an EUI
// collision on the LNS.
func GatewayEUI(runID, label string) string {
	return euiFromRunID(runID + "-" + label)
}

// SagaLorawanGatewayFor returns an eui-mode gateway builder whose EUI (and name)
// are label-scoped, so several gateways can coexist in one journey. Otherwise
// identical to SagaLorawanGateway.
func SagaLorawanGatewayFor(label string) func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
	return func(runID, templateID, routeGroupID string) *AssetCreateBuilder {
		b := SagaLorawanGateway(runID, templateID, routeGroupID)
		b.spec.Name = fmt.Sprintf("saga-lorawan-gateway-%s-%s", label, runID)
		b.spec.AssetUUID = GatewayEUI(runID, label)
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

// euiFromRunID maps the runID to a stable uppercase 16-hex EUI (pad/truncate).
func euiFromRunID(runID string) string {
	h := strings.ToUpper(hex.EncodeToString([]byte(runID)))
	if len(h) >= 16 {
		return h[:16]
	}
	return (h + "0000000000000000")[:16]
}
