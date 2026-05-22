package trigger_execute

import (
	"fmt"
	"time"

	"triggers/src/modules/events/application/constants"
	"triggers/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * NewConsumer creates and starts a NATS consumer for trigger execution events.
 *
 * Following Hexagonal Architecture:
 * - Consumer (Interface Layer) only receives messages and calls service
 * - Service (Application Layer) handles all business logic and message lifecycle
 *
 * Subject pattern: trigger.{triggerId}.execute
 * Example: trigger.507f1f77bcf86cd799439011.execute
 *
 * The consumer processes trigger execution events by:
 * 1. Receiving the event payload from NATS
 * 2. Fetching the trigger configuration from the database/cache
 * 3. Resolving placeholders in the trigger config using event payload
 * 4. Executing the trigger (HTTP call, email send, webhook, etc.)
 *
 * Parameters:
 *   - bus: The NATS bus instance for connecting to the message broker
 *   - eventService: The EventServicePort interface for processing trigger executions
 *
 * Configuration from environment:
 *   - service_name: Used to create unique consumer and queue group names
 *   - nats_batch_size: Number of messages to process per batch (default: 10)
 *   - nats_fetch_timeout: Timeout in seconds for fetching messages (default: 30)
 */
func NewConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-trigger-executor", serviceName)
	queueGroup := fmt.Sprintf("%s-TRIGGER-EXEC-GROUP", serviceName)

	natsBatchSize, err := config.GetIntValue("nats_batch_size")
	if err != nil {
		natsBatchSize = constants.DefaultBatchSize
	}

	natsFetchTimeout, err := config.GetIntValue("nats_fetch_timeout")
	if err != nil {
		natsFetchTimeout = constants.DefaultFetchTimeoutSeconds
	}

	logger.Info(fmt.Sprintf("[CONSUMER:TriggerExecute] Starting %s with retry/DLQ support", consumerName))
	logger.Info(fmt.Sprintf("[CONSUMER:TriggerExecute] Batch size: %d, Fetch timeout: %d seconds", natsBatchSize, natsFetchTimeout))

	_, err = bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		BatchSize:    natsBatchSize,
		FetchTimeout: time.Duration(natsFetchTimeout) * time.Second,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "triggers",
			EventType:   EventType,
		},

		BatchMessageHandlerV2: func(messages []*natsModel.Message) {
			eventService.ProcessTriggerExecutionBatch(messages)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:TriggerExecute] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:TriggerExecute] Started successfully")
	logger.Info("[CONSUMER:TriggerExecute] Listening on subject: trigger.ROUTER.execute")
}
