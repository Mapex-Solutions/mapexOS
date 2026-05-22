package ports

import (
	"math/big"
	"time"
)

// RootCAReq carries the inputs for self-signing a root CA.
type RootCAReq struct {
	SubjectCN string
	TTL       time.Duration
}

// IntermediateCAReq carries the inputs for signing an intermediate CA
// with an existing root.
type IntermediateCAReq struct {
	SubjectCN string
	TTL       time.Duration
}

// ServerCertReq carries the inputs for signing a server cert with the
// intermediate CA. SANs is the operator-supplied list (CSV split at the
// HTTP boundary); the signer decides DNS vs IP per entry.
type ServerCertReq struct {
	CN   string
	SANs []string
	TTL  time.Duration
}

// X509SignerPort owns asymmetric key generation + cert signing.
// Stateless. The application layer feeds it decrypted CA material per
// request and discards plaintext after the call returns.
type X509SignerPort interface {
	GenerateRootCA(req RootCAReq) (certPEM, keyPEM []byte, err error)
	GenerateIntermediateCA(rootCertPEM, rootKeyPEM []byte, req IntermediateCAReq) (certPEM, keyPEM []byte, err error)
	SignServerCert(intermediateCertPEM, intermediateKeyPEM []byte, req ServerCertReq) (certPEM []byte, serial *big.Int, err error)
}
