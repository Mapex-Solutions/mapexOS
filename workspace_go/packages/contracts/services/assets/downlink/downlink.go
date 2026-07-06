package downlink

import (
	"encoding/json"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// Downlink is the platform -> device command path. The Asset MS (the owner of
// device commands) publishes a DownlinkEnvelope on a STATIC per-transport
// subject; each edge server (mapexMQTTBroker / mapexLNS) consumes its subject
// and delivers the command over its own transport. The device identity travels
// in the PAYLOAD, never in the subject — the edge is a dumb transport.
//
// The envelope is command-agnostic: CommandType selects the payload schema, so
// future template-defined downlink commands (script-mounted payloads) reuse the
// same envelope with new command types.

// SubjectMQTTDownlink carries downlink commands for MQTT devices. Consumed by
// mapexMQTTBroker, which resolves the device session from the payload identity
// and publishes to the device's MQTT topic.
var SubjectMQTTDownlink = config.Subject("mqtt", "downlink")

// SubjectLoRaWANDownlink carries downlink commands for LoRaWAN devices.
// Consumed by mapexLNS, which schedules the downlink through the network server.
var SubjectLoRaWANDownlink = config.Subject("lorawan", "downlink")

// StreamMQTTDownlink is the JetStream stream retaining MQTT downlink commands
// (ensured by the broker on boot; WorkQueue retention).
var StreamMQTTDownlink = config.StreamName("MQTT", "DOWNLINK")

// StreamLoRaWANDownlink is the JetStream stream retaining LoRaWAN downlink
// commands (ensured by the LNS on boot; WorkQueue retention).
var StreamLoRaWANDownlink = config.StreamName("LORAWAN", "DOWNLINK")

// CommandTypeOTAUpdate identifies an OTA firmware-update command; its payload
// is an OTAUpdateCommand.
const CommandTypeOTAUpdate = "ota_update"

// DownlinkEnvelope is the wire message on the downlink subjects. Payload is the
// command body selected by CommandType, kept raw so each edge forwards it
// without knowing every command schema.
type DownlinkEnvelope struct {
	OrgID       string          `json:"orgId"`
	AssetUUID   string          `json:"assetUUID"`
	CommandType string          `json:"commandType"`
	Payload     json.RawMessage `json:"payload"`
}

// OTAUpdateCommand is the ota_update payload: everything a device needs to
// download, verify, and apply a firmware image, plus where to report status.
type OTAUpdateCommand struct {
	PlanID        string `json:"planId"`
	ExecutionID   string `json:"executionId"`
	TargetVersion string `json:"targetVersion"`
	DownloadURL   string `json:"downloadUrl"`
	Checksum      string `json:"checksum"`
	Size          int64  `json:"size"`
	// Where the device reports status back: ReportTransport is the device's
	// transport ("mqtt"/"lorawan"/"http"); ReportTarget is the MQTT status topic
	// (or the HTTP status URL for polled devices).
	ReportTransport string `json:"reportTransport"`
	ReportTarget    string `json:"reportTarget"`
}

// NewOTAUpdateEnvelope wraps an OTAUpdateCommand in a DownlinkEnvelope for the
// given device identity.
func NewOTAUpdateEnvelope(orgID, assetUUID string, cmd OTAUpdateCommand) (DownlinkEnvelope, error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return DownlinkEnvelope{}, err
	}
	return DownlinkEnvelope{
		OrgID:       orgID,
		AssetUUID:   assetUUID,
		CommandType: CommandTypeOTAUpdate,
		Payload:     body,
	}, nil
}
