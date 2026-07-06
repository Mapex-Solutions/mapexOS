package status

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"assets/src/modules/ota/application/services"
	otaMsg "assets/src/modules/ota/interfaces/message"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

const EventType = "ota.status"

// NewConsumer starts the OTA status consumer. It receives the normalized device
// status advisories the edge (MQTT broker / HTTP gateway) publishes on the
// static SubjectOTAStatusAdvisory and drives the status handler.
func NewConsumer(bus *natsModel.Bus, handler *services.StatusHandler) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")
	consumerName := fmt.Sprintf("%s-ota-status", serviceName)
	queueGroup := fmt.Sprintf("%s-OTA-STATUS-GROUP", serviceName)

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       otaMsg.OTAScheduleStream,
		Subject:      otaMsg.SubjectOTAStatusAdvisory,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		FetchTimeout: 5 * time.Second,
		RetryPolicy:  natsModel.DefaultRetryPolicy(),
		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "assets",
			EventType:   EventType,
		},
		MessageHandlerV2: func(msg *natsModel.Message) {
			var adv otaEvents.OTAStatusAdvisory
			if err := json.Unmarshal(msg.Data, &adv); err != nil {
				msg.Ack() // malformed advisory — drop, don't retry
				return
			}
			_ = handler.HandleStatus(context.Background(), adv)
			msg.Ack()
		},
	})
	if err != nil {
		logger.Error(err, "[CONSUMER:OTAStatus] Failed to start consumer")
		return nil
	}
	logger.Info("[CONSUMER:OTAStatus] Started")
	return consumer
}
