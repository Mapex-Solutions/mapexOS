package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	pkiDtos "mapexVault/src/modules/pki/application/dtos"
	pkiPorts "mapexVault/src/modules/pki/application/ports"
)

func (s *PkiService) validateSignRequest(req *pkiDtos.SignServerRequest) error {
	if req == nil {
		return errors.New("nil request")
	}
	if req.CN == "" {
		return errors.New("cn required")
	}
	if len(req.SANs) == 0 {
		return errors.New("at least one SAN required")
	}
	if req.TTLDays <= 0 {
		return errors.New("ttlDays must be > 0")
	}
	return nil
}

func (s *PkiService) loadDecryptedIntermediateMaterial(ctx context.Context) ([]byte, []byte, error) {
	e, err := s.loadIntermediateEntity(ctx)
	if err != nil {
		return nil, nil, err
	}
	keyPEM, err := s.decryptIntermediatePrivateKey(e)
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt: %w", err)
	}
	return e.CertPEM, keyPEM, nil
}

func (s *PkiService) runSigner(intCertPEM, intKeyPEM []byte, req *pkiDtos.SignServerRequest) (*signedCertResult, error) {
	ttl := time.Duration(req.TTLDays) * 24 * time.Hour
	certPEM, serial, err := s.deps.Signer.SignServerCert(intCertPEM, intKeyPEM, pkiPorts.ServerCertReq{
		CN:   req.CN,
		SANs: req.SANs,
		TTL:  ttl,
	})
	if err != nil {
		return nil, err
	}
	return &signedCertResult{
		certPEM:   certPEM,
		serialHex: fmt.Sprintf("%X", serial),
		notAfter:  time.Now().UTC().Add(ttl),
	}, nil
}
