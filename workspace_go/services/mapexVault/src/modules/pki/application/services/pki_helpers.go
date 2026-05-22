package services

import (
	pkiDtos "mapexVault/src/modules/pki/application/dtos"
	"mapexVault/src/modules/pki/domain/entities"
)

// entityToBundleDTO maps a CertificateAuthority + decrypted priv key
// into the wire DTO. The decryptedKey slice ownership transfers to the
// caller; the application service MUST NOT log it.
func (s *PkiService) entityToBundleDTO(e *entities.CertificateAuthority, decryptedKey []byte) *pkiDtos.IntermediateCABundleResponse {
	return &pkiDtos.IntermediateCABundleResponse{
		CertPEM:       e.CertPEM,
		PrivateKeyPEM: decryptedKey,
		SubjectCN:     e.SubjectCN,
		NotBefore:     e.NotBefore,
		NotAfter:      e.NotAfter,
		Fingerprint:   e.Fingerprint,
	}
}

// entityToCAChainDTO concatenates root + intermediate cert PEMs.
func (s *PkiService) entityToCAChainDTO(root, intermediate *entities.CertificateAuthority) *pkiDtos.CAChainResponse {
	chain := append([]byte{}, root.CertPEM...)
	chain = append(chain, intermediate.CertPEM...)
	return &pkiDtos.CAChainResponse{Chain: chain}
}
