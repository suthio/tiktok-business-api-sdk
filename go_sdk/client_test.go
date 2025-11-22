package tiktok

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-access-token")
	assert.NotNil(t, client)
	assert.Equal(t, "test-access-token", client.accessToken)
	assert.Equal(t, "https://business-api.tiktok.com", client.baseURL)
	assert.NotNil(t, client.httpClient)
}

func TestNewClientWithConfig(t *testing.T) {
	customHTTPClient := &http.Client{}
	client := NewClientWithConfig("test-access-token", "https://custom-url.com", customHTTPClient)
	assert.NotNil(t, client)
	assert.Equal(t, "test-access-token", client.accessToken)
	assert.Equal(t, "https://custom-url.com", client.baseURL)
	assert.Equal(t, customHTTPClient, client.httpClient)
}

func TestNewClientWithConfig_DefaultValues(t *testing.T) {
	client := NewClientWithConfig("test-access-token", "", nil)
	assert.NotNil(t, client)
	assert.Equal(t, "https://business-api.tiktok.com", client.baseURL)
	assert.NotNil(t, client.httpClient)
}

func TestClient_Get_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test/path", r.URL.Path)
		// Access-Token should now be in header, not query parameter
		assert.Equal(t, "test-token", r.Header.Get("Access-Token"))
		assert.Equal(t, "test-value", r.URL.Query().Get("test-param"))

		response := Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("request-123"),
			Data:      json.RawMessage(`{"test":"data"}`),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClientWithConfig("test-token", server.URL, nil)
	params := url.Values{}
	params.Set("test-param", "test-value")

	resp, err := client.Get(context.Background(), "/test/path", params)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), *resp.Code)
	assert.Equal(t, "Success", *resp.Message)
}

func TestClient_Get_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Code:      ptrInt64(40000),
			Message:   ptrString("Invalid request"),
			RequestID: ptrString("request-456"),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClientWithConfig("test-token", server.URL, nil)
	_, err := client.Get(context.Background(), "/test/path", nil)
	require.Error(t, err)

	errResp, ok := err.(*ErrorResponse)
	require.True(t, ok)
	assert.Equal(t, int64(40000), errResp.Code)
	assert.Equal(t, "Invalid request", errResp.Message)
}

func TestClient_Post_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/test/path", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Read and verify request body
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)
		assert.Equal(t, "test-value", reqBody["test-field"])

		response := Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("request-789"),
			Data:      json.RawMessage(`{"result":"created"}`),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClientWithConfig("test-token", server.URL, nil)
	body := map[string]string{"test-field": "test-value"}

	resp, err := client.Post(context.Background(), "/test/path", nil, body)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), *resp.Code)
}

func TestClient_Put_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		response := Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("request-put"),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClientWithConfig("test-token", server.URL, nil)
	resp, err := client.Put(context.Background(), "/test/path", nil, map[string]string{"test": "data"})
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestClient_Delete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		response := Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("request-delete"),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClientWithConfig("test-token", server.URL, nil)
	resp, err := client.Delete(context.Background(), "/test/path", nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// Helper functions
func ptrInt64(i int64) *int64 {
	return &i
}

func ptrString(s string) *string {
	return &s
}
