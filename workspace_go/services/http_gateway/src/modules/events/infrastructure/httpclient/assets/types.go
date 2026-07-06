package assets

import (
	downlink "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// OTAJobsAdapter implements ports.OTAJobsPort over the Asset MS internal API.
type OTAJobsAdapter struct {
	client *httpclient.HTTPClient
}

// assetsAPIResponse is the standard response envelope from the Asset MS
// internal API. Data is null when the device has no pending job.
type assetsAPIResponse struct {
	Status int                        `json:"status"`
	Errors []string                   `json:"errors"`
	Data   *downlink.OTAUpdateCommand `json:"data"`
}
