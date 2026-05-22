package ram

import (
	"sync/atomic"

	"assets/src/modules/mqttcerts/domain/entities"
)

// InMemoryCAStore holds the intermediate CA bundle in process memory
// after the OnMount lifecycle fetches it from mapexVault. atomic.Pointer
// makes both Get and Set lock-free; hot-swap-safe for future rotation.
type InMemoryCAStore struct {
	ptr atomic.Pointer[entities.CertificateAuthorityRAM]
}
