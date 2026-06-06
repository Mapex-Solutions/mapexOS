package services

import (
	kekDi "mapexVault/src/modules/kek/application/di"
)

// KekService implements the KekServicePort. Stateless from a caching
// perspective: every GetKEK loads + decrypts on demand and discards plaintext
// after use. KEK documents are seeded into Mongo by the mongodb-init container;
// the service itself does no bootstrapping.
type KekService struct {
	deps kekDi.KekServiceDI
}
