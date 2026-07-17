package bootstrap

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"

	constants "mapexVault/src/modules/credentials/application/constants"
	message "mapexVault/src/modules/credentials/interfaces/message"
)

// InitNATS registers NATS connection with Bus, Fanout, Publisher, and ScheduleManager providers.
// Vault uses NATS for: publishing credential refresh/revoke events (MAPEX-VAULT stream),
// scheduled credential refresh via VAULT-SCHEDULE stream (schedules + fired consumer).
func InitNATS(c *dig.Container) {
	natsCfg := config.GetNatsConfig()
	c.Provide(func() *natsModel.Client {
		nc, err := natsModel.New(natsCfg)
		if err != nil {
			logger.Panic("Failed to connect to NATS: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] Connected to NATS")
		return nc
	}, container.Name("core"))

	// Bus for JetStream
	c.Provide(func(params struct {
		container.In
		Client *natsModel.Client `name:"core"`
	}) *natsModel.Bus {
		return natsModel.NewBus(params.Client)
	}, container.Name("core"))

	// Publisher interface
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.Publisher {
		return params.Bus
	}, container.Name("core"))

	// Fanout interface (for credential invalidation broadcast)
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.Fanout {
		return params.Bus
	}, container.Name("core"))

	// Reconcile leader-election KV store — the dedicated lease bucket that elects
	// the single pod running the credential reseed sweep. Short TTL is only a
	// backstop; the helper fails over on a stale revision. Replicas default to 1
	// (single-node safe); a NATS cluster can raise it via NATS_KV_REPLICAS.
	c.Provide(func(params struct {
		container.In
		Client *natsModel.Client `name:"core"`
	}) natsModel.KeyValueStore {
		replicas, _ := config.GetIntValue("nats_kv_replicas")
		store, err := params.Client.CreateKeyValue(natsModel.KVConfig{
			Bucket:   constants.VaultLeaderBucket,
			Replicas: replicas,
			TTL:      20 * time.Second,
		})
		if err != nil {
			logger.Panic("[INFRA:NATS] Failed to create reconcile leader KV bucket: " + err.Error())
		}
		logger.Info("[INFRA:NATS] Reconcile leader KV bucket ready (TTL=20s)")
		return store
	}, container.Name("vault-leader"))

	// ScheduleManager interface + the VAULT-SCHEDULE stream (file storage so the
	// per-credential refresh timers survive NATS restarts; one subject per
	// credential). The reseed safety-net now runs on the reconcile leader ticker,
	// not on a self-republishing stream.
	c.Provide(func(params struct {
		container.In
		Bus *natsModel.Bus `name:"core"`
	}) natsModel.ScheduleManager {
		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:              constants.VaultScheduleStreamName,
			Description:       "Vault credential refresh schedule",
			Subjects:          []string{message.VaultScheduleSubjectPattern},
			Storage:           jetstream.FileStorage,
			AllowMsgSchedules: true,
			Retention:         jetstream.LimitsPolicy,
			MaxAge:            30 * 24 * time.Hour,
		}); err != nil {
			logger.Panic("Failed to create vault schedule stream: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] Vault schedule stream ready (file storage, AllowMsgSchedules)")

		return params.Bus
	}, container.Name("core"))

}
