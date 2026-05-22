package bootstrap

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"

	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// ConsumerRegistry holds references to all NATS consumers for shutdown.
type ConsumerRegistry struct {
	consumers []*natsModel.Consumer
}

// Register adds a consumer to the registry.
func (r *ConsumerRegistry) Register(c *natsModel.Consumer) {
	if c != nil {
		r.consumers = append(r.consumers, c)
	}
}

// StopAll stops all registered consumers.
func (r *ConsumerRegistry) StopAll() {
	for _, c := range r.consumers {
		c.Stop()
	}
}

// WalkerDrainer is implemented by RuntimeService to wait for active walkers.
type WalkerDrainer interface {
	WaitForActiveWalkers()
}

// InitShutdown registers graceful shutdown hooks for all infrastructure
// components via the DIG container. Hooks are executed by priority:
//
//	P0 — Fiber HTTP (stop accepting, drain in-flight requests)
//	P1 — Consumers (stop receiving new messages)
//	P2 — Walker drain (wait for active walkers to finish current step)
//	P5 — Connections: MongoDB, Redis, NATS (concurrent)
func InitShutdown(c *dig.Container, sm *shutdown.ShutdownManager, app *fiber.App, registry *ConsumerRegistry, drainer WalkerDrainer) {
	c.Invoke(func(params struct {
		container.In
		Mongo *mongoManager.MongoManager
		Redis *redisModel.RedisClient `name:"shared"`
		NATS  *natsModel.Client       `name:"core"`
	}) {
		// P0: HTTP — stop accepting new requests, drain in-flight
		sm.RegisterFunc("fiber", 0, func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		})

		// P1: Consumers — stop receiving new messages
		sm.RegisterFunc("consumers", 1, func(ctx context.Context) error {
			registry.StopAll()
			return nil
		})

		// P2: Walker drain — wait for active walkers to finish current step + checkpoint
		sm.RegisterFunc("walker-drain", 2, func(ctx context.Context) error {
			drainer.WaitForActiveWalkers()
			return nil
		})

		// P5: Connections (same priority — run concurrently)
		sm.RegisterFunc("mongodb", 5, func(ctx context.Context) error {
			return params.Mongo.Close(ctx)
		})

		sm.RegisterFunc("redis", 5, func(_ context.Context) error {
			return params.Redis.Close()
		})

		sm.RegisterFunc("nats", 5, func(_ context.Context) error {
			params.NATS.Close()
			return nil
		})
	})
}
