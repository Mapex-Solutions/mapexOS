package services

import (
	"context"

	"mapexVault/src/modules/kek/domain/entities"
)

// loadByContext fetches the KEK record, mapping a missing record to
// ErrKEKNotSeeded so the HTTP layer can return 503 for the boot retry loop.
func (s *KekService) loadByContext(ctx context.Context, keyContext string) (*entities.EncryptionKey, error) {
	e, err := s.deps.Repository.FindByContext(ctx, keyContext)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrKEKNotSeeded
	}
	return e, nil
}

// decryptKEK returns the plaintext KEK (a 64-hex string). The caller MUST NOT
// log or persist it.
func (s *KekService) decryptKEK(e *entities.EncryptionKey) ([]byte, error) {
	return s.deps.Envelope.Decrypt(e.EncryptedDEK, e.DekNonce, e.EncryptedKey, e.KeyNonce)
}
