package services

import (
	"time"

	pkiDi "mapexVault/src/modules/pki/application/di"
)

// PkiService implements the PkiServicePort. Stateless from a caching
// perspective — every Get/Sign decrypts the CA on demand and discards
// plaintext after use. CA documents are seeded into Mongo by the
// mongodb-init container; the service itself does no bootstrapping.
type PkiService struct {
	deps pkiDi.PkiServiceDI
}

// signedCertResult bundles the signer's output for the orchestrator.
type signedCertResult struct {
	certPEM   []byte
	serialHex string
	notAfter  time.Time
}
