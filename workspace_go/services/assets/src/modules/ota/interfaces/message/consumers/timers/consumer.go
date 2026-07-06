package timers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/application/services"
	otaMsg "assets/src/modules/ota/interfaces/message"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

const EventType = "ota.timer"

// timerPayload is the scheduled-message body (planId/firmwareId in the PAYLOAD,
// never the subject).
type timerPayload struct {
	PlanId     string `json:"planId"`
	FirmwareId string `json:"firmwareId"`
}

// NewConsumer starts the OTA timers consumer. It receives the fired
// start/close/scan/abandon timer messages (delivered to the static ota.timer.*
// subjects) and drives the reconciler-timers (and the firmware abandon-check).
// QueueGroup so exactly one pod handles each.
func NewConsumer(bus *natsModel.Bus, timers *services.ReconcilerTimers, firmware ports.OTAFirmwareServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")
	consumerName := fmt.Sprintf("%s-ota-timers", serviceName)
	queueGroup := fmt.Sprintf("%s-OTA-TIMERS-GROUP", serviceName)

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:          otaMsg.OTAScheduleStream,
		Subject:         config.Subject("ota", "timer") + ".>",
		Durable:         consumerName,
		QueueGroup:      queueGroup,
		FetchTimeout:    5 * time.Second,
		DuplicateWindow: 10 * time.Second,
		RetryPolicy:     natsModel.DefaultRetryPolicy(),
		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "assets",
			EventType:   EventType,
		},
		MessageHandlerV2: func(msg *natsModel.Message) {
			ctx := context.Background()
			var p timerPayload
			_ = json.Unmarshal(msg.Data, &p)

			switch msg.Subject {
			case otaMsg.SubjectOTATimerStart:
				_ = timers.OnStart(ctx, p.PlanId)
			case otaMsg.SubjectOTATimerClose:
				_ = timers.OnClose(ctx, p.PlanId)
			case otaMsg.SubjectOTATimerScan:
				_ = timers.OnScanTick(ctx)
			case otaMsg.SubjectOTATimerAbandon:
				_ = firmware.HandleFirmwareAbandon(ctx, p.FirmwareId)
			}
			msg.Ack()
		},
	})
	if err != nil {
		logger.Error(err, "[CONSUMER:OTATimers] Failed to start consumer")
		return nil
	}
	logger.Info("[CONSUMER:OTATimers] Started")
	return consumer
}
