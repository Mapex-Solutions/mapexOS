package mapexvault_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	mqttPorts "assets/src/modules/mqttcerts/application/ports"
	"assets/src/modules/mqttcerts/domain/entities"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

var _ mqttPorts.MapexVaultClientPort = (*MapexVaultClient)(nil)

// NewMapexVaultClient builds the HTTP adapter against mapexVault on
// the platform's gokit httpclient. Reads `mapex_vault_url` and
// `mapex_vault_api_key` from config. The vault API key is intentionally
// separate from this service's own `internal_api_key` so each side can
// rotate independently.
func NewMapexVaultClient() mqttPorts.MapexVaultClientPort {
	base, _ := config.GetStringValue("mapex_vault_url")
	key, _ := config.GetStringValue("mapex_vault_api_key")
	return &MapexVaultClient{
		client: httpclient.New(httpclient.Config{
			BaseURL: base,
			APIKey:  key,
			Timeout: 10 * time.Second,
		}),
	}
}

// FetchIntermediateCABundle calls mapexVault's internal endpoint and
// returns the decrypted CA in a RAM-only entity. The plaintext priv
// key is in the response — caller (the OnMount goroutine) sets it on
// the InMemoryCAStore; nothing else holds a reference.
//
// Uses client.Raw so this adapter can map specific status codes
// (401/503/5xx) onto the typed sentinels the bootstrap retry loop
// inspects. The gokit client's Get/Post collapse non-2xx into a single
// generic error, which would erase that contract.
func (m *MapexVaultClient) FetchIntermediateCABundle(ctx context.Context) (*entities.CertificateAuthorityRAM, error) {
	resp, err := m.client.Raw(ctx, http.MethodGet, EndpointIntermediateCABundle, nil)
	if err != nil {
		return nil, fmt.Errorf("transport: %w", err)
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusOK:
	case resp.StatusCode == http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case resp.StatusCode == http.StatusServiceUnavailable:
		return nil, ErrCANotReady
	case resp.StatusCode >= 500:
		return nil, fmt.Errorf("%w (status=%d)", ErrTransient, resp.StatusCode)
	default:
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status=%d body=%s", resp.StatusCode, string(body))
	}
	var wire intermediateCABundleWire
	if err := json.NewDecoder(resp.Body).Decode(&wire); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	logger.Info(fmt.Sprintf("[INFRA:MapexVaultClient] CA bundle fetched subjectCN=%s notAfter=%s",
		wire.SubjectCN, wire.NotAfter.Format(time.RFC3339)))
	return &entities.CertificateAuthorityRAM{
		CertPEM:       wire.CertPEM,
		PrivateKeyPEM: wire.PrivateKeyPEM,
		SubjectCN:     wire.SubjectCN,
		NotBefore:     wire.NotBefore,
		NotAfter:      wire.NotAfter,
		Fingerprint:   wire.Fingerprint,
	}, nil
}
