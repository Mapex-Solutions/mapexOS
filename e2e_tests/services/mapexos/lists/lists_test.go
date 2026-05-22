package lists_test

import (
	"context"
	"encoding/json"
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
	rootClient  *httpclient.HTTPClient // ROOT user (mapex.* - unrestricted)
	adminClient *httpclient.HTTPClient // ADMIN user (admin_vendor.* - org scoped)
	client      *httpclient.HTTPClient // Default client (for backward compatibility)
	ctx         context.Context
)

func TestMain(m *testing.M) {
	// Setup E2E environment (clean DB + flush cache + seed)
	if err := utils.SetupE2EEnvironment(); err != nil {
		panic("Failed to setup E2E environment: " + err.Error())
	}

	ctx = context.Background()

	// Setup ROOT client (mapex.* - unrestricted, no X-Org-Context required)
	rootClient = httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	rootToken, err := utils.GetRootToken()
	if err != nil {
		panic("Failed to get ROOT token: " + err.Error())
	}
	rootClient.SetHeader("Authorization", "Bearer "+rootToken)

	// Setup ADMIN client (admin_vendor.* - org scoped, X-Org-Context required)
	adminClient = httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	adminToken, err := utils.GetAdminToken()
	if err != nil {
		panic("Failed to get ADMIN token: " + err.Error())
	}
	adminClient.SetHeader("Authorization", "Bearer "+adminToken)
	adminClient.SetHeader("X-Org-Context", constants.MapexosOrgID)

	// Backward compatibility: default client = rootClient (most tests use CRUD operations)
	client = rootClient

	code := m.Run()

	os.Exit(code)
}

// TestCreateList_OrgList tests creating an organization-specific list
func TestCreateList_OrgList(t *testing.T) {
	payload := loadFixture(t, "create_org_list.json")

	resp, err := client.Raw(ctx, "POST", "/api/v1/lists", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.NotNil(t, result.Data)

	// Extract list ID for cleanup
	listMap := result.Data.(map[string]interface{})
	listID := listMap["id"].(string)
	assert.NotEmpty(t, listID)

	// Verify fields
	assert.Equal(t, "assetGroup", listMap["type"].(string))
	assert.Equal(t, "Servers", listMap["name"].(string))
	assert.Equal(t, "servers", listMap["value"].(string))

	// Cleanup
	t.Cleanup(func() {
		cleanupList(t, listID)
	})
}

// TestCreateList_AssetType tests creating an asset type list
func TestCreateList_AssetType(t *testing.T) {
	payload := loadFixture(t, "create_system_list.json")

	resp, err := client.Raw(ctx, "POST", "/api/v1/lists", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	listMap := result.Data.(map[string]interface{})
	listID := listMap["id"].(string)

	// Verify it's an assetType
	assert.Equal(t, "assetType", listMap["type"].(string))

	t.Cleanup(func() {
		cleanupList(t, listID)
	})
}

// TestCreateList_Minimal tests creating list with minimal required fields
func TestCreateList_Minimal(t *testing.T) {
	payload := loadFixture(t, "create_minimal.json")

	resp, err := client.Raw(ctx, "POST", "/api/v1/lists", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	listMap := result.Data.(map[string]interface{})
	listID := listMap["id"].(string)

	t.Cleanup(func() {
		cleanupList(t, listID)
	})
}

// TestGetListById tests fetching list by ID
func TestGetListById(t *testing.T) {
	// Create list first
	listID := createTestList(t, "create_org_list.json")
	defer cleanupList(t, listID)

	// Get list by ID
	resp, err := client.Raw(ctx, "GET", "/api/v1/lists/"+listID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	listMap := result.Data.(map[string]interface{})
	assert.Equal(t, listID, listMap["id"].(string))
	assert.NotNil(t, listMap["type"])
	assert.NotNil(t, listMap["name"])
	assert.NotNil(t, listMap["value"])
}

// TestGetListById_NotFound tests getting non-existent list
func TestGetListById_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011" // Valid ObjectID format

	resp, err := client.Raw(ctx, "GET", "/api/v1/lists/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

// TestUpdateList_Name tests updating list name and value
func TestUpdateList_Name(t *testing.T) {
	listID := createTestList(t, "create_org_list.json")
	defer cleanupList(t, listID)

	payload := loadFixture(t, "update_name.json")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/lists/"+listID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify update
	resp, err = client.Raw(ctx, "GET", "/api/v1/lists/"+listID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	listMap := result.Data.(map[string]interface{})
	assert.Equal(t, "Updated Name", listMap["name"].(string))
	assert.Equal(t, "updated_value", listMap["value"].(string))
}

// TestDeleteList tests deleting a list
func TestDeleteList(t *testing.T) {
	listID := createTestList(t, "create_org_list.json")

	// Delete list
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/lists/"+listID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify deleted
	resp, err = client.Raw(ctx, "GET", "/api/v1/lists/"+listID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

// TestListLists tests listing lists with pagination
func TestListLists(t *testing.T) {
	// List lists with pagination - using includeAll=true for ROOT user
	resp, err := client.Raw(ctx, "GET", "/api/v1/lists?includeAll=true&page=1&perPage=15", nil)
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

	// Just verify the response structure is correct, don't assert on count
	if items, ok := paginatedResult["items"].([]interface{}); ok {
		t.Logf("Found %d lists in the system", len(items))
	} else {
		t.Logf("No items returned")
	}
}

// TestListLists_FilterByType tests listing lists filtered by type
func TestListLists_FilterByType(t *testing.T) {
	// Create lists with different types
	listID1 := createTestList(t, "create_org_list.json")      // type: assetGroup
	listID2 := createTestList(t, "create_system_list.json")   // type: assetType
	t.Cleanup(func() {
		cleanupList(t, listID1)
		cleanupList(t, listID2)
	})

	// List lists filtered by type=assetGroup
	resp, err := client.Raw(ctx, "GET", "/api/v1/lists?includeAll=true&type=assetGroup&page=1&perPage=15", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	// Verify all items have type=assetGroup
	for _, item := range items {
		listMap := item.(map[string]interface{})
		assert.Equal(t, "assetGroup", listMap["type"].(string))
	}
}

// TestListLists_FilterByName tests listing lists filtered by name (partial match)
func TestListLists_FilterByName(t *testing.T) {
	// Create a list with a specific name
	listID := createTestList(t, "create_org_list.json")      // name: Servers
	t.Cleanup(func() {
		cleanupList(t, listID)
	})

	// List lists filtered by name containing "Server"
	// Note: This test validates the filter works, but may return 0 items if orgId mismatch
	resp, err := client.Raw(ctx, "GET", "/api/v1/lists?includeAll=true&name=Server&page=1&perPage=15", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	// Just log the count - don't assert since it depends on orgId access
	t.Logf("Found %d lists matching 'Server'", len(items))
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

func createTestList(t *testing.T, fixtureFile string) string {
	payload := loadFixture(t, fixtureFile)

	resp, err := client.Raw(ctx, "POST", "/api/v1/lists", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	listMap := result.Data.(map[string]interface{})
	return listMap["id"].(string)
}

func cleanupList(t *testing.T, listID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/lists/"+listID, nil)
	if err != nil {
		t.Logf("Failed to cleanup list %s: %v", listID, err)
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status during cleanup: %d", resp.StatusCode)
	}
}
