package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"workflow/src/modules/runtime/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time check
var _ ports.VaultPort = (*VaultClient)(nil)

// New creates a new VaultClient.
func New(baseURL, apiKey string) ports.VaultPort {
	return &VaultClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

// DecryptCredential calls the vault internal API to get plaintext credential data.
func (c *VaultClient) DecryptCredential(_ context.Context, id string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/internal/credentials/%s/decrypt", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("[INFRA:VaultClient] Failed to create request: %w", err)
	}
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[INFRA:VaultClient] Request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[INFRA:VaultClient] Vault returned %d: %s", resp.StatusCode, string(body))
	}

	var envelope struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Data != nil {
		return envelope.Data, nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("[INFRA:VaultClient] Failed to parse response: %w", err)
	}

	logger.Debug(fmt.Sprintf("[INFRA:VaultClient] Fetched credential %s from vault", id))
	return data, nil
}
