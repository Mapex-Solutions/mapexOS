package handlers

import (
	"fmt"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// GetKEK returns the decrypted KEK for the :context path param to the calling
// service (assets MS / LNS at boot).
//
// The HTTP status code is part of the contract — callers map it to behavior:
//   - 200: KEK returned (in the response envelope's data)
//   - 503: KEK not yet seeded in Mongo (caller retries)
//   - 500: any other failure (decrypt, mongo transport, etc.)
//
// response.Custom preserves the mapped status while keeping the standard
// {status, errors, data} envelope.
func (h *KekInternalHandler) GetKEK(c *web.Ctx) error {
	keyContext := c.Params("context")
	logger.Info("[HANDLER:KekInternal] GetKEK: request received context=" + keyContext)
	resp, err := h.service.GetKEK(c.UserContext(), keyContext)
	if err != nil {
		status := mapServiceErrToStatus(err)
		logger.Warn(fmt.Sprintf("[HANDLER:KekInternal] GetKEK: err=%v status=%d", err, status))
		return response.Custom(c, status, []string{err.Error()})
	}
	logger.Info("[HANDLER:KekInternal] GetKEK: ok context=" + resp.Context)
	return response.Success(c, resp)
}
