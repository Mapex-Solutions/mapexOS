package configMod

import (
	"workflow/src/modules/archiver"
	"workflow/src/modules/definitions"
	fetch_options "workflow/src/modules/fetch_options"
	"workflow/src/modules/engine"
	"workflow/src/modules/instances"
	"workflow/src/modules/plugins"
	"workflow/src/modules/runtime"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/common"
)

// Modules defines the order and configuration of all modules to be initialized.
// The order is important as some modules depend on others.
//
// Initialization Order:
//  1. Definitions: Workflow definition CRUD (no dependencies)
//  2. Engine: Condition evaluation + value resolution (pure computation, no I/O)
//  3. Runtime: DAG execution engine (depends on engine + definitions)
//  4. Archiver: WORKFLOW-STATE consumer → MongoDB BulkWrite (depends on runtime KV + same collection)
var Modules = []common.ModuleConfig{

	// Definitions module: CRUD for workflow definitions + MinIO storage
	{
		Name:             "definitions",
		Lazy:             false,
		InitRepositories: definitions.InitRepositories,
		InitServices:     definitions.InitServices,
		InitInterfaces:   definitions.InitInterfaces,
	},

	// Plugins module: CRUD for plugin manifests + TieredCache (L0→L1→MongoDB) + NATS Fanout
	{
		Name:             "plugins",
		Lazy:             false,
		InitRepositories: plugins.InitRepositories,
		InitServices:     plugins.InitServices,
		InitInterfaces:   plugins.InitInterfaces,
	},

	// Engine module: Condition evaluator + value resolver (pure, no repos/interfaces)
	{
		Name:         "engine",
		Lazy:         false,
		InitServices: engine.InitServices,
	},

	// Instances module: CRUD for workflow instances (depends on runtime entities)
	{
		Name:             "instances",
		Lazy:             false,
		InitRepositories: instances.InitRepositories,
		InitServices:     instances.InitServices,
		InitInterfaces:   instances.InitInterfaces,
	},

	// Runtime module: DAG execution engine (depends on engine + definitions)
	{
		Name:             "runtime",
		Lazy:             false,
		InitRepositories: runtime.InitRepositories,
		InitServices:     runtime.InitServices,
		InitInterfaces:   runtime.InitInterfaces,
	},

	// Archiver module: Consumes WORKFLOW-STATE → BulkWrite MongoDB → cleanup KV
	{
		Name:             "archiver",
		Lazy:             false,
		InitRepositories: archiver.InitRepositories,
		InitServices:     archiver.InitServices,
		InitInterfaces:   archiver.InitInterfaces,
	},

	// FetchOptions module: Proxy for plugin dynamic dropdowns (vault → provider HTTP)
	{
		Name:           "fetch_options",
		Lazy:           false,
		InitServices:   fetch_options.InitServices,
		InitInterfaces: fetch_options.InitInterfaces,
	},
}
