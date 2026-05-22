package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	mqttPorts "assets/src/modules/mqttcerts/application/ports"
	"assets/src/modules/mqttcerts/domain/entities"
)

var _ mqttPorts.X509SignerPort = (*X509Signer)(nil)

func NewX509Signer() mqttPorts.X509SignerPort {
	return &X509Signer{}
}

// SignDeviceCert produces a fresh ECDSA P-256 keypair, signs a client
// cert with the in-RAM intermediate CA, and returns PEMs + serial +
// SHA-256 fingerprint. Subject CN format `{orgId}:{assetUUID}` matches
// the broker plugin's username parser.
func (s *X509Signer) SignDeviceCert(
	ca *entities.CertificateAuthorityRAM,
	subjectCN string,
	ttlDays int,
) ([]byte, []byte, *big.Int, string, error) {
	intCert, intKey, err := decodeCertAndKey(ca.CertPEM, ca.PrivateKeyPEM)
	if err != nil {
		return nil, nil, nil, "", fmt.Errorf("decode CA: %w", err)
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, "", fmt.Errorf("generate key: %w", err)
	}
	serial, err := randomSerial()
	if err != nil {
		return nil, nil, nil, "", err
	}
	now := time.Now().UTC()
	ttl := time.Duration(ttlDays) * 24 * time.Hour
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: subjectCN},
		NotBefore:    now,
		NotAfter:     now.Add(ttl),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, intCert, &priv.PublicKey, intKey)
	if err != nil {
		return nil, nil, nil, "", fmt.Errorf("sign device cert: %w", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, _ := x509.MarshalECPrivateKey(priv)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	sum := sha256.Sum256(certDER)
	return certPEM, keyPEM, serial, fmt.Sprintf("%x", sum[:]), nil
}

// ---- helpers ----

func randomSerial() (*big.Int, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, limit)
}

func decodeCertAndKey(certPEM, keyPEM []byte) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	cb, _ := pem.Decode(certPEM)
	if cb == nil {
		return nil, nil, fmt.Errorf("decode cert pem")
	}
	cert, err := x509.ParseCertificate(cb.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse cert: %w", err)
	}
	kb, _ := pem.Decode(keyPEM)
	if kb == nil {
		return nil, nil, fmt.Errorf("decode key pem")
	}
	key, err := x509.ParseECPrivateKey(kb.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse key: %w", err)
	}
	return cert, key, nil
}
