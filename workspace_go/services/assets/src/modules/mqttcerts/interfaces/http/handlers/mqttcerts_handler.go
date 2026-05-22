package handlers

import (
	dtos "assets/src/modules/mqttcerts/application/dtos"
	domConsts "assets/src/modules/mqttcerts/domain/constants"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/gofiber/fiber/v2"
)

// IssueCert — POST /api/v1/mqtt_certs
func (h *MqttCertsHandler) IssueCert(c *fiber.Ctx) error {
	var req dtos.IssueCertRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	rc, _ := c.Locals("requestContext").(*reqCtx.RequestContext)
	resp, err := h.service.IssueCert(c.UserContext(), rc, &req)
	if err != nil {
		// 409 when the asset already has a current cert and force was not set.
		// The service returns a typed error; mapping kept simple here.
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

// RevokeCert — DELETE /api/v1/mqtt_certs/:serial
func (h *MqttCertsHandler) RevokeCert(c *fiber.Ctx) error {
	serial := c.Params("serial")
	rc, _ := c.Locals("requestContext").(*reqCtx.RequestContext)
	if err := h.service.RevokeCert(c.UserContext(), rc, serial, string(domConsts.ReasonUserAction)); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ListByAsset — GET /api/v1/mqtt_certs?assetUUID=...
func (h *MqttCertsHandler) ListByAsset(c *fiber.Ctx) error {
	var q dtos.ListRevokedQuery
	if err := c.QueryParser(&q); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	rc, _ := c.Locals("requestContext").(*reqCtx.RequestContext)
	rows, err := h.service.ListRevokedByAsset(c.UserContext(), rc, q.AssetUUID)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rows)
}
