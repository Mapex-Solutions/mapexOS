package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
)

// InitMongo registers MongoDB manager in DIG container.
// The workflow service enables backpressure tracking so the Archiver
// can adapt batch behavior based on MongoDB write latency.
func InitMongo(c *dig.Container) {
	mongoCfg := config.GetMongoConfig()
	mongoCfg.EnableBackpressure = true

	c.Provide(func() *mongoManager.MongoManager {
		m, err := mongoManager.New(mongoCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return m
	})
}
