package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	payloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assettemplates/payloads"
)

// migrationCreateResponse decodes the created plan id the migration create returns.
type migrationCreateResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// migrationsPath is the base route for template-migration plans, under the asset
// templates group.
const migrationsPath = "/api/v1/asset_templates/migrations"

// CreateMigrationPlan POSTs an immediate template-migration plan moving the assets
// read from assetIDKeys off the fromKey template onto the toKey template, and
// publishes the plan id on BagKeyMigrationPlanID. All params are compile-time wiring
// (bag-key selectors), never runtime bag data threaded through the constructor.
//
// Reads (bag):
//   - fromKey       string  source template id
//   - toKey         string  target template id
//   - each assetIDKey  string  an asset id to migrate
//
// Writes (bag):
//   - BagKeyMigrationPlanID  string  Mongo ObjectID hex of the created plan
//
// Compensate: DELETE the plan. It tolerates 404 (never created / already gone) AND
// 409, because an immediate plan is already terminal by rollback time and the API
// refuses to cancel a non-editable plan — both mean "nothing left to undo".
func CreateMigrationPlan(fromKey, toKey string, assetIDKeys ...string) saga.Step {
	return createMigrationPlan(fromKey, toKey, nil, assetIDKeys)
}

// CreateMigrationPlanWithExtra is CreateMigrationPlan plus literal extra asset ids
// added to the plan (compile-time wiring). It backs the partial-failure phase, which
// adds a fabricated, non-existent asset id whose execution is expected to fail.
func CreateMigrationPlanWithExtra(fromKey, toKey string, extraAssetIDs []string, assetIDKeys ...string) saga.Step {
	return createMigrationPlan(fromKey, toKey, extraAssetIDs, assetIDKeys)
}

// createMigrationPlan is the shared implementation for both constructors.
func createMigrationPlan(fromKey, toKey string, extraAssetIDs, assetIDKeys []string) saga.Step {
	return saga.Step{
		Name: "assets/assettemplates.CreateMigrationPlan",
		Do: func(c *saga.Context) error {
			fromID := c.MustGetString(fromKey)
			toID := c.MustGetString(toKey)
			assetIDs := make([]string, 0, len(assetIDKeys)+len(extraAssetIDs))
			for _, k := range assetIDKeys {
				assetIDs = append(assetIDs, c.MustGetString(k))
			}
			assetIDs = append(assetIDs, extraAssetIDs...)

			body := payloads.NewMigrationPlan(c.RunID, fromID, toID, assetIDs).Build()
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodPost, migrationsPath, body)
			if err != nil {
				return fmt.Errorf("create migration plan: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				b, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("create migration plan: unexpected status %d body=%s", resp.StatusCode, string(b))
			}
			var out migrationCreateResponse
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return fmt.Errorf("decode create-migration-plan response: %w", err)
			}
			if out.Data.ID == "" {
				return fmt.Errorf("create migration plan: empty id in response")
			}
			c.Set(BagKeyMigrationPlanID, out.Data.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyMigrationPlanID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Assets.Raw(c.Stdctx, http.MethodDelete, migrationsPath+"/"+id.(string), nil)
			if err != nil {
				return fmt.Errorf("cancel migration plan: %w", err)
			}
			defer resp.Body.Close()
			// 404: never created / already gone. 409: the immediate plan is already
			// terminal and cannot be cancelled — nothing left to undo either way.
			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusConflict {
				return nil
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("cancel migration plan: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
