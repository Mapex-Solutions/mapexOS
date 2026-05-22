package nats

import (
	"time"

	"triggers/src/modules/events/application/ports"
)

// NATSExecutor is an infrastructure adapter that implements the TriggerExecutor port.
//
// Following Hexagonal Architecture, this adapter:
// - Lives in the infrastructure layer
// - Implements the application port interface (ports.TriggerExecutor)
// - Contains concrete implementation details (NATS client, network I/O)
// - Has framework dependencies (nats-io/nats.go)
//
// Config schema (from trigger.config.nats):
//
//	{
//	  "url": "nats://localhost:4222",          // NATS server URL(s)
//	  "subject": "events.user.created",        // Subject to publish to
//	  "message": {                             // Message payload
//	    "event": "{{payload.event}}",
//	    "data": "{{payload.data}}"
//	  },
//	  "username": "user",                      // Optional
//	  "password": "pass",                      // Optional
//	  "token": "mytoken",                      // Optional (alternative to user/pass)
//	  "credsFile": "/path/to/creds",           // Optional (for NKey/JWT auth)
//	  "headers": {                             // Optional message headers
//	    "X-Custom": "value"
//	  }
//	}
type NATSExecutor struct {
	connectTimeout time.Duration
	publishTimeout time.Duration
}

// Compile-time check to ensure NATSExecutor implements ports.TriggerExecutor interface
var _ ports.TriggerExecutor = (*NATSExecutor)(nil)
