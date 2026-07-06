package ports

import (
	"context"

	downlinkContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/downlink"
)

// OTACommandPayload is the OTA command sent to a device via the edge (or served
// to an HTTP device when it polls for its job). Canonical in the shared
// downlink contract (consumed by mapexMQTTBroker / mapexLNS); re-aliased here.
type OTACommandPayload = downlinkContract.OTAUpdateCommand

// EdgeDispatchPort publishes an OTA command to a device over our own edge
// transports. Implemented by an adapter that publishes a DownlinkEnvelope to
// the STATIC subjects `mqtt.downlink` / `lorawan.downlink` — the org +
// assetUUID travel in the message payload, NEVER in the subject (the edge is a
// dumb transport).
type EdgeDispatchPort interface {
	Dispatch(ctx context.Context, protocol, orgID, assetUUID string, payload OTACommandPayload) error
}
