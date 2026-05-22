package handlers

import (
	"github.com/gofiber/fiber/v2"

	"http_gateway/src/bootstrap"
	dsDto "http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/events/application/dtos"
	"http_gateway/src/modules/events/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// ProcessEvent returns a Fiber handler that processes incoming events
// (webhook receiver). Following Hexagonal Architecture, the handler accepts
// the service port interface and delegates business logic to the service.
//
// Expected inputs:
//   - Event data in the request body (parsed as map[string]any).
//   - DataSource resolved by CustomAuthMiddleware and stored under
//     "dataSource" in the Fiber locals.
func ProcessEvent(service ports.EventServicePort, m *bootstrap.HttpGatewayMetrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
func ProcessHeartbeat(service ports.EventServicePort, _ *bootstrap.HttpGatewayMetrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		return c.SendStatus(fiber.StatusNoContent)
	}
}
