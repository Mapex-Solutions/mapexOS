package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	pkiDi "mapexVault/src/modules/pki/application/di"
	pkiDtos "mapexVault/src/modules/pki/application/dtos"
	pkiPorts "mapexVault/src/modules/pki/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time port check.
var _ pkiPorts.PkiServicePort = (*PkiService)(nil)

// ErrCANotBootstrapped is returned when an endpoint is hit before the
// Mongo seed has populated pkiCertificateAuthorities. The HTTP layer
// maps it to 503 so the Assets MS bootstrap retry loop keeps trying
// until the operator runs generate-pki and the mongodb-init container
// loads the seed JSON.
var ErrCANotBootstrapped = errors.New("ca not bootstrapped")

// New constructs the PKI service. Returns the port interface — the
// concrete struct stays unexposed.
func New(deps pkiDi.PkiServiceDI) pkiPorts.PkiServicePort {
	return &PkiService{deps: deps}
}

// GetIntermediateCABundle returns the intermediate cert + decrypted
// private key. Used by Assets MS at boot.
func (s *PkiService) GetIntermediateCABundle(ctx context.Context) (*pkiDtos.IntermediateCABundleResponse, error) {
	logger.Info("[SERVICE:Pki] GetIntermediateCABundle: loading entity from Mongo")
	entity, err := s.loadIntermediateEntity(ctx)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Pki] GetIntermediateCABundle: load entity failed err=%v", err))
		return nil, fmt.Errorf("load intermediate: %w", err)
	}
	logger.Info(fmt.Sprintf("[SERVICE:Pki] GetIntermediateCABundle: entity loaded subjectCN=%s notAfter=%s", entity.SubjectCN, entity.NotAfter.Format("2006-01-02")))
	decryptedKey, err := s.decryptIntermediatePrivateKey(entity)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Pki] GetIntermediateCABundle: decrypt failed err=%v", err))
		return nil, fmt.Errorf("decrypt intermediate key: %w", err)
	}
	logger.Info("[SERVICE:Pki] GetIntermediateCABundle: decrypt ok, returning bundle")
	return s.entityToBundleDTO(entity, decryptedKey), nil
}

// GetCAChain returns root + intermediate cert PEM concatenated.
// Public material — safe to expose anywhere.
func (s *PkiService) GetCAChain(ctx context.Context) (*pkiDtos.CAChainResponse, error) {
	root, err := s.loadRootEntity(ctx)
	if err != nil {
		return nil, fmt.Errorf("load root: %w", err)
	}
	intermediate, err := s.loadIntermediateEntity(ctx)
	if err != nil {
		return nil, fmt.Errorf("load intermediate: %w", err)
	}
	return s.entityToCAChainDTO(root, intermediate), nil
}

// SignServerCert decrypts the intermediate on demand, signs a server
// cert for the requested CN+SANs, returns the cert. Stateless —
// plaintext key material exists in process memory only for this call.
func (s *PkiService) SignServerCert(ctx context.Context, req *pkiDtos.SignServerRequest) (*pkiDtos.SignServerResponse, error) {
	if err := s.validateSignRequest(req); err != nil {
		return nil, err
	}
	intCertPEM, intKeyPEM, err := s.loadDecryptedIntermediateMaterial(ctx)
	if err != nil {
		return nil, fmt.Errorf("load intermediate material: %w", err)
	}
	result, err := s.runSigner(intCertPEM, intKeyPEM, req)
	if err != nil {
		return nil, fmt.Errorf("sign server cert: %w", err)
	}
	return &pkiDtos.SignServerResponse{
		CertPEM:   result.certPEM,
		SerialHex: result.serialHex,
		NotAfter:  time.Now().UTC().Add(time.Duration(req.TTLDays) * 24 * time.Hour),
	}, nil
}
