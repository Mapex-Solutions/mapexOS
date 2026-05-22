package types

// StandardResponse represents standard API response
type StandardResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Errors  []string `json:"errors"`
	Message string   `json:"message,omitempty"`
	Success bool     `json:"success"`
}

// PaginatedResponse represents paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalItems int         `json:"totalItems"`
	TotalPages int         `json:"totalPages"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}
