package services

import (
	"context"
	"fmt"
	"time"

	dsDto "http_gateway/src/modules/datasources/application/dtos"
	"http_gateway/src/modules/events/application/dtos"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// buildOTAAdvisory composes the OTAStatusAdvisory from the resolved DataSource
// and the device's report. orgId comes from dataSource (trusted post-auth); the
// device supplies assetUUID/executionId/status/progress/error in the body.
func (s *EventService) buildOTAAdvisory(dataSource *dsDto.DataSourceResponse, report *dtos.OTAStatusRequestDTO) (otaEvents.OTAStatusAdvisory, error) {
	if dataSource == nil || dataSource.OrgId == nil {
		return otaEvents.OTAStatusAdvisory{}, &customErrors.ServerCustomError{
			Code: status.NOT_FOUND, Errors: []string{"dataSource with orgId is required"},
		}
	}
	if report == nil {
		return otaEvents.OTAStatusAdvisory{}, &customErrors.ServerCustomError{
			Code: status.UNPROCESSABLE_ENTITY, Errors: []string{"status report body is required"},
		}
	}
	return otaEvents.OTAStatusAdvisory{
		OrgID:          dataSource.OrgId.Hex(),
		AssetUUID:      report.AssetUUID,
		OTAExecutionID: report.ExecutionID,
		Status:         report.Status,
		Progress:       report.Progress,
		Error:          report.Error,
		Message:        report.Message,
		Timestamp:      time.Now().UTC(),
	}, nil
}

// publishOTAAdvisory publishes the advisory on the STATIC inbound subject the
// Asset MS consumes (JetStream — the assets-owned OTA stream retains it).
func (s *EventService) publishOTAAdvisory(ctx context.Context, adv otaEvents.OTAStatusAdvisory) error {
	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{
		Ctx:     ctx,
		Subject: otaEvents.SubjectOTAStatusAdvisory,
		Data:    adv,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] OTA status publish failed: assetUUID=%s executionId=%s", adv.AssetUUID, adv.OTAExecutionID))
		return &customErrors.ServerCustomError{Code: status.INTERNAL_SERVER_ERROR, Errors: []string{err.Error()}}
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Event] OTA status published: assetUUID=%s executionId=%s status=%s", adv.AssetUUID, adv.OTAExecutionID, adv.Status))
	return nil
}

// validateOTAJobQuery checks the resolved DataSource and the assetUUID query
// before relaying the poll to the Asset MS.
func (s *EventService) validateOTAJobQuery(dataSource *dsDto.DataSourceResponse, assetUUID string) error {
	if dataSource == nil || dataSource.OrgId == nil {
		return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"dataSource with orgId is required"}}
	}
	if assetUUID == "" {
		return &customErrors.ServerCustomError{Code: status.UNPROCESSABLE_ENTITY, Errors: []string{"assetUUID query parameter is required"}}
	}
	return nil
}
