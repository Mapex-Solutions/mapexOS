package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	otaPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/ota/payloads"
)

type planCreateResponse struct {
	Data struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}

// CreatePlan creates the OTA plan (source template → firmware's target template)
// over the given assets and publishes the plan id. The firmware id, source
// template id, and asset ids are read from the bag at execution time, so no
// runtime value is threaded through the constructor — sourceTemplateIDKey and
// assetIDKeys are compile-time bag-key selectors.
//
// Reads (bag):
//   - BagKeyFirmwareID       string    firmware id (set by InitFirmware)
//   - sourceTemplateIDKey    string    source template id (set by CreateTemplateWithLabel)
//   - assetIDKeys...         string    one asset id per key (set by CreateAsset*)
//
// Writes (bag):
//   - BagKeyPlanID  string  the new plan id
//
// Compensate: DELETE /api/v1/ota/plans/{id} (cancel), read back from the bag,
// 404-tolerant.
func CreatePlan(sourceTemplateIDKey string, assetIDKeys ...string) saga.Step {
	return saga.Step{
		Name: "assets/ota.CreatePlan",
		Do: func(c *saga.Context) error {
			firmwareID := c.MustGetString(BagKeyFirmwareID)
			sourceTemplateID := c.MustGetString(sourceTemplateIDKey)
			assetIDs := make([]string, 0, len(assetIDKeys))
			for _, k := range assetIDKeys {
				assetIDs = append(assetIDs, c.MustGetString(k))
			}

			spec := otaPayloads.NewOTAPlan(c.RunID, firmwareID, sourceTemplateID, assetIDs).Build()
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, "/api/v1/ota/plans", spec)
			if err != nil {
				return fmt.Errorf("create plan: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create plan: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
			var out planCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create-plan response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create plan: empty id in response")
			}
			c.Set(BagKeyPlanID, out.Data.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyPlanID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, "/api/v1/ota/plans/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("cancel plan: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("cancel plan: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
