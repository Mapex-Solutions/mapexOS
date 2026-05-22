package events

import (
	"fmt"
)

type EvenIdentification struct {
	// Data source id
	Ds *string `json:"ds" validate:"omitempty,mongoid"`

	// Slug for data source
	Sl *string `json:"sl" validate:"omitempty,min=5"`
}

// Transform validates the EvenIdentificationDto and ensures that either 'ds' or 'sl' is provided.
//
// Parameters:
// - None
//
// Returns:
// - error: An error if either 'ds' or 'sl' is not provided, or nil if both are provided.
func (d *EvenIdentification) Transform() error {
	if d.Ds == nil && d.Sl == nil {
		return fmt.Errorf("one of them 'ds' or 'sl' must be provided")
	}

	return nil
}
