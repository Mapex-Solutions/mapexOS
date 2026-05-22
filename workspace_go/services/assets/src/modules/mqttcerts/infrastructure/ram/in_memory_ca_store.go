package ram

import (
	mqttPorts "assets/src/modules/mqttcerts/application/ports"
	"assets/src/modules/mqttcerts/domain/entities"
)

var _ mqttPorts.CAStorePort = (*InMemoryCAStore)(nil)

// NewInMemoryCAStore returns an empty store; IsReady() is false until Set().
func NewInMemoryCAStore() mqttPorts.CAStorePort {
	return &InMemoryCAStore{}
}

func (s *InMemoryCAStore) Set(ca *entities.CertificateAuthorityRAM) {
	s.ptr.Store(ca)
}

func (s *InMemoryCAStore) Get() *entities.CertificateAuthorityRAM {
	return s.ptr.Load()
}

func (s *InMemoryCAStore) IsReady() bool {
	return s.ptr.Load() != nil
}
