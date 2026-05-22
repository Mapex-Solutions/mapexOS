package ports

import (
	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
)

// HeartbeatEvent is a re-export of the cross-service contract so application
// files can reference the type without reaching into interfaces/message.
// Authoritative definition lives in packages/contracts/services/assets/healthmonitor.
type HeartbeatEvent = contracts.HeartbeatEvent

// PresenceAdvisory is a re-export of the cross-service contract so the
// presence handler can reference the type without reaching into
// packages/contracts directly.
type PresenceAdvisory = contracts.PresenceAdvisory
