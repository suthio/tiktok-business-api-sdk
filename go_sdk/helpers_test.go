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

func TestAddPagination(t *testing.T) {
	t.Run("with both page and page_size", func(t *testing.T) {
		params := url.Values{}
		page := int64(2)
		pageSize := int64(50)
		pagination := &PaginationParams{
			Page:     &page,
			PageSize: &pageSize,
		}

		AddPagination(params, pagination)

		assert.Equal(t, "2", params.Get("page"))
		assert.Equal(t, "50", params.Get("page_size"))
	})

	t.Run("with only page", func(t *testing.T) {
		params := url.Values{}
		page := int64(1)
		pagination := &PaginationParams{
			Page: &page,
		}

		AddPagination(params, pagination)

		assert.Equal(t, "1", params.Get("page"))
		assert.Empty(t, params.Get("page_size"))
	})

	t.Run("with only page_size", func(t *testing.T) {
		params := url.Values{}
		pageSize := int64(100)
		pagination := &PaginationParams{
			PageSize: &pageSize,
		}

		AddPagination(params, pagination)

		assert.Empty(t, params.Get("page"))
		assert.Equal(t, "100", params.Get("page_size"))
	})

	t.Run("with nil pagination", func(t *testing.T) {
		params := url.Values{}
		AddPagination(params, nil)

		assert.Empty(t, params.Get("page"))
		assert.Empty(t, params.Get("page_size"))
	})

	t.Run("with empty pagination", func(t *testing.T) {
		params := url.Values{}
		pagination := &PaginationParams{}

		AddPagination(params, pagination)

		assert.Empty(t, params.Get("page"))
		assert.Empty(t, params.Get("page_size"))
	})
}

func TestAddJSONParam(t *testing.T) {
	t.Run("with valid struct", func(t *testing.T) {
		params := url.Values{}
		value := struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		}{
			Field1: "test",
			Field2: 123,
		}

		err := AddJSONParam(params, "test_key", value)

		require.NoError(t, err)
		expected := `{"field1":"test","field2":123}`
		assert.JSONEq(t, expected, params.Get("test_key"))
	})

	t.Run("with string slice", func(t *testing.T) {
		params := url.Values{}
		value := []string{"val1", "val2", "val3"}

		err := AddJSONParam(params, "items", value)

		require.NoError(t, err)
		expected := `["val1","val2","val3"]`
		assert.JSONEq(t, expected, params.Get("items"))
	})

	t.Run("with empty string slice", func(t *testing.T) {
		params := url.Values{}
		value := []string{}

		err := AddJSONParam(params, "items", value)

		require.NoError(t, err)
		assert.Empty(t, params.Get("items"))
	})

	t.Run("with nil value", func(t *testing.T) {
		params := url.Values{}

		err := AddJSONParam(params, "test_key", nil)

		require.NoError(t, err)
		assert.Empty(t, params.Get("test_key"))
	})

	t.Run("with invalid value (channel)", func(t *testing.T) {
		params := url.Values{}
		value := make(chan int) // channels cannot be marshaled to JSON

		err := AddJSONParam(params, "test_key", value)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal test_key")
	})
}

func TestAddStringSlice(t *testing.T) {
	t.Run("with non-empty slice", func(t *testing.T) {
		params := url.Values{}
		values := []string{"apple", "banana", "cherry"}

		err := AddStringSlice(params, "fruits", values)

		require.NoError(t, err)
		expected := `["apple","banana","cherry"]`
		assert.JSONEq(t, expected, params.Get("fruits"))
	})

	t.Run("with empty slice", func(t *testing.T) {
		params := url.Values{}
		values := []string{}

		err := AddStringSlice(params, "fruits", values)

		require.NoError(t, err)
		assert.Empty(t, params.Get("fruits"))
	})

	t.Run("with nil slice", func(t *testing.T) {
		params := url.Values{}
		var values []string

		err := AddStringSlice(params, "fruits", values)

		require.NoError(t, err)
		assert.Empty(t, params.Get("fruits"))
	})

	t.Run("with single item", func(t *testing.T) {
		params := url.Values{}
		values := []string{"single"}

		err := AddStringSlice(params, "items", values)

		require.NoError(t, err)
		expected := `["single"]`
		assert.JSONEq(t, expected, params.Get("items"))
	})
}

