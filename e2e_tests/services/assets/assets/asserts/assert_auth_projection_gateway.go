package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// gatewayAuthEnvelope decodes the gateway-relevant subset of the internal
// asset-auth response (the platform {data: ...} envelope).
type gatewayAuthEnvelope struct {
	Data struct {
		Type    string `json:"type"`
		Lorawan *struct {
			Kind    string `json:"kind"`
			Gateway *struct {
				AuthMode          string `json:"authMode"`
				FrequencyPlanID   string `json:"frequencyPlanId"`
				CurrentCertSerial string `json:"currentCertSerial"`
			} `json:"gateway"`
		} `json:"lorawan"`
	} `json:"data"`
}

// AssertAuthProjectionGateway fetches the LNS Gateway Server's L3 source
// (GET /internal/asset-auth/:assetUUID) and asserts the asset is published as a
// LoRaWAN gateway: type=lorawan, lorawan.kind=gateway, with the gateway block
// carrying the expected frequency plan. This proves the assets service projects a
// gateway asset in the exact shape the LNS reads on connect.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetUUID  string  set by CreateAsset
func AssertAuthProjectionGateway(wantFrequencyPlanID string) saga.Assert {
	return saga.Assert{
		Name: "assets/assets.AssertAuthProjectionGateway",
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)
			headers := map[string]string{"X-API-Key": constants.InternalApiKey}
			resp, err := c.Clients.Assets.RawWithHeaders(c.Stdctx, http.MethodGet, "/internal/asset-auth/"+uuid, nil, headers)
			if err != nil {
				return fmt.Errorf("get asset-auth %s: %w", uuid, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("asset-auth %s: status %d", uuid, resp.StatusCode)
			}
			var env gatewayAuthEnvelope
			if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
				return fmt.Errorf("decode asset-auth %s: %w", uuid, err)
			}
			if env.Data.Type != "lorawan" {
				return fmt.Errorf("asset-auth %s: type = %q, want lorawan", uuid, env.Data.Type)
			}
			lw := env.Data.Lorawan
			if lw == nil || lw.Kind != "gateway" {
				return fmt.Errorf("asset-auth %s: expected lorawan.kind=gateway, got %+v", uuid, lw)
			}
			if lw.Gateway == nil {
				return fmt.Errorf("asset-auth %s: gateway block missing", uuid)
			}
			if lw.Gateway.FrequencyPlanID != wantFrequencyPlanID {
				return fmt.Errorf("asset-auth %s: frequencyPlanId = %q, want %q", uuid, lw.Gateway.FrequencyPlanID, wantFrequencyPlanID)
			}
			return nil
		},
	}
}

// AssertAuthProjectionGatewayCert fetches the LNS Gateway Server's L3 source and
// asserts the cert-mode gateway projection reflects the freshly-issued cert:
// type=lorawan, lorawan.kind=gateway, gateway.authMode=cert, the expected
// frequency plan, and currentCertSerial equal to the serial IssueGatewayCert
// captured on the bag. This proves the cert reflection (mqttcerts -> SetCurrentCert
// -> auth projection) lands in the exact shape the LNS pins on connect.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetUUID        string  set by CreateAsset
//   - assetSteps.BagKeyAssetCertSerial  string  set by IssueGatewayCert
func AssertAuthProjectionGatewayCert(wantFrequencyPlanID string) saga.Assert {
	return saga.Assert{
		Name: "assets/assets.AssertAuthProjectionGatewayCert",
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)
			wantSerial := c.MustGetString(assetSteps.BagKeyAssetCertSerial)
			headers := map[string]string{"X-API-Key": constants.InternalApiKey}
			resp, err := c.Clients.Assets.RawWithHeaders(c.Stdctx, http.MethodGet, "/internal/asset-auth/"+uuid, nil, headers)
			if err != nil {
				return fmt.Errorf("get asset-auth %s: %w", uuid, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("asset-auth %s: status %d", uuid, resp.StatusCode)
			}
			var env gatewayAuthEnvelope
			if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
				return fmt.Errorf("decode asset-auth %s: %w", uuid, err)
			}
			if env.Data.Type != "lorawan" {
				return fmt.Errorf("asset-auth %s: type = %q, want lorawan", uuid, env.Data.Type)
			}
			lw := env.Data.Lorawan
			if lw == nil || lw.Kind != "gateway" {
				return fmt.Errorf("asset-auth %s: expected lorawan.kind=gateway, got %+v", uuid, lw)
			}
			if lw.Gateway == nil {
				return fmt.Errorf("asset-auth %s: gateway block missing", uuid)
			}
			if lw.Gateway.AuthMode != "cert" {
				return fmt.Errorf("asset-auth %s: authMode = %q, want cert", uuid, lw.Gateway.AuthMode)
			}
			if lw.Gateway.FrequencyPlanID != wantFrequencyPlanID {
				return fmt.Errorf("asset-auth %s: frequencyPlanId = %q, want %q", uuid, lw.Gateway.FrequencyPlanID, wantFrequencyPlanID)
			}
			if lw.Gateway.CurrentCertSerial != wantSerial {
				return fmt.Errorf("asset-auth %s: currentCertSerial = %q, want %q (the issued serial)", uuid, lw.Gateway.CurrentCertSerial, wantSerial)
			}
			return nil
		},
	}
}
