package configMod

import (
	"triggers/src/modules/events"
	"triggers/src/modules/triggers"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/common"
)

// Modules defines the order and configuration of all modules to be initialized.
// The order is important as some modules depend on others.
//
// Initialization Order:
//  1. Core modules: triggers (HTTP API)
//  2. Event consumers: events (NATS consumer, depends on triggers for event processing)
var Modules = []common.ModuleConfig{
	// Core HTTP module (no dependencies on other Triggers service modules)
	{
		Name:             "triggers",
		Lazy:             false,
		InitRepositories: triggers.InitRepositories,
		InitServices:     triggers.InitServices,
		InitInterfaces:   triggers.InitInterfaces,
	},

	// NATS event consumer (processes events related to triggers)
	{
		Name:             "events",
		Lazy:             false,
		InitRepositories: nil, // Event consumers don't have repositories
		InitServices:     events.InitServices,
		InitInterfaces:   events.InitInterfaces, // Registers NATS consumer
	},
}
