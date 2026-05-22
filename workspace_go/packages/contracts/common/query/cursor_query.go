package query

import "time"

// CursorQueryDTO contains standard query parameters for cursor-based pagination.
// This is optimized for large datasets where offset-based pagination is inefficient.
//
// Usage:
//
//	type EventsRawQuery struct {
//	    query.CursorQueryDTO
//	    // ... specific fields
//	}
type CursorQueryDTO struct {
	// Cursor is the timestamp to start from (RFC3339 format)
	// If empty, starts from the beginning (newest or oldest based on sort)
	Cursor *time.Time `query:"cursor"`

	// Direction specifies pagination direction: "next" (default) or "prev"
	// - "next": Get items AFTER cursor (older items in DESC order)
	// - "prev": Get items BEFORE cursor (newer items in DESC order)
	Direction *string `query:"direction"`

	// Limit specifies the maximum number of items to return
	// Default: 20, Max: 100
	Limit *int `query:"limit"`

	// SortAsc specifies sort direction
	// - false (default): DESC (newest first)
	// - true: ASC (oldest first)
	SortAsc *bool `query:"sortAsc"`

	// IncludeChildren specifies whether to include child organizations hierarchically
	// When true, uses pathKey range query to include all descendants
	// Default: false
	IncludeChildren *bool `query:"includeChildren"`
}

// GetCursor returns the cursor timestamp or nil if not provided.
func (q *CursorQueryDTO) GetCursor() *time.Time {
	return q.Cursor
}

// GetDirection returns the pagination direction with default value of "next".
func (q *CursorQueryDTO) GetDirection() string {
	if q.Direction == nil || *q.Direction == "" {
		return "next"
	}
	return *q.Direction
}

// GetLimit returns the limit with default value of 20.
// Ensures limit is between 1 and 100.
func (q *CursorQueryDTO) GetLimit() int {
	if q.Limit == nil || *q.Limit < 1 {
		return 20
	}
	if *q.Limit > 100 {
		return 100
	}
	return *q.Limit
}

// GetSortAsc returns whether to sort ascending.
// Default: false (DESC - newest first)
func (q *CursorQueryDTO) GetSortAsc() bool {
	if q.SortAsc == nil {
		return false
	}
	return *q.SortAsc
}

// GetIncludeChildren returns whether to include child organizations.
// Default: false
func (q *CursorQueryDTO) GetIncludeChildren() bool {
	if q.IncludeChildren == nil {
		return false
	}
	return *q.IncludeChildren
}
