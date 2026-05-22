package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
)

const (
	tokenFileName      = ".test-token"
	rootTokenFileName  = ".test-token-root"
	adminTokenFileName = ".test-token-admin"
)

// LoginRequest represents the login payload
type LoginRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	KeepConnected bool   `json:"keepConnected"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Status int      `json:"status"`
	Errors []string `json:"errors"`
	Data   struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	} `json:"data"`
}

// DoLogin performs login with email/password and returns the JWT token
func DoLogin(email, password string) (string, error) {
	authURL := getAuthURL()

	payload := LoginRequest{
		Email:         email,
		Password:      password,
		KeepConnected: true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login payload: %w", err)
	}

	req, err := http.NewRequest("POST", authURL+"/auth/login", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	if loginResp.Status != 200 || loginResp.Data.AccessToken == "" {
		return "", fmt.Errorf("login failed: status=%d, token empty=%v", loginResp.Status, loginResp.Data.AccessToken == "")
	}

	return loginResp.Data.AccessToken, nil
}

// GetRootToken gets token for ROOT user (root@mapex.global)
// ROOT user has mapex.* permission and can query without X-Org-Context header
// Use this for CRUD tests that should PASS
func GetRootToken() (string, error) {
	return getOrRefreshTokenFor(constants.RootUserEmail, constants.RootUserPassword, rootTokenFileName)
}

// GetAdminToken gets token for ADMIN user (admin@mapex.global)
// ADMIN user has admin_vendor.* permission and REQUIRES X-Org-Context header
// Use this for middleware/permission tests (both PASS and DENY scenarios)
func GetAdminToken() (string, error) {
	return getOrRefreshTokenFor(constants.AdminUserEmail, constants.AdminUserPassword, adminTokenFileName)
}

// GetOrRefreshToken gets token from file or performs fresh login
// Default behavior: uses ADMIN user for backward compatibility
func GetOrRefreshToken() (string, error) {
	return GetAdminToken()
}

// getOrRefreshTokenFor is the generic implementation
func getOrRefreshTokenFor(email, password, tokenFile string) (string, error) {
	// Try to read existing token
	token, err := readTokenFromFileByName(tokenFile)
	if err == nil && token != "" {
		// Token exists, verify if it's still valid
		if isTokenValid(token) {
			return token, nil
		}
	}

	// Token doesn't exist or is invalid, do fresh login
	token, err = DoLogin(email, password)
	if err != nil {
		return "", err
	}

	// Save token to file
	if err := saveTokenToFileByName(token, tokenFile); err != nil {
		// Log warning but don't fail - we have the token
		fmt.Printf("Warning: failed to save token to file: %v\n", err)
	}

	return token, nil
}

// saveTokenToFile saves the token to a file (backward compatibility)
func saveTokenToFile(token string) error {
	return saveTokenToFileByName(token, tokenFileName)
}

// saveTokenToFileByName saves the token to a specific file
func saveTokenToFileByName(token, fileName string) error {
	tokenPath := getTokenFilePathByName(fileName)

	// Create directory if it doesn't exist
	dir := filepath.Dir(tokenPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	// Write token to file
	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// readTokenFromFile reads the token from file (backward compatibility)
func readTokenFromFile() (string, error) {
	return readTokenFromFileByName(tokenFileName)
}

// readTokenFromFileByName reads the token from a specific file
func readTokenFromFileByName(fileName string) (string, error) {
	tokenPath := getTokenFilePathByName(fileName)

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// isTokenValid checks if token is still valid by making a test request
func isTokenValid(token string) bool {
	// Make a simple request to verify token works
	req, err := http.NewRequest("GET", constants.MapexosURL+"/api/v1/organizations", nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Token is valid if we get 200 or any success code
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// getTokenFilePath returns the path to the token file (backward compatibility)
func getTokenFilePath() string {
	return getTokenFilePathByName(tokenFileName)
}

// getTokenFilePathByName returns the path to a specific token file
func getTokenFilePathByName(fileName string) string {
	// Use temp directory
	return filepath.Join(os.TempDir(), "mapexos-e2e-tests", fileName)
}

// getAuthURL returns the auth service URL (mapexos service)
func getAuthURL() string {
	// Auth endpoints are in mapexos service
	return constants.MapexosURL
}

// GenerateAdminToken is deprecated - use GetAdminToken() or GetRootToken() instead
// Keeping for backward compatibility
func GenerateAdminToken() (string, error) {
	return GetAdminToken()
}

// GetAdminUserID returns the fixed ADMIN user ID from constants
// No need to perform login - the ID is deterministic from seed script
func GetAdminUserID() (string, error) {
	return constants.AdminUserID, nil
}

// GetRootUserID returns the fixed ROOT user ID from constants
// No need to perform login - the ID is deterministic from seed script
func GetRootUserID() (string, error) {
	return constants.RootUserID, nil
}
