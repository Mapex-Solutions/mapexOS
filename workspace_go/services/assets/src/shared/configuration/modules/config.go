package configMod

import (
	"assets/src/modules/assets"
	"assets/src/modules/assettemplates"
	"assets/src/modules/healthmonitor"
	"assets/src/modules/mqttcerts"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/common"
)

// Modules defines the order and configuration of all modules to be initialized.
// The order is important as some modules depend on others.
//
// Initialization Order:
//  1. Core modules: assettemplates (no dependencies)
//  2. Assets module (depends on assettemplates for templates)
//     The Mosquitto auth callout lives here as an internal route — see
//     interfaces/http/routes/internal_routes.go — so no separate auth module.
//  3. Health Monitor module (depends on assets)
//  4. MqttCerts module (depends on assets + mapexVault)
var Modules = []common.ModuleConfig{
	// Core modules (no dependencies on other Assets service modules)
	{
		Name:             "assettemplates",
		Lazy:             false,
		InitRepositories: assettemplates.InitRepositories,
		InitServices:     assettemplates.InitServices,
		InitInterfaces:   assettemplates.InitInterfaces,
	},

	// Assets module — owns the asset entity AND the L3 read-model
	// fallback (GET /internal/assets/:assetUUID) that the
	// mapex-mqtt-broker plugin's TieredCache hits on L1+L2 miss.
	// Auth decisions live entirely inside the broker plugin off the
	// returned AssetReadModel — no HTTP auth callout.
	{
		Name:             "assets",
		Lazy:             false,
		InitRepositories: assets.InitRepositories,
		InitServices:     assets.InitServices,
		InitInterfaces:   assets.InitInterfaces,
	},

	// Health Monitor module (sensor online/offline detection)
	// Depends on: assets (repository), Redis (state), NATS (heartbeat + schedule)
	{
		Name:             "healthmonitor",
		Lazy:             false,
		InitRepositories: healthmonitor.InitRepositories,
		InitServices:     healthmonitor.InitServices,
		InitInterfaces:   healthmonitor.InitInterfaces,
	},

	// MqttCerts module — device cert lifecycle (issue + revoke + list).
	// Bootstraps the CA bundle from mapexVault on OnMount with a retry
	// goroutine on failure (caReady flag gates external endpoints).
	{
		Name:             "mqttcerts",
		Lazy:             false,
		InitRepositories: mqttcerts.InitRepositories,
		InitServices:     mqttcerts.InitServices,
		InitInterfaces:   mqttcerts.InitInterfaces,
	},
}
