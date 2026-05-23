package routegroups_test

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
	rootClient  *httpclient.HTTPClient
	adminClient *httpclient.HTTPClient
	client      *httpclient.HTTPClient
	ctx         context.Context
)

func TestMain(m *testing.M) {
	if err := utils.SetupE2EEnvironment(); err != nil {
		panic("Failed to setup E2E environment: " + err.Error())
	}

	ctx = context.Background()

	rootClient = httpclient.New(httpclient.Config{BaseURL: constants.RouterURL})
	rootToken, err := utils.GetRootToken()
	if err != nil {
		panic("Failed to get ROOT token: " + err.Error())
	}
	rootClient.SetHeader("Authorization", "Bearer "+rootToken)
	// The router service follows the same coverage-middleware pattern as
	// mapexos: every CRUD endpoint demands X-Org-Context even when the
	// bearer carries the wildcard role.
	rootClient.SetHeader("X-Org-Context", constants.MapexosOrgID)

	adminClient = httpclient.New(httpclient.Config{BaseURL: constants.RouterURL})
	adminToken, err := utils.GetAdminToken()
	if err != nil {
		panic("Failed to get ADMIN token: " + err.Error())
	}
	adminClient.SetHeader("Authorization", "Bearer "+adminToken)
	adminClient.SetHeader("X-Org-Context", constants.MapexosOrgID)

	client = rootClient

	code := m.Run()

	os.Exit(code)
}

func TestCreateRouteGroup(t *testing.T) {
	payload := loadFixture(t, "create_routegroup.json")

	resp, err := client.Raw(ctx, "POST", "/api/v1/route_groups", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.NotNil(t, result.Data)

	routeGroupMap := result.Data.(map[string]interface{})
	routeGroupID := routeGroupMap["id"].(string)
	assert.NotEmpty(t, routeGroupID)

	assert.Equal(t, "API Routes v1", routeGroupMap["name"].(string))
	assert.Equal(t, true, routeGroupMap["enabled"].(bool))
	assert.Equal(t, "1.0.0", routeGroupMap["version"].(string))

	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID)
	})
}

func TestGetRouteGroupById(t *testing.T) {
	routeGroupID := createTestRouteGroup(t, "create_routegroup.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/route_groups/"+routeGroupID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	routeGroupMap := result.Data.(map[string]interface{})
	assert.Equal(t, routeGroupID, routeGroupMap["id"].(string))
	assert.Equal(t, "API Routes v1", routeGroupMap["name"].(string))
}

func TestUpdateRouteGroup(t *testing.T) {
	routeGroupID := createTestRouteGroup(t, "create_routegroup.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID)
	})

	updatePayload := loadFixture(t, "update_routegroup.json")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/route_groups/"+routeGroupID, updatePayload)
	require.NoError(t, err)
	// The router PATCH endpoint settled on 200 OK; older builds returned
	// 201 Created, accept both for build skew.
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated,
		"expected 200 or 201, got %d", resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	routeGroupMap := result.Data.(map[string]interface{})
	assert.Equal(t, "API Routes v1 Updated", routeGroupMap["name"].(string))
}

