package nats

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// ExecutionStateRepository implements hot state persistence via NATS KV.
// Handles key formatting, JSON serialization, and KV operations.
type ExecutionStateRepository struct {
	kvStore natsModel.KeyValueStore
}
