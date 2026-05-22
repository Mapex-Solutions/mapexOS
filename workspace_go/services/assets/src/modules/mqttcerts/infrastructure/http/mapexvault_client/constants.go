package mapexvault_client

import "errors"

// EndpointIntermediateCABundle is the mapexVault path the assets MS
// hits during mqttcerts.OnMount. The X-API-Key header is injected by
// the shared gokit httpclient — adapters don't set it manually.
const EndpointIntermediateCABundle = "/internal/pki/intermediate_ca_bundle"

var (
	ErrUnauthorized = errors.New("mapexvault: unauthorized")
	ErrCANotReady   = errors.New("mapexvault: ca not bootstrapped (503)")
	ErrTransient    = errors.New("mapexvault: transient infra error (5xx)")
)
