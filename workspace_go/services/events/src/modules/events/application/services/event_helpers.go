package services

import (
	ctx "context"
	"encoding/json"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// setTenantContext extracts orgId, pathKey and eventTrackerId from raw JSON
// BEFORE full unmarshal. Ensures DLQ messages always carry tenant context
// even when the full DTO unmarshal fails. The tenantContext struct is
// declared in event_types.go.
func setTenantContext(msg *natsModel.Message) {
	var tc tenantContext
	if err := json.Unmarshal(msg.Data, &tc); err != nil {
		return
	}
	msg.OrgId = tc.OrgId
	msg.PathKey = tc.PathKey
	msg.EventTrackerId = tc.EventTrackerId
}

// getRetentionDays fetches the retention days for a specific table and
// organization. Delegates to RetentionService which handles caching and
// fallback internally.
func (s *EventService) getRetentionDays(c ctx.Context, orgId, tableName string) (uint16, error) {
	return s.deps.RetentionService.GetRetentionDays(c, orgId, tableName)
}

// headersToMetadata converts NATS message headers to a metadata map.
// Single-value headers are stored as strings, multi-value as arrays.
func (s *EventService) headersToMetadata(headers map[string][]string) map[string]interface{} {
	metadata := make(map[string]interface{})
	for k, v := range headers {
		if len(v) == 1 {
			metadata[k] = v[0]
		} else {
			metadata[k] = v
		}
	}
	return metadata
}

// toRawJSON converts a JSON string (from ClickHouse) to json.RawMessage
// so it serializes as a JSON object in the API response, not as an escaped
// string.
func toRawJSON(s string) json.RawMessage {
	if s == "" {
		return nil
	}
	return json.RawMessage(s)
}
