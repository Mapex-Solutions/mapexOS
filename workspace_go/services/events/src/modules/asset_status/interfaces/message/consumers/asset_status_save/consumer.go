package asset_status_save

import (
	"fmt"
	"time"

	"events/src/modules/asset_status/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer starts the asset_status save consumer on EVENTS-ASSET-STATUS.
//
// Durable name uses hyphens only — NATS rejects dots in durable names.
// Queue group load-balances work across events MS replicas. RetryPolicy +
// DLQPolicy use platform defaults.
func NewConsumer(bus *natsModel.Bus, service ports.AssetStatusServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-asset-status-save", serviceName)
	queueGroup := fmt.Sprintf("%s-EVENTS-ASSET-STATUS-GROUP", serviceName)

	natsBatchSize, _ := config.GetIntValue("nats_batch_size")
	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:AssetStatusSave] Starting %s", consumerName))

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
			service.ProcessAssetStatusBatch(messages)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:AssetStatusSave] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:AssetStatusSave] Started successfully")
}
