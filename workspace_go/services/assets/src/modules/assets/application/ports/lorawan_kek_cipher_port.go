package ports

// LorawanKEKCipherPort is the in-RAM holder for the LoRaWAN device-keys KEK
// fetched from mapexVault at OnMount. It owns the envelope service built from
// the KEK and performs local encryption, so the plaintext KEK never leaves this
// process boundary. atomic.Pointer-backed impl lives in infrastructure/ram.
type LorawanKEKCipherPort interface {
	// SetKey builds the envelope service from the 64-hex KEK and marks the cipher
	// ready. Returns an error if the KEK is not valid hex / wrong length.
	SetKey(kekHex string) error

	// IsReady reports whether a KEK has been loaded.
	IsReady() bool

	// Encrypt envelope-encrypts plaintext with the in-RAM KEK, returning the four
	// envelope fields persisted alongside the device.
	Encrypt(plaintext []byte) (encryptedDEK, dekNonce, encryptedKey, keyNonce []byte, err error)
}
