package dto

// Common response DTOs used across the application

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Legacy - deprecated, use PaginatedResponse[T] from pagination.go instead