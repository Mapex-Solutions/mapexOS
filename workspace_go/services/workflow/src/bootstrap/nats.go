package bootstrap

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"

	runtimeConstants "workflow/src/modules/runtime/application/constants"
)

// InitNATS registers NATS Core connection with Bus, Fanout, Subscriber, and KV providers.
// Workflow service uses a single NATS connection (Core) for JetStream consumers, publishers,
// and a KV bucket for workflow instance hot state during execution.
func InitNATS(c *dig.Container) {
	// Initialize NATS Core Connection (JetStream, Domain Events)
	// Used for: Workflow trigger consumers, transition workers, fanout events, KV state
	natsCoreCfg := config.GetNatsConfig()
	c.Provide(func() *natsModel.Client {
		nc, err := natsModel.New(natsCoreCfg)
		if err != nil {
			logger.Panic("Failed to connect to NATS Core: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] Connected to NATS Core (JetStream)")
		return nc
	}, container.Name("core"))

	// Provide NATS Core Bus for JetStream consumers
	c.Provide(func(params struct {
		container.In
		Client *natsModel.Client `name:"core"`
	}) *natsModel.Bus {
		return natsModel.NewBus(params.Client)
	}, container.Name("core"))

	// Provide NATS Core interfaces for DI (services depend on interfaces, not concrete types)
	// Fanout interface - used by DefinitionsService for cache invalidation events
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.Fanout {
		return params.Bus
	}, container.Name("core"))

	// Subscriber interface - used by consumers for JetStream subscriptions
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.Subscriber {
		return params.Bus
	}, container.Name("core"))

	// Publisher interface - used by RuntimeService for resume/callback publishing
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.Publisher {
		return params.Bus
	}, container.Name("core"))

	// Provide NATS KV store for workflow instance hot state.
	// Creates the WORKFLOW-INSTANCES bucket (file-backed, survives NATS restarts).
	// Runtime writes state per-step (CAS), Archiver reads on completion then deletes.
	c.Provide(func(params struct {
		container.In
		Client *natsModel.Client `name:"core"`
	}) natsModel.KeyValueStore {
		kvStore, err := params.Client.CreateKeyValue(natsModel.KVConfig{
			Bucket:      runtimeConstants.KVBucketName,
			Description: runtimeConstants.KVBucketDescription,
			Replicas:    runtimeConstants.KVReplicas,
		})
		if err != nil {
			logger.Panic("Failed to create NATS KV bucket: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] NATS KV bucket ready: " + runtimeConstants.KVBucketName)
		return kvStore
	})

	// Provide NATS ScheduleManager interface + create WORKFLOW-SCHEDULE stream.
	// File storage: schedules survive NATS restarts — no reconciliation needed.
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.ScheduleManager {
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:              runtimeConstants.ScheduleStreamName,
			Description:       "Workflow schedule (long-retention, scheduled msgs)",
			Subjects:          []string{runtimeConstants.ScheduleSubjectPattern},
			Storage:           jetstream.FileStorage,
			AllowMsgSchedules: true,
			Retention:         jetstream.LimitsPolicy,
			MaxAge:            30 * 24 * time.Hour,
			Duplicates:        2 * time.Minute,
		}); err != nil {
			logger.Panic("Failed to create workflow schedule stream: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] Workflow schedule stream ready (file storage, AllowMsgSchedules)")
		return params.Bus
	}, container.Name("core"))
}
