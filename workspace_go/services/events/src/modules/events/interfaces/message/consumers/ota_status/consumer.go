package ota_status

import (
	"fmt"
	"time"

	"events/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer starts the NATS consumer for OTA status-history advisories.
//
// Following Hexagonal Architecture, the consumer only wires the JetStream
// subscription and forwards each batch to the application service, which owns
// parsing, mapping, the bulk ClickHouse insert, and the Ack/Nack/Reject
// lifecycle. It subscribes to SubjectOTAStatus on StreamOTAStatus and uses
// BatchMessageHandlerV2 for efficient bulk inserts into events_ota_status.
func NewConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-ota-status", serviceName)
	queueGroup := fmt.Sprintf("%s-OTA-STATUS-GROUP", serviceName)

	natsBatchSize, _ := config.GetIntValue("nats_batch_size")
	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:OTAStatus] Starting %s with retry/DLQ support", consumerName))

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
			eventService.ProcessOTAStatusBatch(messages)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:OTAStatus] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:OTAStatus] Started successfully with retry/DLQ support")
}
