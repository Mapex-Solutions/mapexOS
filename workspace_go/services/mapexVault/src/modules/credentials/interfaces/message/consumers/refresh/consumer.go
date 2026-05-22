package refresh

import (
	"fmt"

	"mapexVault/src/modules/credentials/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewConsumer creates and starts a NATS consumer for vault refresh events.
//
// Per-credential timers fire on VAULT-SCHEDULE; the handler delegates to
// HandleRefreshMessage which refreshes the token and re-arms the next timer
// before acking, keeping the schedule self-sustaining.
func NewConsumer(bus *natsModel.Bus, service ports.CredentialServicePort) *natsModel.Consumer {
	logger.Info(fmt.Sprintf("[CONSUMER:VaultRefresh] Starting %s", Durable))

	consumer, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:    Stream,
		Subject:   Subject,
		Durable:   Durable,
		BatchSize: 1,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: "mapexVault",
			ServiceType: "vault",
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			service.HandleRefreshMessage(msg)
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:VaultRefresh] Failed to start consumer")
		return nil
	}

	logger.Info("[CONSUMER:VaultRefresh] Started successfully")
	return consumer
}
