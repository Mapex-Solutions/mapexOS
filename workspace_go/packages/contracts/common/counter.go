package common

// CounterResponse is the standard response for counter endpoints.
// Used by all services that implement the /counter endpoint.
type CounterResponse struct {
	Count int64 `json:"count"`
}
