package bootstrap

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"

	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	chManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/manager"
	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// InitShutdown registers graceful shutdown hooks for all infrastructure
// components via the DIG container. Hooks are executed by priority:
//
//	P0 — Fiber HTTP (stop accepting, drain in-flight requests)
//	P5 — Connections: ClickHouse, MongoDB, Redis, NATS (concurrent)
func InitShutdown(c *dig.Container, sm *shutdown.ShutdownManager, app *fiber.App) {
	c.Invoke(func(params struct {
		container.In
		ClickHouse *chManager.ClickHouseManager
		Mongo      *mongoManager.MongoManager
		RedisApp   *redisModel.RedisClient `name:"app"`
		RedisShd   *redisModel.RedisClient `name:"shared"`
		NATS       *natsModel.Client
	}) {
		// P0: HTTP — stop accepting new requests, drain in-flight
		sm.RegisterFunc("fiber", 0, func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		})

		// P5: Connections (same priority — run concurrently)
		sm.RegisterFunc("clickhouse", 5, func(_ context.Context) error {
			return params.ClickHouse.Close()
		})

		sm.RegisterFunc("mongodb", 5, func(ctx context.Context) error {
			return params.Mongo.Close(ctx)
		})

		sm.RegisterFunc("redis-app", 5, func(_ context.Context) error {
			return params.RedisApp.Close()
		})

		sm.RegisterFunc("redis-shared", 5, func(_ context.Context) error {
			return params.RedisShd.Close()
		})

		sm.RegisterFunc("nats", 5, func(_ context.Context) error {
			params.NATS.Close()
			return nil
		})
	})
}
