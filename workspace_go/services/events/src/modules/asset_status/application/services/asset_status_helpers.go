package services

import (
	"encoding/json"
	"fmt"
	"time"

	"events/src/modules/asset_status/domain/entities"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// setTenantContext extracts orgId/pathKey/eventTrackerId from the raw JSON
// before the full unmarshal so DLQ routing keeps tenant context even on
// validation failures. Mirrors events module's helper of the same name.
func setTenantContext(msg *natsModel.Message) {
	var tc struct {
		OrgId          string `json:"orgId"`
		PathKey        string `json:"pathKey"`
		EventTrackerId string `json:"eventTrackerId"`
	}
	if err := json.Unmarshal(msg.Data, &tc); err != nil {
		return
	}
	msg.OrgId = tc.OrgId
	msg.PathKey = tc.PathKey
	msg.EventTrackerId = tc.EventTrackerId
}

// parsePersistencePayload decodes a single NATS message body into the
// wire-format struct and validates that required fields are present.
// Returns an empty payload and a non-nil error on any failure — the handler
// translates that into msg.Reject(...) for DLQ.
func parsePersistencePayload(data []byte) (persistencePayload, error) {
	var p persistencePayload
	if err := json.Unmarshal(data, &p); err != nil {
		return p, fmt.Errorf("invalid JSON: %w", err)
	}
	if p.OrgId == "" {
		return p, fmt.Errorf("missing required field 'orgId'")
	}
	if p.AssetUUID == "" {
		return p, fmt.Errorf("missing required field 'assetUUID'")
	}
	if p.EventId == "" {
		return p, fmt.Errorf("missing required field 'eventId'")
	}
	if p.EventType == "" {
		return p, fmt.Errorf("missing required field 'eventType'")
	}
	if p.Created == "" {
		return p, fmt.Errorf("missing required field 'created'")
	}
	return p, nil
}

// payloadToEntity builds the ClickHouse row from the wire payload. Parses
// created (RFC3339Nano) — fatal on parse failure so the caller can Reject.
func payloadToEntity(p persistencePayload) (*entities.AssetStatusEvent, error) {
	created, err := time.Parse(time.RFC3339Nano, p.Created)
	if err != nil {
		return nil, fmt.Errorf("invalid created timestamp %q: %w", p.Created, err)
	}
	return &entities.AssetStatusEvent{
		Created:          created.UTC(),
		OrgId:            p.OrgId,
		PathKey:          p.PathKey,
		AssetUUID:        p.AssetUUID,
		AssetName:        p.AssetName,
		EventId:          p.EventId,
		EventType:        p.EventType,
		LastSeenAt:       p.LastSeenAt,
		ThresholdMinutes: p.ThresholdMinutes,
		MissCount:        p.MissCount,
	}, nil
}
