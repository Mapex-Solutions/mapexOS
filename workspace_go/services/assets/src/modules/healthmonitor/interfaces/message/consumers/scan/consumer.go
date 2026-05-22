package scan

import (
	"context"
	"fmt"
	"time"

	"assets/src/modules/healthmonitor/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer creates and starts a NATS consumer for health monitor scan events.
// Uses QueueGroup for org-level distribution across pods.
func NewConsumer(bus *natsModel.Bus, service ports.HealthMonitorServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")
	consumerName := fmt.Sprintf("%s-health-scanner", serviceName)
	queueGroup := fmt.Sprintf("%s-HEALTH-SCAN-GROUP", serviceName)

	logger.Info(fmt.Sprintf("[CONSUMER:HealthScanner] Starting %s", consumerName))

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		FetchTimeout: 5 * time.Second,

		// DuplicateWindow MUST be < scan interval (60s) to avoid re-schedule dedup.
		// Without this, the wrapper forces stream Duplicates=15m (subscribe.go:33),
		// silently discarding every re-schedule within 15 minutes of the previous one.
		DuplicateWindow: 10 * time.Second,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "assets",
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			logger.Info(fmt.Sprintf("[CONSUMER:HealthScanner] Scan message received — triggering RunScan: subject=%s ts=%s",
				msg.Subject, msg.Timestamp.Format("15:04:05")))
			service.RunScan(context.Background())
			msg.Ack()
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:HealthScanner] Failed to start consumer")
		return nil
	}

	logger.Info("[CONSUMER:HealthScanner] Started successfully")
	return consumer
}
