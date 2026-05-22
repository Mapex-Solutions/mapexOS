package template_l2sync

import (
	"context"
	"encoding/json"
	"fmt"

	"assets/src/modules/assettemplates/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// retryPayload is the wire shape of a template retry message.
type retryPayload struct {
	TemplateID string `json:"templateId"`
}

// NewConsumer wires the L2 sync fallback consumer for the
// assettemplates module onto the MAPEXOS-L2-WRITES stream, filtered to
// `mapexos.l2_writes.template`. Mirrors the asset_l2sync consumer
// shape (work-queue + DLQ + exponential retry) — only the QueueGroup
// and payload shape differ.
func NewConsumer(bus *natsModel.Bus, service ports.AssetTemplateServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	logger.Info(fmt.Sprintf("[CONSUMER:TemplateL2Sync] Starting %s with DLQ support", Durable))

	_, err := bus.StartConsumer(natsModel.ConsumerOptions{
		Stream:     Stream,
		Subject:    Subject,
		Durable:    Durable,
		QueueGroup: QueueGroup,

		RetryPolicy: natsModel.DefaultRetryPolicy(),

		DLQPolicy: &natsModel.DLQPolicy{
			ServiceName: serviceName,
			ServiceType: "assets",
			EventType:   EventType,
		},

		MessageHandlerV2: func(msg *natsModel.Message) {
			var p retryPayload
			if err := json.Unmarshal(msg.Data, &p); err != nil {
				logger.Warn(fmt.Sprintf("[CONSUMER:TemplateL2Sync] payload unmarshal failed: %v", err))
				msg.Reject("malformed l2 retry payload")
				return
			}
			if p.TemplateID == "" {
				msg.Reject("l2 retry payload missing templateId")
				return
			}
			if err := service.ProcessL2WriteRetry(context.Background(), p.TemplateID); err != nil {
				logger.Warn(fmt.Sprintf("[CONSUMER:TemplateL2Sync] retry failed for %s: %v", p.TemplateID, err))
				msg.Nack(err)
				return
			}
			msg.Ack()
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:TemplateL2Sync] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:TemplateL2Sync] Started successfully with DLQ support")
}
