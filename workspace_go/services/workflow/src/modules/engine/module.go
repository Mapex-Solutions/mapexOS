package engine

import (
	service "workflow/src/modules/engine/application/services"

	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitServices registers the engine services in the DIG container.
// Engine is a pure computation module — no repositories or interfaces needed.
func InitServices() {
	c := container.GetContainer()
	c.Provide(service.New)
	logger.Info("[MODULE:Engine] Services registered")
}
