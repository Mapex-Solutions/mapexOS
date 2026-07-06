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

// CreateAssetWithLabel POSTs the AssetCreate built by the supplied payload fn and
// writes the returned id/uuid under LABEL-scoped bag keys (AssetIDKey(label) /
// AssetUUIDKey(label)) instead of the shared keys, so a journey can provision
// several assets without their bag entries or Compensate teardown colliding.
// Compensate deletes only this label's asset (idempotent, 404-tolerant).
//
// Reads (bag):
//   - templateSteps.BagKeyTemplateID   string  set by CreateTemplate*
//   - rgSteps.BagKeyRouteGroupID       string  set by CreateRouteGroup
//
// Writes (bag):
//   - AssetIDKey(label)    string  Mongo ObjectID hex of the new asset
//   - AssetUUIDKey(label)  string  device id consumed by presence/ingestion
func CreateAssetWithLabel(fn PayloadFn, label string) saga.Step {
	return saga.Step{
		Name: "assets/assets.CreateAsset[" + label + "]",
		Do: func(c *saga.Context) error {
			// Template and route group are optional: a LoRaWAN gateway asset requires
			// neither (the gateway-only presence journey omits both). Absent keys
			// resolve to "" and the payload builder decides whether to set the field.
			templateID := optBagString(c, templateSteps.BagKeyTemplateID)
			routeGroupID := optBagString(c, rgSteps.BagKeyRouteGroupID)
			spec := fn(c.RunID, templateID, routeGroupID).Build()

			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/assets", spec)
			if err != nil {
				return fmt.Errorf("create asset %q: %w", label, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create asset %q: unexpected status %d body=%s", label, resp.StatusCode, string(body))
			}
			var out assetCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create-asset %q response: %w", label, err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create asset %q: empty id in response", label)
			}
			c.Set(AssetIDKey(label), out.Data.ID)
			c.Set(AssetUUIDKey(label), out.Data.AssetUUID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(AssetIDKey(label))
			if !ok {
				return nil
			}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, "/api/v1/assets/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("delete asset %q: %w", label, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("delete asset %q: unexpected status %d", label, resp.StatusCode)
			}
			return nil
		},
	}
}

// optBagString returns the string at key, or "" when the key is absent (or not a
// string). Used for the optional template/route-group inputs.
func optBagString(c *saga.Context, key string) string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
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
