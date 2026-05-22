package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
)

// ConnectivityPayloadFn returns the AssetCreateBuilder a
// CreateConnectivityAsset variant POSTs. Unlike PayloadFn (which only
// receives runID + templateID + a single routeGroupID), this signature
// receives the full saga.Context so the builder can resolve as many
// dependencies as the journey wires through the bag — online RG,
// offline RG, datasource id, etc.
type ConnectivityPayloadFn func(c *saga.Context) *payloads.AssetCreateBuilder

// CreateConnectivityAsset POSTs the AssetCreate built by the supplied
// closure and publishes id/uuid/password on the bag. Used by the
// connectivity-action journeys whose payloads bind multiple RG ids and
// (for HTTP) a datasource id.
//
// Writes (bag):
//   - BagKeyAssetID            string  Mongo ObjectID hex of the new asset
//   - BagKeyAssetUUID          string  device id consumed by the presence pipeline
//   - BagKeyAssetMqttPassword  string  plaintext password (only when the variant sets one)
//
// Compensate: DELETE /api/v1/assets/{id}. Idempotent.
func CreateConnectivityAsset(fn ConnectivityPayloadFn) saga.Step {
	return saga.Step{
		Name: "assets/assets.CreateConnectivityAsset",
		Do: func(c *saga.Context) error {
			builder := fn(c)
			spec := builder.Build()

			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/assets", spec)
			if err != nil {
				return fmt.Errorf("create connectivity asset: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create connectivity asset: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out assetCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create connectivity asset response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create connectivity asset: empty id in response")
			}
			c.Set(BagKeyAssetID, out.Data.ID)
			c.Set(BagKeyAssetUUID, out.Data.AssetUUID)
			if pwd := builder.MqttPassword(); pwd != "" {
				c.Set(BagKeyAssetMqttPassword, pwd)
			}
			return nil
		},
		Compensate: func(c *saga.Context) error {
			if _, ok := c.Get(BagKeyAssetDeleted); ok {
				return nil
			}
			id, ok := c.Get(BagKeyAssetID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, "/api/v1/assets/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete asset: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete asset: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
