package workflow_state

import (
	"fmt"
	"time"

	"workflow/src/modules/archiver/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer starts the WORKFLOW-STATE batch consumer.
// Uses BatchMessageHandlerV2 for per-message ACK/NACK control within batches.
func NewConsumer(bus *natsModel.Bus, service ports.ArchiverServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-workflow-state-archiver", serviceName)
	queueGroup := fmt.Sprintf("%s-WORKFLOW-STATE-GROUP", serviceName)

	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:WorkflowState] Starting %s with batch processing", consumerName))

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		BatchSize:    BatchSize,
		FetchTimeout: time.Duration(natsFetchTimeout) * time.Second,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: DLQServiceType,
			EventType:   EventType,
		},

		BatchMessageHandlerV2: func(messages []*natsModel.Message) {
			service.ProcessStateBatch(messages)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:WorkflowState] Failed to start consumer")
		return nil
	}

	logger.Info("[CONSUMER:WorkflowState] Started successfully with batch processing")
	return consumer
}
