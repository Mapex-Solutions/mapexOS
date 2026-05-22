package nats

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// L2WritesPublisherAdapter implements L2WritesPublisherPort using
// the gokit NATS Bus. Msg-Id is forwarded via the Nats-Msg-Id header
// so JetStream's native dedup (5s window on stream MAPEXOS-L2-WRITES)
// coalesces rapid successive failures on the same entity.
type L2WritesPublisherAdapter struct {
	bus *natsModel.Bus
}
