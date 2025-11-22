package tiktok

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHTTPClient implements the HTTPClient interface for testing
type mockHTTPClient struct{}

func (m *mockHTTPClient) Do(_ interface{}) (interface{}, error) {
	return nil, nil
}

func TestResponse_Marshaling(t *testing.T) {
	t.Run("unmarshal valid response", func(t *testing.T) {
		jsonData := `{
			"code": 0,
			"message": "OK",
			"request_id": "req-123",
			"data": {"key": "value"}
		}`

		var resp Response
		err := json.Unmarshal([]byte(jsonData), &resp)

		require.NoError(t, err)
		assert.Equal(t, int64(0), *resp.Code)
		assert.Equal(t, "OK", *resp.Message)
		assert.Equal(t, "req-123", *resp.RequestID)
		assert.NotNil(t, resp.Data)
	})

	t.Run("unmarshal response with missing fields", func(t *testing.T) {
		jsonData := `{"code": 0}`

		var resp Response
		err := json.Unmarshal([]byte(jsonData), &resp)

		require.NoError(t, err)
		assert.Equal(t, int64(0), *resp.Code)
		assert.Nil(t, resp.Message)
		assert.Nil(t, resp.RequestID)
	})

	t.Run("marshal response", func(t *testing.T) {
		code := int64(0)
		msg := "Success"
		reqID := "test-req-id"
		data := json.RawMessage(`{"test": "data"}`)

		resp := Response{
			Code:      &code,
			Message:   &msg,
			RequestID: &reqID,
			Data:      data,
		}

		jsonBytes, err := json.Marshal(resp)
		require.NoError(t, err)

		var unmarshaled Response
		err = json.Unmarshal(jsonBytes, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, int64(0), *unmarshaled.Code)
		assert.Equal(t, "Success", *unmarshaled.Message)
		assert.Equal(t, "test-req-id", *unmarshaled.RequestID)
	})
}

func TestErrorResponse_Error(t *testing.T) {
	t.Run("error message", func(t *testing.T) {
		err := &ErrorResponse{
			Code:      40001,
			Message:   "Invalid access token",
			RequestID: "req-error-123",
		}

		assert.Equal(t, "Invalid access token", err.Error())
	})

	t.Run("empty error message", func(t *testing.T) {
		err := &ErrorResponse{
			Code:      40000,
			Message:   "",
			RequestID: "req-error-456",
		}

		assert.Equal(t, "", err.Error())
	})

	t.Run("error implements error interface", func(t *testing.T) {
		var err error = &ErrorResponse{
			Code:    50000,
			Message: "Server error",
		}

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Server error")
	})
}

func TestErrorResponse_Marshaling(t *testing.T) {
	t.Run("unmarshal error response", func(t *testing.T) {
		jsonData := `{
			"code": 40100,
			"message": "Authentication failed",
			"request_id": "req-auth-error"
		}`

		var errResp ErrorResponse
		err := json.Unmarshal([]byte(jsonData), &errResp)

		require.NoError(t, err)
		assert.Equal(t, int64(40100), errResp.Code)
		assert.Equal(t, "Authentication failed", errResp.Message)
		assert.Equal(t, "req-auth-error", errResp.RequestID)
	})

	t.Run("marshal error response", func(t *testing.T) {
		errResp := ErrorResponse{
			Code:      40400,
			Message:   "Resource not found",
			RequestID: "req-404",
		}

		jsonBytes, err := json.Marshal(errResp)
		require.NoError(t, err)

		var unmarshaled ErrorResponse
		err = json.Unmarshal(jsonBytes, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, int64(40400), unmarshaled.Code)
		assert.Equal(t, "Resource not found", unmarshaled.Message)
		assert.Equal(t, "req-404", unmarshaled.RequestID)
	})
}

