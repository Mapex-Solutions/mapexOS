package services

// NodeValidationError represents a validation failure for a specific node.
// Returned by ValidateNodes (one entry per failing node) and surfaced to
// the application layer to assemble a `customErrors.ValidationError`.
type NodeValidationError struct {
	NodeID   string   `json:"nodeId"`
	NodeType string   `json:"nodeType"`
	Errors   []string `json:"errors"`
}
