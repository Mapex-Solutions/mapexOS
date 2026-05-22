package configMod

import (
	"events/src/modules/asset_status"
	"events/src/modules/events"
	"events/src/modules/retention"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/common"
)

// Modules defines the order and configuration of all modules to be initialized.
// For the Events service:
//  - retention:     Manages retention policies per organization (MongoDB).
//  - events:        Consumes events from NATS and stores in ClickHouse.
//  - asset_status:  Persists asset connectivity transitions to ClickHouse
//                   (EVENTS-ASSET-STATUS → asset_status_history) and serves
//                   the connectivity_history HTTP query.
//
// IMPORTANT: retention MUST be initialized BEFORE events because
// EventService depends on RetentionServicePort for TTL resolution.
var Modules = []common.ModuleConfig{
	{
		Name:             "retention",
		Lazy:             false,
		InitRepositories: retention.InitRepositories,
		InitServices:     retention.InitServices,
		InitInterfaces:   retention.InitInterfaces,
	},
	{
		Name:             "events",
		Lazy:             false,
		InitRepositories: events.InitRepositories,
		InitServices:     events.InitServices,
		InitInterfaces:   events.InitInterfaces,
	},
	{
		Name:             "asset_status",
		Lazy:             false,
		InitRepositories: asset_status.InitRepositories,
		InitServices:     asset_status.InitServices,
		InitInterfaces:   asset_status.InitInterfaces,
	},
}
