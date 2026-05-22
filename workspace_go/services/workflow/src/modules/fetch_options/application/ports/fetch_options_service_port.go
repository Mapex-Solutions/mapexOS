package ports

import (
	"context"

	"workflow/src/modules/fetch_options/application/types"
)

// FetchOptionsServicePort defines the contract for fetchOptions proxy operations.
type FetchOptionsServicePort interface {
	// FetchOptions decrypts credential via vault, loads manifest, makes HTTP call to provider.
	FetchOptions(ctx context.Context, credentialId, pluginId, resourceKey string, dependsOn map[string]string) ([]types.FetchOptionsItem, error)
}
