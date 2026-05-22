package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	atContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assettemplates"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// parseListNameUpdatedEvent unmarshals the inbound NATS payload. It rejects
// the message in place when the payload is malformed so the caller can exit
// early without touching the dispatcher.
func (s *AssetTemplateService) parseListNameUpdatedEvent(msg *natsModel.Message) (atContract.ListNameUpdatedEvent, bool) {
	var event atContract.ListNameUpdatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		msg.Reject(fmt.Sprintf("invalid list name updated event: %s", err))
		return event, false
	}
	return event, true
}

// logListNameUpdatedReceived emits the audit log right after parsing so we
// always see what was attempted, even when downstream dispatch fails.
func (s *AssetTemplateService) logListNameUpdatedReceived(event atContract.ListNameUpdatedEvent) {
	logger.Info(fmt.Sprintf("[SERVICE:AssetTemplate] List name updated event received - listId: %s, type: %s, newName: %s",
		event.ListId, event.ListType, event.NewName))
}

// dispatchListNameUpdate routes to the matching denormalization based on
// ListType. Returns skipped=true for non-classification list types so the
// caller can Ack without treating it as a failure.
func (s *AssetTemplateService) dispatchListNameUpdate(c ctx.Context, event atContract.ListNameUpdatedEvent) (skipped bool, err error) {
	switch event.ListType {
	case "asset_manufacturer":
		return false, s.UpdateManufacturerName(c, event.ListId, event.NewName)
	case "asset_model":
		return false, s.UpdateModelName(c, event.ListId, event.NewName)
	case "asset_category":
		return false, s.UpdateCategoryName(c, event.ListId, event.NewName)
	default:
		return true, nil
	}
}

// completeListNameUpdated finalizes the message lifecycle: Ack on success or
// skip, Nack on dispatch error, and a per-outcome log line.
func (s *AssetTemplateService) completeListNameUpdated(msg *natsModel.Message, event atContract.ListNameUpdatedEvent, skipped bool, err error) {
	if skipped {
		logger.Info(fmt.Sprintf("[SERVICE:AssetTemplate] Ignoring non-classification list type: %s", event.ListType))
		msg.Ack()
		return
	}
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:AssetTemplate] Failed to update %s name for listId: %s", event.ListType, event.ListId))
		msg.Nack(err)
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:AssetTemplate] %s name updated successfully for listId: %s", event.ListType, event.ListId))
	msg.Ack()
}
