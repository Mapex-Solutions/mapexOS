package inline

import (
	"context"
	"fmt"
	"strings"
	"time"

	"workflow/src/modules/runtime/domain/entities"
)

/*
 * LOG EXECUTOR
 * Emits a log entry with interpolated message. Outputs ["out"].
 * Supports simple template interpolation: ${state.field}, ${event.field}
 */

// LogExecutor emits a log entry with interpolated message and outputs ["out"].
// Supports simple template interpolation: ${state.field}, ${event.field}.
type LogExecutor struct{}

// NewLogExecutor creates a new LogExecutor.
func NewLogExecutor() entities.NodeExecutor {
	return &LogExecutor{}
}

// NodeType returns "core/log".
func (e *LogExecutor) NodeType() string {
	return "core/log"
}

// Execute interpolates the log message template with state and event data,
// producing a LogEntry in the result.
func (e *LogExecutor) Execute(_ context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.LogNodeConfig)
	if !ok || cfg == nil {
		return &entities.NodeExecutionResult{
			OutputHandles: []string{"out"},
		}, nil
	}

	level := entities.LogLevel(cfg.Level)
	if level == "" {
		level = entities.LogInfo
	}

	message := interpolateMessage(cfg.Message, execCtx.State, execCtx.EventPayload)

	return &entities.NodeExecutionResult{
		OutputHandles: []string{"out"},
		LogEntries: []entities.LogEntry{
			{
				Level:     level,
				Message:   message,
				Timestamp: time.Now(),
				NodeID:    execCtx.NodeID,
				NodeType:  execCtx.NodeType,
			},
		},
	}, nil
}

func interpolateMessage(template string, state, eventPayload map[string]interface{}) string {
	if !strings.Contains(template, "${") {
		return template
	}

	result := template

	for key, val := range state {
		token := "${state." + key + "}"
		if strings.Contains(result, token) {
			result = strings.ReplaceAll(result, token, fmt.Sprintf("%v", val))
		}
	}

	for key, val := range eventPayload {
		token := "${event." + key + "}"
		if strings.Contains(result, token) {
			result = strings.ReplaceAll(result, token, fmt.Sprintf("%v", val))
		}
	}

	return result
}
