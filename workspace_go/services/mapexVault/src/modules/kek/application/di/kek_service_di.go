package di

import (
	kekPorts "mapexVault/src/modules/kek/application/ports"
	domainRepos "mapexVault/src/modules/kek/domain/repositories"

	"go.uber.org/dig"
)

// KekServiceDI is the DI struct for the KEK bounded context. Every field is a
// port interface, never a concrete driver, per /go-arch §6.
type KekServiceDI struct {
	dig.In
	Repository domainRepos.EncryptionKeyRepository
	Envelope   kekPorts.EnvelopePort
}
