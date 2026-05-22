package asset_l2sync

import (
	"context"
	"encoding/json"
	"fmt"

	"assets/src/modules/assets/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// retryPayload is the wire shape of a retry message. Producer side
// marshals it from `asset_handler_l2sync.publishL2Retry` — only the
// id travels; the consumer re-fetches current state from Mongo.
type retryPayload struct {
	AssetID string `json:"assetId"`
}

// NewConsumer wires the L2 sync fallback consumer onto the
// MAPEXOS-L2-WRITES stream, filtered to the asset subject.
// Failure-retry policy uses the gokit default (exponential with
// DLQ on exhaustion). On a successful retry the service publishes
// the existing FANOUT invalidation so caches downstream refresh.
func NewConsumer(bus *natsModel.Bus, assetService ports.AssetServicePort) {
	serviceName, _ := config.GetStringValue("service_name")

	logger.Info(fmt.Sprintf("[CONSUMER:AssetL2Sync] Starting %s with DLQ support", Durable))

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
				logger.Warn(fmt.Sprintf("[CONSUMER:AssetL2Sync] payload unmarshal failed: %v", err))
				msg.Reject("malformed l2 retry payload")
				return
			}
			if p.AssetID == "" {
				msg.Reject("l2 retry payload missing assetId")
				return
			}
			if err := assetService.ProcessL2WriteRetry(context.Background(), p.AssetID); err != nil {
				logger.Warn(fmt.Sprintf("[CONSUMER:AssetL2Sync] retry failed for %s: %v", p.AssetID, err))
				msg.Nack(err)
				return
			}
			msg.Ack()
		},
	})

	if err != nil {
		logger.Error(err, "[CONSUMER:AssetL2Sync] Failed to start consumer")
		return
	}

	logger.Info("[CONSUMER:AssetL2Sync] Started successfully with DLQ support")
}
