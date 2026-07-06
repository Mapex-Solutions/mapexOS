package repositories

import (
	"context"

	"events/src/modules/events/domain/entities"
)

// OTAStatusEventRepository defines storage operations for OTA status-history
// events (the events_ota_status projection).
type OTAStatusEventRepository interface {
	SaveOTAStatusEvent(ctx context.Context, event *entities.OTAStatusEvent) error
	SaveOTAStatusEventBatch(ctx context.Context, events []*entities.OTAStatusEvent) error
}
