package assets_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/types"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"
)

var (
	client           *httpclient.HTTPClient
	internalClient   *httpclient.HTTPClient
	routerClient     *httpclient.HTTPClient
	ctx              context.Context
	templateID       string                       // Will be set in TestMain
	testOrgID        = constants.MapexosOrgID // Seed root organization from mongodb-init
	testCategoryID   = "670a4cde48e006e3f95e8eb3"
	testAssetTypeID  = "670a4cde48e006e3f95e8eb4"
	testRouteGroupID string // Will be set in TestMain (route group seeded by test)
)

func TestMain(m *testing.M) {
	// Setup E2E environment (clean DB + flush cache + seed)
	if err := utils.SetupE2EEnvironment(); err != nil {
		panic("Failed to setup E2E environment: " + err.Error())
	}

	// Setup
	ctx = context.Background()
	client = httpclient.New(httpclient.Config{BaseURL: constants.AssetsURL})
	internalClient = httpclient.New(httpclient.Config{BaseURL: constants.AssetsURL})

	// Generate admin token for tests
	token, err := utils.GenerateAdminToken()
	if err != nil {
		panic("Failed to generate admin token: " + err.Error())
	}
	client.SetHeader("Authorization", "Bearer "+token)

	// Set organization context for client (Mapex vendor organization)
	client.SetHeader("X-Org-Context", testOrgID)

	// Setup internal client with API Key (default API Key for internal communication)
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "5230c2e2-e245-468d-89e8-94154cf520d0" // Default API Key
	}
	internalClient.SetHeader("X-API-Key", apiKey)

	// Set organization context for internal client
	internalClient.SetHeader("X-Org-Context", testOrgID)

	// Setup router client for creating prerequisite route group
	routerClient = httpclient.New(httpclient.Config{BaseURL: constants.RouterURL})
	routerClient.SetHeader("Authorization", "Bearer "+token)
	routerClient.SetHeader("X-Org-Context", testOrgID)

	// Create a route group (assets require RouteGroupIds min=1)
	testRouteGroupID = createTestRouteGroup()

	// Create a template for testing assets
	templateID = createTestTemplate()

	// Run tests
	code := m.Run()

	// Cleanup template & route group
	cleanupTemplate(templateID)
	cleanupRouteGroup(testRouteGroupID)

	os.Exit(code)
}

