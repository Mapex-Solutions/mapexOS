package rabbitmq

import (
	"time"

	"triggers/src/modules/events/application/ports"
)

// RabbitMQExecutor is an infrastructure adapter that implements the TriggerExecutor port.
//
// Following Hexagonal Architecture, this adapter:
// - Lives in the infrastructure layer
// - Implements the application port interface (ports.TriggerExecutor)
// - Contains concrete implementation details (AMQP client, network I/O)
// - Has framework dependencies (rabbitmq/amqp091-go)
//
// Config schema (from trigger.config.rabbitmq):
//
//	{
//	  "host": "rabbitmq.example.com",
//	  "port": 5672,
//	  "username": "guest",          // Optional (default: guest)
//	  "password": "guest",          // Optional (default: guest)
//	  "vhost": "/",                 // Optional (default: /)
//	  "exchange": "my-exchange",    // Optional (direct publish to queue if empty)
//	  "exchangeType": "direct",     // Optional (direct, fanout, topic, headers)
//	  "routingKey": "my-key",       // Routing key for exchange or queue name
//	  "queue": "my-queue",          // Queue name (if publishing directly to queue)
//	  "message": {                  // Message payload
//	    "data": "{{payload.data}}",
//	    "timestamp": "{{payload.timestamp}}"
//	  },
//	  "useTLS": true                // Optional, enables TLS connection
//	}
type RabbitMQExecutor struct {
	connectTimeout time.Duration
	publishTimeout time.Duration
}

// Compile-time check to ensure RabbitMQExecutor implements ports.TriggerExecutor interface
var _ ports.TriggerExecutor = (*RabbitMQExecutor)(nil)
