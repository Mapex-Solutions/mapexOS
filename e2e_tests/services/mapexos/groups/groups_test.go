package groups_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

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
	testOrgID   string // Mapexos organization ID from constants
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

	// Use Mapexos organization ID from constants (pre-seeded with fixed ID)
	testOrgID = constants.MapexosOrgID

	code := m.Run()

	os.Exit(code)
}

func TestCreateGroup_OrgGroup(t *testing.T) {
	payload := loadFixture(t, "create_org_group.json", testOrgID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/groups", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	groupID := groupMap["id"].(string)
	assert.Equal(t, "Engineering Team", groupMap["name"].(string))
	assert.False(t, groupMap["isSystem"].(bool))

	t.Cleanup(func() {
		cleanupGroup(t, groupID)
	})
}

func TestCreateGroup_SystemGroup(t *testing.T) {
	payload := loadFixture(t, "create_system_group.json", "")

	resp, err := client.Raw(ctx, "POST", "/api/v1/groups", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	groupID := groupMap["id"].(string)
	assert.True(t, groupMap["isSystem"].(bool))

	t.Cleanup(func() {
		cleanupGroup(t, groupID)
	})
}

func TestCreateGroup_Minimal(t *testing.T) {
	payload := loadFixture(t, "create_minimal.json", testOrgID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/groups", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	groupID := groupMap["id"].(string)

	t.Cleanup(func() {
		cleanupGroup(t, groupID)
	})
}

func TestCreateGroup_NoOrgIdForNonSystem(t *testing.T) {
	payload := map[string]interface{}{
		"name":     "Invalid Group",
		"enabled":  true,
		"isSystem": false,
		// Missing orgId - should fail
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/groups", payload)
	require.NoError(t, err)
	utils.AssertBadRequest(t, resp)
}

func TestGetGroupById(t *testing.T) {
	groupID := createTestGroup(t, "create_org_group.json", testOrgID)
	defer cleanupGroup(t, groupID)

	resp, err := client.Raw(ctx, "GET", "/api/v1/groups/"+groupID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	assert.Equal(t, groupID, groupMap["id"].(string))
	assert.NotNil(t, groupMap["name"])
}

func TestGetGroupById_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011"

	resp, err := client.Raw(ctx, "GET", "/api/v1/groups/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

func TestUpdateGroup_Name(t *testing.T) {
	groupID := createTestGroup(t, "create_org_group.json", testOrgID)
	defer cleanupGroup(t, groupID)

	payload := loadFixture(t, "update_name.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/groups/"+groupID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/groups/"+groupID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	assert.Equal(t, "Engineering Team Updated", groupMap["name"].(string))
}

func TestUpdateGroup_Description(t *testing.T) {
	groupID := createTestGroup(t, "create_org_group.json", testOrgID)
	defer cleanupGroup(t, groupID)

	payload := loadFixture(t, "update_description.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/groups/"+groupID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/groups/"+groupID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	assert.Equal(t, "Updated description for the team", groupMap["description"].(string))
}

func TestUpdateGroup_Disable(t *testing.T) {
	groupID := createTestGroup(t, "create_org_group.json", testOrgID)
	defer cleanupGroup(t, groupID)

	payload := loadFixture(t, "update_disable.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/groups/"+groupID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/groups/"+groupID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	assert.False(t, groupMap["enabled"].(bool))
}

func TestUpdateGroup_Full(t *testing.T) {
	groupID := createTestGroup(t, "create_org_group.json", testOrgID)
	defer cleanupGroup(t, groupID)

	payload := loadFixture(t, "update_full.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/groups/"+groupID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/groups/"+groupID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	assert.Equal(t, "New Team Name", groupMap["name"].(string))
	assert.Equal(t, "Completely updated group", groupMap["description"].(string))
}

func TestDeleteGroup(t *testing.T) {
	groupID := createTestGroup(t, "create_org_group.json", testOrgID)

	resp, err := client.Raw(ctx, "DELETE", "/api/v1/groups/"+groupID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify deleted
	resp, err = client.Raw(ctx, "GET", "/api/v1/groups/"+groupID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

func TestListGroups(t *testing.T) {
	// Create a test group in Mapexos organization
	groupID := createTestGroup(t, "create_org_group.json", testOrgID)
	defer cleanupGroup(t, groupID)

	// Add delay to ensure MongoDB commit completes
	time.Sleep(100 * time.Millisecond)

	// Request all accessible groups using includeAll=true
	resp, err := client.Raw(ctx, "GET", "/api/v1/groups?includeAll=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Handle paginated response: {items: [...], pagination: {...}}
	dataMap := result.Data.(map[string]interface{})
	groups := dataMap["items"].([]interface{})

	// Debug: log complete response payload
	t.Logf("=== COMPLETE RESPONSE PAYLOAD ===")
	responseJSON, _ := json.MarshalIndent(result.Data, "", "  ")
	t.Logf("%s", string(responseJSON))

	// Debug: log how many groups returned and their IDs
	t.Logf("\n=== GROUPS SUMMARY ===")
	t.Logf("LIST returned %d groups. Looking for groupID: %s", len(groups), groupID)
	for i, g := range groups {
		groupMap := g.(map[string]interface{})
		var id string
		if idVal, ok := groupMap["id"]; ok && idVal != nil {
			id = idVal.(string)
		} else if idVal, ok := groupMap["ID"]; ok && idVal != nil {
			id = idVal.(string)
		}
		var name interface{} = "<nil>"
		if nameVal, ok := groupMap["name"]; ok {
			name = nameVal
		}
		t.Logf("  Group %d: ID=%s, Name=%v", i+1, id, name)
	}

	// Should return at least the created group (may include others from previous tests/seeds)
	assert.GreaterOrEqual(t, len(groups), 1)

	// Verify the created group is in the results
	foundCreatedGroup := false
	for _, g := range groups {
		groupMap := g.(map[string]interface{})
		// Try both "id" (lowercase) and "ID" (uppercase) for compatibility
		var groupIDStr string
		if id, ok := groupMap["id"]; ok && id != nil {
			groupIDStr = id.(string)
		} else if id, ok := groupMap["ID"]; ok && id != nil {
			groupIDStr = id.(string)
		}

		if groupIDStr == groupID {
			foundCreatedGroup = true
			break
		}
	}
	assert.True(t, foundCreatedGroup, "Created group should be in the list")
}

// Helper functions

func loadFixture(t *testing.T, filename string, orgID string) map[string]interface{} {
	data, err := os.ReadFile("fixtures/" + filename)
	require.NoError(t, err)

	content := string(data)
	fmt.Printf("\n[DEBUG loadFixture] BEFORE replacement: orgID='%s', content contains {{ORG_ID}}: %v\n",
		orgID, strings.Contains(content, "{{ORG_ID}}"))

	if orgID != "" {
		content = strings.ReplaceAll(content, "{{ORG_ID}}", orgID)
	}

	fmt.Printf("[DEBUG loadFixture] AFTER replacement: content=%s\n", content)

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(content), &payload)
	require.NoError(t, err)

	fmt.Printf("[DEBUG loadFixture] PAYLOAD orgId field: %v\n\n", payload["orgId"])

	return payload
}

func createTestGroup(t *testing.T, fixtureFile string, orgID string) string {
	payload := loadFixture(t, fixtureFile, orgID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/groups", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	groupMap := result.Data.(map[string]interface{})
	return groupMap["id"].(string)
}

func cleanupGroup(t *testing.T, groupID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/groups/"+groupID, nil)
	if err != nil {
		t.Logf("Failed to cleanup group %s: %v", groupID, err)
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status during cleanup: %d", resp.StatusCode)
	}
}

func createTestOrganization() string {
	// Use Mapexos vendor organization ID from constants (pre-seeded with fixed ID)
	payload := map[string]interface{}{
		"name":        "Test Organization E2E Groups",
		"type":        "customer",
		"parentOrgId": constants.MapexosOrgID,
		"enabled":     true,
		"address": map[string]interface{}{
			"city":    "Test City",
			"state":   "TC",
			"country": "Test Country",
			"zipCode": "12345",
		},
		"phone": "+1234567890",
		"authConfig": map[string]interface{}{
			"providerType": "internal",
		},
		"accessPolicy": map[string]interface{}{
			"rolePolicy":   "merge",
			"defaultScope": "local",
		},
	}

	resp, err := rootClient.Raw(ctx, "POST", "/api/v1/organizations", payload)
	if err != nil {
		panic("Failed to create test organization: " + err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("Failed to create test organization: got status %d", resp.StatusCode))
	}

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic("Failed to parse organization response: " + err.Error())
	}

	if result.Data == nil {
		panic(fmt.Sprintf("Organization creation returned nil data: success=%v, message=%s", result.Success, result.Message))
	}

	orgMap := result.Data.(map[string]interface{})
	return orgMap["id"].(string)
}

func cleanupOrganization(orgID string) {
	rootClient.Raw(ctx, "DELETE", "/api/v1/organizations/"+orgID, nil)
}
