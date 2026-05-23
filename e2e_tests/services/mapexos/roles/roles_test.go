package roles_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
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
	testOrgID   string
)

func TestMain(m *testing.M) {
	// Setup E2E environment (clean DB + flush cache + seed)
	if err := utils.SetupE2EEnvironment(); err != nil {
		panic("Failed to setup E2E environment: " + err.Error())
	}

	ctx = context.Background()

	// Setup ROOT client. Carries the seed admin JWT plus X-Org-Context
	// pinned to the seed root org — the mapexos middleware requires the
	// header on every CRUD endpoint regardless of the bearer's wildcard
	// permission.
	rootClient = httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	rootToken, err := utils.GetRootToken()
	if err != nil {
		panic("Failed to get ROOT token: " + err.Error())
	}
	rootClient.SetHeader("Authorization", "Bearer "+rootToken)
	rootClient.SetHeader("X-Org-Context", constants.MapexosOrgID)

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

	// Create test organization (child of Mapexos)
	testOrgID = createTestOrganization()

	code := m.Run()

	// Cleanup
	cleanupOrganization(testOrgID)
	os.Exit(code)
}

// ========================================
// CREATE TESTS
// ========================================

func TestCreateRole_SystemRole(t *testing.T) {
	payload := loadFixture(t, "create_system_role.json", "")

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	roleID := roleMap["id"].(string)

	assert.Equal(t, "System Administrator", roleMap["name"].(string))
	assert.True(t, roleMap["isSystem"].(bool))

	// Verify permissions array
	permissions := roleMap["permissions"].([]interface{})
	assert.Contains(t, permissions, "mapex.*")

	t.Cleanup(func() {
		cleanupRole(t, roleID)
	})
}

