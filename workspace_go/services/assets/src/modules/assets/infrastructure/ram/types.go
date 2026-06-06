package ram

import (
	"sync/atomic"

	envelopeUtil "github.com/Mapex-Solutions/mapexGoKit/utils/envelope"
)

// InMemoryKEKCipher holds the envelope service built from the LoRaWAN device-keys
// KEK fetched at OnMount. atomic.Pointer keeps load/store lock-free and
// hot-swap-safe; the plaintext KEK never leaves this process boundary.
type InMemoryKEKCipher struct {
	svc atomic.Pointer[envelopeUtil.EnvelopeService]
}
