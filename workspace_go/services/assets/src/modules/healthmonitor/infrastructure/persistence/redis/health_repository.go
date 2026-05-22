package redis

import (
	"context"
	"time"

	"assets/src/modules/healthmonitor/application/ports"

	redisModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redis"
)

// Compile-time check
var _ ports.HealthRepository = (*healthRepository)(nil)

// New creates a new Redis health repository using the existing RedisClient.
func New(client *redisModel.RedisClient) ports.HealthRepository {
	return &healthRepository{client: client}
}

/**
 * Hot path — called per heartbeat
 */

// UpdateLastSeen updates the last-seen timestamp for an asset in the sorted set.
func (r *healthRepository) UpdateLastSeen(ctx context.Context, orgId string, assetUUID string, ts time.Time) error {
	key := RedisKeyLastSeen + orgId
	return r.client.ZAdd(ctx, key, redisModel.TimeToScore(ts), assetUUID)
}

// ResetMissCounter removes the miss counter for an asset.
func (r *healthRepository) ResetMissCounter(ctx context.Context, orgId string, assetUUID string) error {
	key := RedisKeyMiss + orgId
	return r.client.HDel(ctx, key, assetUUID)
}

// IsAlerted checks if an asset is currently in the alerted (offline) set.
func (r *healthRepository) IsAlerted(ctx context.Context, orgId string, assetUUID string) (bool, error) {
	key := RedisKeyAlerted + orgId
	return r.client.SIsMember(ctx, key, assetUUID)
}

// RemoveAlerted atomically removes an asset from the alerted (offline) set.
// Returns (true, nil) only for the caller whose SREM actually removed the member
// (Redis returns count=1). Concurrent heartbeats for the same asset get (false, nil)
// and must NOT re-run the offline→online transition.
func (r *healthRepository) RemoveAlerted(ctx context.Context, orgId string, assetUUID string) (bool, error) {
	key := RedisKeyAlerted + orgId
	removed, err := r.client.SRemN(ctx, key, assetUUID)
	if err != nil {
		return false, err
	}
	return removed == 1, nil
}

// RegisterOrg adds an org to the active orgs set.
func (r *healthRepository) RegisterOrg(ctx context.Context, orgId string) error {
	return r.client.SAdd(ctx, RedisKeyOrgs, orgId)
}

// IsKnownOnline checks if the asset has already been confirmed online via heartbeat.
func (r *healthRepository) IsKnownOnline(ctx context.Context, orgId string, assetUUID string) (bool, error) {
	key := RedisKeyKnown + orgId
	return r.client.SIsMember(ctx, key, assetUUID)
}

// MarkKnownOnline adds the asset to the known-online set.
// This is a one-time flag to avoid redundant MongoDB writes on every heartbeat.
func (r *healthRepository) MarkKnownOnline(ctx context.Context, orgId string, assetUUID string) error {
	key := RedisKeyKnown + orgId
	return r.client.SAdd(ctx, key, assetUUID)
}

/**
 * Scanner — called per scan cycle
 */

// FindStale returns assetUUIDs whose last-seen timestamp is before the cutoff.
func (r *healthRepository) FindStale(ctx context.Context, orgId string, cutoff time.Time, offset int64, limit int64) ([]string, error) {
	key := RedisKeyLastSeen + orgId
	return r.client.ZRangeByScore(ctx, key, "-inf", redisModel.FormatCutoff(cutoff), offset, limit)
}

// IncrementMiss increments the miss counter for an asset and returns the new count.
func (r *healthRepository) IncrementMiss(ctx context.Context, orgId string, assetUUID string) (int64, error) {
	key := RedisKeyMiss + orgId
	return r.client.HIncrBy(ctx, key, assetUUID, 1)
}

// MarkAlerted adds an asset to the alerted set (prevents duplicate alerts).
func (r *healthRepository) MarkAlerted(ctx context.Context, orgId string, assetUUID string) error {
	key := RedisKeyAlerted + orgId
	return r.client.SAdd(ctx, key, assetUUID)
}

