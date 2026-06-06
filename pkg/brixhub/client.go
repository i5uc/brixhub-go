package brixhub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Client is the BrixHub API client
type Client struct {
	apiKey     string
	baseURL    string
	userAgent  string
	httpClient *http.Client
	
	// Rate limit info from last response
	LastRateLimit *RateLimit
}

// ClientOption configures the Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL (for testing)
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithUserAgent sets a custom User-Agent
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// NewClient creates a new BrixHub API client
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	
	c := &Client{
		apiKey:    apiKey,
		baseURL:   BaseURL,
		userAgent: DefaultUserAgent,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	
	for _, opt := range opts {
		opt(c)
	}
	
	return c, nil
}

// doRequest performs an HTTP request and handles the response
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}
	
	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse rate limit headers
	c.LastRateLimit = &RateLimit{}
	if limit := resp.Header.Get("X-RateLimit-Limit-Day"); limit != "" {
		c.LastRateLimit.LimitDay, _ = strconv.Atoi(limit)
	}
	if remaining := resp.Header.Get("X-RateLimit-Remaining-Day"); remaining != "" {
		c.LastRateLimit.RemainingDay, _ = strconv.Atoi(remaining)
	}
	if limitMin := resp.Header.Get("X-RateLimit-Limit-Min"); limitMin != "" {
		c.LastRateLimit.LimitMin, _ = strconv.Atoi(limitMin)
	}
	
	// Check for error status codes
	if resp.StatusCode >= 400 {
		return parseError(resp, respBody)
	}
	
	// Parse successful response
	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	
	if apiResp.Status >= 400 {
		return &APIError{
			Type:    ErrInternal,
			Message: apiResp.Message,
			Status:  apiResp.Status,
		}
	}
	
	// Unmarshal data into result
	if result != nil && apiResp.Data != nil {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	}
	
	return nil
}

// GetRemainingQuota returns the remaining daily quota from the last response
func (c *Client) GetRemainingQuota() int {
	if c.LastRateLimit != nil {
		return c.LastRateLimit.RemainingDay
	}
	return -1
}

// Health checks the API health status (no authentication required)
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	// Health endpoint doesn't require authentication
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}
	
	return &health, nil
}