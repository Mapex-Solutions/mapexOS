package router

import (
	"context"
	"fmt"
	"strings"

	"assets/src/modules/assets/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewRouteGroupAdapter creates a new RouteGroupAdapter.
//
// Parameters:
//   - client: Configured HTTP client for Router service communication
//
// Returns:
//   - RouteGroupPort: The port interface implementation
func NewRouteGroupAdapter(client *httpclient.HTTPClient) ports.RouteGroupPort {
	return &RouteGroupAdapter{
		client: client,
	}
}

// Compile-time check to ensure RouteGroupAdapter implements RouteGroupPort interface.
var _ ports.RouteGroupPort = (*RouteGroupAdapter)(nil)

// GetNamesByIds fetches RouteGroup names from Router service by their IDs.
//
// This method calls the Router service's internal API endpoint to retrieve
// RouteGroup details and extracts the names. It gracefully handles errors by
// returning an empty array if the API call fails.
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - ids: Array of RouteGroup IDs to fetch names for
//
// Returns:
//   - []string: Array of RouteGroup names (empty array if fetch fails)
//   - error: nil on success, error if critical failure
func (a *RouteGroupAdapter) GetNamesByIds(ctx context.Context, ids []string) ([]string, error) {
	if len(ids) == 0 {
		return []string{}, nil
	}

	// Build query string: comma-separated IDs
	idsParam := strings.Join(ids, ",")

	// Call Router service internal API
	var apiResponse RouterAPIResponse
	endpoint := "/api/internal/v1/routegroups?ids=" + idsParam
	err := a.client.Get(ctx, endpoint, &apiResponse)

	// If API call fails, return empty array (non-critical operation)
	if err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:HTTPClient] Failed to fetch route group names: %v", err))
		return []string{}, nil
	}

	// Extract names from response data
	names := make([]string, 0, len(apiResponse.Data))
	for _, rg := range apiResponse.Data {
		if rg.Name != nil && *rg.Name != "" {
			names = append(names, *rg.Name)
		}
	}

	return names, nil
}

// GetRouterKindsByIds fetches, for each RouteGroup id, the list of router
// kinds present in the group. Empty ids returns an empty map. If the API
// call fails, returns an empty map (non-critical — caller's validation
// degrades open, backend already enforces via router-side skip).
//
// Parameters:
//   - ctx: Context for controlling cancellation and timeouts
//   - ids: Array of RouteGroup IDs to inspect
//
// Returns:
//   - map[string][]string: groupId -> list of router kinds. Missing groups
//     are absent from the map.
//   - error: nil on success, error on critical failure.
func (a *RouteGroupAdapter) GetRouterKindsByIds(ctx context.Context, ids []string) (map[string][]string, error) {
	if len(ids) == 0 {
		return map[string][]string{}, nil
	}

	idsParam := strings.Join(ids, ",")

	var apiResponse RouterAPIResponse
	endpoint := "/api/internal/v1/routegroups?ids=" + idsParam
	err := a.client.Get(ctx, endpoint, &apiResponse)
	if err != nil {
		logger.Warn(fmt.Sprintf("[INFRA:HTTPClient] Failed to fetch route group kinds: %v", err))
		return map[string][]string{}, nil
	}

	result := make(map[string][]string, len(apiResponse.Data))
	for _, rg := range apiResponse.Data {
		if rg.ID == nil || *rg.ID == "" {
			continue
		}
		kinds := make([]string, 0, len(rg.Routers))
		for _, r := range rg.Routers {
			kinds = append(kinds, r.Kind)
		}
		result[*rg.ID] = kinds
	}

	return result, nil
}
