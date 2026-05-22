package services

import (
	ctx "context"
	"fmt"
	"mapexIam/src/modules/lists/domain/entities"

	listContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assettemplates"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// publishListNameUpdatedEvent publishes a NATS event when a list name changes.
// Subject + payload shape are owned by the cross-service contract package
// so consumers (assets/assettemplates) can decode without local duplication.
//
// Event structure: { listId, listType, newName }.
func (s *ListService) publishListNameUpdatedEvent(c ctx.Context, list *entities.List) {
	subject := listContract.ListNameUpdatedSubject
	event := map[string]interface{}{
		"listId":   list.ID.Hex(),
		"listType": list.Type,
		"newName":  list.Name,
	}
	publishConfig := natsModel.PublishConfig{
		Ctx:     c,
		Subject: subject,
		Data:    event,
	}
	if err := s.deps.NatsBus.Publish(publishConfig); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:List] Failed to publish list name update event to subject: %s", subject))
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:List] Published list name update event to subject: %s for listId: %s", subject, list.ID.Hex()))
}
