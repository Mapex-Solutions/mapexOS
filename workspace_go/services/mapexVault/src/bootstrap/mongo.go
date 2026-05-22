package bootstrap

import (
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
)

// InitMongo registers MongoDB manager in DIG container.
func InitMongo(c *dig.Container) {
	mongoCfg := config.GetMongoConfig()

	c.Provide(func() *mongoManager.MongoManager {
		m, err := mongoManager.New(mongoCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return m
	})
}
