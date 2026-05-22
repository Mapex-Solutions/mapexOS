package di

import (
	"mapexVault/src/modules/credentials/domain/repositories"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	envelope "github.com/Mapex-Solutions/mapexGoKit/utils/envelope"

	"go.uber.org/dig"
)

// CredentialServiceDependenciesInjection aggregates all dependencies for CredentialService.
type CredentialServiceDependenciesInjection struct {
	dig.In

	// CredentialRepo provides credential persistence (MongoDB)
	CredentialRepo repositories.CredentialRepository

	// ConnectionRepo provides connection persistence (MongoDB)
	ConnectionRepo repositories.ConnectionRepository

	// Encryption provides envelope encryption (Master Key → DEK → Data)
	Encryption *envelope.EnvelopeService

	// Publisher provides NATS JetStream publishing for vault events
	Publisher natsModel.Publisher `name:"core"`

	// ScheduleManager provides NATS JetStream scheduling for credential refresh
	ScheduleManager natsModel.ScheduleManager `name:"core"`
}
