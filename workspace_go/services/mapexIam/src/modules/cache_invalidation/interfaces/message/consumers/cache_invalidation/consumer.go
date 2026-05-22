package cache_invalidation

import (
	"mapexIam/src/modules/cache_invalidation/application/ports"

	contractsCacheInvalidation "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/cache_invalidation"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer wires the NATS cache-invalidation consumer and delegates every message
// to the application-layer service. No business logic lives here — this function only
// binds transport (NATS stream/subject/durable/queue-group/DLQ) to the service port.
func NewConsumer(bus *natsModel.Bus, service ports.CacheInvalidationServicePort) *natsModel.Consumer {
	serviceName := "mapexIam"

	logger.Info("[CONSUMER:CacheInvalidation] Starting centralized cache invalidation consumer")

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:     contractsCacheInvalidation.Stream,
		Subject:    contractsCacheInvalidation.Subject,
		Durable:    "mapexos-cache-invalidation-consumer",
		QueueGroup: serviceName + "-CACHE-INVALIDATION-GROUP",

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "mapex-iam",
			EventType:   contractsCacheInvalidation.EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			service.HandleEvent(msg)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:CacheInvalidation] Failed to start cache invalidation consumer")
		return nil
	}

	logger.Info("[CONSUMER:CacheInvalidation] Started successfully with DLQ support")
	return consumer
}
