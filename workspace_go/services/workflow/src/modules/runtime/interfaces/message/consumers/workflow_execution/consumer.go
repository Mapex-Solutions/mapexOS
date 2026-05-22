package workflow_execution

import (
	"fmt"
	"time"

	"workflow/src/modules/runtime/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer creates and starts a NATS consumer for the WORKFLOW-EXECUTION stream.
//
// Handles three delivery modes (dispatched by "mode" field in payload):
//   - newInstance: creates a new execution from an instance config
//   - signal: delivers a signal to a waiting execution
//   - signalOrStart: tries signal, falls back to newInstance
func NewConsumer(bus *natsModel.Bus, service ports.RuntimeServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")

	consumerName := fmt.Sprintf("%s-workflow-execution", serviceName)
	queueGroup := fmt.Sprintf("%s-WORKFLOW-EXECUTION-GROUP", serviceName)

	natsFetchTimeout, _ := config.GetIntValue("nats_fetch_timeout")

	logger.Info(fmt.Sprintf("[CONSUMER:WorkflowExecution] Starting %s with 3-mode dispatch", consumerName))

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		FetchTimeout: time.Duration(natsFetchTimeout) * time.Second,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: DLQServiceType,
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			preview := string(msg.Data)
			if len(preview) > 400 {
				preview = preview[:400] + "...(truncated)"
			}
			logger.Debug(fmt.Sprintf("[CONSUMER:WorkflowExecution] msg received subject=%s size=%d preview=%s",
				msg.Subject, len(msg.Data), preview))
			service.HandleExecution(msg)
			logger.Debug(fmt.Sprintf("[CONSUMER:WorkflowExecution] msg returned from HandleExecution subject=%s", msg.Subject))
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:WorkflowExecution] Failed to start consumer")
		return nil
	}

	logger.Info("[CONSUMER:WorkflowExecution] Started successfully with 3-mode dispatch")
	return consumer
}
