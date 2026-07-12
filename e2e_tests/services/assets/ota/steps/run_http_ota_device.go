package steps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	downlink "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	assetSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/assets/assets/steps"
	dsPayloads "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/payloads"
	dsSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/http_gateway/datasources/steps"
)

type otaJobResponse struct {
	Data downlink.OTAUpdateCommand `json:"data"`
}

// RunHttpOtaDevice simulates an HTTP device: it polls the gateway for its OTA
// job (data-source-authenticated via ?ds= + the API key), downloads + verifies
// the firmware, records the executionId, and POSTs the status progression back.
// The HTTP device is never pushed — no subscribe, no presence gate — so this runs
// after CreatePlan without any ordering race.
//
// Reads (bag):
//   - dsSteps.BagKeyDataSourceID / dsSteps.BagKeyDataSourceApiKey  set by CreateDataSource
//   - assetSteps.BagKeyAssetUUID                                   set by CreateConnectivityAsset
//
// Writes (bag):
//   - BagKeyExecutionID  string
//
// Compensate: no-op.
func RunHttpOtaDevice() saga.Step {
	return saga.Step{
		Name: "assets/ota.RunHttpOtaDevice",
		Do: func(c *saga.Context) error {
			dsID := c.MustGetString(dsSteps.BagKeyDataSourceID)
			apiKey := c.MustGetString(dsSteps.BagKeyDataSourceApiKey)
			uuid := c.MustGetString(assetSteps.BagKeyAssetUUID)
			headers := map[string]string{dsPayloads.SagaApiKeyHeaderName: apiKey}

			cmd, err := pollOtaJob(c, dsID, uuid, headers)
			if err != nil {
				return err
			}
			if err := downloadAndVerifyFirmware(c.Stdctx, cmd); err != nil {
				return err
			}
			c.Set(BagKeyExecutionID, cmd.ExecutionID)

			for _, s := range otaStatusProgression {
				report := otaDtos.OTAStatusRequestDTO{
					AssetUUID:   uuid,
					ExecutionID: cmd.ExecutionID,
					Status:      s.status,
					Progress:    s.progress,
				}
				resp, err := c.Clients.Gateway.RawWithHeaders(c.Stdctx, http.MethodPost, "/api/v1/ota/status?ds="+dsID, report, headers)
				if err != nil {
					return fmt.Errorf("post ota status %q: %w", s.status, err)
				}
				if resp.StatusCode < 200 || resp.StatusCode >= 300 {
					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()
					return fmt.Errorf("post ota status %q: unexpected status %d body=%s", s.status, resp.StatusCode, string(body))
				}
				resp.Body.Close()
				time.Sleep(150 * time.Millisecond)
			}
			return nil
		},
		Compensate: func(_ *saga.Context) error { return nil },
	}
}

// pollOtaJob polls GET /api/v1/ota/jobs until the gateway returns a command
// (200) rather than "nothing actionable" (204), or the timeout elapses.
func pollOtaJob(c *saga.Context, dsID, uuid string, headers map[string]string) (downlink.OTAUpdateCommand, error) {
	deadline := time.Now().Add(60 * time.Second)
	path := "/api/v1/ota/jobs?ds=" + dsID + "&assetUUID=" + uuid
	for {
		resp, err := c.Clients.Gateway.RawWithHeaders(c.Stdctx, http.MethodGet, path, nil, headers)
		if err == nil {
			switch resp.StatusCode {
			case http.StatusNoContent:
				resp.Body.Close() // nothing actionable yet — keep polling
			case http.StatusOK:
				var out otaJobResponse
				decErr := json.NewDecoder(resp.Body).Decode(&out)
				resp.Body.Close()
				if decErr != nil {
					return downlink.OTAUpdateCommand{}, fmt.Errorf("decode ota job response: %w", decErr)
				}
				if out.Data.ExecutionID == "" {
					return downlink.OTAUpdateCommand{}, fmt.Errorf("ota job: empty command in 200 response")
				}
				return out.Data, nil
			default:
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				return downlink.OTAUpdateCommand{}, fmt.Errorf("get ota job: unexpected status %d body=%s", resp.StatusCode, string(body))
			}
		}
		if time.Now().After(deadline) {
			return downlink.OTAUpdateCommand{}, fmt.Errorf("get ota job: no command for %s within 60s", uuid)
		}
		select {
		case <-c.Stdctx.Done():
			return downlink.OTAUpdateCommand{}, fmt.Errorf("get ota job poll cancelled: %w", c.Stdctx.Err())
		case <-time.After(time.Second):
		}
	}
}
