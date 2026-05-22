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
	"net"
	"time"

	pkiPorts "mapexVault/src/modules/pki/application/ports"
)

// X509Signer is the concrete crypto adapter. Stateless: every method
// generates fresh keys and signs ad hoc. Used by the mapexVault pki
// service for on-demand operations (intermediate CA generation when
// bootstrapping; server cert signing for the broker provisioning
// script).
type X509Signer struct{}

var _ pkiPorts.X509SignerPort = (*X509Signer)(nil)

func NewX509Signer() pkiPorts.X509SignerPort {
	return &X509Signer{}
}

// GenerateRootCA produces a self-signed ECDSA P-256 root CA cert + key.
func (s *X509Signer) GenerateRootCA(req pkiPorts.RootCAReq) ([]byte, []byte, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate root key: %w", err)
	}
	serial, err := randomSerial()
	if err != nil {
		return nil, nil, err
	}
	now := time.Now().UTC()
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: req.SubjectCN},
		NotBefore:             now,
		NotAfter:              now.Add(req.TTL),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("sign root cert: %w", err)
	}
	return encodeCertPEM(der), encodeECPrivKeyPEM(priv), nil
}

// GenerateIntermediateCA signs an intermediate CA with the supplied root.
func (s *X509Signer) GenerateIntermediateCA(rootCertPEM, rootKeyPEM []byte, req pkiPorts.IntermediateCAReq) ([]byte, []byte, error) {
	rootCert, rootKey, err := decodeCertAndKey(rootCertPEM, rootKeyPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("decode root: %w", err)
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate intermediate key: %w", err)
	}
	serial, err := randomSerial()
	if err != nil {
		return nil, nil, err
	}
	now := time.Now().UTC()
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: req.SubjectCN},
		NotBefore:             now,
		NotAfter:              now.Add(req.TTL),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, rootCert, &priv.PublicKey, rootKey)
	if err != nil {
		return nil, nil, fmt.Errorf("sign intermediate cert: %w", err)
	}
	return encodeCertPEM(der), encodeECPrivKeyPEM(priv), nil
}

// SignServerCert produces a server cert signed by the supplied intermediate.
func (s *X509Signer) SignServerCert(intCertPEM, intKeyPEM []byte, req pkiPorts.ServerCertReq) ([]byte, *big.Int, error) {
	intCert, intKey, err := decodeCertAndKey(intCertPEM, intKeyPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("decode intermediate: %w", err)
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate server key: %w", err)
	}
	serial, err := randomSerial()
	if err != nil {
		return nil, nil, err
	}
	now := time.Now().UTC()
	dnsNames, ipAddrs := splitSANs(req.SANs)
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: req.CN},
		NotBefore:    now,
		NotAfter:     now.Add(req.TTL),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     dnsNames,
		IPAddresses:  ipAddrs,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, intCert, &priv.PublicKey, intKey)
	if err != nil {
		return nil, nil, fmt.Errorf("sign server cert: %w", err)
	}
	return encodeCertPEM(der), serial, nil
}

// CertFingerprintSHA256 exposes a fingerprint helper for the application
// layer to populate the entity's Fingerprint field.
func CertFingerprintSHA256(certPEM []byte) string {
	b, _ := pem.Decode(certPEM)
	if b == nil {
		return ""
	}
	sum := sha256.Sum256(b.Bytes)
	return fmt.Sprintf("%x", sum[:])
}

// ---- helpers ----

func randomSerial() (*big.Int, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, limit)
}

func encodeCertPEM(der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

func encodeECPrivKeyPEM(priv *ecdsa.PrivateKey) []byte {
	der, _ := x509.MarshalECPrivateKey(priv)
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der})
}

func decodeCertAndKey(certPEM, keyPEM []byte) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	cb, _ := pem.Decode(certPEM)
	if cb == nil {
		return nil, nil, fmt.Errorf("decode cert pem: empty block")
	}
	cert, err := x509.ParseCertificate(cb.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse cert: %w", err)
	}
	kb, _ := pem.Decode(keyPEM)
	if kb == nil {
		return nil, nil, fmt.Errorf("decode key pem: empty block")
	}
	key, err := x509.ParseECPrivateKey(kb.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse key: %w", err)
	}
	return cert, key, nil
}

func splitSANs(sans []string) ([]string, []net.IP) {
	var dns []string
	var ips []net.IP
	for _, s := range sans {
		if ip := net.ParseIP(s); ip != nil {
			ips = append(ips, ip)
			continue
		}
		dns = append(dns, s)
	}
	return dns, ips
}
