package auth_test

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
)

var (
	client *httpclient.HTTPClient
	ctx    context.Context
)

func TestMain(m *testing.M) {
	// Setup
	ctx = context.Background()
	client = httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})

	// Run tests (Auth tests don't need token in TestMain since they test login itself)
	code := m.Run()

	os.Exit(code)
}

// TestLogin_Valid tests successful login with valid credentials
func TestLogin_Valid(t *testing.T) {
	payload := loadFixture(t, "login_valid.json")

	resp, err := client.Raw(ctx, "POST", "/auth/login", payload)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	dataMap := result.Data.(map[string]interface{})
	assert.NotNil(t, dataMap["access_token"])
	assert.NotNil(t, dataMap["refresh_token"])
	assert.NotNil(t, dataMap["user"])

	// Verify user object
	userMap := dataMap["user"].(map[string]interface{})
	assert.NotNil(t, userMap["id"])
	assert.NotNil(t, userMap["email"])
}

// TestLogin_InvalidEmail tests login with invalid email format
func TestLogin_InvalidEmail(t *testing.T) {
	payload := loadFixture(t, "login_invalid_email.json")

	resp, err := client.Raw(ctx, "POST", "/auth/login", payload)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestLogin_WrongPassword tests login with wrong password
func TestLogin_WrongPassword(t *testing.T) {
	payload := loadFixture(t, "login_wrong_password.json")

	resp, err := client.Raw(ctx, "POST", "/auth/login", payload)
	require.NoError(t, err)
	// Should return 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestLogin_ShortPassword tests login with password shorter than 8 chars
func TestLogin_ShortPassword(t *testing.T) {
	payload := loadFixture(t, "login_short_password.json")

	resp, err := client.Raw(ctx, "POST", "/auth/login", payload)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestLogout tests logout functionality
func TestLogout(t *testing.T) {
	// First, login to get a token
	payload := loadFixture(t, "login_valid.json")
	resp, err := client.Raw(ctx, "POST", "/auth/login", payload)
	require.NoError(t, err)

	var loginResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResult)
	require.NoError(t, err)

	dataMap := loginResult.Data.(map[string]interface{})
	accessToken := dataMap["access_token"].(string)

	// Now logout with the token
	clientWithAuth := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	clientWithAuth.SetHeader("Authorization", "Bearer "+accessToken)

	resp, err = clientWithAuth.Raw(ctx, "POST", "/auth/logout", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestRefreshToken tests token refresh functionality
func TestRefreshToken(t *testing.T) {
	// First, login to get tokens
	payload := loadFixture(t, "login_valid.json")
	resp, err := client.Raw(ctx, "POST", "/auth/login", payload)
	require.NoError(t, err)

	var loginResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResult)
	require.NoError(t, err)

	dataMap := loginResult.Data.(map[string]interface{})
	accessToken := dataMap["access_token"].(string)
	refreshToken := dataMap["refresh_token"].(string)

	// Now refresh the token - refresh token must be passed in X-Refresh-Token header
	clientWithAuth := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	clientWithAuth.SetHeader("Authorization", "Bearer "+accessToken)

	// Use POSTWithHeaders to add the X-Refresh-Token header
	headers := map[string]string{
		"X-Refresh-Token": refreshToken,
	}

	resp, err = clientWithAuth.RawWithHeaders(ctx, "POST", "/auth/refresh", nil, headers)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var refreshResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&refreshResult)
	require.NoError(t, err)

	// Verify new tokens are returned
	refreshDataMap := refreshResult.Data.(map[string]interface{})
	assert.NotNil(t, refreshDataMap["access_token"])
	assert.NotNil(t, refreshDataMap["refresh_token"])
}

// TestGetMyCoverage tests getting user's organization coverage
func TestGetMyCoverage(t *testing.T) {
	// First, login to get a token
	payload := loadFixture(t, "login_valid.json")
	resp, err := client.Raw(ctx, "POST", "/auth/login", payload)
	require.NoError(t, err)

	var loginResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResult)
	require.NoError(t, err)

	dataMap := loginResult.Data.(map[string]interface{})
	accessToken := dataMap["access_token"].(string)

	// Now get coverage
	clientWithAuth := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	clientWithAuth.SetHeader("Authorization", "Bearer "+accessToken)

	resp, err = clientWithAuth.Raw(ctx, "GET", "/auth/users/me/coverage", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var coverageResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&coverageResult)
	require.NoError(t, err)

	// Verify coverage data structure
	assert.NotNil(t, coverageResult.Data)

	// Coverage should be an array of organization IDs or organization objects
	t.Logf("Coverage data: %+v", coverageResult.Data)
}

// TestGetMyCoverage_Unauthorized tests coverage endpoint without authentication
func TestGetMyCoverage_Unauthorized(t *testing.T) {
	resp, err := client.Raw(ctx, "GET", "/auth/users/me/coverage", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
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
