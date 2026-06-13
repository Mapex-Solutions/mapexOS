package ports

// EnvelopePort wraps the existing envelope encryption primitive
// (Master Key → DEK → payload). Hides infrastructure from the
// application layer.
type EnvelopePort interface {
	Encrypt(plaintext []byte) (encryptedDEK, dekNonce, encryptedData, dataNonce []byte, err error)
	Decrypt(encryptedDEK, dekNonce, encryptedData, dataNonce []byte) (plaintext []byte, err error)
}
