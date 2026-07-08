package ports

import (
	"context"

	v1lists "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/lists"
)

// ListsClientPort resolves a classification slug into an org-scoped list id via
// mapexIam's internal resolve endpoint, creating the list when it does not yet
// exist. The install calls it once per classification level.
type ListsClientPort interface {
	// Resolve returns the id of the org-scoped list matching the request,
	// creating it (never global/system) when absent.
	Resolve(ctx context.Context, req v1lists.ListResolveRequest) (string, error)
}
