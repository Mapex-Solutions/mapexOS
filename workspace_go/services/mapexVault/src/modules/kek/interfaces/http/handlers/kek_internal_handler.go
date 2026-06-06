package handlers

import (
	"fmt"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"github.com/gofiber/fiber/v2"
)

// GetKEK returns the decrypted KEK for the :context path param to the calling
// service (assets MS / LNS at boot).
//
// Status code mapping:
//   - 200: KEK returned
//   - 503: KEK not yet seeded in Mongo (caller maps to retry)
//   - 500: any other failure (decrypt, mongo transport, etc.)
func (h *KekInternalHandler) GetKEK(c *fiber.Ctx) error {
	keyContext := c.Params("context")
	logger.Info("[HANDLER:KekInternal] GetKEK: request received context=" + keyContext)
	resp, err := h.service.GetKEK(c.UserContext(), keyContext)
	if err != nil {
		status := mapServiceErrToStatus(err)
		logger.Warn(fmt.Sprintf("[HANDLER:KekInternal] GetKEK: err=%v status=%d", err, status))
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	logger.Info("[HANDLER:KekInternal] GetKEK: ok context=" + resp.Context)
	return c.JSON(resp)
}
