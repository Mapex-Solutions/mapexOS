package org_created

import (
	"fmt"
	"time"

	"events/src/modules/retention/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer creates and starts a NATS consumer that listens for organization created events.
// When an organization is created, it automatically creates 8 default retention policy documents.
//
// Parameters:
//   - bus: The NATS bus instance for connecting to the message broker
//   - retentionService: The RetentionServicePort interface for creating default policies
func NewConsumer(bus *natsModel.Bus, retentionService ports.RetentionServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-retention-org-created", serviceName)
	queueGroup := fmt.Sprintf("%s-RETENTION-ORG-CREATED-GROUP", serviceName)

	natsBatchSize, _ := config.GetIntValue("nats_batch_size")
	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:RetentionOrgCreated] Starting %s with retry/DLQ support", consumerName))

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
			for _, msg := range messages {
				if err := retentionService.HandleOrgCreatedEvent(msg); err != nil {
					msg.Nack(err)
					continue
				}
				msg.Ack()
			}
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:RetentionOrgCreated] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:RetentionOrgCreated] Started successfully with retry/DLQ support")
}
