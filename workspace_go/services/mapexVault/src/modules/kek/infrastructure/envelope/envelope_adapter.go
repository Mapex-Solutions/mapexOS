package envelope

import (
	kekPorts "mapexVault/src/modules/kek/application/ports"

	envelopeUtil "github.com/Mapex-Solutions/mapexGoKit/utils/envelope"
)

// EnvelopeAdapter wraps the shared *envelope.EnvelopeService behind the kek
// module's EnvelopePort. Hides the concrete driver from the application
// layer. The kek module only decrypts: KEKs are envelope-encrypted at
// seed time by the mongodb-init container.
type EnvelopeAdapter struct {
	svc *envelopeUtil.EnvelopeService
}

var _ kekPorts.EnvelopePort = (*EnvelopeAdapter)(nil)

func NewEnvelopeAdapter(svc *envelopeUtil.EnvelopeService) kekPorts.EnvelopePort {
	return &EnvelopeAdapter{svc: svc}
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
