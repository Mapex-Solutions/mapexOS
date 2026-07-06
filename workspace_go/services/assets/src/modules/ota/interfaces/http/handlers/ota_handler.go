package handlers

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/ota/application/dtos"
	"assets/src/modules/ota/application/ports"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// ── Firmware ──

// InitFirmwareUpload creates a firmware artifact and returns a presigned PUT URL.
func InitFirmwareUpload(service ports.OTAFirmwareServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		body, _ := requestValidation.GetDTO[*otaDtos.FirmwareInitRequest](c, "bodyDTO")
		ret, err := service.InitUpload(ctx, rc, body)
		if err != nil {
			return err
		}
		return response.Created(c, ret)
	}
}

// CompleteFirmwareUpload confirms the stored object and marks the firmware READY.
func CompleteFirmwareUpload(service ports.OTAFirmwareServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		params, _ := requestValidation.GetDTO[*dtos.OTAFirmwareIdDto](c, "paramsDTO")
		if err := service.CompleteUpload(ctx, rc, params.FirmwareId); err != nil {
			return err
		}
		return response.Success(c, map[string]bool{"success": true})
	}
}

// ── Plans ──

// CreatePlan creates a scheduled OTA plan.
func CreatePlan(service ports.OTAPlanServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		body, _ := requestValidation.GetDTO[*otaDtos.OTAPlanCreateRequest](c, "bodyDTO")
		ret, err := service.CreatePlan(ctx, rc, body)
		if err != nil {
			return err
		}
		return response.Created(c, ret)
	}
}

// ListPlans returns a paginated, filtered list of OTA plans.
func ListPlans(service ports.OTAPlanServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		query, _ := requestValidation.GetDTO[*dtos.OTAPlanQueryDTO](c, "queryDTO")
		ret, err := service.ListPlans(ctx, rc, query)
		if err != nil {
			return err
		}
		return response.Success(c, ret)
	}
}

// GetPlan returns a single OTA plan by id.
func GetPlan(service ports.OTAPlanServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		params, _ := requestValidation.GetDTO[*dtos.OTAPlanIdDto](c, "paramsDTO")
		ret, err := service.GetPlan(ctx, rc, params.PlanId)
		if err != nil {
			return err
		}
		return response.Success(c, ret)
	}
}

// UpdatePlan partially updates an OTA plan.
func UpdatePlan(service ports.OTAPlanServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		params, _ := requestValidation.GetDTO[*dtos.OTAPlanIdDto](c, "paramsDTO")
		body, _ := requestValidation.GetDTO[*dtos.OTAPlanUpdateRequest](c, "bodyDTO")
		ret, err := service.UpdatePlan(ctx, rc, params.PlanId, body)
		if err != nil {
			return err
		}
		return response.Success(c, ret)
	}
}

// DeletePlan cancels/deletes an OTA plan.
func DeletePlan(service ports.OTAPlanServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		params, _ := requestValidation.GetDTO[*dtos.OTAPlanIdDto](c, "paramsDTO")
		if err := service.DeletePlan(ctx, rc, params.PlanId); err != nil {
			return err
		}
		return response.Success(c, map[string]bool{"success": true})
	}
}

// DownloadFirmware returns a presigned GET URL for the plan's firmware artifact.
func DownloadFirmware(service ports.OTAPlanServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		params, _ := requestValidation.GetDTO[*dtos.OTAPlanIdDto](c, "paramsDTO")
		ret, err := service.DownloadFirmware(ctx, rc, params.PlanId)
		if err != nil {
			return err
		}
		return response.Success(c, ret)
	}
}

// ListExecutions returns the per-device executions of a plan.
func ListExecutions(service ports.OTAPlanServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		rc, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found", nil)
		}
		params, _ := requestValidation.GetDTO[*dtos.OTAPlanIdDto](c, "paramsDTO")
		query, _ := requestValidation.GetDTO[*dtos.OTAExecutionQueryDTO](c, "queryDTO")
		ret, err := service.ListExecutions(ctx, rc, params.PlanId, query)
		if err != nil {
			return err
		}
		return response.Success(c, ret)
	}
}

// ── Internal (MS-to-MS) ──

// GetPendingJob serves the HTTP device poll relayed by the HTTP gateway: the
// asset's pending OTA command with a fresh download URL, or data=null when the
// device has nothing actionable (the gateway translates null into 204).
func GetPendingJob(jobs ports.DeviceJobPort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()
		assetUUID := c.Query("assetUUID")
		if assetUUID == "" {
			return response.BadRequest(c, []string{"assetUUID query parameter is required"})
		}
		job, err := jobs.PendingJob(ctx, assetUUID)
		if err != nil {
			return err
		}
		return response.Success(c, job)
	}
}
