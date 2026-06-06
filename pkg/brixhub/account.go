package brixhub

import (
	"context"
	"fmt"
	"strconv"
)

// GetAccountInfo retrieves account information and quota status
func (c *Client) GetAccountInfo(ctx context.Context) (*AccountInfo, error) {
	var info AccountInfo
	if err := c.doRequest(ctx, "GET", "/me", nil, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetUsage retrieves detailed usage history (Pro+ plans only)
func (c *Client) GetUsage(ctx context.Context, limit, offset int) (*UsageResponse, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	
	endpoint := fmt.Sprintf("/usage?limit=%s&offset=%s", 
		strconv.Itoa(limit), 
		strconv.Itoa(offset))
	
	var usage UsageResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &usage); err != nil {
		return nil, err
	}
	
	return &usage, nil
}