package di

import (
	"assets/src/modules/assettemplates/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// AssetTemplateConsumerDependenciesInjection aggregates all dependencies required
// by AssetTemplate consumers.
//
// Following Hexagonal Architecture principles:
//   - Messaging: NatsBus (Core) for JetStream subscriptions
//   - Application: AssetTemplateService for handling business logic
//
// This struct uses dig.In to enable automatic dependency injection by uber/dig.
//
// LIST NAME UPDATED CONSUMER:
//   - CoreBus subscribes to list name updated events via JetStream
//   - Used for denormalized name synchronization across services
type AssetTemplateConsumerDependenciesInjection struct {
	dig.In

	// CoreBus provides NATS Bus for JetStream consumers with DLQ support
	CoreBus *natsModel.Bus `name:"core"`

	// AssetTemplateService handles asset template business logic
	AssetTemplateService ports.AssetTemplateServicePort
}
