package mapexvault

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"assets/src/modules/assets/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

var _ ports.LorawanKEKClientPort = (*KEKClient)(nil)

// NewKEKClient builds the HTTP adapter against mapexVault on the platform's
// gokit httpclient. Reads `mapex_vault_url` + `mapex_vault_api_key` from config
// (the vault API key is separate from this service's own internal_api_key so
// each side rotates independently).
func NewKEKClient() ports.LorawanKEKClientPort {
	base, _ := config.GetStringValue("mapex_vault_url")
	key, _ := config.GetStringValue("mapex_vault_api_key")
	return &KEKClient{
		client: httpclient.New(httpclient.Config{
			BaseURL: base,
			APIKey:  key,
			Timeout: 10 * time.Second,
		}),
	}
}

// FetchKEK calls mapexVault's internal KEK endpoint for a context and returns
// the decrypted 64-hex KEK. Uses client.Raw so specific status codes map onto
// typed sentinels the bootstrap retry loop inspects (the gokit Get/Post collapse
// non-2xx into a single generic error, which would erase that contract).
func (c *KEKClient) FetchKEK(ctx context.Context, keyContext string) (string, error) {
	resp, err := c.client.Raw(ctx, http.MethodGet, EndpointKEKBase+keyContext, nil)
	if err != nil {
		return "", fmt.Errorf("transport: %w", err)
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusOK:
	case resp.StatusCode == http.StatusUnauthorized:
		return "", ErrUnauthorized
	case resp.StatusCode == http.StatusServiceUnavailable:
		return "", ErrKEKNotReady
	case resp.StatusCode >= 500:
		return "", fmt.Errorf("%w (status=%d)", ErrTransient, resp.StatusCode)
	default:
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status=%d body=%s", resp.StatusCode, string(body))
	}
	var wire kekWire
	if err := json.NewDecoder(resp.Body).Decode(&wire); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}
	logger.Info("[INFRA:MapexVault] KEK fetched context=" + wire.Context)
	return wire.Kek, nil
}
