package tiktok

import (
	"encoding/json"
)

// Response represents the common response structure for TikTok Business API
type Response struct {
	Code      *int64          `json:"code,omitempty"`
	Message   *string         `json:"message,omitempty"`
	RequestID *string         `json:"request_id,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Code      int64  `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	return e.Message
}

// ClientConfig represents the configuration for the TikTok Business API client
type ClientConfig struct {
	BaseURL    string
	HTTPClient interface {
		Do(req interface{}) (interface{}, error)
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		BaseURL: "https://business-api.tiktok.com",
	}
}

// PageInfo represents common pagination information used across all API responses
type PageInfo struct {
	Page        int64 `json:"page"`
	PageSize    int64 `json:"page_size"`
	TotalNumber int64 `json:"total_number"`
	TotalPage   int64 `json:"total_page"`
}

// PaginationParams represents common pagination parameters
type PaginationParams struct {
	Page     *int64
	PageSize *int64
}
