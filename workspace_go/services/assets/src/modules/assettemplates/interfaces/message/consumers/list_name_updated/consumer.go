package list_name_updated

import (
	"fmt"

	"assets/src/modules/assettemplates/application/ports"

	atContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assettemplates"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer creates and starts a durable NATS consumer for list name updates.
// Uses StartConsumer with DLQPolicy for unified DLQ compliance.
//
// This file contains only consumer wiring. The dispatch logic (JSON
// unmarshal + ListType switch) lives inside the AssetTemplateService via
// the HandleListNameUpdated port method.
func NewConsumer(bus *natsModel.Bus, assetTemplateService ports.AssetTemplateServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-list-name-updated", serviceName)
	queueGroup := fmt.Sprintf("%s-LIST-NAME-GROUP", serviceName)

	logger.Info(fmt.Sprintf("[CONSUMER:ListNameUpdated] Starting %s with DLQ support", consumerName))

	_, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:     atContract.ListNameUpdatedStream,
		Subject:    atContract.ListNameUpdatedSubject,
		Durable:    consumerName,
		QueueGroup: queueGroup,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "assets",
			EventType:   atContract.ListNameUpdatedEventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			assetTemplateService.HandleListNameUpdated(msg)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:ListNameUpdated] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:ListNameUpdated] Started successfully with DLQ support")
}
