package assets

import (
	"context"
	"fmt"
	"net/url"

	"http_gateway/src/modules/events/application/ports"

	downlink "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewOTAJobsAdapter creates the adapter over a configured HTTP client.
func NewOTAJobsAdapter(client *httpclient.HTTPClient) ports.OTAJobsPort {
	return &OTAJobsAdapter{client: client}
}

// Compile-time check that the adapter implements the port.
var _ ports.OTAJobsPort = (*OTAJobsAdapter)(nil)

// FetchPendingJob calls the Asset MS internal API for the asset's pending OTA
// command. Returns nil (no error) when the device has nothing actionable.
func (a *OTAJobsAdapter) FetchPendingJob(ctx context.Context, assetUUID string) (*downlink.OTAUpdateCommand, error) {
	endpoint := "/internal/ota/jobs?assetUUID=" + url.QueryEscape(assetUUID)

	var apiResponse assetsAPIResponse
	if err := a.client.Get(ctx, endpoint, &apiResponse); err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:OTAJobs] pending-job fetch failed: assetUUID=%s", assetUUID))
		return nil, err
	}
	return apiResponse.Data, nil
}
