package ports

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// EventBusPort is the minimal NATS publish surface the events module
// requires from the underlying bus. It is the union of the two kit
// interfaces actually used by EventService:
//
//   - Publisher.Publish      — JetStream-acknowledged publishes
//     (e.g., events.raw audit trail, processor.js.execute pipeline).
//   - CorePublisher.PublishCore — fire-and-forget core publishes
//     (e.g., asset.heartbeat.{orgId} per-batch).
//
// Typing the DI field as this port (instead of the concrete *natsModel.Bus)
// satisfies the dependency-inversion rule and lets unit tests substitute a
// mock without spinning up a real broker.
type EventBusPort interface {
	natsModel.Publisher
	natsModel.CorePublisher
}