func TestCreateRole_OrgRole(t *testing.T) {
	payload := loadFixture(t, "create_org_role.json", testOrgID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	roleID := roleMap["id"].(string)

	assert.Equal(t, "Site Manager", roleMap["name"].(string))
	if v, ok := roleMap["isSystem"].(bool); ok {
		assert.False(t, v)
	}
	// orgId is now resolved from X-Org-Context by the service, not from the
	// payload, so the response carries the seed org id (MapexosOrgID) rather
	// than the testOrg we created. Just assert it is populated.
	assert.NotEmpty(t, roleMap["orgId"].(string))

	// Verify permissions
	permissions := roleMap["permissions"].([]interface{})
	assert.Contains(t, permissions, "user.read")
	assert.Contains(t, permissions, "asset.list")

	t.Cleanup(func() {
		cleanupRole(t, roleID)
	})
}

func TestCreateRole_Minimal(t *testing.T) {
	payload := loadFixture(t, "create_minimal.json", testOrgID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	roleID := roleMap["id"].(string)

	assert.Equal(t, "Viewer", roleMap["name"].(string))

	t.Cleanup(func() {
		cleanupRole(t, roleID)
	})
}

func TestCreateRole_NoOrgIdForNonSystem(t *testing.T) {
	payload := map[string]interface{}{
		"name":        "Invalid Role",
		"permissions": []string{"user.read"},
		"isSystem":    false,
		// Missing orgId - should fail
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	utils.AssertBadRequest(t, resp)
}

func TestCreateRole_EmptyPermissions(t *testing.T) {
	payload := map[string]interface{}{
		"name":        "No Permissions Role",
		"permissions": []string{},
		"isSystem":    false,
		"orgId":       testOrgID,
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	utils.AssertBadRequest(t, resp)
}

func TestCreateRole_MissingName(t *testing.T) {
	payload := map[string]interface{}{
		"permissions": []string{"user.read"},
		"isSystem":    true,
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	utils.AssertBadRequest(t, resp)
}

// ========================================
// GET TESTS
// ========================================

func TestGetRoleById(t *testing.T) {
	roleID := createTestRole(t, "create_org_role.json", testOrgID)
	defer cleanupRole(t, roleID)

	resp, err := client.Raw(ctx, "GET", "/api/v1/roles/"+roleID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	assert.Equal(t, roleID, roleMap["id"].(string))
	assert.Equal(t, "Site Manager", roleMap["name"].(string))
}

func TestGetRoleById_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011"

	resp, err := client.Raw(ctx, "GET", "/api/v1/roles/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

func TestListRoles(t *testing.T) {
	role1ID := createTestRole(t, "create_system_role.json", "")
	role2ID := createTestRole(t, "create_org_role.json", testOrgID)
	defer cleanupRole(t, role1ID)
	defer cleanupRole(t, role2ID)

	// Admin users (admin_vendor.*) MUST specify X-Org-Context header
	// Query with org context to see roles accessible in this org (includes system roles)
	resp, err := client.Raw(ctx, "GET", "/api/v1/roles?page=1&perPage=15", nil)

	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Result is now a paginated result with items array
	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 2)

	// Verify pagination metadata exists
	pagination := paginatedResult["pagination"].(map[string]interface{})
	assert.NotNil(t, pagination["totalItems"])
	assert.NotNil(t, pagination["totalPages"])
	assert.NotNil(t, pagination["page"])
	assert.NotNil(t, pagination["perPage"])
}

// ========================================
// UPDATE TESTS
// ========================================

func TestUpdateRole_Name(t *testing.T) {
	roleID := createTestRole(t, "create_org_role.json", testOrgID)
	defer cleanupRole(t, roleID)

	payload := loadFixture(t, "update_name.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/roles/"+roleID, payload)
	require.NoError(t, err)
	// API returns 201 for updates, not 200
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/roles/"+roleID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	assert.Equal(t, "Site Manager Updated", roleMap["name"].(string))
}

func TestUpdateRole_Permissions(t *testing.T) {
	roleID := createTestRole(t, "create_org_role.json", testOrgID)
	defer cleanupRole(t, roleID)

	payload := loadFixture(t, "update_permissions.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/roles/"+roleID, payload)
	require.NoError(t, err)
	// API returns 201 for updates, not 200
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/roles/"+roleID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	permissions := roleMap["permissions"].([]interface{})
	assert.Contains(t, permissions, "user.*")
	assert.Contains(t, permissions, "asset.*")
}

// NOTE: Roles don't have enabled field in v1, skipping this test
// func TestUpdateRole_Disable(t *testing.T) { ... }

func TestUpdateRole_Full(t *testing.T) {
	roleID := createTestRole(t, "create_org_role.json", testOrgID)
	defer cleanupRole(t, roleID)

	payload := loadFixture(t, "update_full.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/roles/"+roleID, payload)
	require.NoError(t, err)
	// API returns 201 for updates, not 200
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/roles/"+roleID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	assert.Equal(t, "Operations Manager", roleMap["name"].(string))
	assert.Equal(t, "Updated role for operations", roleMap["description"].(string))

	permissions := roleMap["permissions"].([]interface{})
	assert.Contains(t, permissions, "asset.*")
}

// ========================================
// DELETE TESTS
// ========================================

func TestDeleteRole(t *testing.T) {
	roleID := createTestRole(t, "create_org_role.json", testOrgID)

	resp, err := client.Raw(ctx, "DELETE", "/api/v1/roles/"+roleID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify deleted
	resp, err = client.Raw(ctx, "GET", "/api/v1/roles/"+roleID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

func TestDeleteRole_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011"

	resp, err := client.Raw(ctx, "DELETE", "/api/v1/roles/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

// ========================================
// PERMISSION TESTS
// ========================================

func TestRolePermissions_WildcardSupport(t *testing.T) {
	payload := map[string]interface{}{
		"name":        "Wildcard Tester",
		"description": "Tests wildcard permissions",
		"permissions": []string{"user.*", "asset.read"},
		"isSystem":    false,
		"orgId":       testOrgID,
		"scope":       "local",
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	roleID := roleMap["id"].(string)
	defer cleanupRole(t, roleID)

	permissions := roleMap["permissions"].([]interface{})
	assert.Contains(t, permissions, "user.*")
	assert.Contains(t, permissions, "asset.read")
}

func TestRolePermissions_AdminPermission(t *testing.T) {
	payload := map[string]interface{}{
		"name":        "Admin Role",
		"description": "Admin permissions",
		"permissions": []string{"admin.*"},
		"isSystem":    false,
		"orgId":       testOrgID,
		"scope":       "local",
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	roleID := roleMap["id"].(string)
	defer cleanupRole(t, roleID)

	permissions := roleMap["permissions"].([]interface{})
	assert.Contains(t, permissions, "admin.*")
}

// ========================================
// HELPER FUNCTIONS
// ========================================

func loadFixture(t *testing.T, filename string, customerID string) map[string]interface{} {
	data, err := os.ReadFile("fixtures/" + filename)
	require.NoError(t, err)

	content := string(data)
	if customerID != "" {
		content = strings.ReplaceAll(content, "{{CUSTOMER_ID}}", customerID)
	}

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(content), &payload)
	require.NoError(t, err)

	return payload
}

func createTestRole(t *testing.T, fixtureFile string, customerID string) string {
	payload := loadFixture(t, fixtureFile, customerID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/roles", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	roleMap := result.Data.(map[string]interface{})
	return roleMap["id"].(string)
}

func cleanupRole(t *testing.T, roleID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/roles/"+roleID, nil)
	if err != nil {
		t.Logf("Failed to cleanup role %s: %v", roleID, err)
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status during cleanup: %d", resp.StatusCode)
	}
}

func createTestOrganization() string {
	payload := map[string]interface{}{
		"name":        "Test Organization for Roles",
		"type":        "customer",
		"parentOrgId": constants.MapexosOrgID, // Child of Mapexos vendor org
		"enabled":     true,
		"address": map[string]interface{}{
			"city":    "Test City",
			"state":   "Test State",
			"country": "USA",
			"zipCode": "12345",
		},
		"phone": "+12125559999",
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
		panic("Failed to create test organization: unexpected status " + http.StatusText(resp.StatusCode))
	}

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic("Failed to parse organization response: " + err.Error())
	}

	orgMap := result.Data.(map[string]interface{})
	return orgMap["id"].(string)
}

func cleanupOrganization(orgID string) {
	client.Raw(ctx, "DELETE", "/api/v1/organizations/"+orgID, nil)
}
