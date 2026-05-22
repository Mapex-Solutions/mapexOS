package http

import (
	"net/http"

	"triggers/src/modules/events/application/ports"
)

// HTTPExecutor is an infrastructure adapter that implements the TriggerExecutor port.
//
// Following Hexagonal Architecture, this adapter:
// - Lives in the infrastructure layer
// - Implements the application port interface (ports.TriggerExecutor)
// - Contains concrete implementation details (HTTP client, network I/O)
// - Has framework dependencies (net/http, encoding/json)
//
// Config schema:
//
//	{
//	  "endpoint": "https://api.example.com/endpoint",
//	  "method": "POST",           // GET, POST, PUT, PATCH, DELETE
//	  "headers": {
//	    "Content-Type": "application/json",
//	    "Authorization": "Bearer {{token}}"
//	  },
//	  "body": {
//	    "sensor": "{{payload.sensorId}}",
//	    "value": "{{payload.value}}"
//	  },
//	  "timeout": 30               // Optional, in milliseconds
//	}
type HTTPExecutor struct {
	client *http.Client
}

// Compile-time check to ensure HTTPExecutor implements ports.TriggerExecutor interface
var _ ports.TriggerExecutor = (*HTTPExecutor)(nil)
