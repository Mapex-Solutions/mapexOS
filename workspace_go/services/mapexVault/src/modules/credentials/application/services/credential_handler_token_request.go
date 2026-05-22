package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	templatereplace "github.com/Mapex-Solutions/mapexGoKit/utils/templatereplace"
)

/**
 * Token Request Execution
 */

// executeTokenRequest builds and executes an HTTP request from a TokenRequestConfig.
// Resolves {{credential.*}} templates, serializes body per content type, extracts tokens from response.
func (s *CredentialService) executeTokenRequest(credential *entities.Credential, config *entities.TokenRequestConfig, data map[string]interface{}) (*tokenResponse, error) {
	if config == nil || config.Url == "" {
		return nil, fmt.Errorf("token request config is nil or has no URL")
	}

	// Step 1 — Resolve templates
	templateCtx := map[string]interface{}{"credential": data}

	resolvedBody := templatereplace.Resolve(config.Body, templateCtx)
	resolvedHeaders := resolveStringMap(config.Headers, templateCtx)
	resolvedQuery := resolveStringMap(config.QueryParams, templateCtx)

	// Default method and content type
	method := config.Method
	if method == "" {
		method = "POST"
	}
	contentType := config.ContentType
	if contentType == "" {
		contentType = "application/json"
	}

	// If body is nil and JSON, send full credential data (backward compat)
	if resolvedBody == nil && contentType == "application/json" {
		resolvedBody = data
	}

	// Step 2 — Build HTTP request
	var bodyReader io.Reader
	if contentType == "application/x-www-form-urlencoded" {
		form := url.Values{}
		for k, v := range flattenToStringMap(resolvedBody) {
			form.Set(k, v)
		}
		bodyReader = bytes.NewReader([]byte(form.Encode()))
	} else {
		bodyBytes, err := json.Marshal(resolvedBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, config.Url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	for k, v := range resolvedHeaders {
		req.Header.Set(k, v)
	}

	if len(resolvedQuery) > 0 {
		query := req.URL.Query()
		for k, v := range resolvedQuery {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}

	// Step 3 — Execute
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	// Step 4 — Parse response
	var respBody map[string]interface{}
	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Step 5 — Extract tokens
	accessToken := extractByPath(respBody, config.AccessTokenPath)
	refreshToken := ""
	if config.RefreshTokenPath != "" {
		refreshToken = extractByPath(respBody, config.RefreshTokenPath)
	}
	expiresAt := resolveTokenExpirationFromResponse(respBody, accessToken, config.ExpiresInPath)

	logger.Info(fmt.Sprintf("[SERVICE:Credential] Token request: method=%s url=%s status=%d", method, config.Url, resp.StatusCode))

	// Step 6 — Return
	return &tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		RawResponse:  respBody,
	}, nil
}

/**
 * Token Update
 */

// updateCredentialTokens encrypts updated tokens, updates MongoDB, publishes events and schedule.
func (s *CredentialService) updateCredentialTokens(credential *entities.Credential, data map[string]interface{}, resp *tokenResponse) error {
	if resp.AccessToken != "" {
		data["accessToken"] = resp.AccessToken
	}
	if resp.RefreshToken != "" {
		data["refreshToken"] = resp.RefreshToken
	}

	env, err := encryptData(s.deps.Encryption, data)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated tokens: %w", err)
	}

	now := time.Now()
	id := credential.ID.Hex()
	_, err = s.deps.CredentialRepo.FindByIdAndUpdate(context.Background(), &id, model.Map{
		"encryptedDEK":    env.EncryptedDEK,
		"dekNonce":        env.DEKNonce,
		"encryptedData":   env.EncryptedData,
		"dataNonce":       env.DataNonce,
		"tokenExpiresAt":  resp.ExpiresAt,
		"lastRefreshedAt": now,
		"refreshError":    "",
		"status":          string(entities.CredentialStatusActive),
		"updated":         now,
	})
	if err != nil {
		return fmt.Errorf("failed to update credential after token request: %w", err)
	}

	s.publishVaultEvent(credential.ID.Hex(), "refreshed")
	s.publishRefreshSchedule(credential.ID.Hex(), credential.Type, resp.ExpiresAt)

	return nil
}

/**
 * Private Helpers
 */

// resolveStringMap applies templatereplace.ResolveString to each value in a map[string]string.
func resolveStringMap(m map[string]string, contexts map[string]interface{}) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = templatereplace.ResolveString(v, contexts)
	}
	return result
}

// flattenToStringMap converts map[string]interface{} to map[string]string for form-urlencoded serialization.
func flattenToStringMap(m interface{}) map[string]string {
	raw, ok := m.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
