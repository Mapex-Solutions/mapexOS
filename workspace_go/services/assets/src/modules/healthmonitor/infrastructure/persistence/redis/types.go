package redis

import (
	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// healthRepository is the Redis-backed adapter implementing ports.HealthRepository.
// It holds the RedisClient used by all repository methods defined in health_repository.go.
type healthRepository struct {
	client *redisModel.RedisClient
}
