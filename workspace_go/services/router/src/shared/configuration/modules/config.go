package configMod

import (
	"router/src/modules/events"
	"router/src/modules/routegroups"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/common"
)

// Modules defines the order and configuration of all modules to be initialized.
// The order is important as some modules depend on others.
//
// Initialization Order:
//  1. Core modules: routegroups (HTTP API)
//  2. Event consumers: events (NATS consumer, depends on routegroups for event processing)
var Modules = []common.ModuleConfig{
	// Core HTTP module (no dependencies on other Router service modules)
	{
		Name:             "routegroups",
		Lazy:             false,
		InitRepositories: routegroups.InitRepositories,
		InitServices:     routegroups.InitServices,
		InitInterfaces:   routegroups.InitInterfaces,
	},

	// NATS event consumer (processes events related to route groups)
	// No InitRepositories - uses TieredCache injected via DI from main.go
	{
		Name:           "events",
		Lazy:           false,
		InitServices:   events.InitServices,
		InitInterfaces: events.InitInterfaces, // Registers NATS consumers (WorkQueue + FANOUT)
	},
}
