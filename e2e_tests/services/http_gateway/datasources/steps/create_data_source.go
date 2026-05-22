// Package steps holds saga steps that exercise the http_gateway
// datasources module HTTP endpoints.
package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/payloads"
)

type dataSourceCreateResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// CreateDataSource POSTs the canonical SagaHttpDataSource payload to
// the http_gateway and publishes the id + apiKey on the bag for the
// HTTP heartbeat / publish steps to consume.
//
// Writes (bag):
//   - BagKeyDataSourceID      string  Mongo ObjectID hex of the new datasource
//   - BagKeyDataSourceApiKey  string  plaintext apiKey embedded on auth.apiKey.key
//
// Compensate: DELETE /api/v1/data_sources/{id}. The id is read back
// from the bag rather than captured in a closure so Compensate stays
// idempotent and the Step value can be safely reused across runs.
func CreateDataSource() saga.Step {
	return saga.Step{
		Name: "http_gateway/datasources.CreateDataSource",
		Do: func(c *saga.Context) error {
			builder := payloads.SagaHttpDataSource(c.RunID)
			spec := builder.Build()

			resp, err := c.Clients.Gateway.Raw(c.Stdctx, http.MethodPost, "/api/v1/data_sources", spec)
			if err != nil {
				return fmt.Errorf("create data source: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create data source: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out dataSourceCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create data source response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create data source: empty id in response")
			}
			c.Set(BagKeyDataSourceID, out.Data.ID)
			c.Set(BagKeyDataSourceApiKey, builder.ApiKey())
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyDataSourceID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Gateway.Raw(c.Stdctx, http.MethodDelete, "/api/v1/data_sources/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete data source: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete data source: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
