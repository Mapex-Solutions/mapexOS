package ports

import "context"

// L2WritesPublisherPort publishes a retry hint to the L2 writes stream
// when the synchronous L2 (MinIO) write fails. The hint carries only
// the entity id — the consumer re-fetches the current state from
// Mongo before retrying the write, so a stale payload cannot
// overwrite newer state.
//
// msgId enables NATS-native dedup with a 5s window: multiple rapid
// failures on the same entity coalesce into a single retry message.
//
// The port is platform-shared (assets + assettemplates both depend
// on it) because the publisher concern is the same — only the
// subject and msgId prefix vary per module.
type L2WritesPublisherPort interface {
	PublishRetry(ctx context.Context, subject, msgId string, payload []byte) error
}
