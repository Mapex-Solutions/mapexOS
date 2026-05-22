package workflow_resume

import (
	"fmt"
	"time"

	"workflow/src/modules/runtime/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * NewConsumer creates and starts a NATS consumer for workflow resume messages.
 *
 * Following Hexagonal Architecture:
 * - Consumer (Interface Layer) only receives messages and calls service
 * - Service (Application Layer) handles all business logic and message lifecycle
 *
 * This consumer subscribes to the "workflow.resume.>" subject on the "WORKFLOW-RESUME" stream
 * and processes messages individually via MessageHandlerV2.
 *
 * Parameters:
 *   - bus: The NATS bus instance for connecting to the message broker
 *   - service: The RuntimeServicePort interface for handling resume messages
 *
 * Configuration from environment:
 *   - service_name: Used to create unique consumer and queue group names
 *   - nats_fetch_timeout: Timeout in seconds for fetching messages
 */
func NewConsumer(bus *natsModel.Bus, service ports.RuntimeServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-workflow-resume", serviceName)
	queueGroup := fmt.Sprintf("%s-WORKFLOW-RESUME-GROUP", serviceName)

	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:WorkflowResume] Starting %s with retry/DLQ support", consumerName))

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		FetchTimeout: time.Duration(natsFetchTimeout) * time.Second,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: DLQServiceType,
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			service.HandleResume(msg)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:WorkflowResume] Failed to start consumer")
		return nil
	}

	logger.Info("[CONSUMER:WorkflowResume] Started successfully with retry/DLQ support")
	return consumer
}
