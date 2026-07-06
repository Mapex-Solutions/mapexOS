package nats

import (
	"context"

	"assets/src/modules/ota/application/ports"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// statusHistoryPublisher implements ports.StatusHistoryPublisherPort by
// publishing the status advisory to the events history subject (→ ClickHouse).
type statusHistoryPublisher struct {
	pub natsModel.Publisher
}

// NewStatusHistoryPublisher returns a StatusHistoryPublisherPort over the NATS
// publisher.
func NewStatusHistoryPublisher(pub natsModel.Publisher) ports.StatusHistoryPublisherPort {
	return &statusHistoryPublisher{pub: pub}
}

func (p *statusHistoryPublisher) Publish(ctx context.Context, advisory otaEvents.OTAStatusAdvisory) error {
	return p.pub.Publish(natsModel.PublishConfig{
		Ctx:     ctx,
		Subject: otaEvents.SubjectOTAStatus,
		Data:    advisory,
	})
}

var _ ports.StatusHistoryPublisherPort = (*statusHistoryPublisher)(nil)
