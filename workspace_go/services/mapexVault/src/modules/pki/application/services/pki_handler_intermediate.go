package services

import (
	"context"

	domainConsts "mapexVault/src/modules/pki/domain/constants"
	"mapexVault/src/modules/pki/domain/entities"
)

func (s *PkiService) loadRootEntity(ctx context.Context) (*entities.CertificateAuthority, error) {
	return s.loadEntityByKind(ctx, domainConsts.CAKindRoot)
}

func (s *PkiService) loadIntermediateEntity(ctx context.Context) (*entities.CertificateAuthority, error) {
	return s.loadEntityByKind(ctx, domainConsts.CAKindIntermediate)
}

func (s *PkiService) loadEntityByKind(ctx context.Context, kind domainConsts.CAKind) (*entities.CertificateAuthority, error) {
	e, err := s.deps.CARepository.FindByKind(ctx, kind)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrCANotBootstrapped
	}
	return e, nil
}

// decryptIntermediatePrivateKey returns the plaintext priv key PEM.
// The caller MUST discard the returned slice as soon as it has been
// consumed (signing or HTTP response) — never log, never persist.
func (s *PkiService) decryptIntermediatePrivateKey(e *entities.CertificateAuthority) ([]byte, error) {
	return s.deps.Envelope.Decrypt(e.EncryptedDEK, e.DekNonce, e.EncryptedKey, e.KeyNonce)
}
