package common

// CursorInfo contains cursor pagination metadata for HTTP responses.
// This struct is used to provide pagination information in API responses
// that use cursor-based pagination strategy.
type CursorInfo struct {
	// Next cursor value to fetch the next page.
	// Empty string if there are no more pages forward.
	Next string `json:"next"`

	// Previous cursor value to fetch the previous page.
	// Empty string if there are no pages backward.
	Previous string `json:"previous"`

	// HasNext indicates whether there are more pages forward.
	HasNext bool `json:"hasNext"`

	// HasPrevious indicates whether there are pages backward.
	HasPrevious bool `json:"hasPrevious"`
}
