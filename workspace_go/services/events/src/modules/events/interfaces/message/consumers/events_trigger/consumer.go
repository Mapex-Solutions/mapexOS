package events_trigger

import (
	"fmt"
	"time"

	"events/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * NewConsumer creates and starts a NATS consumer for trigger execution history events.
 *
 * Following Hexagonal Architecture:
 * - Consumer (Interface Layer) only receives messages and calls service
 * - Service (Application Layer) handles all business logic and message lifecycle
 *
 * This consumer subscribes to the "events.trigger" subject on the "EVENTS-TRIGGER" stream
 * and processes messages using BatchMessageHandlerV2 for efficient bulk ClickHouse inserts.
 *
 * Parameters:
 *   - bus: The NATS bus instance for connecting to the message broker
 *   - eventService: The EventServicePort interface for processing trigger events
 *
 * Configuration from environment:
 *   - service_name: Used to create unique consumer and queue group names
 *   - nats_batch_size: Number of messages to fetch per NATS batch
 *   - nats_fetch_timeout: Timeout in seconds for fetching messages
 */
func NewConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-events-trigger", serviceName)
	queueGroup := fmt.Sprintf("%s-EVENTS-TRIGGER-GROUP", serviceName)

	natsBatchSize, _ := config.GetIntValue("nats_batch_size")
	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:EventsTrigger] Starting %s with retry/DLQ support", consumerName))

	_, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		BatchSize:    natsBatchSize,
		FetchTimeout: time.Duration(natsFetchTimeout) * time.Second,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "events",
			EventType:   EventType,
		},

		BatchMessageHandlerV2: func(messages []*natsModel.Message) {
			eventService.ProcessTriggerEventBatch(messages)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:EventsTrigger] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:EventsTrigger] Started successfully with retry/DLQ support")
}
