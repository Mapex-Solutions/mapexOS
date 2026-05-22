package router

import (
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// RouteGroupAdapter implements RouteGroupPort using HTTP calls to Router service.
//
// This adapter follows Hexagonal Architecture by:
//   - Implementing the application port interface
//   - Encapsulating all HTTP/infrastructure details
//   - Keeping the application layer clean from HTTP concerns
type RouteGroupAdapter struct {
	client *httpclient.HTTPClient
}

// RouteGroupResponse represents a single route group from Router service API.
// This is an infrastructure-specific type for HTTP API response parsing.
type RouteGroupResponse struct {
	ID      *string        `json:"id,omitempty"`
	Name    *string        `json:"name,omitempty"`
	Routers []RouterResult `json:"routers,omitempty"`
}

// RouterResult is the minimal shape needed from each router entry in a route
// group — just the kind, for validation purposes.
type RouterResult struct {
	Kind string `json:"kind"`
}

// RouterAPIResponse is the standard API response wrapper from Router service.
// Used when fetching route group names via internal API.
type RouterAPIResponse struct {
	Status int                  `json:"status"`
	Errors []string             `json:"errors"`
	Data   []RouteGroupResponse `json:"data"`
}
