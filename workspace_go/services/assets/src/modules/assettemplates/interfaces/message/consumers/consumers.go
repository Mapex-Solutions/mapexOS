package consumers

import (
	"assets/src/modules/assettemplates/application/ports"
	"assets/src/modules/assettemplates/interfaces/message/consumers/list_name_updated"
	"assets/src/modules/assettemplates/interfaces/message/consumers/template_l2sync"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * Consumers barrel file - exports all consumer initialization functions
 *
 * Following Hexagonal Architecture:
 * - Consumers (Interface Layer) only receive messages and call service
 * - Service (Application Layer) handles all business logic
 */

// NewListNameUpdatedConsumer creates the list name updated consumer
func NewListNameUpdatedConsumer(bus *natsModel.Bus, assetTemplateService ports.AssetTemplateServicePort) {
	list_name_updated.NewConsumer(bus, assetTemplateService)
}

// NewTemplateL2SyncConsumer creates the template L2 write retry consumer.
// Drains MAPEXOS-L2-WRITES for the `mapexos.l2_writes.template` subject
// and forwards each retry to AssetTemplateService.ProcessL2WriteRetry.
func NewTemplateL2SyncConsumer(bus *natsModel.Bus, assetTemplateService ports.AssetTemplateServicePort) {
	template_l2sync.NewConsumer(bus, assetTemplateService)
}
