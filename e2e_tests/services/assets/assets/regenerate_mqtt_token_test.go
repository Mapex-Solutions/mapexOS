package assets_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/types"
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// parseJWTClaims decodes the payload segment of a NATS user JWT and
// returns the claim map. It does NOT verify the signature — tests use
// it to compare jti / exp across rotations, not to authenticate.
func parseJWTClaims(t *testing.T, jwtStr string) map[string]any {
	t.Helper()
	parts := strings.Split(jwtStr, ".")
	require.Len(t, parts, 3, "jwt must have header.payload.signature")
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err, "decode jwt payload")
	var claims map[string]any
	require.NoError(t, json.Unmarshal(payload, &claims), "unmarshal jwt payload")
	return claims
}

// createMqttAssetForRegenerate creates an MQTT-protocol asset using
// inline payload (the shared create_asset.json fixture predates the
// NKey-JWT migration and carries the removed `password` field). Returns
// the new asset's id and the one-shot credential string emitted by the
// platform.
func createMqttAssetForRegenerate(t *testing.T, runIDSuffix string) (assetID string, credential string) {
	t.Helper()
	payload := map[string]any{
		"name":            fmt.Sprintf("regenerate-test-%s", runIDSuffix),
		"enabled":         true,
		"assetUUID":       fmt.Sprintf("regen-uuid-%s", runIDSuffix),
		"assetTemplateId": templateID,
		"orgId":           testOrgID,
		"routeGroupIds":   []string{testRouteGroupID},
		"protocol": map[string]any{
			"type": "mqtt",
			"mqtt": map[string]any{
				"clientId": fmt.Sprintf("regen-client-%s", runIDSuffix),
				"username": fmt.Sprintf("regen-user-%s", runIDSuffix),
			},
		},
		"mqtt_token_ttl": "1y",
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/assets", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "create asset")

	var result types.StandardResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	resp.Body.Close()

	data := result.Data.(map[string]any)
	assetID, _ = data["id"].(string)
	require.NotEmpty(t, assetID, "create response missing id")

	credential, _ = data["mqtt_credential"].(string)
	require.NotEmpty(t, credential, "create response missing mqtt_credential — assets MS not on NKey-JWT build?")

	return assetID, credential
}

// TestRegenerate_200 covers the happy path: create an MQTT asset, hit
// the regenerate endpoint, and assert the response carries a fresh
// credential along with the user-pubkey + expiry metadata.
func TestRegenerate_200(t *testing.T) {
	assetID, _ := createMqttAssetForRegenerate(t, "200")
	t.Cleanup(func() { cleanupAsset(t, assetID) })

	resp, err := client.Raw(ctx, "POST", "/api/v1/assets/"+assetID+"/regenerate_mqtt_token", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "regenerate must return 200")

	var result types.StandardResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	resp.Body.Close()

	data, ok := result.Data.(map[string]any)
	require.True(t, ok, "response.data must decode to a map")
	assert.NotEmpty(t, data["mqtt_credential"], "credential present in response")
	assert.NotEmpty(t, data["mqtt_token_expires_at"], "expires-at present in response")
}

// TestRegenerate_NoToken_401 builds an anonymous httpclient (no
// Authorization header) and asserts the platform rejects the
// regenerate request as unauthorized.
func TestRegenerate_NoToken_401(t *testing.T) {
	anon := httpclient.New(httpclient.Config{BaseURL: constants.AssetsURL})
	// X-Org-Context is informational; the auth middleware fires first
	// and returns 401 before any org-scoped logic runs.
	resp, err := anon.Raw(ctx, "POST", "/api/v1/assets/0000000000000000deadbeef/regenerate_mqtt_token", nil)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "anonymous regenerate must be 401")
}

// TestRegenerate_404 calls the endpoint with an authenticated client
// against a synthetic non-existent id. The asset-load step in the
// service surfaces NOT_FOUND.
func TestRegenerate_404(t *testing.T) {
	resp, err := client.Raw(ctx, "POST", "/api/v1/assets/0000000000000000deadbeef/regenerate_mqtt_token", nil)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "regenerate on missing asset must be 404")
}

// TestRegenerate_NewCredentialDifferent asserts the regenerated JWT is
// distinct from the original — both as a raw string and at the JTI
// claim level (the platform persists the new JTI so any stale bearer
// fails the device-side jti-match invariant on refresh).
func TestRegenerate_NewCredentialDifferent(t *testing.T) {
	assetID, originalCredential := createMqttAssetForRegenerate(t, "diff")
	t.Cleanup(func() { cleanupAsset(t, assetID) })

	resp, err := client.Raw(ctx, "POST", "/api/v1/assets/"+assetID+"/regenerate_mqtt_token", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result types.StandardResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	resp.Body.Close()

	data := result.Data.(map[string]any)
	rotated, _ := data["mqtt_credential"].(string)
	require.NotEmpty(t, rotated)

	assert.NotEqual(t, originalCredential, rotated, "raw JWT string must change after regenerate")

	originalClaims := parseJWTClaims(t, originalCredential)
	rotatedClaims := parseJWTClaims(t, rotated)
	assert.NotEqual(t, originalClaims["jti"], rotatedClaims["jti"], "jti must change after regenerate")
}
