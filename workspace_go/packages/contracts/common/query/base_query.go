package query

// BaseQueryDTO contains standard query parameters for all list endpoints.
// These fields are MANDATORY for ANY "/" list endpoint across ALL microservices.
//
// Usage:
//
//	type AssetQueryDTO struct {
//	    query.BaseQueryDTO
//	    // ... specific fields
//	}
type BaseQueryDTO struct {
	// Projection specifies which fields to return (comma-separated)
	// Example: "name,type,status"
	Projection *string `query:"projection"`

	// Page specifies the current page number (1-based)
	// Default: 1
	Page *int `query:"page"`

	// PerPage specifies the number of items per page
	// Default: 20
	PerPage *int `query:"perPage"`

	// Sort specifies the sort order (field:direction)
	// Example: "created:desc", "name:asc", "updated:desc"
	// IMPORTANT: Uses "created" and "updated" (NOT "createdAt")
	// Default: "created:desc"
	Sort *string `query:"sort"`

	// IncludeChildren specifies whether to include child organizations hierarchically
	// When true, uses pathKey range query to include all descendants
	// Default: false
	IncludeChildren *bool `query:"includeChildren"`
}

// GetPage returns the page number with default value of 1.
// Ensures page is always >= 1.
func (q *BaseQueryDTO) GetPage() int {
	if q.Page == nil || *q.Page < 1 {
		return 1
	}
	return *q.Page
}

// GetPerPage returns the items per page with default value of 20.
// Ensures perPage is always >= 1.
func (q *BaseQueryDTO) GetPerPage() int {
	if q.PerPage == nil || *q.PerPage < 1 {
		return 20
	}
	return *q.PerPage
}

// GetSort returns the sort order with default value of "created:desc".
//
// IMPORTANT: Sort field names use "created" and "updated" (with capital letter),
// NOT "createdAt" or "updatedAt".
//
// Examples:
//   - "created:desc" (default)
//   - "created:asc"
//   - "updated:desc"
//   - "name:asc"
func (q *BaseQueryDTO) GetSort() string {
	if q.Sort == nil || *q.Sort == "" {
		return "created:desc"
	}
	return *q.Sort
}

// GetIncludeChildren returns whether to include child organizations.
// Default: false
func (q *BaseQueryDTO) GetIncludeChildren() bool {
	if q.IncludeChildren == nil {
		return false
	}
	return *q.IncludeChildren
}
