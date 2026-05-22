package constants

// EventType is the DLQ metadata tag for this consumer. Surfaced in the DLQ
// pipeline so operators can filter dead letters by origin.
//
// NOTE: NATS stream/subject constants intentionally live only in the consumer
// package (interfaces/message/consumers/asset_status_save/constants.go) — they
// are consumer-local infrastructure details.
const EventType = "asset-status"

// TableName is the ClickHouse table this module persists into and queries
// from. Used by the repository and by retention-policy TTL application.
const TableName = "asset_status_history"
