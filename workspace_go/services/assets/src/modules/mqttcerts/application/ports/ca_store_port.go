package ports

import "assets/src/modules/mqttcerts/domain/entities"

// CAStorePort is the in-RAM hot-swap-safe holder for the intermediate
// CA bundle. atomic.Pointer-backed impl lives in infrastructure/ram.
type CAStorePort interface {
	Set(ca *entities.CertificateAuthorityRAM)
	Get() *entities.CertificateAuthorityRAM
	IsReady() bool
}
