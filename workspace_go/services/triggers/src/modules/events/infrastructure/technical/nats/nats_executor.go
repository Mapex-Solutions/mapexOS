package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"triggers/src/modules/events/application/ports"

	"github.com/nats-io/nats.go"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewNATSExecutor creates a new NATS trigger executor adapter.
func NewNATSExecutor() ports.TriggerExecutor {
	return &NATSExecutor{
		connectTimeout: 10 * time.Second,
		publishTimeout: 5 * time.Second,
	}
}

// Execute publishes a message to a NATS server based on the trigger configuration.
//
// This is the concrete implementation of the ports.TriggerExecutor interface.
//
// Steps:
// 1. Extract NATS config (url, subject, credentials, message)
// 2. Build connection options
// 3. Connect to server
// 4. Publish message
// 5. Flush and close connection
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - config: Trigger configuration (nats field) with placeholders already resolved
//
// Returns:
//   - error: If NATS publish fails
func (e *NATSExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	// Extract nats config
	natsConfig, ok := config["nats"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("NATS trigger config missing 'nats' field")
	}

	// Extract URL (required) — accepts both "url" and "server" field names.
	// The contracts NatsConfig struct uses json:"server", but executor tests use "url".
	// Supporting both ensures compatibility between typed TriggerConfig and manual maps.
	url, ok := natsConfig["url"].(string)
	if !ok || url == "" {
		url, ok = natsConfig["server"].(string)
	}
	if !ok || url == "" {
		return fmt.Errorf("NATS trigger config missing required field 'url' (or 'server')")
	}

	// Extract subject (required)
	subject, ok := natsConfig["subject"].(string)
	if !ok || subject == "" {
		return fmt.Errorf("NATS trigger config missing required field 'subject'")
	}

	// Extract message payload
	message, err := e.extractMessage(natsConfig)
	if err != nil {
		return fmt.Errorf("failed to extract NATS message: %w", err)
	}

	// Build connection options
	opts := []nats.Option{
		nats.Timeout(e.connectTimeout),
		nats.Name("triggers-executor"),
	}

	// Extract authentication options
	username, _ := natsConfig["username"].(string)
	password, _ := natsConfig["password"].(string)
	token, _ := natsConfig["token"].(string)
	credsFile, _ := natsConfig["credsFile"].(string)

	if username != "" && password != "" {
		opts = append(opts, nats.UserInfo(username, password))
	} else if token != "" {
		opts = append(opts, nats.Token(token))
	} else if credsFile != "" {
		opts = append(opts, nats.UserCredentials(credsFile))
	}

	// Check context cancellation before connecting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	logger.Info(fmt.Sprintf("[INFRA:NATSExecutor] Connecting to %s", url))

	// Connect to NATS
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:NATSExecutor] Connection failed: %s", url))
		return fmt.Errorf("NATS connection failed: %w", err)
	}
	defer nc.Close()

	// Check context cancellation before publishing
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	logger.Info(fmt.Sprintf("[INFRA:NATSExecutor] Publishing to subject: %s", subject))

	// Extract headers if present
	headers, hasHeaders := natsConfig["headers"].(map[string]interface{})

	if hasHeaders && len(headers) > 0 {
		// Publish with headers using nats.Msg
		msg := &nats.Msg{
			Subject: subject,
			Data:    message,
			Header:  make(nats.Header),
		}

		for key, value := range headers {
			if strVal, ok := value.(string); ok {
				msg.Header.Add(key, strVal)
			}
		}

		err = nc.PublishMsg(msg)
	} else {
		// Simple publish without headers
		err = nc.Publish(subject, message)
	}

	if err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:NATSExecutor] Publish failed: %s", subject))
		return fmt.Errorf("NATS publish failed: %w", err)
	}

	// Flush to ensure message is sent
	err = nc.FlushTimeout(e.publishTimeout)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:NATSExecutor] Flush failed: %s", subject))
		return fmt.Errorf("NATS flush failed: %w", err)
	}

	logger.Info(fmt.Sprintf("[INFRA:NATSExecutor] Published successfully to subject: %s", subject))
	return nil
}

// GetType returns the trigger type this executor handles.
func (e *NATSExecutor) GetType() string {
	return "nats"
}

// extractMessage extracts and serializes the message payload from config.
func (e *NATSExecutor) extractMessage(config map[string]interface{}) ([]byte, error) {
	message, exists := config["message"]
	if !exists {
		return []byte(""), nil
	}

	switch m := message.(type) {
	case string:
		return []byte(m), nil
	case map[string]interface{}:
		return json.Marshal(m)
	default:
		return json.Marshal(m)
	}
}
