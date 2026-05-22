package envelope

import (
	pkiPorts "mapexVault/src/modules/pki/application/ports"

	envelopeUtil "github.com/Mapex-Solutions/mapexGoKit/utils/envelope"
)

// EnvelopeAdapter wraps the existing *envelope.EnvelopeService behind
// the pki module's EnvelopePort. Hides the concrete driver type from
// the application layer per /go-arch §6.
type EnvelopeAdapter struct {
	svc *envelopeUtil.EnvelopeService
}

var _ pkiPorts.EnvelopePort = (*EnvelopeAdapter)(nil)

func NewEnvelopeAdapter(svc *envelopeUtil.EnvelopeService) pkiPorts.EnvelopePort {
	return &EnvelopeAdapter{svc: svc}
}

func (a *EnvelopeAdapter) Encrypt(plaintext []byte) ([]byte, []byte, []byte, []byte, error) {
	env, err := a.svc.Encrypt(plaintext)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return env.EncryptedDEK, env.DEKNonce, env.EncryptedData, env.DataNonce, nil
}

func (a *EnvelopeAdapter) Decrypt(encryptedDEK, dekNonce, encryptedData, dataNonce []byte) ([]byte, error) {
	env := &envelopeUtil.EncryptedEnvelope{
		EncryptedDEK:  encryptedDEK,
		DEKNonce:      dekNonce,
		EncryptedData: encryptedData,
		DataNonce:     dataNonce,
	}
	return a.svc.Decrypt(env)
}
