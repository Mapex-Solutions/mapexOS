package ports

// EnvelopePort wraps the existing envelope encryption primitive
// (Master Key → DEK → payload). Hides infrastructure from the
// application layer per /go-arch §6.
type EnvelopePort interface {
	Encrypt(plaintext []byte) (encryptedDEK, dekNonce, encryptedData, dataNonce []byte, err error)
	Decrypt(encryptedDEK, dekNonce, encryptedData, dataNonce []byte) (plaintext []byte, err error)
}