func TestDoGet(t *testing.T) {
	type TestResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	t.Run("successful get request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Contains(t, r.URL.Path, "/test/path")

			testData := TestResponse{
				ID:   "test-123",
				Name: "Test Name",
			}

			responseData, _ := json.Marshal(testData)
			response := Response{
				Code:      ptrInt64(0),
				Message:   ptrString("Success"),
				RequestID: ptrString("req-123"),
				Data:      json.RawMessage(responseData),
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClientWithConfig("test-token", server.URL, nil)
		params := url.Values{}

		var result TestResponse
		err := DoGet(context.Background(), client, "/test/path", params, &result)

		require.NoError(t, err)
		assert.Equal(t, "test-123", result.ID)
		assert.Equal(t, "Test Name", result.Name)
	})

	t.Run("get request with query params", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "value1", r.URL.Query().Get("key1"))
			assert.Equal(t, "value2", r.URL.Query().Get("key2"))

			testData := TestResponse{ID: "test-456"}
			responseData, _ := json.Marshal(testData)
			response := Response{
				Code: ptrInt64(0),
				Data: json.RawMessage(responseData),
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClientWithConfig("test-token", server.URL, nil)
		params := url.Values{}
		params.Set("key1", "value1")
		params.Set("key2", "value2")

		var result TestResponse
		err := DoGet(context.Background(), client, "/test/path", params, &result)

		require.NoError(t, err)
		assert.Equal(t, "test-456", result.ID)
	})

	t.Run("unmarshaling error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := Response{
				Code: ptrInt64(0),
				Data: json.RawMessage(`{"invalid json`),
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClientWithConfig("test-token", server.URL, nil)
		params := url.Values{}

		var result TestResponse
		err := DoGet(context.Background(), client, "/test/path", params, &result)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal response")
	})
}

func TestDoPost(t *testing.T) {
	type TestRequest struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	type TestResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	t.Run("successful post request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Contains(t, r.URL.Path, "/test/post")

			var reqBody TestRequest
			json.NewDecoder(r.Body).Decode(&reqBody)
			assert.Equal(t, "test-name", reqBody.Name)
			assert.Equal(t, 42, reqBody.Value)

			testData := TestResponse{
				ID:     "created-123",
				Status: "SUCCESS",
			}

			responseData, _ := json.Marshal(testData)
			response := Response{
				Code:      ptrInt64(0),
				Message:   ptrString("Created"),
				RequestID: ptrString("req-456"),
				Data:      json.RawMessage(responseData),
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClientWithConfig("test-token", server.URL, nil)
		requestBody := TestRequest{
			Name:  "test-name",
			Value: 42,
		}

		var result TestResponse
		err := DoPost(context.Background(), client, "/test/post", requestBody, &result)

		require.NoError(t, err)
		assert.Equal(t, "created-123", result.ID)
		assert.Equal(t, "SUCCESS", result.Status)
	})

	t.Run("post with empty body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			testData := TestResponse{
				ID:     "empty-123",
				Status: "OK",
			}
			responseData, _ := json.Marshal(testData)
			response := Response{
				Code: ptrInt64(0),
				Data: json.RawMessage(responseData),
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClientWithConfig("test-token", server.URL, nil)

		var result TestResponse
		err := DoPost(context.Background(), client, "/test/post", nil, &result)

		require.NoError(t, err)
		assert.Equal(t, "empty-123", result.ID)
	})

	t.Run("unmarshaling error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := Response{
				Code: ptrInt64(0),
				Data: json.RawMessage(`invalid`),
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClientWithConfig("test-token", server.URL, nil)

		var result TestResponse
		err := DoPost(context.Background(), client, "/test/post", nil, &result)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal response")
	})
}
