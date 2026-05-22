package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// InitNATS registers NATS client, Bus, and Publisher interface in DIG container.
func InitNATS(c *dig.Container) {
	natsCfg := config.GetNatsConfig()
	c.Provide(func() *natsModel.Client {
		nc, err := natsModel.New(natsCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return nc
	})

	// Provide the NATS Bus (concrete type)
	c.Provide(natsModel.NewBus)

	// Provide the NATS Publisher interface (for Hexagonal Architecture testability)
	c.Provide(func(bus *natsModel.Bus) natsModel.Publisher {
		return bus
	})

	// Provide CorePublisher interface (fire-and-forget batch publishing)
	c.Provide(func(bus *natsModel.Bus) natsModel.CorePublisher {
		return bus
	})
}
