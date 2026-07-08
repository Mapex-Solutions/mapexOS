package listsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"assets/src/modules/assettemplates/application/ports"

	v1lists "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/lists"
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

var _ ports.ListsClientPort = (*ListsClient)(nil)

// NewListsClient builds the outbound HTTP adapter to mapexIam's internal lists
// resolve endpoint, authenticated with the shared internal_api_key (injected as
// the X-API-Key header by the client wrapper).
func NewListsClient() ports.ListsClientPort {
	base, _ := config.GetStringValue("mapexiam_url")
	key, _ := config.GetStringValue("internal_api_key")
	return &ListsClient{
		client: httpclient.New(httpclient.Config{
			BaseURL: base,
			APIKey:  key,
			Timeout: 10 * time.Second,
		}),
	}
}

// Resolve POSTs the classification request to mapexIam and returns the resolved
// (found-or-created) org-scoped list id, unwrapping the standard {data} envelope.
func (c *ListsClient) Resolve(ctx context.Context, req v1lists.ListResolveRequest) (string, error) {
	resp, err := c.client.Raw(ctx, http.MethodPost, "/internal/lists/resolve", req)
	if err != nil {
		return "", fmt.Errorf("transport: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("lists resolve status=%d body=%s", resp.StatusCode, string(body))
	}

	var env struct {
		Data v1lists.ListResolveResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return "", fmt.Errorf("decode resolve response: %w", err)
	}
	return env.Data.Id, nil
}
