package handlers

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"http_gateway/src/bootstrap"
	dsDto "http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/events/application/dtos"
	"http_gateway/src/modules/events/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// ProcessEvent returns a Fiber handler that processes incoming events
// (webhook receiver). Following Hexagonal Architecture, the handler accepts
// the service port interface and delegates business logic to the service.
//
// Expected inputs:
//   - Event data in the request body (parsed as map[string]any).
//   - DataSource resolved by CustomAuthMiddleware and stored under
//     "dataSource" in the Fiber locals.
func ProcessEvent(service ports.EventServicePort, m *bootstrap.HttpGatewayMetrics) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		m.EventPayloadSize.Observe(float64(len(c.Body())))

		var event map[string]any
		if err := c.BodyParser(&event); err != nil {
			return err
		}

		dataSource, errGetDs := requestValidation.GetDTO[*dsDto.DataSourceResponse](c, "dataSource")
		if errGetDs != nil {
			return errGetDs
		}

		retData, err := service.ProcessEvent(ctx, event, dataSource)
		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// ProcessHeartbeat returns a Fiber handler for POST /api/v1/heartbeat?ds={dataSourceId}.
//
// Body shape: { "assetUUID": "<v>" } (validated upstream by ValidationMiddleware
// using HeartbeatRequestDTO). The CustomAuthMiddleware (mounted on the route)
// resolves the DataSource and stores it in c.Locals("dataSource"); this
// handler retrieves both the body and the DataSource from c.Locals, then
// delegates publishing to the EventService. orgId and pathKey come from the
// resolved DataSource — never from the body — so a compromised body cannot
// spoof a different tenant. The metrics arg is kept for signature symmetry
// with ProcessEvent; heartbeat metrics are emitted from the service layer.
func ProcessHeartbeat(service ports.EventServicePort, _ *bootstrap.HttpGatewayMetrics) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		body, errBody := requestValidation.GetDTO[*dtos.HeartbeatRequestDTO](c, "bodyDTO")
		if errBody != nil {
			return errBody
		}

		dataSource, errGetDs := requestValidation.GetDTO[*dsDto.DataSourceResponse](c, "dataSource")
		if errGetDs != nil {
			return errGetDs
		}

		if err := service.ProcessHeartbeat(ctx, dataSource, body.AssetUUID); err != nil {
			return err
		}
		return response.Success(c, map[string]bool{"success": true})
	}
}

// ProcessOTAStatus returns a Fiber handler for POST /api/v1/ota/status?ds={dataSourceId}.
//
// Body shape: { assetUUID, executionId, status, progress, error?, message? }
// (validated upstream via OTAStatusRequestDTO). The CustomAuthMiddleware
// resolves the DataSource; orgId comes from it — never from the body — so a
// compromised body cannot spoof a different tenant. The service normalizes the
// report into the shared OTAStatusAdvisory and publishes it for the Asset MS.
func ProcessOTAStatus(service ports.EventServicePort, _ *bootstrap.HttpGatewayMetrics) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		body, errBody := requestValidation.GetDTO[*dtos.OTAStatusRequestDTO](c, "bodyDTO")
		if errBody != nil {
			return errBody
		}

		dataSource, errGetDs := requestValidation.GetDTO[*dsDto.DataSourceResponse](c, "dataSource")
		if errGetDs != nil {
			return errGetDs
		}

		if err := service.ProcessOTAStatus(ctx, dataSource, body); err != nil {
			return err
		}
		return response.Success(c, map[string]bool{"success": true})
	}
}

// GetOTAJob returns a Fiber handler for GET /api/v1/ota/jobs?ds={dataSourceId}&assetUUID={assetUUID}.
//
// The device polls for its pending OTA job; the gateway relays the poll to the
// Asset MS internal API (which mints a fresh presigned download URL). Responds
// 204 when the device has nothing actionable.
func GetOTAJob(service ports.EventServicePort, _ *bootstrap.HttpGatewayMetrics) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		dataSource, errGetDs := requestValidation.GetDTO[*dsDto.DataSourceResponse](c, "dataSource")
		if errGetDs != nil {
			return errGetDs
		}

		job, err := service.GetOTAJob(ctx, dataSource, c.Query("assetUUID"))
		if err != nil {
			return err
		}
		if job == nil {
			return response.Custom(c, status.NO_CONTENT, nil)
		}
		return response.Success(c, job)
	}
}
