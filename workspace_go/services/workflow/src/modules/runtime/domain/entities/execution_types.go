package entities

import (
	"time"

	sharedTypes "workflow/src/shared/types"
)

// PathEntry records one step in the execution path (DAG visualization).
// Each node visited during execution produces a PathEntry with timing and status info.
type PathEntry struct {
	NodeID       string     `bson:"nodeId"`
	NodeType     string     `bson:"nodeType"`
	Status       string     `bson:"status"`
	EnteredAt    time.Time  `bson:"enteredAt"`
	ExitedAt     *time.Time `bson:"exitedAt,omitempty"`
	DurationMs   int64      `bson:"durationMs"`
	OutputHandle string     `bson:"outputHandle,omitempty"`
	Error        *string    `bson:"error,omitempty"`
}

// ExecutionError is an alias to the shared type.
// Defined in src/shared/types/ because it is used by multiple modules (runtime, instances, archiver).
type ExecutionError = sharedTypes.ExecutionError

// LogLevel represents the severity level for step log entries.
type LogLevel string

const (
	// LogDebug is the lowest severity level.
	LogDebug LogLevel = "debug"
	// LogInfo is the default severity level.
	LogInfo LogLevel = "info"
	// LogWarn indicates a potential issue.
	LogWarn LogLevel = "warn"
	// LogError indicates a node execution error.
	LogError LogLevel = "error"
)

// LogEntry is a structured log entry emitted by a node executor.
// Pure domain entity — kept in-memory on NodeExecutionResult.LogEntries.
// Wire-format (NATS/ClickHouse) is owned by the contracts package; this
// entity carries only persistence/in-memory shape.
type LogEntry struct {
	Level     LogLevel               `bson:"level"`
	Message   string                 `bson:"message"`
	Timestamp time.Time              `bson:"timestamp"`
	NodeID    string                 `bson:"nodeId"`
	NodeType  string                 `bson:"nodeType"`
	Data      map[string]interface{} `bson:"data,omitempty"`
}
