package events

import (
	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	dsContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/http_gateway/datasources"
)

// ProcessorExecuteDataSource is the minimized DataSource descriptor carried inside
// a ProcessorExecutePayload. Only the fields required by js-executor to resolve the
// target Asset are included; the remaining DataSource metadata (name, description,
// pathKey) is intentionally omitted because js-executor reads it from the Asset
// cache (source of truth).
//
// Field names and types MUST match the original ad-hoc map shape exactly — any
// change here is a cross-service contract change and requires updating the TS
// counterpart in workspace_js/packages/schemas/.
type ProcessorExecuteDataSource struct {
	// OrgId is the organization identifier used as the Asset cache key prefix
	// ({orgId}/{assetUUID}).
	OrgId *common.ObjectID `json:"orgId"`

	// AssetBind describes how to resolve the Asset UUID from the event body.
	AssetBind *dsContracts.AssetBind `json:"assetBind"`
}

// ProcessorExecutePayload is the cross-service NATS payload published by
// http_gateway on SubjectProcessorJSExecute and consumed by js-executor
// (workspace_js).
//
// This struct mirrors, 1:1, the historical ad-hoc map[string]any payload. It
// exists to make the contract explicit and type-safe on the Go side.
// Reciprocity on the TS side is tracked as a separate follow-up.
type ProcessorExecutePayload struct {
	// SourceType identifies the gateway that produced the event. For
	// http_gateway this is always the literal "http".
	SourceType string `json:"sourceType"`

	// DataSource is the minimized DataSource descriptor required by js-executor.
	DataSource ProcessorExecuteDataSource `json:"dataSource"`

	// Event is the raw event body posted by the external system.
	Event map[string]any `json:"event"`

	// EventTrackerId is the UUID minted at ingestion to correlate the event
	// across the http_gateway -> js-executor -> router -> events pipeline.
	EventTrackerId string `json:"eventTrackerId"`
}
