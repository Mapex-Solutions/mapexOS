package plugin_fanout

import (
	"fmt"

	"workflow/src/modules/plugins/application/constants"
	"workflow/src/modules/plugins/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * NewConsumer wires the NATS FANOUT subscription for plugin cache invalidation
 * and delegates every message to the application-layer service. No business
 * logic lives here — this function only binds transport (FANOUT stream/subject)
 * to the service port.
 *
 * FANOUT Pattern:
 *  - Each workflow pod receives a copy of the message (no queue group)
 *  - Used to invalidate TieredCache (L0/L1) across all replicas
 *  - Ephemeral subscription (not durable) — created fresh on each startup
 */
func NewConsumer(bus natsModel.Fanout, service ports.PluginServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	logger.Info(fmt.Sprintf("[CONSUMER:PluginFanout] Starting FANOUT subscription: %s -> %s", serviceName, constants.FanoutPluginSubject))

	_, err := bus.SubscribeFanout(
		constants.FanoutStreamName,
		serviceName,
		constants.FanoutPluginSubject,
		func(data []byte) error {
			service.HandleFanoutEvent(&natsModel.Message{Data: data})
			return nil
		},
	)
	if err != nil {
		logger.Error(err, "[CONSUMER:PluginFanout] Failed to start FANOUT subscription")
		return
	}

	logger.Info("[CONSUMER:PluginFanout] FANOUT subscription started successfully")
}
