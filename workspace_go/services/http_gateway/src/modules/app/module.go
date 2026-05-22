package appModule

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	configMod "http_gateway/src/shared/configuration/modules"
)

// InitModule initializes all modules in 3 phases:
// 1. InitRepositories - Registers all repositories in DIG container
// 2. InitServices - Registers all services in DIG container
// 3. InitInterfaces - Registers HTTP routes and consumers
func InitModule(c *fiber.App) {

	logger.Info("[APP:MODULE] Initializing modules")

	// Show initialization order
	for i, mod := range configMod.Modules {
		logger.Info(fmt.Sprintf("[APP:MODULE] %d. %s", i+1, mod.Name))
	}

	// Phase 1: Initialize all repositories
	logger.Info("[APP:MODULE] Phase 1 - Initializing Repositories")

	for _, mod := range configMod.Modules {
		if !mod.Lazy && mod.InitRepositories != nil {
			logger.Info(fmt.Sprintf("[MODULE:%s] Initializing repositories...", mod.Name))
			mod.InitRepositories()
		}
	}

	// Phase 2: Initialize all services
	logger.Info("[APP:MODULE] Phase 2 - Initializing Services")

	for _, mod := range configMod.Modules {
		if !mod.Lazy && mod.InitServices != nil {
			logger.Info(fmt.Sprintf("[MODULE:%s] Initializing services...", mod.Name))
			mod.InitServices()
		}
	}

	// Phase 3: Initialize all interfaces (HTTP + Consumers)
	logger.Info("[APP:MODULE] Phase 3 - Initializing Interfaces (HTTP Routes & Consumers)")

	for _, mod := range configMod.Modules {
		if !mod.Lazy && mod.InitInterfaces != nil {
			logger.Info(fmt.Sprintf("[MODULE:%s] Initializing interfaces...", mod.Name))
			mod.InitInterfaces()
		}
	}

	logger.Info("[APP:MODULE] All modules initialized successfully")
}
