package redis

import (
	"context"
	"time"

	"assets/src/modules/ota/application/ports"
	"assets/src/modules/ota/domain/entities"

	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// liveStateAdapter implements ports.LiveStatePort by writing the current
// execution state to Redis (the live read model for the status screen).
type liveStateAdapter struct {
	rc  *redisModel.RedisClient
	ttl time.Duration
}

// NewLiveStateAdapter returns a LiveStatePort over the Redis client.
func NewLiveStateAdapter(rc *redisModel.RedisClient, ttl time.Duration) ports.LiveStatePort {
	return &liveStateAdapter{rc: rc, ttl: ttl}
}

func (a *liveStateAdapter) SetExecutionLive(ctx context.Context, exec *entities.OTAExecution) error {
	key := "ota:exec:" + exec.PlanID.Hex() + ":" + exec.AssetID.Hex()
	return a.rc.SetEx(ctx, key, exec, a.ttl)
}

var _ ports.LiveStatePort = (*liveStateAdapter)(nil)
