package configMod

import (
	"mapexVault/src/modules/credentials"
	"mapexVault/src/modules/pki"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/common"
)

// Modules defines the order and configuration of all modules to be initialized.
var Modules = []common.ModuleConfig{
	{
		Name:             "credentials",
		Lazy:             false,
		InitRepositories: credentials.InitRepositories,
		InitServices:     credentials.InitServices,
		InitInterfaces:   credentials.InitInterfaces,
	},
	{
		Name:             "pki",
		Lazy:             false,
		InitRepositories: pki.InitRepositories,
		InitServices:     pki.InitServices,
		InitInterfaces:   pki.InitInterfaces,
	},
}
