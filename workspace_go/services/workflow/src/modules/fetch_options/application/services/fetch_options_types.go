package services

import (
	"workflow/src/modules/fetch_options/application/di"
)

// FetchOptionsService handles fetchOptions proxy requests.
// Decrypts credential via vault, loads plugin manifest, makes HTTP call to provider.
type FetchOptionsService struct {
	deps di.FetchOptionsServiceDI
}
