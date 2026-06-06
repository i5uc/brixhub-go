package brixhub

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorType represents the type of API error
type ErrorType string

const (
	// Error types
	ErrBadRequest       ErrorType = "bad_request"
	ErrUnauthorized     ErrorType = "unauthorized"
	ErrExpired          ErrorType = "expired"
	ErrPlanLimited      ErrorType = "plan_limited"
	ErrQuotaExceeded    ErrorType = "quota_exceeded"
	ErrRateLimited      ErrorType = "rate_limited"
	ErrInternal         ErrorType = "internal"
	ErrServiceUnavailable ErrorType = "service_unavailable"
)

// APIError represents a BrixHub API error
type APIError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Status  int       `json:"status"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("brixhub API error [%d - %s]: %s", e.Status, e.Type, e.Message)
}

// IsUnauthorized returns true if the error is an authentication error
func (e *APIError) IsUnauthorized() bool {
	return e.Type == ErrUnauthorized || e.Type == ErrExpired
}

// IsRateLimit returns true if the error is a rate limit error
func (e *APIError) IsRateLimit() bool {
	return e.Type == ErrQuotaExceeded || e.Type == ErrRateLimited
}

// IsPlanLimited returns true if the feature requires a higher plan
func (e *APIError) IsPlanLimited() bool {
	return e.Type == ErrPlanLimited
}

// parseError parses an error response from the API
func parseError(resp *http.Response, body []byte) error {
	apiErr := &APIError{
		Status: resp.StatusCode,
	}
	
	// Try to parse the error type from body
	var errResp struct {
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
		Message string `json:"message"`
	}
	
	if err := json.Unmarshal(body, &errResp); err == nil {
		if errResp.Error.Type != "" {
			apiErr.Type = ErrorType(errResp.Error.Type)
			apiErr.Message = errResp.Error.Message
		} else {
			apiErr.Message = errResp.Message
			// Infer type from status code
			switch resp.StatusCode {
			case 400:
				apiErr.Type = ErrBadRequest
			case 401:
				apiErr.Type = ErrUnauthorized
			case 403:
				apiErr.Type = ErrPlanLimited
			case 429:
				apiErr.Type = ErrRateLimited
			case 500:
				apiErr.Type = ErrInternal
			case 503:
				apiErr.Type = ErrServiceUnavailable
			}
		}
	} else {
		apiErr.Message = string(body)
	}
	
	return apiErr
}