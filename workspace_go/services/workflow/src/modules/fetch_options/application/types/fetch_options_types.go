package types

// FetchOptionsItem is a single option returned by the fetch-options proxy.
type FetchOptionsItem struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}
