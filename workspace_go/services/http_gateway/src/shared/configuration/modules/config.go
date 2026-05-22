package configMod

import (
	datasources "http_gateway/src/modules/datasources"
	events "http_gateway/src/modules/events"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/common"
)

// Modules defines the order and configuration of all modules to be initialized.
// The order is important as some modules depend on others.
//
// Initialization Order:
//  1. Core modules: datasources (manages data source configurations)
//  2. Event processing: events (depends on datasources for authentication/validation)
var Modules = []common.ModuleConfig{
	// Core module (no dependencies on other HTTP Gateway modules)
	{
		Name:             "datasources",
		Lazy:             false,
		InitRepositories: datasources.InitRepositories,
		InitServices:     datasources.InitServices,
		InitInterfaces:   datasources.InitInterfaces,
	},

	// Event webhook receiver (depends on datasources for auth validation)
	{
		Name:             "events",
		Lazy:             false,
		InitRepositories: nil, // Event handlers don't have repositories
		InitServices:     events.InitServices,
		InitInterfaces:   events.InitInterfaces,
	},
}
