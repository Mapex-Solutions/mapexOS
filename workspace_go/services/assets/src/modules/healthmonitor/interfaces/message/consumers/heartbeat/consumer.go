package heartbeat

import (
	"fmt"
	"time"

	"assets/src/modules/healthmonitor/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer creates and starts a NATS consumer for asset heartbeat events.
func NewConsumer(bus *natsModel.Bus, service ports.HealthMonitorServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")
	consumerName := fmt.Sprintf("%s-health-heartbeat", serviceName)

	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")
	if natsFetchTimeout <= 0 {
		natsFetchTimeout = 1
	}

	logger.Info(fmt.Sprintf("[CONSUMER:HealthHeartbeat] Starting %s", consumerName))

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		FetchTimeout: time.Duration(natsFetchTimeout) * time.Second,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "assets",
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			service.HandleHeartbeat(msg)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:HealthHeartbeat] Failed to start consumer")
		return nil
	}

	logger.Info("[CONSUMER:HealthHeartbeat] Started successfully")
	return consumer
}