func TestDefaultConfig(t *testing.T) {
	t.Run("returns default configuration", func(t *testing.T) {
		config := DefaultConfig()

		assert.NotNil(t, config)
		assert.Equal(t, "https://business-api.tiktok.com", config.BaseURL)
		assert.Nil(t, config.HTTPClient)
	})

	t.Run("returns new instance each time", func(t *testing.T) {
		config1 := DefaultConfig()
		config2 := DefaultConfig()

		assert.NotSame(t, config1, config2)
		assert.Equal(t, config1.BaseURL, config2.BaseURL)
	})
}

func TestClientConfig(t *testing.T) {
	t.Run("can create custom config", func(t *testing.T) {
		config := &ClientConfig{
			BaseURL: "https://custom-api.example.com",
		}

		assert.Equal(t, "https://custom-api.example.com", config.BaseURL)
		assert.Nil(t, config.HTTPClient)
	})

	t.Run("can set custom HTTP client", func(t *testing.T) {
		mockClient := &mockHTTPClient{}

		config := &ClientConfig{
			BaseURL:    "https://api.example.com",
			HTTPClient: mockClient,
		}

		assert.NotNil(t, config.HTTPClient)
		assert.Equal(t, mockClient, config.HTTPClient)
	})
}

func TestPageInfo_Marshaling(t *testing.T) {
	t.Run("unmarshal page info", func(t *testing.T) {
		jsonData := `{
			"page": 2,
			"page_size": 50,
			"total_number": 235,
			"total_page": 5
		}`

		var pageInfo PageInfo
		err := json.Unmarshal([]byte(jsonData), &pageInfo)

		require.NoError(t, err)
		assert.Equal(t, int64(2), pageInfo.Page)
		assert.Equal(t, int64(50), pageInfo.PageSize)
		assert.Equal(t, int64(235), pageInfo.TotalNumber)
		assert.Equal(t, int64(5), pageInfo.TotalPage)
	})

	t.Run("marshal page info", func(t *testing.T) {
		pageInfo := PageInfo{
			Page:        1,
			PageSize:    100,
			TotalNumber: 450,
			TotalPage:   5,
		}

		jsonBytes, err := json.Marshal(pageInfo)
		require.NoError(t, err)

		var unmarshaled PageInfo
		err = json.Unmarshal(jsonBytes, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, int64(1), unmarshaled.Page)
		assert.Equal(t, int64(100), unmarshaled.PageSize)
		assert.Equal(t, int64(450), unmarshaled.TotalNumber)
		assert.Equal(t, int64(5), unmarshaled.TotalPage)
	})

	t.Run("zero values", func(t *testing.T) {
		pageInfo := PageInfo{}

		jsonBytes, err := json.Marshal(pageInfo)
		require.NoError(t, err)

		expected := `{"page":0,"page_size":0,"total_number":0,"total_page":0}`
		assert.JSONEq(t, expected, string(jsonBytes))
	})
}

func TestPaginationParams(t *testing.T) {
	t.Run("with both values", func(t *testing.T) {
		page := int64(3)
		pageSize := int64(25)

		params := PaginationParams{
			Page:     &page,
			PageSize: &pageSize,
		}

		assert.NotNil(t, params.Page)
		assert.NotNil(t, params.PageSize)
		assert.Equal(t, int64(3), *params.Page)
		assert.Equal(t, int64(25), *params.PageSize)
	})

	t.Run("with nil values", func(t *testing.T) {
		params := PaginationParams{}

		assert.Nil(t, params.Page)
		assert.Nil(t, params.PageSize)
	})

	t.Run("with only page", func(t *testing.T) {
		page := int64(1)

		params := PaginationParams{
			Page: &page,
		}

		assert.NotNil(t, params.Page)
		assert.Nil(t, params.PageSize)
		assert.Equal(t, int64(1), *params.Page)
	})

	t.Run("with only page_size", func(t *testing.T) {
		pageSize := int64(10)

		params := PaginationParams{
			PageSize: &pageSize,
		}

		assert.Nil(t, params.Page)
		assert.NotNil(t, params.PageSize)
		assert.Equal(t, int64(10), *params.PageSize)
	})
}
