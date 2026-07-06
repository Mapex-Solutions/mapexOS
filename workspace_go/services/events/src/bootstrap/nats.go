package bootstrap

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/dig"

	otaEvents "github.com/Mapex-Solutions/MapexOS/contracts/services/ota/events"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// InitNATS registers NATS client and Bus in DIG container.
func InitNATS(c *dig.Container) {
	natsCfg := config.GetNatsConfig()
	c.Provide(func() *natsModel.Client {
		nc, err := natsModel.New(natsCfg)
		if err != nil {
			logger.Panic(err.Error())
		}
		return nc
	})

	// Provide the NATS Bus for publishers and consumers
	c.Provide(natsModel.NewBus)

	// Ensure the OTA status-history stream exists. The Asset MS publishes each
	// per-device OTA status transition on events.ota.status; this service's
	// ota_status consumer drains it into events_ota_status. WorkQueue retention
	// drops each row once acked (persisted to ClickHouse); MaxAge caps drainage
	// backlog. EnsureStream is a no-op when the stream already exists.
	if err := c.Invoke(func(bus *natsModel.Bus) {
		if err := bus.EnsureStream(jetstream.StreamConfig{
			Name:        otaEvents.StreamOTAStatus,
			Description: "OTA per-device status history ingestion (events_ota_status)",
			Subjects:    []string{otaEvents.SubjectOTAStatus},
			Storage:     jetstream.FileStorage,
			Retention:   jetstream.WorkQueuePolicy,
			MaxAge:      30 * 24 * time.Hour,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to ensure OTA status stream")
		} else {
			logger.Info("[INFRA:NATS] OTA status stream ready (WorkQueue, events_ota_status)")
		}
	}); err != nil {
		logger.Error(err, "[INFRA:NATS] Failed to invoke OTA status stream ensure")
	}
}