func TestDeleteRouteGroup(t *testing.T) {
	routeGroupID := createTestRouteGroup(t, "create_routegroup.json")

	resp, err := client.Raw(ctx, "DELETE", "/api/v1/route_groups/"+routeGroupID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	getResp, err := client.Raw(ctx, "GET", "/api/v1/route_groups/"+routeGroupID, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestListRouteGroups_BasicPagination(t *testing.T) {
	resp, err := client.Raw(ctx, "GET", "/api/v1/route_groups?page=1&perPage=10", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})

	if pagination, ok := paginatedResult["pagination"].(map[string]interface{}); ok {
		assert.NotNil(t, pagination["page"])
		assert.NotNil(t, pagination["perPage"])
		assert.Equal(t, float64(1), pagination["page"].(float64))
		assert.Equal(t, float64(10), pagination["perPage"].(float64))
	}

	if items, ok := paginatedResult["items"].([]interface{}); ok {
		t.Logf("Found %d route groups in the system", len(items))
	}
}

func TestListRouteGroups_FilterByName(t *testing.T) {
	routeGroupID1 := createTestRouteGroup(t, "create_routegroup.json")
	routeGroupID2 := createTestRouteGroup(t, "create_routegroup_versioned.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID1)
		cleanupRouteGroup(t, routeGroupID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/route_groups?name=API%20Routes", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		routeGroupMap := item.(map[string]interface{})
		name := routeGroupMap["name"].(string)
		assert.Contains(t, name, "API Routes")
	}
}

func TestListRouteGroups_FilterByEnabled(t *testing.T) {
	routeGroupID1 := createTestRouteGroup(t, "create_routegroup.json")
	routeGroupID2 := createTestRouteGroup(t, "create_routegroup_versioned.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID1)
		cleanupRouteGroup(t, routeGroupID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/route_groups?enabled=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		routeGroupMap := item.(map[string]interface{})
		assert.Equal(t, true, routeGroupMap["enabled"].(bool))
	}
}

func TestListRouteGroups_FilterByVersion(t *testing.T) {
	routeGroupID1 := createTestRouteGroup(t, "create_routegroup.json")
	routeGroupID2 := createTestRouteGroup(t, "create_routegroup_versioned.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID1)
		cleanupRouteGroup(t, routeGroupID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/route_groups?version=1.0.0", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		routeGroupMap := item.(map[string]interface{})
		if version, ok := routeGroupMap["version"].(string); ok {
			assert.Equal(t, "1.0.0", version)
		}
	}
}

func TestListRouteGroups_MultipleFilters(t *testing.T) {
	routeGroupID1 := createTestRouteGroup(t, "create_routegroup.json")
	routeGroupID2 := createTestRouteGroup(t, "create_routegroup_versioned.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID1)
		cleanupRouteGroup(t, routeGroupID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/route_groups?name=API&enabled=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		routeGroupMap := item.(map[string]interface{})
		assert.Contains(t, routeGroupMap["name"].(string), "API")
		assert.Equal(t, true, routeGroupMap["enabled"].(bool))
	}
}

func TestListRouteGroups_Projection(t *testing.T) {
	routeGroupID := createTestRouteGroup(t, "create_routegroup.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/route_groups?projection=name,enabled", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	if len(items) > 0 {
		routeGroupMap := items[0].(map[string]interface{})
		assert.NotNil(t, routeGroupMap["id"])
		assert.NotNil(t, routeGroupMap["name"])
		assert.NotNil(t, routeGroupMap["enabled"])
	}
}

func TestListRouteGroups_WithOrgContext(t *testing.T) {
	routeGroupID := createTestRouteGroup(t, "create_routegroup.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID)
	})

	resp, err := adminClient.Raw(ctx, "GET", "/api/v1/route_groups", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		routeGroupMap := item.(map[string]interface{})
		if orgId, ok := routeGroupMap["orgId"].(string); ok {
			assert.Equal(t, constants.MapexosOrgID, orgId)
		}
	}
}

func TestListRouteGroups_RootUser(t *testing.T) {
	routeGroupID := createTestRouteGroup(t, "create_routegroup.json")
	t.Cleanup(func() {
		cleanupRouteGroup(t, routeGroupID)
	})

	resp, err := rootClient.Raw(ctx, "GET", "/api/v1/route_groups", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	assert.NotNil(t, paginatedResult["items"])
}

func loadFixture(t *testing.T, filename string) map[string]interface{} {
	data, err := os.ReadFile("fixtures/" + filename)
	require.NoError(t, err, "Failed to load fixture: "+filename)

	var payload map[string]interface{}
	err = json.Unmarshal(data, &payload)
	require.NoError(t, err, "Failed to parse fixture: "+filename)

	return payload
}

func createTestRouteGroup(t *testing.T, fixtureFile string) string {
	payload := loadFixture(t, fixtureFile)

	resp, err := client.Raw(ctx, "POST", "/api/v1/route_groups", payload)
	require.NoError(t, err, "Failed to create test route group")
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected 201 Created")

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to parse create response")

	routeGroupMap := result.Data.(map[string]interface{})
	routeGroupID := routeGroupMap["id"].(string)
	require.NotEmpty(t, routeGroupID, "Route group ID should not be empty")

	return routeGroupID
}

func cleanupRouteGroup(t *testing.T, routeGroupID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/route_groups/"+routeGroupID, nil)
	if err != nil {
		t.Logf("Failed to cleanup route group %s: %v", routeGroupID, err)
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status code during cleanup: %d", resp.StatusCode)
	}
}
