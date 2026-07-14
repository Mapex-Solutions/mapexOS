package handlers

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/assets/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetLorawanByDevAddr returns a handler that resolves the LoRaWAN devices sharing a
// DevAddr into their auth projections. The LNS calls it on an uplink to hydrate the
// session for the network address the uplink carries; the DevAddr is not unique, so
// the response is a list the LNS disambiguates against the session keys.
//
// Security: Protected by API Key authentication (X-API-Key header).
//
// Parameters:
//   - service: The AssetServicePort interface for asset business operations
//
// Returns:
//   - A handler function that resolves devices by DevAddr
func GetLorawanByDevAddr(service ports.AssetServicePort) web.Handler {
	return func(c *web.Ctx) error {
		devAddr := c.Params("devAddr")
		if devAddr == "" {
			return response.BadRequest(c, []string{"devAddr is required"})
		}

		devices, err := service.GetLorawanDevicesByDevAddr(c.UserContext(), devAddr)
		if err != nil {
			return err
		}

		return response.Success(c, devices)
	}
}
