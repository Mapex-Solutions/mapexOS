package services

import (
	"testing"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

/**
 * setTenantContext — unit tests
 */

func TestSetTenantContext(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		wantOrgId     string
		wantPathKey   string
		wantTrackerId string
	}{
		{
			name:          "extracts all tenant fields",
			payload:       `{"orgId":"org-123","pathKey":"000001","eventTrackerId":"tracker-abc"}`,
			wantOrgId:     "org-123",
			wantPathKey:   "000001",
			wantTrackerId: "tracker-abc",
		},
		{
			name:          "extracts tenant fields ignoring unknown fields",
			payload:       `{"orgId":"org-456","pathKey":"000002","eventTrackerId":"tracker-def","execution":{"invalid":true},"extra":999}`,
			wantOrgId:     "org-456",
			wantPathKey:   "000002",
			wantTrackerId: "tracker-def",
		},
		{
			name:          "partial fields — only orgId",
			payload:       `{"orgId":"org-789"}`,
			wantOrgId:     "org-789",
			wantPathKey:   "",
			wantTrackerId: "",
		},
		{
			name:          "invalid JSON — all fields empty",
			payload:       `{invalid`,
			wantOrgId:     "",
			wantPathKey:   "",
			wantTrackerId: "",
		},
		{
			name:          "empty JSON object",
			payload:       `{}`,
			wantOrgId:     "",
			wantPathKey:   "",
			wantTrackerId: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := natsModel.NewTestMessage([]byte(tt.payload), 0, nil)

			setTenantContext(msg)

			if msg.OrgId != tt.wantOrgId {
				t.Fatalf("OrgId: got %q, want %q", msg.OrgId, tt.wantOrgId)
			}
			if msg.PathKey != tt.wantPathKey {
				t.Fatalf("PathKey: got %q, want %q", msg.PathKey, tt.wantPathKey)
			}
			if msg.EventTrackerId != tt.wantTrackerId {
				t.Fatalf("EventTrackerId: got %q, want %q", msg.EventTrackerId, tt.wantTrackerId)
			}
		})
	}
}

/**
 * DLQ Tenant Context per stream — verifies that when the full DTO unmarshal
 * fails, the Message already carries orgId/pathKey/eventTrackerId for DLQ.
 *
 * Each sub-test sends a payload with valid tenant fields but broken data
 * that causes the specific stream's unmarshal/validate to fail.
 */

func TestProcessMessage_DLQTenantContext(t *testing.T) {
	// Zero-value service — failure paths never touch deps
	svc := &EventService{}

	const (
		orgId     = "0000000000000000000aa001"
		pathKey   = "000001"
		trackerID = "abc12345-1111-2222-3333-444455556666"
	)

	// Payloads with valid tenant fields but broken DTO-specific fields:
	// - Raw/JsExec: missing required fields → validator.UnmarshalAndValidate fails
	// - Router/BusinessRule/Trigger/Workflow/EventStore: "created":"NOT-A-DATE" → json.Unmarshal fails on time.Time
	tests := []struct {
		name    string
		payload string
		process func(idx int, msg *natsModel.Message) string // returns action
	}{
		{
			name:    "raw — missing required fields",
			payload: `{"orgId":"` + orgId + `","pathKey":"` + pathKey + `","eventTrackerId":"` + trackerID + `"}`,
			process: func(idx int, msg *natsModel.Message) string {
				return svc.processRawEventMessage(idx, msg).action
			},
		},
		{
			name:    "jsexec — missing required fields",
			payload: `{"orgId":"` + orgId + `","pathKey":"` + pathKey + `","eventTrackerId":"` + trackerID + `"}`,
			process: func(idx int, msg *natsModel.Message) string {
				return svc.processJsExecEventMessage(idx, msg).action
			},
		},
		{
			name:    "router — invalid created time",
			payload: `{"orgId":"` + orgId + `","pathKey":"` + pathKey + `","eventTrackerId":"` + trackerID + `","created":"NOT-A-DATE"}`,
			process: func(idx int, msg *natsModel.Message) string {
				return svc.processRouterEventMessage(idx, msg).action
			},
		},
		{
			name:    "businessrule — invalid created time",
			payload: `{"orgId":"` + orgId + `","pathKey":"` + pathKey + `","eventTrackerId":"` + trackerID + `","created":"NOT-A-DATE"}`,
			process: func(idx int, msg *natsModel.Message) string {
				return svc.processBusinessRuleEventMessage(idx, msg).action
			},
		},
		{
			name:    "trigger — invalid created time",
			payload: `{"orgId":"` + orgId + `","pathKey":"` + pathKey + `","eventTrackerId":"` + trackerID + `","created":"NOT-A-DATE"}`,
			process: func(idx int, msg *natsModel.Message) string {
				return svc.processTriggerEventMessage(idx, msg).action
			},
		},
		{
			name:    "workflow — invalid created time",
			payload: `{"orgId":"` + orgId + `","pathKey":"` + pathKey + `","eventTrackerId":"` + trackerID + `","created":"NOT-A-DATE"}`,
			process: func(idx int, msg *natsModel.Message) string {
				return svc.processWorkflowEventMessage(idx, msg).action
			},
		},
		{
			name:    "eventstore — invalid created time",
			payload: `{"orgId":"` + orgId + `","pathKey":"` + pathKey + `","eventTrackerId":"` + trackerID + `","created":"NOT-A-DATE"}`,
			process: func(idx int, msg *natsModel.Message) string {
				return svc.processEventStoreMessage(idx, msg).action
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := natsModel.NewTestMessage([]byte(tt.payload), 0, nil)

			action := tt.process(0, msg)

			// Must reject (unmarshal/validation failure)
			if action != "reject" {
				t.Fatalf("action: got %q, want %q", action, "reject")
			}

			// Tenant context must be populated even on reject
			if msg.OrgId != orgId {
				t.Fatalf("OrgId: got %q, want %q", msg.OrgId, orgId)
			}
			if msg.PathKey != pathKey {
				t.Fatalf("PathKey: got %q, want %q", msg.PathKey, pathKey)
			}
			if msg.EventTrackerId != trackerID {
				t.Fatalf("EventTrackerId: got %q, want %q", msg.EventTrackerId, trackerID)
			}
		})
	}
}
