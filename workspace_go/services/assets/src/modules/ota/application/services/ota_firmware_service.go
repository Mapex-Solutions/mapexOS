package services

import (
	"context"
	"time"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/entities"

	otaDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/dtos"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	mapper "github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// NewFirmwareService returns a firmware service over the given dependencies.
func NewFirmwareService(deps FirmwareServiceDeps) ports.OTAFirmwareServicePort {
	return &FirmwareService{deps: deps}
}

// InitUpload creates the PENDING_UPLOAD artifact, mints a presigned PUT URL, and
// schedules the abandon-check. The firmwareId is a fresh ObjectId, so the object
// key is always new (immutable artifacts, no overwrite).
func (s *FirmwareService) InitUpload(ctx context.Context, rc *reqCtx.RequestContext, req *otaDtos.FirmwareInitRequest) (*otaDtos.FirmwareInitResponse, error) {
	orgID, err := s.resolveOrg(rc)
	if err != nil {
		return nil, err
	}

	// Map request → entity via the shared mapper (Filename, Size, and the
	// string→ObjectId conversion of TargetTemplateID); then set the computed
	// fields the mapper cannot derive.
	fw, err := mapper.DtoToEntityWithOptions[otaDtos.FirmwareInitRequest, entities.Firmware](
		req, mapper.MapperOptions{StringToObjectId: true},
	)
	if err != nil || fw.TargetTemplateID.IsZero() {
		return nil, badRequest("invalid targetTemplateId")
	}

	id := model.NewObjectID()
	now := time.Now()
	fw.ID = id
	fw.OrgID = orgID
	fw.ObjectKey = orgID.Hex() + "/" + fw.TargetTemplateID.Hex() + "/" + id.Hex() + ".bin"
	fw.Checksum = req.SHA256
	fw.ChecksumAlgorithm = checksumAlgorithmSHA256
	fw.Status = entities.FirmwarePendingUpload
	fw.Created = now
	fw.Updated = now

	if _, err := s.deps.FirmwareRepo.Create(ctx, fw); err != nil {
		return nil, err
	}

	uploadURL, err := s.deps.Store.PresignPut(ctx, fw.ObjectKey, req.SHA256, s.deps.PresignTTL)
	if err != nil {
		return nil, err
	}

	// Best-effort: schedule the abandon-check (purged on finalize).
	_ = s.deps.Scheduler.ScheduleAbandonCheck(id.Hex(), now.Add(s.deps.AbandonTTL))

	return &otaDtos.FirmwareInitResponse{FirmwareID: id.Hex(), UploadURL: uploadURL}, nil
}

// CompleteUpload confirms the object is in storage (size + checksum) and flips
// the artifact to READY. It returns success ONLY after the object is confirmed.
func (s *FirmwareService) CompleteUpload(ctx context.Context, rc *reqCtx.RequestContext, firmwareID string) error {
	fw, err := s.deps.FirmwareRepo.FindById(ctx, &firmwareID)
	if err != nil {
		return err
	}
	if fw == nil {
		return notFound("Firmware not found")
	}
	if fw.Status != entities.FirmwarePendingUpload {
		return badRequest("firmware is not pending upload")
	}

	size, checksum, err := s.deps.Store.Stat(ctx, fw.ObjectKey)
	if err != nil {
		return unprocessable("firmware object not found in storage")
	}
	if size != fw.Size {
		return unprocessable("uploaded size does not match the declared size")
	}
	if checksum != "" && fw.Checksum != "" && checksum != fw.Checksum {
		return unprocessable("uploaded checksum does not match the declared checksum")
	}

	payload := map[string]any{"status": string(entities.FirmwareReady), "updated": time.Now()}
	if _, err := s.deps.FirmwareRepo.FindByIdAndUpdate(ctx, &firmwareID, payload); err != nil {
		return err
	}
	_ = s.deps.Scheduler.PurgeAbandonCheck(firmwareID)
	return nil
}

// HandleFirmwareAbandon is fired by the abandon-check timer. If the artifact was
// never finalized, it deletes any orphan object and marks it ABANDONED;
// otherwise it is a no-op (the schedule was superseded by finalize).
func (s *FirmwareService) HandleFirmwareAbandon(ctx context.Context, firmwareID string) error {
	fw, err := s.deps.FirmwareRepo.FindById(ctx, &firmwareID)
	if err != nil {
		return err
	}
	if fw == nil || fw.Status != entities.FirmwarePendingUpload {
		return nil
	}
	_ = s.deps.Store.Delete(ctx, fw.ObjectKey) // best-effort (object may not exist)
	payload := map[string]any{"status": string(entities.FirmwareAbandoned), "updated": time.Now()}
	_, err = s.deps.FirmwareRepo.FindByIdAndUpdate(ctx, &firmwareID, payload)
	return err
}

func (s *FirmwareService) resolveOrg(rc *reqCtx.RequestContext) (model.ObjectId, error) {
	return resolveOrgID(rc)
}

func badRequest(msg string) error {
	return &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{msg}}
}

func notFound(msg string) error {
	return &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{msg}}
}

func unprocessable(msg string) error {
	return &customErrors.ServerCustomError{Code: httpStatus.UNPROCESSABLE_ENTITY, Errors: []string{msg}}
}

// Compile-time check that FirmwareService implements the port.
var _ ports.OTAFirmwareServicePort = (*FirmwareService)(nil)
