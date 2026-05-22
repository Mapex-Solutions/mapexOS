package consumers

import (
	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/interfaces/message/consumers/asset_l2sync"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * Consumers barrel — exports module-level consumer initialization
 * functions registered from the assets module.go.
 *
 * Following Hexagonal Architecture:
 *   - Consumers (interfaces layer) only receive messages and call service
 *   - Service (application layer) handles all business logic
 */

// NewAssetL2SyncConsumer creates the asset L2 write retry consumer.
// Drains the MAPEXOS-L2-WRITES stream for the `mapexos.l2_writes.asset`
// subject and forwards each retry to AssetService.ProcessL2WriteRetry.
func NewAssetL2SyncConsumer(bus *natsModel.Bus, service ports.AssetServicePort) {
	asset_l2sync.NewConsumer(bus, service)
}
