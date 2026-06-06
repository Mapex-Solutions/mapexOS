package ram

import (
	"errors"

	"assets/src/modules/assets/application/ports"

	envelopeUtil "github.com/Mapex-Solutions/mapexGoKit/utils/envelope"
)

var _ ports.LorawanKEKCipherPort = (*InMemoryKEKCipher)(nil)

// NewInMemoryKEKCipher returns an empty cipher; IsReady() is false until SetKey().
func NewInMemoryKEKCipher() ports.LorawanKEKCipherPort {
	return &InMemoryKEKCipher{}
}

// SetKey builds the envelope service from the 64-hex KEK and stores it.
func (c *InMemoryKEKCipher) SetKey(kekHex string) error {
	svc, err := envelopeUtil.New(kekHex)
	if err != nil {
		return err
	}
	c.svc.Store(svc)
	return nil
}

func (c *InMemoryKEKCipher) IsReady() bool {
	return c.svc.Load() != nil
}

// Encrypt envelope-encrypts plaintext with the in-RAM KEK, returning the four
// envelope fields.
func (c *InMemoryKEKCipher) Encrypt(plaintext []byte) (encryptedDEK, dekNonce, encryptedKey, keyNonce []byte, err error) {
	svc := c.svc.Load()
	if svc == nil {
		return nil, nil, nil, nil, errors.New("kek cipher not ready")
	}
	env, err := svc.Encrypt(plaintext)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return env.EncryptedDEK, env.DEKNonce, env.EncryptedData, env.DataNonce, nil
}
