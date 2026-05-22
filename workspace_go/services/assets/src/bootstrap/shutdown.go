package bootstrap

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"

	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// InitShutdown registers graceful shutdown hooks for all infrastructure
// components via the DIG container. Hooks are executed by priority:
//
//	P0 — Fiber HTTP (stop accepting, drain in-flight requests)
//	P5 — Connections: MongoDB, Redis, NATS Core (concurrent)
//
// Single NATS connection (core) — see bootstrap/nats.go for the rationale.
func InitShutdown(c *dig.Container, sm *shutdown.ShutdownManager, app *fiber.App) {
	c.Invoke(func(params struct {
		container.In
		Mongo    *mongoManager.MongoManager
		Redis    *redisModel.RedisClient `name:"shared"`
		NATSCore *natsModel.Client       `name:"core"`
	}) {
		// P0: HTTP — stop accepting new requests, drain in-flight
		sm.RegisterFunc("fiber", 0, func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		})

		// P5: Connections (same priority — run concurrently)
		sm.RegisterFunc("mongodb", 5, func(ctx context.Context) error {
			return params.Mongo.Close(ctx)
		})

		sm.RegisterFunc("redis", 5, func(_ context.Context) error {
			return params.Redis.Close()
		})

		sm.RegisterFunc("nats-core", 5, func(_ context.Context) error {
			params.NATSCore.Close()
			return nil
		})
	})
}
