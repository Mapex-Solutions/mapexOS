package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"

	eventsPorts "http_gateway/src/modules/events/application/ports"
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

	// Adapter: the events module DI struct requires the EventBusPort
	// interface. The concrete *natsModel.Bus satisfies it, but dig cannot
	// inject a concrete type into an interface field automatically — this
	// provider bridges the two.
	c.Provide(func(b *natsModel.Bus) eventsPorts.EventBusPort {
		return b
	})
}
