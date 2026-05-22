package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"triggers/src/modules/events/application/ports"

	amqp "github.com/rabbitmq/amqp091-go"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewRabbitMQExecutor creates a new RabbitMQ trigger executor adapter.
func NewRabbitMQExecutor() ports.TriggerExecutor {
	return &RabbitMQExecutor{
		connectTimeout: 10 * time.Second,
		publishTimeout: 5 * time.Second,
	}
}

// Execute publishes a message to a RabbitMQ broker based on the trigger configuration.
//
// This is the concrete implementation of the ports.TriggerExecutor interface.
//
// Steps:
// 1. Extract RabbitMQ config (host, port, credentials, exchange, queue, message)
// 2. Build AMQP URL
// 3. Connect to broker
// 4. Open channel
// 5. Publish message
// 6. Close connection
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - config: Trigger configuration (rabbitmq field) with placeholders already resolved
//
// Returns:
//   - error: If RabbitMQ publish fails
func (e *RabbitMQExecutor) Execute(ctx context.Context, config map[string]interface{}) error {
	// Extract rabbitmq config
	rabbitConfig, ok := config["rabbitmq"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("RabbitMQ trigger config missing 'rabbitmq' field")
	}

	// Extract host (required)
	host, ok := rabbitConfig["host"].(string)
	if !ok || host == "" {
		return fmt.Errorf("RabbitMQ trigger config missing required field 'host'")
	}

	// Extract port (optional, default 5672)
	port := 5672
	if portVal, exists := rabbitConfig["port"]; exists {
		switch p := portVal.(type) {
		case float64:
			port = int(p)
		case int:
			port = p
		}
	}

	// Extract credentials (optional, default guest/guest)
	username := "guest"
	if u, ok := rabbitConfig["username"].(string); ok && u != "" {
		username = u
	}

	password := "guest"
	if p, ok := rabbitConfig["password"].(string); ok && p != "" {
		password = p
	}

	// Extract vhost (optional, default /)
	vhost := "/"
	if v, ok := rabbitConfig["vhost"].(string); ok && v != "" {
		vhost = v
	}

	// Extract exchange settings
	exchange, _ := rabbitConfig["exchange"].(string)
	exchangeType := "direct"
	if et, ok := rabbitConfig["exchangeType"].(string); ok && et != "" {
		exchangeType = et
	}

	// Extract routing key or queue
	routingKey, _ := rabbitConfig["routingKey"].(string)
	queue, _ := rabbitConfig["queue"].(string)

	// Either routingKey or queue must be provided
	if routingKey == "" && queue == "" {
		return fmt.Errorf("RabbitMQ trigger config requires either 'routingKey' or 'queue'")
	}

	// If no exchange, use queue as routing key (direct queue publish)
	if exchange == "" && routingKey == "" {
		routingKey = queue
	}

	// Check TLS setting
	useTLS := false
	if tlsVal, exists := rabbitConfig["useTLS"]; exists {
		if b, ok := tlsVal.(bool); ok {
			useTLS = b
		}
	}

	// Extract message payload
	message, err := e.extractMessage(rabbitConfig)
	if err != nil {
		return fmt.Errorf("failed to extract RabbitMQ message: %w", err)
	}

	// Build AMQP URL
	protocol := "amqp"
	if useTLS {
		protocol = "amqps"
	}
	amqpURL := fmt.Sprintf("%s://%s:%s@%s:%d%s", protocol, username, password, host, port, vhost)

	// Check context cancellation before connecting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	logger.Info(fmt.Sprintf("[INFRA:RabbitMQExecutor] Connecting to %s://%s:%d%s", protocol, host, port, vhost))

	// Connect to RabbitMQ
	conn, err := amqp.DialConfig(amqpURL, amqp.Config{
		Dial: amqp.DefaultDial(e.connectTimeout),
	})
	if err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:RabbitMQExecutor] Connection failed: %s://%s:%d", protocol, host, port))
		return fmt.Errorf("RabbitMQ connection failed: %w", err)
	}
	defer conn.Close()

	// Open channel
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(err, "[INFRA:RabbitMQExecutor] Failed to open channel")
		return fmt.Errorf("RabbitMQ channel open failed: %w", err)
	}
	defer ch.Close()

	// Declare exchange if specified
	if exchange != "" {
		err = ch.ExchangeDeclare(
			exchange,     // name
			exchangeType, // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			logger.Error(err, fmt.Sprintf("[INFRA:RabbitMQExecutor] Failed to declare exchange: %s", exchange))
			return fmt.Errorf("RabbitMQ exchange declare failed: %w", err)
		}
	}

	// Check context cancellation before publishing
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	logger.Info(fmt.Sprintf("[INFRA:RabbitMQExecutor] Publishing to exchange: %s, routing key: %s", exchange, routingKey))

	// Create publish context with timeout
	pubCtx, cancel := context.WithTimeout(ctx, e.publishTimeout)
	defer cancel()

	// Publish message
	err = ch.PublishWithContext(
		pubCtx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[INFRA:RabbitMQExecutor] Publish failed: %s/%s", exchange, routingKey))
		return fmt.Errorf("RabbitMQ publish failed: %w", err)
	}

	logger.Info(fmt.Sprintf("[INFRA:RabbitMQExecutor] Published successfully to exchange: %s, routing key: %s", exchange, routingKey))
	return nil
}

// GetType returns the trigger type this executor handles.
func (e *RabbitMQExecutor) GetType() string {
	return "rabbitmq"
}

// extractMessage extracts and serializes the message payload from config.
func (e *RabbitMQExecutor) extractMessage(config map[string]interface{}) ([]byte, error) {
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
