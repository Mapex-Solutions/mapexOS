package memberships_test

import (
	"context"
	"encoding/json"
	"fmt"
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
	testOrgID   string // Test organization
	rootUserID  string // ROOT user ID (68f5bbce1aef22967c3ebb33)
	adminUserID string // ADMIN user ID (68f5bbce1aef22967c3ebb34)
	testGroupID string // Test group
	testRoleID  string // Test role
	testRole2ID string // Second test role
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

	// Get user IDs from constants (deterministic from seed script)
	rootUserID = constants.RootUserID
	adminUserID = constants.AdminUserID

	// Setup test resources (using ROOT client for unrestricted CRUD)
	testOrgID = createTestOrganization()

	testRoleID = createTestRole("Test Role 1", testOrgID)
	testRole2ID = createTestRole("Test Role 2", testOrgID)
	testGroupID = createTestGroup()

	code := m.Run()

	// Cleanup test resources
	cleanupGroup(testGroupID)
	cleanupRole(testRoleID)
	cleanupRole(testRole2ID)
	cleanupOrganization(testOrgID)

	os.Exit(code)
}

func TestCreateMembership_UserLocal(t *testing.T) {
	payload := loadFixture(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")

	resp, err := client.Raw(ctx, "POST", "/api/v1/memberships", payload)
	require.NoError(t, err)

	// If not created, print error for debugging
	if resp.StatusCode != http.StatusCreated {
		var errResp types.ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		t.Fatalf("Create membership failed: status=%d, errors=%v", resp.StatusCode, errResp.Errors)
	}

	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	membershipID := membershipMap["id"].(string)
	assert.Equal(t, "user", membershipMap["assigneeType"].(string))
	assert.Equal(t, adminUserID, membershipMap["assigneeId"].(string))
	assert.Equal(t, testOrgID, membershipMap["orgId"].(string))
	assert.Equal(t, "local", membershipMap["scope"].(string))
	assert.True(t, membershipMap["enabled"].(bool))

	// Verify roleIds array
	roleIds := membershipMap["roleIds"].([]interface{})
	assert.Equal(t, 1, len(roleIds))
	assert.Equal(t, testRoleID, roleIds[0].(string))

	t.Cleanup(func() {
		cleanupMembership(t, membershipID)
	})
}

func TestCreateMembership_UserRecursive(t *testing.T) {
	payload := loadFixture(t, "create_user_membership_recursive.json", adminUserID, testOrgID, testRoleID, "")

	resp, err := client.Raw(ctx, "POST", "/api/v1/memberships", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	membershipID := membershipMap["id"].(string)
	assert.Equal(t, "recursive", membershipMap["scope"].(string))

	t.Cleanup(func() {
		cleanupMembership(t, membershipID)
	})
}

func TestCreateMembership_Group(t *testing.T) {
	payload := loadFixture(t, "create_group_membership.json", testGroupID, testOrgID, testRoleID, "")

	resp, err := client.Raw(ctx, "POST", "/api/v1/memberships", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	membershipID := membershipMap["id"].(string)
	assert.Equal(t, "group", membershipMap["assigneeType"].(string))
	assert.Equal(t, testGroupID, membershipMap["assigneeId"].(string))

	t.Cleanup(func() {
		cleanupMembership(t, membershipID)
	})
}

func TestCreateMembership_MultipleRoles(t *testing.T) {
	payload := loadFixture(t, "create_multiple_roles.json", adminUserID, testOrgID, testRoleID, testRole2ID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/memberships", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	membershipID := membershipMap["id"].(string)

	// Verify roleIds array has 2 roles
	roleIds := membershipMap["roleIds"].([]interface{})
	assert.Equal(t, 2, len(roleIds))

	t.Cleanup(func() {
		cleanupMembership(t, membershipID)
	})
}

func TestGetMembershipById(t *testing.T) {
	membershipID := createTestMembership(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")
	defer cleanupMembership(t, membershipID)

	resp, err := client.Raw(ctx, "GET", "/api/v1/memberships/"+membershipID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	assert.Equal(t, membershipID, membershipMap["id"].(string))
	assert.NotNil(t, membershipMap["assigneeType"])
	assert.NotNil(t, membershipMap["assigneeId"])
}

func TestGetMembershipById_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011"

	resp, err := client.Raw(ctx, "GET", "/api/v1/memberships/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

func TestUpdateMembership_Scope(t *testing.T) {
	membershipID := createTestMembership(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")
	defer cleanupMembership(t, membershipID)

	payload := loadFixture(t, "update_scope.json", "", "", "", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/memberships/"+membershipID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/memberships/"+membershipID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	assert.Equal(t, "recursive", membershipMap["scope"].(string))
}

func TestUpdateMembership_Disable(t *testing.T) {
	membershipID := createTestMembership(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")
	defer cleanupMembership(t, membershipID)

	payload := loadFixture(t, "update_disable.json", "", "", "", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/memberships/"+membershipID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/memberships/"+membershipID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	assert.False(t, membershipMap["enabled"].(bool))
}

func TestUpdateMembership_Roles(t *testing.T) {
	membershipID := createTestMembership(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")
	defer cleanupMembership(t, membershipID)

	payload := loadFixture(t, "update_roles.json", "", "", testRole2ID, "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/memberships/"+membershipID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/memberships/"+membershipID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	roleIds := membershipMap["roleIds"].([]interface{})
	assert.Equal(t, 1, len(roleIds))
	assert.Equal(t, testRole2ID, roleIds[0].(string))
}

func TestDeleteMembership(t *testing.T) {
	membershipID := createTestMembership(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")

	resp, err := client.Raw(ctx, "DELETE", "/api/v1/memberships/"+membershipID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify deleted
	resp, err = client.Raw(ctx, "GET", "/api/v1/memberships/"+membershipID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

func TestListMemberships(t *testing.T) {
	membershipID1 := createTestMembership(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")
	membershipID2 := createTestMembership(t, "create_group_membership.json", testGroupID, testOrgID, testRoleID, "")
	defer cleanupMembership(t, membershipID1)
	defer cleanupMembership(t, membershipID2)

	// Request all memberships using includeAll=true
	resp, err := client.Raw(ctx, "GET", "/api/v1/memberships?includeAll=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Handle paginated response: {items: [...], pagination: {...}}
	dataMap := result.Data.(map[string]interface{})
	memberships := dataMap["items"].([]interface{})
	assert.GreaterOrEqual(t, len(memberships), 2)
}

func TestListMemberships_FilterByUser(t *testing.T) {
	membershipID := createTestMembership(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")
	defer cleanupMembership(t, membershipID)

	// Filter by userId
	resp, err := client.Raw(ctx, "GET", fmt.Sprintf("/api/v1/memberships?userId=%s&includeAll=true", adminUserID), nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	dataMap := result.Data.(map[string]interface{})
	memberships := dataMap["items"].([]interface{})
	assert.GreaterOrEqual(t, len(memberships), 1)

	// Verify all results are for the test user
	for _, item := range memberships {
		membership := item.(map[string]interface{})
		if membership["assigneeType"].(string) == "user" {
			assert.Equal(t, adminUserID, membership["assigneeId"].(string))
		}
	}
}

func TestGetMeCoverage(t *testing.T) {
	// Create a membership for the test user
	membershipID := createTestMembership(t, "create_user_membership_local.json", adminUserID, testOrgID, testRoleID, "")
	defer cleanupMembership(t, membershipID)

	resp, err := client.Raw(ctx, "GET", "/api/v1/me/coverage", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	coverageMap := result.Data.(map[string]interface{})
	assert.NotNil(t, coverageMap["userId"])

	customers := coverageMap["customers"].([]interface{})
	assert.NotNil(t, customers)
}

// Helper functions

func loadFixture(t *testing.T, filename string, assigneeID, orgID, roleID1, roleID2 string) map[string]interface{} {
	data, err := os.ReadFile("fixtures/" + filename)
	require.NoError(t, err)

	content := string(data)
	if assigneeID != "" {
		content = strings.ReplaceAll(content, "{{USER_ID}}", assigneeID)
		content = strings.ReplaceAll(content, "{{GROUP_ID}}", assigneeID)
	}
	if orgID != "" {
		content = strings.ReplaceAll(content, "{{ORG_ID}}", orgID)
	}
	if roleID1 != "" {
		content = strings.ReplaceAll(content, "{{ROLE_ID}}", roleID1)
		content = strings.ReplaceAll(content, "{{ROLE_ID_1}}", roleID1)
	}
	if roleID2 != "" {
		content = strings.ReplaceAll(content, "{{ROLE_ID_2}}", roleID2)
	}

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(content), &payload)
	require.NoError(t, err)

	return payload
}

func createTestMembership(t *testing.T, fixtureFile string, assigneeID, orgID, roleID1, roleID2 string) string {
	payload := loadFixture(t, fixtureFile, assigneeID, orgID, roleID1, roleID2)

	resp, err := client.Raw(ctx, "POST", "/api/v1/memberships", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	membershipMap := result.Data.(map[string]interface{})
	return membershipMap["id"].(string)
}

func cleanupMembership(t *testing.T, membershipID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/memberships/"+membershipID, nil)
	if err != nil {
		t.Logf("Failed to cleanup membership %s: %v", membershipID, err)
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status during cleanup: %d", resp.StatusCode)
	}
}

func createTestOrganization() string {
	// Use Mapexos vendor organization ID from constants (seeded with fixed ID)
	payload := map[string]interface{}{
		"name":        "Test Organization E2E Memberships",
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

	orgMap := result.Data.(map[string]interface{})
	return orgMap["id"].(string)
}

func cleanupOrganization(orgID string) {
	rootClient.Raw(ctx, "DELETE", "/api/v1/organizations/"+orgID, nil)
}

func createTestRole(name string, orgID string) string {
	payload := map[string]interface{}{
		"name":        name,
		"permissions": []string{"user.read", "user.list"},
		"orgId":       orgID,
		"isSystem":    false,
		"scope":       "local",
		"pathKey":     "",
	}

	resp, err := rootClient.Raw(ctx, "POST", "/api/v1/roles", payload)
	if err != nil {
		panic("Failed to create test role: " + err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("Failed to create test role: got status %d", resp.StatusCode))
	}

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic("Failed to parse role response: " + err.Error())
	}

	roleMap := result.Data.(map[string]interface{})
	return roleMap["id"].(string)
}

func cleanupRole(roleID string) {
	rootClient.Raw(ctx, "DELETE", "/api/v1/roles/"+roleID, nil)
}

func createTestGroup() string {
	payload := map[string]interface{}{
		"name":    "Test Group Membership",
		"enabled": true,
		"orgId":   testOrgID,
		"roleIds": []string{testRoleID},
	}

	resp, err := rootClient.Raw(ctx, "POST", "/api/v1/groups", payload)
	if err != nil {
		panic("Failed to create test group: " + err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("Failed to create test group: got status %d", resp.StatusCode))
	}

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic("Failed to parse group response: " + err.Error())
	}

	groupMap := result.Data.(map[string]interface{})
	return groupMap["id"].(string)
}

func cleanupGroup(groupID string) {
	rootClient.Raw(ctx, "DELETE", "/api/v1/groups/"+groupID, nil)
}
