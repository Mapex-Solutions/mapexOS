package mapexvault

import "errors"

// EndpointKEKBase is the mapexVault path prefix the assets MS hits during the
// asset service OnMount; the KEK context is appended. The X-API-Key header is
// injected by the shared gokit httpclient.
const EndpointKEKBase = "/internal/kek/"

var (
	ErrUnauthorized = errors.New("mapexvault: unauthorized")
	ErrKEKNotReady  = errors.New("mapexvault: kek not seeded (503)")
	ErrTransient    = errors.New("mapexvault: transient infra error (5xx)")
)
