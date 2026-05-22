package handlers

import (
	"fmt"

	pkiDtos "mapexVault/src/modules/pki/application/dtos"

	"github.com/gofiber/fiber/v2"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// GetIntermediateCABundle returns the intermediate cert + decrypted
// priv key to the calling service (Assets MS only).
//
// Status code mapping:
//   - 200: bundle returned
//   - 503: CA not yet bootstrapped in Mongo (caller maps to retry)
//   - 500: any other failure (decrypt, mongo transport, etc.)
func (h *PkiInternalHandler) GetIntermediateCABundle(c *fiber.Ctx) error {
	logger.Info("[HANDLER:PkiInternal] GetIntermediateCABundle: request received")
	bundle, err := h.service.GetIntermediateCABundle(c.UserContext())
	if err != nil {
		status := mapServiceErrToStatus(err)
		logger.Warn(fmt.Sprintf("[HANDLER:PkiInternal] GetIntermediateCABundle: err=%v status=%d", err, status))
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	logger.Info(fmt.Sprintf("[HANDLER:PkiInternal] GetIntermediateCABundle: ok subjectCN=%s", bundle.SubjectCN))
	return c.JSON(bundle)
}

// GetCAChain returns root + intermediate concatenated as PEM. Public material.
func (h *PkiInternalHandler) GetCAChain(c *fiber.Ctx) error {
	logger.Info("[HANDLER:PkiInternal] GetCAChain: request received")
	chain, err := h.service.GetCAChain(c.UserContext())
	if err != nil {
		status := mapServiceErrToStatus(err)
		logger.Warn(fmt.Sprintf("[HANDLER:PkiInternal] GetCAChain: err=%v status=%d", err, status))
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(chain)
}

// SignServer signs a server cert per the request body using the intermediate CA.
func (h *PkiInternalHandler) SignServer(c *fiber.Ctx) error {
	var req pkiDtos.SignServerRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn(fmt.Sprintf("[HANDLER:PkiInternal] SignServer: bad body err=%v", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	logger.Info(fmt.Sprintf("[HANDLER:PkiInternal] SignServer: request received cn=%s sans=%v ttlDays=%d", req.CN, req.SANs, req.TTLDays))
	resp, err := h.service.SignServerCert(c.UserContext(), &req)
	if err != nil {
		status := mapServiceErrToStatus(err)
		logger.Warn(fmt.Sprintf("[HANDLER:PkiInternal] SignServer: err=%v status=%d", err, status))
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}
