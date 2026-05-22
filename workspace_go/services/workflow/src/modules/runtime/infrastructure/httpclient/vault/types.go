package vault

import "net/http"

// VaultClient implements VaultPort by calling the vault MS internal API.
type VaultClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}