// TestCreateAsset_Valid tests creating an asset with all fields
func TestCreateAsset_Valid(t *testing.T) {
	payload := loadFixture(t, "create_asset.json")
	// Inject the template ID and route group ID created in TestMain
	payload["assetTemplateId"] = templateID
	payload["routeGroupIds"] = []string{testRouteGroupID}

	resp, err := client.Raw(ctx, "POST", "/api/v1/assets", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.NotNil(t, result.Data)

	// Extract asset ID for cleanup
	assetMap := result.Data.(map[string]interface{})
	assetID := assetMap["id"].(string)
	assert.NotEmpty(t, assetID)

	// Verify fields
	assert.Equal(t, "IoT Device 001", assetMap["name"].(string))
	assert.Equal(t, "ABC123DEF456", assetMap["assetUUID"].(string))
	if v, ok := assetMap["enabled"].(bool); ok {
		assert.True(t, v)
	}

	// Cleanup
	t.Cleanup(func() {
		cleanupAsset(t, assetID)
	})
}

// TestCreateAsset_Minimal tests creating asset with minimal required fields
func TestCreateAsset_Minimal(t *testing.T) {
	payload := loadFixture(t, "create_minimal.json")
	payload["assetTemplateId"] = templateID
	payload["routeGroupIds"] = []string{testRouteGroupID}

	resp, err := client.Raw(ctx, "POST", "/api/v1/assets", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assetMap := result.Data.(map[string]interface{})
	assetID := assetMap["id"].(string)

	t.Cleanup(func() {
		cleanupAsset(t, assetID)
	})
}

// TestGetAssetById tests fetching asset by ID
func TestGetAssetById(t *testing.T) {
	// Create asset first
	assetID := createTestAsset(t, "create_asset.json")
	defer cleanupAsset(t, assetID)

	// Get asset by ID
	resp, err := client.Raw(ctx, "GET", "/api/v1/assets/"+assetID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assetMap := result.Data.(map[string]interface{})
	assert.Equal(t, assetID, assetMap["id"].(string))
	assert.NotNil(t, assetMap["name"])
	assert.NotNil(t, assetMap["assetUUID"])
}

// TestGetAssetById_NotFound tests getting non-existent asset
func TestGetAssetById_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011" // Valid ObjectID format

	resp, err := client.Raw(ctx, "GET", "/api/v1/assets/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

// TestUpdateAsset_Name tests updating asset name and description
func TestUpdateAsset_Name(t *testing.T) {
	assetID := createTestAsset(t, "create_asset.json")
	defer cleanupAsset(t, assetID)

	payload := loadFixture(t, "update_name.json")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/assets/"+assetID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify update
	resp, err = client.Raw(ctx, "GET", "/api/v1/assets/"+assetID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assetMap := result.Data.(map[string]interface{})
	assert.Equal(t, "Updated Asset Name", assetMap["name"].(string))
}

// TestDeleteAsset tests deleting an asset
func TestDeleteAsset(t *testing.T) {
	assetID := createTestAsset(t, "create_asset.json")

	// Delete asset
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/assets/"+assetID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify deleted
	resp, err = client.Raw(ctx, "GET", "/api/v1/assets/"+assetID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

// TestListAssets tests listing assets with pagination
func TestListAssets(t *testing.T) {
	// List assets with pagination - using includeAll=true for ROOT user
	resp, err := client.Raw(ctx, "GET", "/api/v1/assets?includeAll=true&page=1&perPage=15", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Result is now a paginated result with items array
	paginatedResult := result.Data.(map[string]interface{})

	// Verify pagination metadata exists
	if pagination, ok := paginatedResult["pagination"].(map[string]interface{}); ok {
		assert.NotNil(t, pagination["page"])
		assert.NotNil(t, pagination["perPage"])
	}

	// Just verify the response structure is correct
	if items, ok := paginatedResult["items"].([]interface{}); ok {
		t.Logf("Found %d assets in the system", len(items))
	} else {
		t.Logf("No items returned")
	}
}

// TestListAssets_FilterByCategory tests listing assets filtered by category
func TestListAssets_FilterByCategory(t *testing.T) {
	// Create assets
	assetID1 := createTestAsset(t, "create_asset.json")
	t.Cleanup(func() {
		cleanupAsset(t, assetID1)
	})

	// List assets filtered by category
	resp, err := client.Raw(ctx, "GET", fmt.Sprintf("/api/v1/assets?includeAll=true&category=%s&page=1&perPage=15", testCategoryID), nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	// Just verify we got results
	t.Logf("Found %d assets with category filter", len(items))
}

// Helper functions

func loadFixture(t *testing.T, filename string) map[string]interface{} {
	data, err := os.ReadFile("fixtures/" + filename)
	require.NoError(t, err, "Failed to load fixture: "+filename)

	var payload map[string]interface{}
	err = json.Unmarshal(data, &payload)
	require.NoError(t, err, "Failed to parse fixture: "+filename)

	return payload
}

func createTestAsset(t *testing.T, fixtureFile string) string {
	payload := loadFixture(t, fixtureFile)
	payload["assetTemplateId"] = templateID
	payload["routeGroupIds"] = []string{testRouteGroupID}

	resp, err := client.Raw(ctx, "POST", "/api/v1/assets", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assetMap := result.Data.(map[string]interface{})
	return assetMap["id"].(string)
}

func cleanupAsset(t *testing.T, assetID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/assets/"+assetID, nil)
	if err != nil {
		t.Logf("Failed to cleanup asset %s: %v", assetID, err)
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status during cleanup: %d", resp.StatusCode)
	}
}

func createTestTemplate() string {
	payload := map[string]interface{}{
		"name":             "Test Template",
		"enabled":          true,
		"manufacturerName": "Test Corp",
		"modelName":        "TEST-001",
		"assetIdPath":      "payload.deviceId",
		"scriptValidator":  "function validate(data) { return true; }",
		"scriptConversion": "function convert(data) { return data; }",
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/asset_templates", payload)
	if err != nil {
		panic("Failed to create test template: " + err.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("Failed to create test template: status %d", resp.StatusCode))
	}

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic("Failed to parse template response: " + err.Error())
	}

	templateMap := result.Data.(map[string]interface{})
	return templateMap["id"].(string)
}

func cleanupTemplate(templateID string) {
	_, _ = client.Raw(ctx, "DELETE", "/api/v1/asset_templates/"+templateID, nil)
}

// createTestRouteGroup provisions a save_event route group on the router
// service so the assets test has a valid RouteGroupId to satisfy the
// AssetCreate validate:"required,min=1" rule. The same shape is used by
// the saga journeys (services/router/routegroups/payloads).
func createTestRouteGroup() string {
	payload := map[string]interface{}{
		"name":    "assets-e2e-save-event",
		"version": "1.0.0",
		"enabled": true,
		"routers": []map[string]interface{}{
			{
				"kind":      "save_event",
				"saveEvent": map[string]interface{}{},
			},
		},
	}

	resp, err := routerClient.Raw(ctx, "POST", "/api/v1/route_groups", payload)
	if err != nil {
		panic("Failed to create test route group: " + err.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("Failed to create test route group: status %d", resp.StatusCode))
	}

	var result types.StandardResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic("Failed to parse route group response: " + err.Error())
	}

	rgMap := result.Data.(map[string]interface{})
	return rgMap["id"].(string)
}

func cleanupRouteGroup(routeGroupID string) {
	_, _ = routerClient.Raw(ctx, "DELETE", "/api/v1/route_groups/"+routeGroupID, nil)
}
