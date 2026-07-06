package asserts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
)

// AssertAuthProjectionGatewayKey fetches the LNS Gateway Server's L3 source and
// asserts the key-mode gateway projection: type=lorawan, lorawan.kind=gateway,
// gateway.authMode=key, the expected frequency plan, and a non-empty apiKeyHash
// that bcrypt-matches the known token. This proves the assets service hashes the
// operator token and projects only the hash in the shape the LNS bcrypt-compares
// on connect.
//
// Reads (bag):
//   - assetSteps.BagKeyAssetUUID  string  set by CreateAsset
func AssertAuthProjectionGatewayKey(wantFrequencyPlanID string) saga.Assert {
	return saga.Assert{
		Name: "assets/assets.AssertAuthProjectionGatewayKey",
		Check: func(c *saga.Context) error {
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)
			headers := map[string]string{"X-API-Key": constants.InternalApiKey}
			resp, err := c.Clients.Assets.RawWithHeaders(c.Stdctx, http.MethodGet, "/internal/asset_auth/"+uuid, nil, headers)
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
			if lw.Gateway.AuthMode != "key" {
				return fmt.Errorf("asset-auth %s: authMode = %q, want key", uuid, lw.Gateway.AuthMode)
			}
			if lw.Gateway.FrequencyPlanID != wantFrequencyPlanID {
				return fmt.Errorf("asset-auth %s: frequencyPlanId = %q, want %q", uuid, lw.Gateway.FrequencyPlanID, wantFrequencyPlanID)
			}
			if lw.Gateway.APIKeyHash == "" {
				return fmt.Errorf("asset-auth %s: apiKeyHash is empty, want the bcrypt hash of the token", uuid)
			}
			if err := bcrypt.CompareHashAndPassword([]byte(lw.Gateway.APIKeyHash), []byte(assetPayloads.SagaGatewayAPIKey)); err != nil {
				return fmt.Errorf("asset-auth %s: apiKeyHash does not match the registered token: %w", uuid, err)
			}
			return nil
		},
	}
}
