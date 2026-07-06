package services

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"time"

	"events/src/modules/events/domain/entities"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// processOTAStatusMessage parses one OTAStatusAdvisory off the wire and maps it
// to an OTAStatusEvent for the bulk history insert. The advisory is a plain
// cross-service contract (no validation tags), so a JSON unmarshal failure is
// the only reject path.
func (s *EventService) processOTAStatusMessage(idx int, msg *natsModel.Message) messageResult[entities.OTAStatusEvent] {
	setTenantContext(msg)

	var adv otaEvents.OTAStatusAdvisory
	if err := json.Unmarshal(msg.Data, &adv); err != nil {
		return messageResult[entities.OTAStatusEvent]{
			msg: msg, action: "reject",
			rejectReason: fmt.Sprintf("invalid OTA status advisory JSON: %s", err.Error()),
		}
	}

	msg.OrgId = adv.OrgID

	created := adv.Timestamp
	if created.IsZero() {
		created = time.Now().UTC()
	}

	retentionDays, _ := s.getRetentionDays(ctx.Background(), adv.OrgID, "eventsOtaStatus")
	if retentionDays == 0 {
		retentionDays = defaultOTAStatusRetentionDays
	}

	event := &entities.OTAStatusEvent{
		Created:        created,
		OrgId:          adv.OrgID,
		AssetUuid:      adv.AssetUUID,
		PlanId:         adv.PlanID,
		OTAExecutionId: adv.OTAExecutionID,
		Status:         adv.Status,
		Progress:       adv.Progress,
		Error:          adv.Error,
		Message:        adv.Message,
		RetentionDays:  retentionDays,
	}

	return messageResult[entities.OTAStatusEvent]{
		msg: msg, action: "pending", entity: event,
	}
}
