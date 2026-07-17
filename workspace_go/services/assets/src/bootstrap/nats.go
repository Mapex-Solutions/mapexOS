package bootstrap

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/dig"

	atMessage "assets/src/modules/assettemplates/interfaces/message"
	hmMessage "assets/src/modules/healthmonitor/interfaces/message"
	otaMessage "assets/src/modules/ota/interfaces/message"

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
//  2. Edge presence consumer — subscribes to mapexos.presence.advisory
//     published by the edge servers (MQTT broker plugin, LNS Gateway
//     Server) on every CONNECT and DISCONNECT.
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

	// OTA leader-election KV store — the dedicated lease bucket that elects the
	// single pod running the OTA pacing scan. Short TTL is only a backstop; the
	// helper fails over on a stale revision (see infrastructure/nats leaderelection).
	// Replicas default to 1 (safe on a single-node NATS, like every other stream
	// here); a clustered deployment can raise it via NATS_KV_REPLICAS for lease HA.
	c.Provide(func(params struct {
		container.In
		Client *natsModel.Client `name:"core"`
	}) natsModel.KeyValueStore {
		replicas, _ := config.GetIntValue("nats_kv_replicas")
		store, err := params.Client.CreateKeyValue(natsModel.KVConfig{
			Bucket:   otaMessage.OTALeaderBucket,
			Replicas: replicas, // config default 1; the kit also floors any 0 to 1
			TTL:      20 * time.Second,
		})
		if err != nil {
			logger.Panic("[INFRA:NATS] Failed to create OTA leader-election KV bucket: " + err.Error())
		}
		logger.Info("[INFRA:NATS] OTA leader-election KV bucket ready (TTL=20s)")
		return store
	}, container.Name("ota-leader"))

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

		// OTA stream — plan timers (schedule → start/close/scan) via native @at
		// scheduling + inbound device status advisories. STATIC subjects (ota.>);
		// planId/executionId travel in the payload. WorkQueue so pods share via
		// queue group; Duplicates < scan interval to keep the re-schedule loop alive.
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:              otaMessage.OTAScheduleStream,
			Description:       "OTA plan timers (schedule/start/close/scan) + device status advisories",
			Subjects:          []string{config.Subject("ota", "") + ">"},
			Storage:           jetstream.FileStorage,
			Retention:         jetstream.WorkQueuePolicy,
			AllowMsgSchedules: true,
			Duplicates:        10 * time.Second,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to create OTA stream")
		} else {
			logger.Info("[INFRA:NATS] OTA stream ready (WorkQueue, AllowMsgSchedules)")
		}

		// Asset template migration stream — plan start timers via native @at
		// scheduling. STATIC subjects (assettemplates.migration.>); the plan id
		// travels in the payload + MsgId. WorkQueue so pods share via queue group;
		// Duplicates < the shortest timer gap to keep re-schedules alive.
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:              atMessage.MigrationScheduleStream,
			Description:       "Asset template migration plan start timers",
			Subjects:          []string{config.Subject("assettemplates", "migration") + ".>"},
			Storage:           jetstream.FileStorage,
			Retention:         jetstream.WorkQueuePolicy,
			AllowMsgSchedules: true,
			Duplicates:        10 * time.Second,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to create asset template migration stream")
		} else {
			logger.Info("[INFRA:NATS] Asset template migration stream ready (WorkQueue, AllowMsgSchedules)")
		}

		// Edge presence stream — captures every edge server's advisories (the
		// MQTT broker plugin, the LNS Gateway Server) for both CONNECT and
		// DISCONNECT on one shared subject. WorkQueue so scaling out
		// healthmonitor pods spreads the subject across the queue group;
		// 5m max-age keeps the buffer small (presence is best-effort, an
		// older advisory than that would be misleading anyway).
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:        hmMessage.EdgePresenceStreamName,
			Description: "Edge presence advisories (connect/disconnect)",
			Subjects:    []string{hmMessage.EdgePresenceAdvisorySubject},
			Storage:     jetstream.FileStorage,
			Retention:   jetstream.WorkQueuePolicy,
			MaxAge:      5 * time.Minute,
		}); err != nil {
			logger.Error(err, "[INFRA:NATS] Failed to create edge presence stream")
		} else {
			logger.Info("[INFRA:NATS] Edge presence stream ready")
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
