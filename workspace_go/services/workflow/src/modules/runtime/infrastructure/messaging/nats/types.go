package nats

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"go.uber.org/dig"
)

// RuntimePublisherParams declares DIG-injected dependencies for RuntimePublisher.
type RuntimePublisherParams struct {
	dig.In
	Publisher       natsModel.Publisher       `name:"core"`
	ScheduleManager natsModel.ScheduleManager `name:"core"`
}

// RuntimePublisher implements RuntimePublisherPort via NATS JetStream.
// Handles message construction, subject formatting, and publishing.
type RuntimePublisher struct {
	publisher       natsModel.Publisher
	scheduleManager natsModel.ScheduleManager
}
