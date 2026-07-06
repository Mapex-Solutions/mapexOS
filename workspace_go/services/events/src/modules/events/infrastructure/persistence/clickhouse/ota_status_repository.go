package clickhouseRepo

import (
	"context"
	"fmt"

	"events/src/modules/events/domain/entities"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// SaveOTAStatusEvent stores a single OTA status transition in the
// events_ota_status ClickHouse table.
func (r *EventRepositoryClickHouse) SaveOTAStatusEvent(ctx context.Context, event *entities.OTAStatusEvent) error {
	if r.otaStatusEventTable == nil {
		return fmt.Errorf("events_ota_status table model not initialized")
	}

	if err := r.otaStatusEventTable.Insert(ctx, event); err != nil {
		logger.Error(err, "[REPO:Event] Failed to insert OTA status event")
		return fmt.Errorf("failed to insert OTA status event: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] OTA status event saved: execution=%s, status=%s", event.OTAExecutionId, event.Status))
	return nil
}

// SaveOTAStatusEventBatch stores multiple OTA status transitions in the
// events_ota_status ClickHouse table in a single bulk insert.
func (r *EventRepositoryClickHouse) SaveOTAStatusEventBatch(ctx context.Context, events []*entities.OTAStatusEvent) error {
	if len(events) == 0 {
		return nil
	}

	if r.otaStatusEventTable == nil {
		return fmt.Errorf("events_ota_status table model not initialized")
	}

	if err := r.otaStatusEventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save OTA status event batch")
		return fmt.Errorf("failed to save OTA status event batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] OTA status batch saved: %d events", len(events)))
	return nil
}
