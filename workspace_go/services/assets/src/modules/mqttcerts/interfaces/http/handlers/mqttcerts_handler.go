package handlers

import (
	dtos "assets/src/modules/mqttcerts/application/dtos"
	domConsts "assets/src/modules/mqttcerts/domain/constants"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// IssueCert — POST /api/v1/mqtt_certs
func (h *MqttCertsHandler) IssueCert(c *web.Ctx) error {
	var req dtos.IssueCertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, []string{err.Error()})
	}
	rc, _ := c.Locals("requestContext").(*reqCtx.RequestContext)
	resp, err := h.service.IssueCert(c.UserContext(), rc, &req)
	if err != nil {
		// 409 when the asset already has a current cert and force was not set.
		// The service returns a typed error; mapping kept simple here.
		return response.Conflict(c, []string{err.Error()})
	}
	return response.Success(c, resp)
}

// IssueGatewayCert — POST /api/v1/gateway_certs
func (h *MqttCertsHandler) IssueGatewayCert(c *web.Ctx) error {
	var req dtos.IssueCertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, []string{err.Error()})
	}
	rc, _ := c.Locals("requestContext").(*reqCtx.RequestContext)
	resp, err := h.service.IssueGatewayCert(c.UserContext(), rc, &req)
	if err != nil {
		return response.Conflict(c, []string{err.Error()})
	}
	return response.Success(c, resp)
}

// RevokeCert — DELETE /api/v1/mqtt_certs/:serial
func (h *MqttCertsHandler) RevokeCert(c *web.Ctx) error {
	serial := c.Params("serial")
	rc, _ := c.Locals("requestContext").(*reqCtx.RequestContext)
	if err := h.service.RevokeCert(c.UserContext(), rc, serial, string(domConsts.ReasonUserAction)); err != nil {
		return response.Custom(c, web.StatusBadGateway, []string{err.Error()})
	}
	return response.Success(c, map[string]bool{"success": true})
}

// ListByAsset — GET /api/v1/mqtt_certs?assetUUID=...
func (h *MqttCertsHandler) ListByAsset(c *web.Ctx) error {
	var q dtos.ListRevokedQuery
	if err := c.QueryParser(&q); err != nil {
		return response.BadRequest(c, []string{err.Error()})
	}
	rc, _ := c.Locals("requestContext").(*reqCtx.RequestContext)
	rows, err := h.service.ListRevokedByAsset(c.UserContext(), rc, q.AssetUUID)
	if err != nil {
		return response.Custom(c, web.StatusBadGateway, []string{err.Error()})
	}
	return response.Success(c, rows)
}
