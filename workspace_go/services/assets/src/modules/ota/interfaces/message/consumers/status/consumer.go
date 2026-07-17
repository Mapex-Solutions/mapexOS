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
		// Match the timers consumer's window: this consumer shares OTAScheduleStream,
		// and an unset DuplicateWindow makes the kit force the stream's Duplicates to
		// the 15m default, which would break the per-plan start/close dedup.
		DuplicateWindow: 10 * time.Second,
		RetryPolicy:     natsModel.DefaultRetryPolicy(),
		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "assets",
			EventType:   EventType,
		},
		MessageHandlerV2: func(msg *natsModel.Message) {
			var adv otaEvents.OTAStatusAdvisory
			if err := json.Unmarshal(msg.Data, &adv); err != nil {
				logger.Warn(fmt.Sprintf("[CONSUMER:OTAStatus] malformed advisory dropped: %v", err))
				msg.Ack() // malformed advisory — drop, don't retry
				return
			}
			logger.Debug(fmt.Sprintf("[CONSUMER:OTAStatus] advisory received: execId=%s status=%s progress=%d",
				adv.OTAExecutionID, adv.Status, adv.Progress))
			if err := handler.HandleStatus(context.Background(), adv); err != nil {
				logger.Error(err, fmt.Sprintf("[CONSUMER:OTAStatus] handle failed: execId=%s status=%s", adv.OTAExecutionID, adv.Status))
			}
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
