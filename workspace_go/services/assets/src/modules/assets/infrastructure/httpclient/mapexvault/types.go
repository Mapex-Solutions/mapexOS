package mapexvault

import (
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// KEKClient is the concrete HTTP adapter for the LorawanKEKClientPort. Wraps the
// platform's gokit httpclient so transport, X-API-Key injection, and timeout
// policy stay aligned with every other internal-API caller.
type KEKClient struct {
	client *httpclient.HTTPClient
}

// kekWire decodes the on-wire JSON; field names match the Go contract
// KEKResponse json tags exactly.
type kekWire struct {
	Context string `json:"context"`
	Kek     string `json:"kek"`
}
