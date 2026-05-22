package evaluators

import "errors"

/*
 * VALUE RESOLUTION AND TYPE CONVERSION ERRORS
 */

var (
	// ErrInvalidFieldValue is returned when a FieldValue is invalid or nil
	ErrInvalidFieldValue = errors.New("invalid field value configuration")

	// ErrFieldNotFound is returned when a field path doesn't exist in the source
	ErrFieldNotFound = errors.New("field not found in source")

	// ErrInvalidSource is returned when the source type is unknown or source is nil
	ErrInvalidSource = errors.New("invalid value source")

	// ErrInvalidPath is returned when a field path is malformed
	ErrInvalidPath = errors.New("invalid field path")

	// ErrInvalidNodeID is returned when a node_output source has no NodeID
	ErrInvalidNodeID = errors.New("node_output source requires nodeId")

	// ErrInvalidGroupItem is returned when a ConditionGroupItem has unsupported data type
	ErrInvalidGroupItem = errors.New("invalid condition group item data type")

	// ErrEmptyGroup is returned when a ConditionGroup has no items
	ErrEmptyGroup = errors.New("condition group has no items")

	// ErrOperatorNotFound is returned when an operator is not found in the registry
	ErrOperatorNotFound = errors.New("operator not found in registry")
)
