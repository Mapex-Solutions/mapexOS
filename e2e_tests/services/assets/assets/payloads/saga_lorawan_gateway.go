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

// euiFromRunID maps the runID to a stable uppercase 16-hex EUI (pad/truncate).
func euiFromRunID(runID string) string {
	h := strings.ToUpper(hex.EncodeToString([]byte(runID)))
	if len(h) >= 16 {
		return h[:16]
	}
	return (h + "0000000000000000")[:16]
}
