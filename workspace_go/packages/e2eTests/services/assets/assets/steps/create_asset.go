// Package steps holds saga steps that exercise the assets/assets module
// HTTP endpoints.
package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/payloads"
	templateSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/steps"
	rgSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/router/routegroups/steps"
)

type assetCreateResponse struct {
	Data struct {
		ID        string `json:"id"`
		AssetUUID string `json:"assetUUID"`
	} `json:"data"`
}

// PayloadFn returns the AssetCreateBuilder a CreateAsset variant
// sends to the API. Callers pick the variant at journey-wire time
// (password vs cert) without touching the step's machinery.
type PayloadFn func(runID, templateID, routeGroupID string) *payloads.AssetCreateBuilder

// CreateAsset POSTs the canonical SagaMqttTemperatureSensor (password
// mode) payload and publishes id/uuid/password on the bag. Used by
// the password journey; the cert journey calls CreateAssetWith with
// the cert variant.
func CreateAsset() saga.Step {
	return CreateAssetWith(payloads.SagaMqttTemperatureSensor)
}

// CreateAssetWith POSTs the AssetCreate built by the supplied
// payload fn. Lets callers pick between password / cert variants
// without forking the step's HTTP + Compensate logic.
//
// Reads (bag):
//   - templateSteps.BagKeyTemplateID   string  set by CreateTemplate
//   - rgSteps.BagKeyRouteGroupID       string  set by CreateRouteGroup
//
// Writes (bag):
//   - BagKeyAssetID            string  Mongo ObjectID hex of the new asset
//   - BagKeyAssetUUID          string  device id consumed by the presence pipeline
//   - BagKeyAssetMqttPassword  string  plaintext password (only when the variant sets one)
//
// Compensate: DELETE /api/v1/assets/{id}. The id is read back from
// the bag rather than captured in a closure so Compensate stays
// idempotent. Skipped when an explicit DeleteAsset step already ran.
func CreateAssetWith(fn PayloadFn) saga.Step {
	return saga.Step{
		Name: "assets/assets.CreateAsset",
		Do: func(c *saga.Context) error {
			templateID := c.MustGetString(templateSteps.BagKeyTemplateID)
			routeGroupID := c.MustGetString(rgSteps.BagKeyRouteGroupID)
			builder := fn(c.RunID, templateID, routeGroupID)
			spec := builder.Build()

			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/assets", spec)
			if err != nil {
				return fmt.Errorf("create asset: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create asset: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out assetCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create-asset response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create asset: empty id in response")
			}
			c.Set(BagKeyAssetID, out.Data.ID)
			c.Set(BagKeyAssetUUID, out.Data.AssetUUID)
			if pwd := builder.MqttPassword(); pwd != "" {
				c.Set(BagKeyAssetMqttPassword, pwd)
			}
			return nil
		},
		Compensate: func(c *saga.Context) error {
			// Skip when the journey already ran DeleteAsset as an
			// explicit step (the bag carries assetDeleted=true). Avoids
			// a redundant DELETE the broker plugin would otherwise see
			// as a second invalidation event.
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
