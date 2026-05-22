package bootstrap

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/dig"

	hmMessage "assets/src/modules/healthmonitor/interfaces/message"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	assetsAuthContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// InitNATS registers a SINGLE NATS Core connection authenticated as the
// shared 'service' user. This connection serves two roles:
//
//  1. JetStream consumers + publishers + ScheduleManager (telemetry,
//     health monitor scan schedule, fanout invalidation, etc).
//  2. MQTT presence consumer — subscribes to mapexos.mqtt.presence.advisory
//     published by the Mosquitto broker plugin (mapex-broker-plugin) on
//     every device CONNECT and DISCONNECT.
//
// Device CONNECT auth runs entirely inside the mapex-mqtt-broker plugin
// off the AssetReadModel returned by its TieredCache (L1 Pebble → L2
// MinIO → L3 GET /internal/assets/:assetUUID). This NATS init only
// wires the service-account connection used by JetStream consumers +
// fanout; it has nothing to do with device auth.
func InitNATS(c *dig.Container) {
	natsCfg := config.GetNatsConfig()
	c.Provide(func() *natsModel.Client {
		nc, err := natsModel.New(natsCfg)
		if err != nil {
			logger.Panic("[INFRA:NATS] Failed to connect to Core: " + err.Error())
		}
		logger.Info("[INFRA:NATS] Connected to Core (JetStream + presence consumer)")
		return nc
	}, container.Name("core"))

	// Bus for JetStream consumers, publishers, ScheduleManager.
	c.Provide(func(params struct {
		container.In
		Client *natsModel.Client `name:"core"`
	}) *natsModel.Bus {
		return natsModel.NewBus(params.Client)
	}, container.Name("core"))

	// Fanout interface — used by AssetTemplateService, AssetService for event publishing.
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.Fanout {
		return params.Bus
	}, container.Name("core"))

	// Subscriber interface — used by JetStream consumers.
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.Subscriber {
		return params.Bus
	}, container.Name("core"))

	// Publisher — used by health monitor alerts (publishes to ROUTE-GROUPS stream).
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.Publisher {
		return params.Bus
	}, container.Name("core"))

	// Health monitoring streams (created on the same core connection).
	c.Invoke(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) {
		// Asset heartbeat stream — public API for heartbeat producers (JS Executor, future LoRaWAN/CoAP).
		// Retention=work matches the canonical create from the nats-init container so EnsureStream
		// is a no-op when the stream already exists (NATS rejects retention-policy switches).
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:        hmMessage.HeartbeatStreamName,
			Description: "Asset heartbeat ingestion (JS Executor, LoRaWAN, CoAP)",
			Subjects:    []string{hmMessage.HeartbeatSubject},
			Storage:     jetstream.FileStorage,
			Retention:   jetstream.WorkQueuePolicy,
			MaxAge:      1 * time.Hour,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to create asset heartbeat stream")
		} else {
			logger.Info("[INFRA:NATS] Asset heartbeat stream ready")
		}

		// Health monitor scheduling stream — internal scheduler for periodic scanning.
		// WorkQueue retention enables horizontal scaling via QueueGroup.
		// Duplicates: 10s (MUST be < scan interval to avoid re-schedule dedup killing the loop).
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:              hmMessage.ScanStreamName,
			Description:       "Health monitor periodic scan scheduler",
			Subjects:          []string{config.Subject("healthmonitor", "") + ">"},
			Storage:           jetstream.FileStorage,
			Retention:         jetstream.WorkQueuePolicy,
			AllowMsgSchedules: true,
			Duplicates:        10 * time.Second,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to create asset health monitor stream")
		} else {
			logger.Info("[INFRA:NATS] Asset health monitor stream ready (WorkQueue, AllowMsgSchedules, Duplicates=10s)")
		}

		// MQTT presence stream — captures broker plugin advisories for both
		// device CONNECT and DISCONNECT events. WorkQueue so scaling out
		// healthmonitor pods spreads each subject across the queue group;
		// 5m max-age keeps the buffer small (presence is best-effort, an
		// older advisory than that would be misleading anyway).
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:        hmMessage.MqttPresenceStreamName,
			Description: "MQTT broker presence advisories (connect/disconnect)",
			Subjects:    []string{hmMessage.MqttPresenceAdvisorySubject},
			Storage:     jetstream.FileStorage,
			Retention:   jetstream.WorkQueuePolicy,
			MaxAge:      5 * time.Minute,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to create MQTT presence stream")
		} else {
			logger.Info("[INFRA:NATS] MQTT presence stream ready")
		}

		// FANOUT cache invalidation stream — broadcast channel for asset and
		// asset-template invalidation messages. Memory storage for speed +
		// short retention (consumers receive in-flight only). Memory means
		// the stream is wiped on nats-core restart, so we recreate here at
		// service start to stay independent of the nats-init sidecar.
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:        assetsContract.FanoutStreamName,
			Description: "Platform-wide cache invalidation broadcast",
			Subjects:    []string{config.Subject("fanout", "") + ">"},
			Storage:     jetstream.MemoryStorage,
			Retention:   jetstream.LimitsPolicy,
			MaxAge:      5 * time.Minute,
			MaxMsgs:     10000,
			MaxBytes:    10 * 1024 * 1024,
			Discard:     jetstream.DiscardOld,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to create FANOUT stream")
		} else {
			logger.Info("[INFRA:NATS] FANOUT stream ready")
		}

		// L2 writes retry stream — durable fallback for MinIO writes
		// that fail on the synchronous happy path. Both asset and
		// assettemplate modules publish here when their CRUD-time L2
		// write errors; in-module consumers drain the stream and
		// reconcile against current Mongo state. NATS-native Msg-Id
		// dedup with a 5s window coalesces rapid successive failures
		// on the same entity. Defense in depth — nats-init already
		// creates this on the sidecar, but ensuring here lets the
		// service survive a clean MS restart without the sidecar.
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:        assetsAuthContract.L2WritesStreamName,
			Description: "L2 write retry stream for durable cache sync",
			Subjects:    []string{config.Subject("l2_writes", "") + ">"},
			Storage:     jetstream.FileStorage,
			Retention:   jetstream.WorkQueuePolicy,
			MaxAge:      24 * time.Hour,
			Discard:     jetstream.DiscardOld,
			Duplicates:  5 * time.Second,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to create L2 writes stream")
		} else {
			logger.Info("[INFRA:NATS] L2 writes stream ready (WorkQueue, MaxAge=24h, Dupe=5s)")
		}

		// ScheduleManager interface for health monitor scheduling.
		c.Provide(func(params struct {
			container.In
			Bus *natsModel.Bus `name:"core"`
		}) natsModel.ScheduleManager {
			return params.Bus
		}, container.Name("core"))
	})
}
