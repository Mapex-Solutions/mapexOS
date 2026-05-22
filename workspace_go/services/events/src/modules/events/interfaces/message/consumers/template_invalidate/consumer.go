package template_invalidate

import (
	"fmt"

	"events/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * NewConsumer creates and starts the FANOUT subscription for template cache
 * invalidation. Each instance of the events service receives a copy of every
 * message and clears its TieredCache (L0+L1) for the matching key.
 *
 * Pure NewConsumer + lambda — no business logic.
 */
func NewConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	logger.Info(fmt.Sprintf("[CONSUMER:TemplateInvalidate] Starting FANOUT subscription: %s -> %s", serviceName, Subject))

	_, err := bus.SubscribeFanout(Stream, serviceName, Subject, func(data []byte) error {
		msg := &natsModel.Message{Data: data}
		eventService.HandleTemplateInvalidate(msg)
		return nil
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:TemplateInvalidate] Failed to start FANOUT subscription")
		return
	}

	logger.Info("[CONSUMER:TemplateInvalidate] FANOUT subscription started successfully")
}
