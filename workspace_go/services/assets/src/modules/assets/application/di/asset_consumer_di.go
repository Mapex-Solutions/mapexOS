package di

import (
	"assets/src/modules/assets/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// AssetConsumerDependenciesInjection aggregates all dependencies
// required by Asset module consumers.
//
// L2 SYNC CONSUMER:
//   - CoreBus subscribes to mapexos.l2_writes.asset on the
//     MAPEXOS-L2-WRITES stream and forwards each retry hint to
//     AssetService.ProcessL2WriteRetry for reconciliation against
//     current Mongo state.
type AssetConsumerDependenciesInjection struct {
	dig.In

	// CoreBus provides NATS Bus for JetStream consumers with DLQ support
	CoreBus *natsModel.Bus `name:"core"`

	// AssetService handles asset business logic
	AssetService ports.AssetServicePort
}
