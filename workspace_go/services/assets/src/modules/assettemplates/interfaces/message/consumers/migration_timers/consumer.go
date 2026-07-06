package migration_timers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"assets/src/modules/assettemplates/application/ports"
	message "assets/src/modules/assettemplates/interfaces/message"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// EventType routes exhausted migration-timer messages to the DLQ.
const EventType = "assettemplates.migration.timer"

// timerPayload is the scheduled-message body. The plan id travels in the PAYLOAD,
// never in the subject.
type timerPayload struct {
	PlanId string `json:"planId"`
}

// NewConsumer starts the migration timers consumer. It receives fired start-timer
// messages (delivered to the static assettemplates.migration.timer.* subjects) and
// drives the plan runner. QueueGroup so exactly one pod handles each firing.
func NewConsumer(bus *natsModel.Bus, service ports.AssetTemplateServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")
	consumerName := fmt.Sprintf("%s-assettemplates-migration-timers", serviceName)
	queueGroup := fmt.Sprintf("%s-ASSETTEMPLATES-MIGRATION-TIMERS-GROUP", serviceName)

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:          message.MigrationScheduleStream,
		Subject:         config.Subject("assettemplates", "migration.timer") + ".>",
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
			case message.SubjectMigrationTimerStart:
				_ = service.RunMigrationPlan(ctx, p.PlanId)
			}
			msg.Ack()
		},
	})
	if err != nil {
		logger.Error(err, "[CONSUMER:MigrationTimers] Failed to start consumer")
		return nil
	}
	logger.Info("[CONSUMER:MigrationTimers] Started")
	return consumer
}
