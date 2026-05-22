package events_dlq

import (
	"fmt"
	"time"

	"events/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * NewConsumer creates and starts a NATS consumer for Dead Letter Queue events.
 *
 * Following Hexagonal Architecture:
 * - Consumer (Interface Layer) only receives messages and calls service
 * - Service (Application Layer) handles all business logic
 *
 * This consumer subscribes to the "dlq.mapexos" subject on the "MAPEXOS-DLQ" stream
 * and stores failed messages in ClickHouse for debugging and analysis.
 *
 * IMPORTANT:
 * This consumer does NOT use retry/DLQ policy itself.
 * If storage fails, messages are ACKed anyway to prevent infinite loops.
 * The DLQ consumer should NEVER send to DLQ itself.
 *
 * Parameters:
 *   - bus: The NATS bus instance for connecting to the message broker
 *   - eventService: The EventServicePort interface for processing DLQ events
 *
 * Configuration from environment:
 *   - service_name: Used to create unique consumer and queue group names
 *   - nats_batch_size: Number of messages to fetch per NATS batch
 *   - nats_fetch_timeout: Timeout in seconds for fetching messages
 */
func NewConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-events-dlq", serviceName)
	queueGroup := fmt.Sprintf("%s-EVENTS-DLQ-GROUP", serviceName)

	natsBatchSize, _ := config.GetIntValue("nats_batch_size")
	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:EventsDLQ] Starting %s (NO retry/DLQ policy)", consumerName))

	// Use StartConsumer WITHOUT retry/DLQ policy to avoid infinite loops
	// This consumer processes DLQ messages, so it cannot send to DLQ itself
	_, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		BatchSize:    natsBatchSize,
		FetchTimeout: time.Duration(natsFetchTimeout) * time.Second,

		// NO RetryPolicy - if storage fails, just ACK and move on
		// NO DLQPolicy - this IS the DLQ consumer, cannot send to itself

		BatchMessageHandlerV2: func(messages []*natsModel.Message) {
			eventService.ProcessDLQEventBatch(messages)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:EventsDLQ] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:EventsDLQ] Started successfully (NO retry/DLQ policy)")
}
