package ports

// EnvelopePort wraps the envelope decryption primitive (Master Key -> DEK ->
// payload). The kek module only decrypts: KEKs are envelope-encrypted at seed
// time by the mongodb-init container. Hides infrastructure from the application
// layer per /go-arch §6.
type EnvelopePort interface {
	Decrypt(encryptedDEK, dekNonce, encryptedData, dataNonce []byte) (plaintext []byte, err error)
}
