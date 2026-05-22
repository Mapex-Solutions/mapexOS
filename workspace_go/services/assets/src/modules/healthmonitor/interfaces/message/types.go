package message

import (
	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
)

// HeartbeatEvent is a cross-service contract, authoritatively defined in
// packages/contracts/services/assets/healthmonitor. Re-aliased here so
// existing intra-module imports (hmMessage.HeartbeatEvent) keep compiling
// without redefining the wire format.
type HeartbeatEvent = contracts.HeartbeatEvent
