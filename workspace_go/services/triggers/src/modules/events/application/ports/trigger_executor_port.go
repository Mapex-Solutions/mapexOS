package ports

import (
	"context"
)

// TriggerExecutor defines the contract for trigger execution adapters.
//
// Following Hexagonal Architecture, this is a PORT (interface) that defines
// what the application needs from infrastructure adapters.
//
// Infrastructure layer provides concrete implementations (ADAPTERS) for different
// protocols: HTTP, Email, Slack, Teams, MQTT, RabbitMQ, NATS, WebSocket.
//
// Each trigger type implements this interface with its specific execution logic.
type TriggerExecutor interface {
	// Execute performs the trigger action with the resolved configuration.
	//
	// The application layer calls this method after:
	// Fetching trigger configuration
	// Resolving placeholders in the config
	//
	// Infrastructure implementations handle protocol-specific details:
	// - HTTP: Making HTTP requests
	// - Email: Sending SMTP emails
	// - Slack/Teams: Posting webhook messages
	// - MQTT/RabbitMQ/NATS: Publishing messages
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - config: Trigger configuration with placeholders already resolved
	//
	// Returns:
	//   - error: If execution fails
	Execute(ctx context.Context, config map[string]interface{}) error

	// GetType returns the trigger type this executor handles.
	//
	// Examples: "http", "email", "teams", "slack", "mqtt", "rabbitmq", "nats", "websocket"
	//
	// This is used by the registry to map trigger types to executors.
	GetType() string
}

// ExecutorRegistry manages the mapping of trigger types to executor implementations.
//
// This is also a PORT interface that infrastructure will implement.
//
// The registry enables dynamic executor selection based on trigger type,
// following the Factory pattern.
//
// Infrastructure layer provides the concrete registry implementation that
// knows about all available executor adapters.
type ExecutorRegistry interface {
	// GetExecutor retrieves an executor by trigger type.
	//
	// Parameters:
	//   - triggerType: The type of trigger (http, email, teams, slack, etc.)
	//
	// Returns:
	//   - TriggerExecutor: The executor for the trigger type
	//   - bool: True if executor exists, false otherwise
	GetExecutor(triggerType string) (TriggerExecutor, bool)
}
