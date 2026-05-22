package websocket

import (
	"time"

	"triggers/src/modules/events/application/ports"
)

// WebSocketExecutor is an infrastructure adapter that implements the TriggerExecutor port.
//
// Following Hexagonal Architecture, this adapter:
// - Lives in the infrastructure layer
// - Implements the application port interface (ports.TriggerExecutor)
// - Contains concrete implementation details (WebSocket client, network I/O)
// - Has framework dependencies (gorilla/websocket)
//
// Config schema (from trigger.config.websocket):
//
//	{
//	  "url": "wss://ws.example.com/events",    // WebSocket server URL
//	  "message": {                             // Message payload
//	    "type": "event",
//	    "data": "{{payload.data}}"
//	  },
//	  "headers": {                             // Optional connection headers
//	    "Authorization": "Bearer token",
//	    "X-Custom": "value"
//	  },
//	  "messageType": "text"                    // Optional: "text" (default) or "binary"
//	}
type WebSocketExecutor struct {
	connectTimeout   time.Duration
	handshakeTimeout time.Duration
	writeTimeout     time.Duration
}

// Compile-time check to ensure WebSocketExecutor implements ports.TriggerExecutor interface
var _ ports.TriggerExecutor = (*WebSocketExecutor)(nil)
