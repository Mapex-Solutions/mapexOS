package users_test

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

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// TestCreateUser_Internal tests creating user with internal auth
func TestCreateUser_Internal(t *testing.T) {
	basePayload := loadFixture(t, "create_internal.json")

	// Wrap in onboarding payload with membership
	onboardingPayload := make(map[string]interface{})
	for k, v := range basePayload {
		onboardingPayload[k] = v
	}
	onboardingPayload["memberships"] = []map[string]interface{}{
		{
			"orgId": constants.MapexosOrgID,
			"roles": []string{constants.SuperAdminRoleID},
			"scope": "local",
		},
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/onboarding/users", onboardingPayload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.NotNil(t, result.Data)

	// Extract user ID for cleanup (onboarding returns {user: {...}, memberships: [...]})
	onboardingResult := result.Data.(map[string]interface{})
	userMap := onboardingResult["user"].(map[string]interface{})
	userID := userMap["id"].(string)
	assert.NotEmpty(t, userID)

	// Cleanup
	t.Cleanup(func() {
		cleanupUser(t, userID)
	})
}

// TestCreateUser_Google tests creating user with Google OAuth
func TestCreateUser_Google(t *testing.T) {
	basePayload := loadFixture(t, "create_google.json")

	// Wrap in onboarding payload with membership
	onboardingPayload := make(map[string]interface{})
	for k, v := range basePayload {
		onboardingPayload[k] = v
	}
	onboardingPayload["memberships"] = []map[string]interface{}{
		{
			"orgId": constants.MapexosOrgID,
			"roles": []string{constants.SuperAdminRoleID},
			"scope": "local",
		},
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/onboarding/users", onboardingPayload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	onboardingResult := result.Data.(map[string]interface{})
	userMap := onboardingResult["user"].(map[string]interface{})
	userID := userMap["id"].(string)

	t.Cleanup(func() {
		cleanupUser(t, userID)
	})
}

// TestCreateUser_Minimal tests creating user with minimal required fields
func TestCreateUser_Minimal(t *testing.T) {
	basePayload := loadFixture(t, "create_minimal.json")

	// Wrap in onboarding payload with membership
	onboardingPayload := make(map[string]interface{})
	for k, v := range basePayload {
		onboardingPayload[k] = v
	}
	onboardingPayload["memberships"] = []map[string]interface{}{
		{
			"orgId": constants.MapexosOrgID,
			"roles": []string{constants.SuperAdminRoleID},
			"scope": "local",
		},
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/onboarding/users", onboardingPayload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	onboardingResult := result.Data.(map[string]interface{})
	userMap := onboardingResult["user"].(map[string]interface{})
	userID := userMap["id"].(string)

	t.Cleanup(func() {
		cleanupUser(t, userID)
	})
}

// TestCreateUser_InvalidEmail tests validation for invalid email
func TestCreateUser_InvalidEmail(t *testing.T) {
	payload := map[string]interface{}{
		"email":    "invalid-email",
		"password": "Pass123!",
		"authProvider": map[string]interface{}{
			"type": "internal",
		},
		"firstName": "Test",
		"lastName":  "User",
		"enabled":   true,
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/users", payload)
	require.NoError(t, err)
	utils.AssertBadRequest(t, resp)
}

// TestGetUserById tests fetching user by ID
func TestGetUserById(t *testing.T) {
	// Create user first
	userID := createTestUser(t, "create_internal.json")
	defer cleanupUser(t, userID)

	// Get user by ID
	resp, err := client.Raw(ctx, "GET", "/api/v1/users/"+userID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	userMap := result.Data.(map[string]interface{})
	assert.Equal(t, userID, userMap["id"].(string))
	assert.NotNil(t, userMap["email"])
	assert.NotNil(t, userMap["firstName"])
	assert.NotNil(t, userMap["lastName"])
}

// TestGetUserById_NotFound tests getting non-existent user
func TestGetUserById_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011" // Valid ObjectID format

	resp, err := client.Raw(ctx, "GET", "/api/v1/users/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

// TestUpdateUser_Name tests updating user name
func TestUpdateUser_Name(t *testing.T) {
	userID := createTestUser(t, "create_internal.json")
	defer cleanupUser(t, userID)

	payload := loadFixture(t, "update_name.json")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/users/"+userID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify update
	resp, err = client.Raw(ctx, "GET", "/api/v1/users/"+userID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	userMap := result.Data.(map[string]interface{})
	assert.Equal(t, "John Updated", userMap["firstName"].(string))
	assert.Equal(t, "Doe Updated", userMap["lastName"].(string))
}

// TestUpdateUser_Password tests updating user password
func TestUpdateUser_Password(t *testing.T) {
	userID := createTestUser(t, "create_internal.json")
	defer cleanupUser(t, userID)

	payload := loadFixture(t, "update_password.json")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/users/"+userID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify changePasswordNextLogin flag
	resp, err = client.Raw(ctx, "GET", "/api/v1/users/"+userID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	userMap := result.Data.(map[string]interface{})
	assert.True(t, userMap["changePasswordNextLogin"].(bool))
}

// TestUpdateUser_Full tests updating multiple user fields
func TestUpdateUser_Full(t *testing.T) {
	userID := createTestUser(t, "create_internal.json")
	defer cleanupUser(t, userID)

	payload := loadFixture(t, "update_full.json")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/users/"+userID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify all updates
	resp, err = client.Raw(ctx, "GET", "/api/v1/users/"+userID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	userMap := result.Data.(map[string]interface{})
	assert.Equal(t, "john.updated@example.com", userMap["email"].(string))
	assert.Equal(t, "Senior Software Engineer", userMap["jobTitle"].(string))
}

// TestUpdateUser_Disable tests disabling a user
func TestUpdateUser_Disable(t *testing.T) {
	userID := createTestUser(t, "create_internal.json")
	defer cleanupUser(t, userID)

	payload := loadFixture(t, "update_disable.json")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/users/"+userID, payload)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify disabled
	resp, err = client.Raw(ctx, "GET", "/api/v1/users/"+userID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	userMap := result.Data.(map[string]interface{})
	assert.False(t, userMap["enabled"].(bool))
}

// TestDeleteUser tests deleting a user
func TestDeleteUser(t *testing.T) {
	userID := createTestUser(t, "create_internal.json")

	// Delete user
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/users/"+userID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify deleted
	resp, err = client.Raw(ctx, "GET", "/api/v1/users/"+userID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

// TestListUsers tests listing users with pagination
func TestListUsers(t *testing.T) {
	// Create a test user in Mapexos org
	userID := createTestUser(t, "create_internal.json")
	defer cleanupUser(t, userID)

	// ROOT users don't need X-Org-Context header (unrestricted global access)
	// Client is now rootClient by default - no need to set org context

	resp, err := client.Raw(ctx, "GET", "/api/v1/users?page=1&perPage=100", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Result is now a paginated result with items array
	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	// Debug: log how many users returned and their IDs
	t.Logf("LIST returned %d users. Looking for userID: %s", len(items), userID)
	for i, u := range items {
		userMap := u.(map[string]interface{})
		var id string
		if idVal, ok := userMap["id"]; ok && idVal != nil {
			id = idVal.(string)
		} else if idVal, ok := userMap["ID"]; ok && idVal != nil {
			id = idVal.(string)
		}
		t.Logf("  User %d: ID=%s, Email=%v", i+1, id, userMap["Email"])
	}

	// Should return at least the created user (may include others like admin or from previous tests/seeds)
	assert.GreaterOrEqual(t, len(items), 1)

	// Verify the created user is in the results
	foundCreatedUser := false
	for _, u := range items {
		userMap := u.(map[string]interface{})
		// Try both "id" (lowercase) and "ID" (uppercase) for compatibility
		var userIDStr string
		if id, ok := userMap["id"]; ok && id != nil {
			userIDStr = id.(string)
		} else if id, ok := userMap["ID"]; ok && id != nil {
			userIDStr = id.(string)
		}

		if userIDStr == userID {
			foundCreatedUser = true
			break
		}
	}
	assert.True(t, foundCreatedUser, "Created user should be in the list")

	// Verify pagination metadata exists
	pagination := paginatedResult["pagination"].(map[string]interface{})
	assert.NotNil(t, pagination["totalItems"])
	assert.NotNil(t, pagination["page"])
	assert.NotNil(t, pagination["perPage"])
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

func createTestUser(t *testing.T, fixtureFile string) string {
	// Load base user data from fixture
	basePayload := loadFixture(t, fixtureFile)

	// Wrap in onboarding payload with membership
	onboardingPayload := make(map[string]interface{})

	// Copy all user fields from base payload
	for k, v := range basePayload {
		onboardingPayload[k] = v
	}

	// Add membership to Mapexos organization with viewer role
	onboardingPayload["memberships"] = []map[string]interface{}{
		{
			"orgId": constants.MapexosOrgID,
			"roles": []string{constants.SuperAdminRoleID},
			"scope": "local",
		},
	}

	// Use onboarding endpoint instead of direct user creation
	resp, err := client.Raw(ctx, "POST", "/api/v1/onboarding/users", onboardingPayload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Onboarding returns {user: {...}, memberships: [...]}
	onboardingResult := result.Data.(map[string]interface{})
	userMap := onboardingResult["user"].(map[string]interface{})
	return userMap["id"].(string)
}

func cleanupUser(t *testing.T, userID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/users/"+userID, nil)
	if err != nil {
		t.Logf("Failed to cleanup user %s: %v", userID, err)
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status during cleanup: %d", resp.StatusCode)
	}
}
