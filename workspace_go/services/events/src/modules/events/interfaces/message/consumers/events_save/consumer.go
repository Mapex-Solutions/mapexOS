package events_save

import (
	"fmt"
	"time"

	"events/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * NewConsumer creates and starts a NATS consumer for processed events from router service.
 *
 * Following Hexagonal Architecture:
 * - Consumer (Interface Layer) only receives messages and calls service
 * - Service (Application Layer) handles all business logic including EVA field mapping
 *
 * These are processed, validated events from the router service that should be
 * stored in the main events table for analytics and querying (retention: 7-365 days).
 *
 * EVA Field Mapping Flow:
 * 1. Event arrives with assetTemplateId
 * 2. Service fetches template from cache (L0/L1/Fallback)
 * 3. DynamicFields provide fieldId mapping for EVA columns
 * 4. Event stored with eva_number, eva_string, eva_bool, eva_date MAPs
 *
 * Parameters:
 *   - bus: The NATS bus instance for connecting to the message broker
 *   - eventService: The EventServicePort interface for processing events
 *
 * Configuration from environment:
 *   - service_name: Used to create unique consumer and queue group names
 *   - nats_batch_size: Number of messages to fetch per NATS batch (same for bulk insert)
 *   - nats_fetch_timeout: Timeout in seconds for fetching messages
 */
func NewConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-events-save", serviceName)
	queueGroup := fmt.Sprintf("%s-EVENTS-SAVE-GROUP", serviceName)

	natsBatchSize, _ := config.GetIntValue("nats_batch_size")
	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:EventsSave] Starting %s with retry/DLQ support", consumerName))

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
			eventService.ProcessEventStoreBatch(messages)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:EventsSave] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:EventsSave] Started successfully with retry/DLQ support")
}
