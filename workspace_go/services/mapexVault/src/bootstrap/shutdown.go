package bootstrap

import (
	"context"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
	"go.uber.org/dig"

	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	mongoManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// InitShutdown registers graceful shutdown hooks for all infrastructure.
func InitShutdown(c *dig.Container, sm *shutdown.ShutdownManager, app *web.App) {
	c.Invoke(func(params struct {
		container.In
		Mongo  *mongoManager.MongoManager
		Redis  *redisModel.RedisClient `name:"shared"`
		NATS   *natsModel.Client       `name:"core"`
		Leader *natsModel.LeaderElection
	}) {
		sm.RegisterFunc("fiber", 0, func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		})

		// P3: reconcile leader — resign (delete the lease key for a fast handoff)
		// BEFORE the NATS connection closes at P5, so the Delete lands.
		sm.RegisterFunc("vault-leader-election", 3, func(_ context.Context) error {
			params.Leader.Stop()
			return nil
		})

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
