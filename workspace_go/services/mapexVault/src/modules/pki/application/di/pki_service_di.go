package di

import (
	pkiPorts "mapexVault/src/modules/pki/application/ports"
	domainRepos "mapexVault/src/modules/pki/domain/repositories"

	"go.uber.org/dig"
)

// PkiServiceDI is the DI struct for the PKI bounded context.
// Every field is a port interface (or constructor-resolved dependency)
// — never a concrete driver, per /go-arch §6.
//
// Logging happens via the mapexGoKit logger package (function calls,
// not an injected type) — matching the project convention.
type PkiServiceDI struct {
	dig.In
	CARepository domainRepos.CARepository
	Envelope     pkiPorts.EnvelopePort
	Signer       pkiPorts.X509SignerPort
}
