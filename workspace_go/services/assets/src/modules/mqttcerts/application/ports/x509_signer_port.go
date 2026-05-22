package ports

import (
	"math/big"

	"assets/src/modules/mqttcerts/domain/entities"
)

// X509SignerPort signs device certs using the RAM-cached CA.
type X509SignerPort interface {
	SignDeviceCert(
		ca *entities.CertificateAuthorityRAM,
		subjectCN string,
		ttlDays int,
	) (certPEM []byte, keyPEM []byte, serial *big.Int, fingerprint string, err error)
}
