package reconcile

import (
	"context"
	"fmt"
	"time"

	"mapexVault/src/modules/credentials/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer creates and starts a NATS consumer for vault reconcile events.
//
// Uses a QueueGroup so multiple pods share the work: only one pod handles the
// fired reconcile message per interval. The handler invokes RunReconcile which
// re-arms the next timer before acking, keeping the loop self-sustaining.
func NewConsumer(bus *natsModel.Bus, service ports.CredentialServicePort) *natsModel.Consumer {
	serviceName, _ := config.GetStringValue("service_name")
	consumerName := fmt.Sprintf("%s-vault-reconciler", serviceName)
	queueGroup := fmt.Sprintf("%s-VAULT-RECONCILE-GROUP", serviceName)

	logger.Info(fmt.Sprintf("[CONSUMER:VaultReconciler] Starting %s", consumerName))

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:       Stream,
		Subject:      Subject,
		Durable:      consumerName,
		QueueGroup:   queueGroup,
		FetchTimeout: 5 * time.Second,

		// DuplicateWindow MUST be < reconcile interval (1h) to avoid re-schedule dedup.
		// Without this, the wrapper forces stream Duplicates=15m (subscribe.go:33);
		// current 1h interval masks the bug but a shorter interval would kill the loop.
		DuplicateWindow: 10 * time.Second,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "mapexVault",
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			service.RunReconcile(context.Background())
			msg.Ack()
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:VaultReconciler] Failed to start consumer")
		return nil
	}

	logger.Info("[CONSUMER:VaultReconciler] Started successfully")
	return consumer
}
