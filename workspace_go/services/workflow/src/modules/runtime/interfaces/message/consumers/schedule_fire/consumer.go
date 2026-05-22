package schedule_fire

import (
	"fmt"
	"time"

	"workflow/src/modules/runtime/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer creates a consumer on WORKFLOW-SCHEDULE stream filtering workflow.schedule.fired.
// When a NATS Schedule fires, it delivers the message to this subject.
// The consumer delegates to the service, which re-publishes to WORKFLOW-RESUME.
func NewConsumer(bus *natsModel.Bus, service ports.RuntimeServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-schedule-fire", serviceName)
	queueGroup := fmt.Sprintf("%s-WORKFLOW-SCHEDULE-GROUP", serviceName)

	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")
	if natsFetchTimeout <= 0 {
		natsFetchTimeout = 1
	}

	logger.Info(fmt.Sprintf("[CONSUMER:ScheduleFire] Starting %s", consumerName))

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:          Stream,
		Subject:         Subject,
		Durable:         consumerName,
		QueueGroup:      queueGroup,
		FetchTimeout:    time.Duration(natsFetchTimeout) * time.Second,
		DuplicateWindow: 2 * time.Minute, // Match bootstrap — no Msg-Id used on this stream

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: DLQServiceType,
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			service.HandleScheduleFire(msg)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:ScheduleFire] Failed to start consumer")
		return nil
	}

	logger.Info("[CONSUMER:ScheduleFire] Started successfully")
	return consumer
}
