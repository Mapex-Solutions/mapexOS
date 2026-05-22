package plugin_execute

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
 * NewConsumer creates and starts a NATS consumer for workflow plugin execution events.
 *
 * Following Hexagonal Architecture:
 * - Consumer (Interface Layer) only receives messages and calls service
 * - Service (Application Layer) handles all business logic and message lifecycle
 *
 * Subject: trigger.TRIGGER_FROM_WORKFLOW.execute
 *
 * The consumer processes plugin execution events from the Workflow Service:
 * 1. Receives the fully resolved action pipeline (hooks + operation)
 * 2. Executes before hook (if defined) — e.g., pre-auth login
 * 3. Resolves {{before.X}} templates in the operation
 * 4. Executes the operation (HTTP call, MQTT publish, etc.)
 * 5. Executes after hook (if defined)
 * 6. Publishes resume to WORKFLOW-RESUME with result (success/error)
 *
 * Parameters:
 *   - bus: The NATS bus instance for connecting to the message broker
 *   - eventService: The EventServicePort interface for processing plugin executions
 */
func NewConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-workflow-plugin-executor", serviceName)
	queueGroup := fmt.Sprintf("%s-WORKFLOW-PLUGIN-EXEC-GROUP", serviceName)

	natsBatchSize, err := config.GetIntValue("nats_batch_size")
	if err != nil {
		natsBatchSize = constants.DefaultBatchSize
	}

	natsFetchTimeout, err := config.GetIntValue("nats_fetch_timeout")
	if err != nil {
		natsFetchTimeout = constants.DefaultFetchTimeoutSeconds
	}

	logger.Info(fmt.Sprintf("[CONSUMER:WorkflowPluginExecute] Starting %s", consumerName))

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
			eventService.ProcessWorkflowExecutionBatch(messages)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:WorkflowPluginExecute] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:WorkflowPluginExecute] Started successfully")
	logger.Info(fmt.Sprintf("[CONSUMER:WorkflowPluginExecute] Listening on subject: %s", Subject))
}
