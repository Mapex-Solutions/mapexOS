package mqtt

import (
	"time"

	"triggers/src/modules/events/application/ports"
)

// MQTTExecutor is an infrastructure adapter that implements the TriggerExecutor port.
//
// Following Hexagonal Architecture, this adapter:
// - Lives in the infrastructure layer
// - Implements the application port interface (ports.TriggerExecutor)
// - Contains concrete implementation details (MQTT client, network I/O)
// - Has framework dependencies (paho.mqtt.golang)
//
// Config schema (from trigger.config.mqtt):
//
//	{
//	  "broker": "mqtt.example.com",
//	  "port": 1883,
//	  "topic": "sensors/{{payload.sensorId}}/data",
//	  "qos": 1,
//	  "username": "user",        // Optional
//	  "password": "pass",        // Optional
//	  "clientId": "client-123",  // Optional (auto-generated if not provided)
//	  "message": {               // Message payload
//	    "temperature": "{{payload.temperature}}",
//	    "timestamp": "{{payload.timestamp}}"
//	  },
//	  "useTLS": true            // Optional, enables TLS connection
//	}
type MQTTExecutor struct {
	connectTimeout time.Duration
	publishTimeout time.Duration
}

// Compile-time check to ensure MQTTExecutor implements ports.TriggerExecutor interface
var _ ports.TriggerExecutor = (*MQTTExecutor)(nil)
