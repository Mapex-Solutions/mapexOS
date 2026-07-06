package ports

import (
	"context"
	"time"

	dtos "assets/src/modules/ota/application/dtos"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// OTAFirmwareServicePort drives the firmware upload lifecycle (presigned
// direct-to-store upload: the bytes never traverse the Asset MS).
type OTAFirmwareServicePort interface {
	// InitUpload creates a PENDING_UPLOAD artifact, returns a presigned PUT URL,
	// and schedules an abandon-check.
	InitUpload(ctx context.Context, rc *reqCtx.RequestContext, req *otaDtos.FirmwareInitRequest) (*otaDtos.FirmwareInitResponse, error)
	// CompleteUpload confirms the object is stored (size + checksum), flips the
	// artifact to READY, and purges the abandon-check. Returns success ONLY
	// after the object is confirmed in storage.
	CompleteUpload(ctx context.Context, rc *reqCtx.RequestContext, firmwareID string) error
	// HandleFirmwareAbandon is fired by the abandon-check timer: if the artifact
	// was never finalized it is marked ABANDONED and any orphan object deleted.
	HandleFirmwareAbandon(ctx context.Context, firmwareID string) error
}

// FirmwareSchedulerPort schedules (and purges) the per-firmware abandon-check
// timer. The NATS implementation lives in the messaging layer.
type FirmwareSchedulerPort interface {
	ScheduleAbandonCheck(firmwareID string, at time.Time) error
	PurgeAbandonCheck(firmwareID string) error
}

// OTAPlanServicePort drives plan CRUD + the /find list endpoints.
type OTAPlanServicePort interface {
	CreatePlan(ctx context.Context, rc *reqCtx.RequestContext, req *otaDtos.OTAPlanCreateRequest) (*otaDtos.OTAPlanResponse, error)
	ListPlans(ctx context.Context, rc *reqCtx.RequestContext, query *dtos.OTAPlanQueryDTO) (*model.PaginatedResult[otaDtos.OTAPlanResponse], error)
	GetPlan(ctx context.Context, rc *reqCtx.RequestContext, id string) (*otaDtos.OTAPlanResponse, error)
	UpdatePlan(ctx context.Context, rc *reqCtx.RequestContext, id string, req *dtos.OTAPlanUpdateRequest) (*otaDtos.OTAPlanResponse, error)
	DeletePlan(ctx context.Context, rc *reqCtx.RequestContext, id string) error
	ListExecutions(ctx context.Context, rc *reqCtx.RequestContext, planID string, query *dtos.OTAExecutionQueryDTO) (*model.PaginatedResult[otaDtos.OTAExecutionResponse], error)
	// DownloadFirmware mints a short-TTL presigned GET URL for the plan's
	// firmware artifact, after confirming the object still exists in storage.
	// Returns a 404 "firmware not found" when the object is gone.
	DownloadFirmware(ctx context.Context, rc *reqCtx.RequestContext, planID string) (*otaDtos.OTAFirmwareDownloadResponse, error)
}
