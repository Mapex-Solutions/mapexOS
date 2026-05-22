package asset_invalidate

import (
	"fmt"

	"router/src/modules/events/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * NewConsumer creates and starts a NATS FANOUT consumer for asset cache invalidation.
 *
 * FANOUT Pattern:
 * - Each service instance receives a copy of the message (no queue group)
 * - Used for cache invalidation across all replicas
 * - Ephemeral consumer (not durable) - created fresh on each startup
 *
 * TieredCache Architecture:
 *   L0 (RAM): Hot cache - cleared on invalidation
 *   L1 (Disk): Persistent cache - cleared on invalidation
 *   L2 (MinIO): Source of truth - NOT affected
 *
 * Flow:
 * 1. Asset service updates MinIO read model
 * 2. Asset service publishes FANOUT invalidation
 * 3. All service instances receive the message
 * 4. Each instance clears L0+L1 for that assetUUID
 * 5. Next request fetches fresh data from L2 → populates L0/L1
 *
 * Parameters:
 *   - bus: The NATS bus instance for FANOUT subscription
 *   - eventService: The EventServicePort interface for processing invalidations
 */
func NewConsumer(bus *natsModel.Bus, eventService ports.EventServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	logger.Info(fmt.Sprintf("[CONSUMER:AssetInvalidate] Starting FANOUT subscription: %s -> %s", serviceName, Subject))

	_, err := bus.SubscribeFanout(Stream, serviceName, Subject, func(data []byte) error {
		msg := &natsModel.Message{Data: data}
		eventService.ProcessAssetInvalidateBatch([]*natsModel.Message{msg})
		return nil
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:AssetInvalidate] Failed to start FANOUT subscription")
		return
	}

	logger.Info("[CONSUMER:AssetInvalidate] FANOUT subscription started successfully")
}
