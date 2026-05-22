package nats

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// alertPublisher is the NATS-backed adapter implementing ports.AlertPublisherPort.
// It wraps a natsModel.Publisher used by publish methods in alert_publisher.go.
type alertPublisher struct {
	publisher natsModel.Publisher
}