// GetActiveOrgs returns all org IDs that have active sensors.
func (r *healthRepository) GetActiveOrgs(ctx context.Context) ([]string, error) {
	return r.client.SMembers(ctx, RedisKeyOrgs)
}

/**
 * API enrichment — called per GetById / List
 */

// GetLastSeen returns the last-seen timestamp for a single asset.
func (r *healthRepository) GetLastSeen(ctx context.Context, orgId string, assetUUID string) (*time.Time, error) {
	key := RedisKeyLastSeen + orgId
	score, err := r.client.ZScore(ctx, key, assetUUID)
	if err != nil {
		return nil, nil // member not found
	}
	return redisModel.ScoreToTime(score), nil
}

// GetLastSeenBatch returns last-seen timestamps for multiple assets in a single round-trip.
func (r *healthRepository) GetLastSeenBatch(ctx context.Context, orgId string, assetUUIDs []string) (map[string]*time.Time, error) {
	if len(assetUUIDs) == 0 {
		return map[string]*time.Time{}, nil
	}
	key := RedisKeyLastSeen + orgId
	scores, err := r.client.ZMScore(ctx, key, assetUUIDs...)
	if err != nil {
		return nil, err
	}
	return redisModel.ScoresToTimeMap(assetUUIDs, scores), nil
}

// IsAlertedBatch checks if multiple assets are in the alerted set in a single round-trip.
func (r *healthRepository) IsAlertedBatch(ctx context.Context, orgId string, assetUUIDs []string) (map[string]bool, error) {
	if len(assetUUIDs) == 0 {
		return map[string]bool{}, nil
	}
	key := RedisKeyAlerted + orgId
	members := make([]interface{}, len(assetUUIDs))
	for i, uuid := range assetUUIDs {
		members[i] = uuid
	}
	flags, err := r.client.SMIsMember(ctx, key, members...)
	if err != nil {
		return nil, err
	}
	return redisModel.BoolSliceToMap(assetUUIDs, flags), nil
}

/**
 * Presence — called from the CONNECT advisory consumer and the disconnect anti-race check
 */

// SetLastConnectAt records the timestamp of the most recent successful MQTT
// CONNECT for an asset. Stored as unix seconds in a per-org HASH.
// Idempotent — overwrites any prior value.
func (r *healthRepository) SetLastConnectAt(ctx context.Context, orgId string, assetUUID string, ts time.Time) error {
	key := RedisKeyLastConnect + orgId
	return r.client.HSetInt64(ctx, key, assetUUID, ts.Unix())
}

// GetLastConnectAt returns the recorded last-connect timestamp or nil when
// the asset has no recorded CONNECT. The presence consumer treats a nil
// result as "never connected" — disconnect handling drops the message in
// that case to avoid marking offline an asset that was never online.
func (r *healthRepository) GetLastConnectAt(ctx context.Context, orgId string, assetUUID string) (*time.Time, error) {
	key := RedisKeyLastConnect + orgId
	sec, err := r.client.HGetInt64(ctx, key, assetUUID)
	if err != nil {
		return nil, nil // miss — caller treats as "never connected"
	}
	t := time.Unix(sec, 0)
	return &t, nil
}

/**
 * Cleanup
 */

// RemoveAsset removes all health state for an asset (used on asset deletion
// or on Enabled true→false transition). Includes the presence last-connect
// HASH so disabled assets do not leak stale CONNECT timestamps.
func (r *healthRepository) RemoveAsset(ctx context.Context, orgId string, assetUUID string) error {
	if err := r.client.PipelineRemoveFromCollections(ctx,
		RedisKeyLastSeen+orgId,
		RedisKeyMiss+orgId,
		RedisKeyAlerted+orgId,
		assetUUID,
	); err != nil {
		return err
	}
	if err := r.client.SRem(ctx, RedisKeyKnown+orgId, assetUUID); err != nil {
		return err
	}
	return r.client.HDel(ctx, RedisKeyLastConnect+orgId, assetUUID)
}
