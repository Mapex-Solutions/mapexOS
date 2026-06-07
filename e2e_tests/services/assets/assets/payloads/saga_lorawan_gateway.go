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

// euiFromRunID maps the runID to a stable uppercase 16-hex EUI (pad/truncate).
func euiFromRunID(runID string) string {
	h := strings.ToUpper(hex.EncodeToString([]byte(runID)))
	if len(h) >= 16 {
		return h[:16]
	}
	return (h + "0000000000000000")[:16]
}
