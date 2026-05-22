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

	// ScheduleManager interface + VAULT-SCHEDULE and VAULT-RECONCILER streams.
	// File storage: schedules survive NATS restarts.
	//
	// VAULT-SCHEDULE holds the per-credential refresh timers (one subject per
	// credential). VAULT-RECONCILER holds the single self-republishing reconciler
	// timer that acts as a safety net if any per-credential timer is lost.
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

		if err := params.Bus.EnsureStream(jetstream.StreamConfig{
			Name:              constants.VaultReconcilerStreamName,
			Description:       "Vault reconciler loop (credential refresh safety-net)",
			Subjects:          []string{message.VaultReconcilerSubjectPattern},
			Storage:           jetstream.FileStorage,
			AllowMsgSchedules: true,
			Retention:         jetstream.WorkQueuePolicy,
			Duplicates:        10 * time.Second,
		}); err != nil {
			logger.Panic("Failed to create vault reconciler stream: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] Vault reconciler stream ready (file storage, AllowMsgSchedules, Duplicates=10s)")

		return params.Bus
	}, container.Name("core"))

}
