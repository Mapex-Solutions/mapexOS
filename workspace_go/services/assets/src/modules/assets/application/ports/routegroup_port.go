package ports

import "context"

// RouteGroupPort defines the interface for RouteGroup lookup operations.
//
// This port abstracts the external service dependency (Router service)
// following Hexagonal Architecture principles. The application layer
// depends on this interface, not on concrete HTTP implementations.
//
// Implementations:
//   - RouteGroupAdapter (infrastructure/httpclient): HTTP-based implementation
type RouteGroupPort interface {
	// GetNamesByIds retrieves RouteGroup names by their IDs.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - ids: Array of RouteGroup IDs to lookup
	//
	// Returns:
	//   - []string: Array of RouteGroup names (may be empty if lookup fails)
	//   - error: nil on success, error on critical failure
	GetNamesByIds(ctx context.Context, ids []string) ([]string, error)

	// GetRouterKindsByIds retrieves, for each RouteGroup ID, the list of
	// router kinds present in that group. Used by HealthMonitor validation
	// to enforce that only trigger/workflow routers can receive healthStatus
	// events.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - ids: Array of RouteGroup IDs to inspect
	//
	// Returns:
	//   - map[string][]string: groupId -> list of router kinds. Groups not
	//     found in the Router service are absent from the map (caller may
	//     treat absence as "skip validation for this id").
	//   - error: nil on success, error on critical failure
	GetRouterKindsByIds(ctx context.Context, ids []string) (map[string][]string, error)
}
