package types

import "time"

// ExecutionError is a structured error produced by node executors.
// Shared across modules: runtime (produces errors), instances (reads errors), archiver (persists errors).
type ExecutionError struct {
	Code       string    `bson:"code"       json:"code"`
	Message    string    `bson:"message"    json:"message"`
	NodeID     string    `bson:"nodeId"     json:"nodeId"`
	NodeType   string    `bson:"nodeType"   json:"nodeType"`
	Timestamp  time.Time `bson:"timestamp"  json:"timestamp"`
	StackTrace string    `bson:"stackTrace" json:"stackTrace"`
}
