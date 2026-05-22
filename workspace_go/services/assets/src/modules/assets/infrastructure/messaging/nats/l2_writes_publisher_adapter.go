package nats

import (
	"context"
	"fmt"

	"assets/src/modules/assets/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewL2WritesPublisherAdapter builds the adapter from an injected Bus.
// The Bus is wired in bootstrap/nats.go alongside the assets-module
// NATS connection.
func NewL2WritesPublisherAdapter(bus *natsModel.Bus) ports.L2WritesPublisherPort {
	return &L2WritesPublisherAdapter{bus: bus}
}

var _ ports.L2WritesPublisherPort = (*L2WritesPublisherAdapter)(nil)

// PublishRetry pushes a JSON message onto the configured subject
// with `Nats-Msg-Id: <msgId>` set on the headers — this is the
// JetStream-native dedup signal.
func (a *L2WritesPublisherAdapter) PublishRetry(ctx context.Context, subject, msgId string, payload []byte) error {
	err := a.bus.Publish(natsModel.PublishConfig{
		Ctx:     ctx,
		Subject: subject,
		Data:    payload,
		Headers: map[string]string{
			"Nats-Msg-Id": msgId,
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("[REPO:Nats] L2 writes publish failed subject=%s msgId=%s err=%v", subject, msgId, err))
		return err
	}
	return nil
}
