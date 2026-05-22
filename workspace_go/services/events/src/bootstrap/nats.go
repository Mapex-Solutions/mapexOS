package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// InitNATS registers NATS client and Bus in DIG container.
func InitNATS(c *dig.Container) {
	natsCfg := config.GetNatsConfig()
	c.Provide(func() *natsModel.Client {
		nc, err := natsModel.New(natsCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return nc
	})

	// Provide the NATS Bus for publishers and consumers
	c.Provide(natsModel.NewBus)
}
