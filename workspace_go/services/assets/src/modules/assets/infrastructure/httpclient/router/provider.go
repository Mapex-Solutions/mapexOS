package router

import (
	"time"

	"assets/src/modules/assets/application/ports"

	configuration "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// NewRouteGroupPort creates and configures a RouteGroupPort for inter-service communication.
//
// This provider creates a RouteGroupAdapter that implements the RouteGroupPort interface,
// following Hexagonal Architecture principles by returning the port interface instead
// of the concrete implementation.
//
// Configuration values are loaded from the service config:
//   - router_service_url: Base URL of the Router service
//   - internal_api_key: API Key for authenticating internal requests
//
// Returns:
//   - ports.RouteGroupPort: Port interface for RouteGroup lookup operations
func NewRouteGroupPort() ports.RouteGroupPort {
	// Load configuration from environment/config
	routerServiceURL, _ := configuration.GetStringValue("router_service_url")
	apiKey, _ := configuration.GetStringValue("internal_api_key")

	// Create HTTP client configured for Router service
	client := httpclient.New(httpclient.Config{
		BaseURL: routerServiceURL,
		APIKey:  apiKey,
		Timeout: 5 * time.Second, // 5 second timeout for internal API calls
	})

	// Return port interface (not concrete implementation)
	return NewRouteGroupAdapter(client)
}
