package brixhub

import (
	"os"
)

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() (*Client, error) {
	apiKey := os.Getenv("BRIXHUB_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyRequired
	}

	opts := []ClientOption{}
	
	if ua := os.Getenv("BRIXHUB_USER_AGENT"); ua != "" {
		opts = append(opts, WithUserAgent(ua))
	}
	
	if baseURL := os.Getenv("BRIXHUB_BASE_URL"); baseURL != "" {
		opts = append(opts, WithBaseURL(baseURL))
	}

	return NewClient(apiKey, opts...)
}

// ErrAPIKeyRequired is returned when API key is missing
var ErrAPIKeyRequired = NewAPIError("API key is required. Set BRIXHUB_API_KEY environment variable")