package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
)

// InitFiber creates and registers the Fiber instance with global middlewares.
func InitFiber(c *dig.Container) *fiber.App {
	serviceName, _ := config.GetStringValue("service_name")
	serviceVersion, _ := config.GetStringValue("service_version")

	fiberInstance := fiber.New(fiber.Config{
		AppName:      serviceName + " " + serviceVersion,
		ErrorHandler: customErrors.FiberErrorHandler,
	})

	c.Provide(func() *fiber.App {
		return fiberInstance
	})

	// Register /metrics endpoint before global middlewares
	c.Invoke(func(m *MapexosMetrics) {
		m.Registry.RegisterEndpoint(fiberInstance)
	})

	// Global middlewares
	fiberInstance.Use(cors.New())
	fiberInstance.Use(helmet.New())

	return fiberInstance
}
