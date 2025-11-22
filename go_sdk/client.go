package tiktok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client represents the TikTok Business API client
type Client struct {
	baseURL     string
	httpClient  *http.Client
	accessToken string
}

// NewClient creates a new TikTok Business API client
func NewClient(accessToken string) *Client {
	baseURL := "https://business-api.tiktok.com"

	// Check if sandbox mode is enabled via environment variable
	if os.Getenv("TIKTOK_AD_IS_SANDBOX") == "true" {
		baseURL = "https://sandbox-ads.tiktok.com"
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		accessToken: accessToken,
	}
}

// NewClientWithConfig creates a new client with custom configuration
func NewClientWithConfig(accessToken string, baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	if baseURL == "" {
		baseURL = "https://business-api.tiktok.com"
	}
	return &Client{
		baseURL:     baseURL,
		httpClient:  httpClient,
		accessToken: accessToken,
	}
}

// doRequest performs an HTTP request and returns the response
func (c *Client) doRequest(ctx context.Context, method, path string, queryParams url.Values, body interface{}) (*Response, error) {
	// Build URL
	fullURL := c.baseURL + path

	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Access-Token in header (not query parameter)
	req.Header.Set("Access-Token", c.accessToken)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors
	if apiResp.Code != nil && *apiResp.Code != 0 {
		errResp := &ErrorResponse{
			Code:    *apiResp.Code,
			Message: *apiResp.Message,
		}
		if apiResp.RequestID != nil {
			errResp.RequestID = *apiResp.RequestID
		}
		return &apiResp, errResp
	}

	return &apiResp, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, queryParams url.Values) (*Response, error) {
	return c.doRequest(ctx, http.MethodGet, path, queryParams, nil)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, queryParams url.Values, body interface{}) (*Response, error) {
	return c.doRequest(ctx, http.MethodPost, path, queryParams, body)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, queryParams url.Values, body interface{}) (*Response, error) {
	return c.doRequest(ctx, http.MethodPut, path, queryParams, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string, queryParams url.Values) (*Response, error) {
	return c.doRequest(ctx, http.MethodDelete, path, queryParams, nil)
}
