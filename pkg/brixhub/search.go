package brixhub

import (
	"context"
	"fmt"
	"net/url"
)

// Search performs a multi-criteria search
// Returns merged profiles with confidence scores
func (c *Client) Search(ctx context.Context, req *SearchRequest) (*SearchResults, *Meta, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("search request cannot be nil")
	}
	
	var result struct {
		Results []Profile `json:"results"`
	}
	
	var meta Meta
	
	// Use a wrapper to capture meta
	err := c.doRequest(ctx, "POST", "/search", req, &result)
	if err != nil {
		return nil, nil, err
	}
	
	return &SearchResults{Results: result.Results}, &meta, nil
}

// LookupEmail performs a reverse lookup by email
func (c *Client) LookupEmail(ctx context.Context, email string) (*Profile, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	
	endpoint := fmt.Sprintf("/lookup/email/%s", url.PathEscape(email))
	
	var result struct {
		Results []Profile `json:"results"`
	}
	
	if err := c.doRequest(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	if len(result.Results) == 0 {
		return nil, nil
	}
	
	return &result.Results[0], nil
}

// LookupPhone performs a reverse lookup by phone number
// Supports all French formats: 06..., +33..., 0033...
func (c *Client) LookupPhone(ctx context.Context, phone string) (*Profile, error) {
	if phone == "" {
		return nil, fmt.Errorf("phone is required")
	}
	
	endpoint := fmt.Sprintf("/lookup/phone/%s", url.PathEscape(phone))
	
	var result struct {
		Results []Profile `json:"results"`
	}
	
	if err := c.doRequest(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	if len(result.Results) == 0 {
		return nil, nil
	}
	
	return &result.Results[0], nil
}

// LookupIBAN performs a reverse lookup by IBAN
// Returns account holder with contact details and BIC
func (c *Client) LookupIBAN(ctx context.Context, iban string) (*Profile, error) {
	if iban == "" {
		return nil, fmt.Errorf("IBAN is required")
	}
	
	endpoint := fmt.Sprintf("/lookup/iban/%s", url.PathEscape(iban))
	
	var result struct {
		Results []Profile `json:"results"`
	}
	
	if err := c.doRequest(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	if len(result.Results) == 0 {
		return nil, nil
	}
	
	return &result.Results[0], nil
}