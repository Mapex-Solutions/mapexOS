package crypto

// X509Signer is the local-signing adapter for the X509SignerPort.
// Stateless — every Sign call reads the intermediate CA from the
// caller-supplied entities.CertificateAuthorityRAM.
type X509Signer struct{}
