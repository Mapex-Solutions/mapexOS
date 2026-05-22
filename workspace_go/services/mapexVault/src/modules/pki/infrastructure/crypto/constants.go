package crypto

import "crypto/x509"

const (
	KeyAlgorithmName    = "ECDSA-P256"
	SerialBitLength     = 128
)

// SignatureAlgorithm — what x509.CreateCertificate actually uses for
// ECDSA P-256 keys (auto-selected). Kept exported for tests / audit.
const SignatureAlgorithm = x509.ECDSAWithSHA256
