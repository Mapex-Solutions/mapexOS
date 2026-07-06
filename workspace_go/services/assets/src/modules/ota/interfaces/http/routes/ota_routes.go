package routes

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	"assets/src/modules/ota/application/dtos"
	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/interfaces/http/handlers"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/ota"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
)

// RegisterRoutes registers the OTA operator HTTP routes (base path /api/v1/ota).
// ONE handler + ONE router for the OTA domain.
func RegisterRoutes(group web.Router, firmwareSvc ports.OTAFirmwareServicePort, planSvc ports.OTAPlanServicePort) {
	r := swagger.Wrap(group).Tag("OTA")

	/** Firmware */

	initDto := validation.NewValidation(&otaDtos.FirmwareInitRequest{}, nil, nil)
	r.Post("/firmware/init", initDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAFirmwareUpload),
		coverageMw.InjectRequestContext(),
		handlers.InitFirmwareUpload(firmwareSvc),
	).
		Summary("Init firmware upload").
		Description("Creates a firmware artifact and returns a short-TTL presigned PUT URL for direct-to-store upload.").
		Returns(&otaDtos.FirmwareInitResponse{})

	completeDto := validation.NewValidation(nil, nil, &dtos.OTAFirmwareIdDto{})
	r.Post("/firmware/:firmwareId/complete", completeDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAFirmwareUpload),
		coverageMw.InjectRequestContext(),
		handlers.CompleteFirmwareUpload(firmwareSvc),
	).
		Summary("Complete firmware upload").
		Description("Confirms the uploaded object (size + checksum) and marks the firmware READY. Returns 200 only after the object is confirmed in storage.").
		Returns(map[string]bool{})

	/** Plans */

	createDto := validation.NewValidation(&otaDtos.OTAPlanCreateRequest{}, nil, nil)
	r.Post("/plans", createDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAPlanCreate),
		coverageMw.InjectRequestContext(),
		handlers.CreatePlan(planSvc),
	).
		Summary("Create OTA plan").
		Description("Creates a scheduled OTA plan (source template → target firmware) and queues one execution per selected asset.").
		Returns(&otaDtos.OTAPlanResponse{})

	listDto := validation.NewValidation(nil, &dtos.OTAPlanQueryDTO{}, nil)
	r.Get("/plans", listDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAPlanList),
		coverageMw.InjectRequestContext(),
		handlers.ListPlans(planSvc),
	).
		Summary("List OTA plans").
		Description("Returns a paginated, filterable list of OTA plans scoped to the caller's organization.").
		Returns(&model.PaginatedResult[otaDtos.OTAPlanResponse]{})

	getDto := validation.NewValidation(nil, nil, &dtos.OTAPlanIdDto{})
	r.Get("/plans/:planId", getDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAPlanRead),
		coverageMw.InjectRequestContext(),
		handlers.GetPlan(planSvc),
	).
		Summary("Get OTA plan").
		Description("Retrieves a single OTA plan by id, including its live counters.").
		Returns(&otaDtos.OTAPlanResponse{})

	updateDto := validation.NewValidation(&dtos.OTAPlanUpdateRequest{}, nil, &dtos.OTAPlanIdDto{})
	r.Patch("/plans/:planId", updateDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAPlanUpdate),
		coverageMw.InjectRequestContext(),
		handlers.UpdatePlan(planSvc),
	).
		Summary("Update OTA plan").
		Description("Partially updates an OTA plan; only provided fields are changed.").
		Returns(&otaDtos.OTAPlanResponse{})

	deleteDto := validation.NewValidation(nil, nil, &dtos.OTAPlanIdDto{})
	r.Delete("/plans/:planId", deleteDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAPlanDelete),
		coverageMw.InjectRequestContext(),
		handlers.DeletePlan(planSvc),
	).
		Summary("Cancel OTA plan").
		Description("Cancels an OTA plan (marks CANCELED; the close routine finalizes it). The firmware artifact is retained and stays downloadable.").
		Returns(map[string]bool{})

	execDto := validation.NewValidation(nil, &dtos.OTAExecutionQueryDTO{}, &dtos.OTAPlanIdDto{})
	r.Get("/plans/:planId/executions", execDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAPlanRead),
		coverageMw.InjectRequestContext(),
		handlers.ListExecutions(planSvc),
	).
		Summary("List OTA plan executions").
		Description("Returns the paginated per-device executions of a plan (live state).").
		Returns(&model.PaginatedResult[otaDtos.OTAExecutionResponse]{})

	downloadDto := validation.NewValidation(nil, nil, &dtos.OTAPlanIdDto{})
	r.Get("/plans/:planId/firmware/download", downloadDto, swagger.Expose,
		permissionMw.RequirePermission(perms.OTAPlanRead),
		coverageMw.InjectRequestContext(),
		handlers.DownloadFirmware(planSvc),
	).
		Summary("Download OTA plan firmware").
		Description("Mints a short-TTL presigned GET URL for the plan's firmware artifact; returns 404 when the object is no longer in storage.").
		Returns(&otaDtos.OTAFirmwareDownloadResponse{})
}
