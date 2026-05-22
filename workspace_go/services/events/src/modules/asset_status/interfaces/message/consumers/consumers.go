package consumers

import (
	"events/src/modules/asset_status/application/ports"
	"events/src/modules/asset_status/interfaces/message/consumers/asset_status_save"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// NewAssetStatusSaveConsumer is the DIG-invocable entry point wired in
// module.go InitInterfaces. Delegates to the asset_status_save package.
func NewAssetStatusSaveConsumer(bus *natsModel.Bus, service ports.AssetStatusServicePort) {
	asset_status_save.NewConsumer(bus, service)
}
