package crypto

import "crypto/x509"

const (
	KeyAlgorithmName    = "ECDSA-P256"
	SerialBitLength     = 128
)

const SignatureAlgorithm = x509.ECDSAWithSHA256
