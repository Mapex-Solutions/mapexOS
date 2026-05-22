package appModule

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	configMod "triggers/src/shared/configuration/modules"
)

// InitModule initializes all modules in 3 phases:
// InitRepositories - Registers all repositories in DIG container
// InitServices - Registers all services in DIG container
// InitInterfaces - Registers HTTP routes and consumers
func InitModule(c *fiber.App) {

	logger.Info("[MODULE:App] Initializing Modules")

	// Show initialization order
	for i, mod := range configMod.Modules {
		logger.Info(fmt.Sprintf("[MODULE:App] %d. %s", i+1, mod.Name))
	}

	// Phase 1: Initialize all repositories
	logger.Info("[MODULE:App] Initializing Repositories")

	for _, mod := range configMod.Modules {
		if !mod.Lazy && mod.InitRepositories != nil {
			logger.Info(fmt.Sprintf("[MODULE:%s] Initializing repositories...", mod.Name))
			mod.InitRepositories()
		}
	}

	// Phase 2: Initialize all services
	logger.Info("[MODULE:App] Initializing Services")

	for _, mod := range configMod.Modules {
		if !mod.Lazy && mod.InitServices != nil {
			logger.Info(fmt.Sprintf("[MODULE:%s] Initializing services...", mod.Name))
			mod.InitServices()
		}
	}

	// Phase 3: Initialize all interfaces (HTTP + Consumers)
	logger.Info("[MODULE:App] Initializing Interfaces (HTTP Routes & Consumers)")

	for _, mod := range configMod.Modules {
		if !mod.Lazy && mod.InitInterfaces != nil {
			logger.Info(fmt.Sprintf("[MODULE:%s] Initializing interfaces...", mod.Name))
			mod.InitInterfaces()
		}
	}

	logger.Info("[MODULE:App] All modules initialized successfully")
}
