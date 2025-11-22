package tiktok

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// AddPagination adds pagination parameters to url.Values
func AddPagination(params url.Values, pagination *PaginationParams) {
	if pagination == nil {
		return
	}
	if pagination.Page != nil {
		params.Set("page", strconv.FormatInt(*pagination.Page, 10))
	}
	if pagination.PageSize != nil {
		params.Set("page_size", strconv.FormatInt(*pagination.PageSize, 10))
	}
}

// AddJSONParam marshals value to JSON and adds it to params with the given key
func AddJSONParam(params url.Values, key string, value interface{}) error {
	if value == nil {
		return nil
	}

	// Check if value is a slice and is empty
	if slice, ok := value.([]string); ok && len(slice) == 0 {
		return nil
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", key, err)
	}

	// Don't add if the marshaled value is "null"
	if string(jsonBytes) == "null" {
		return nil
	}

	params.Set(key, string(jsonBytes))
	return nil
}

// AddStringSlice adds a string slice as a JSON array parameter
func AddStringSlice(params url.Values, key string, values []string) error {
	if len(values) == 0 {
		return nil
	}
	return AddJSONParam(params, key, values)
}

// DoGet executes a GET request and unmarshals the response into result
// This is a generic helper that handles the common pattern of:
// 1. Calling client.Get()
// 2. Unmarshaling the response data
// 3. Returning any errors
func DoGet[T any](ctx context.Context, client *Client, path string, params url.Values, result *T) error {
	resp, err := client.Get(ctx, path, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Data, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// DoPost executes a POST request and unmarshals the response into result
func DoPost[T any](ctx context.Context, client *Client, path string, body interface{}, result *T) error {
	resp, err := client.Post(ctx, path, nil, body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Data, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
