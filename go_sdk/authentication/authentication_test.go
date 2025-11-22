package authentication

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccessToken(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		assert.Equal(t, http.MethodPost, r.Method)

		// Verify request path
		assert.Equal(t, "/open_api/v1.3/oauth2/access_token/", r.URL.Path)

		// Verify content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var req AccessTokenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "test_app_id", req.AppID)
		assert.Equal(t, "test_auth_code", req.AuthCode)
		assert.Equal(t, "test_secret", req.Secret)

		// Send response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"access_token":             "test_access_token",
				"refresh_token":            "test_refresh_token",
				"expires_in":               86400,
				"refresh_token_expires_in": 31536000,
				"advertiser_ids":           []string{"123456", "789012"},
				"token_type":               "Bearer",
				"scope":                    "user.info.basic",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create API client
	api := NewAPIWithConfig(server.URL, nil)

	// Test GetAccessToken
	req := &AccessTokenRequest{
		AppID:    "test_app_id",
		AuthCode: "test_auth_code",
		Secret:   "test_secret",
	}

	resp, err := api.GetAccessToken(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "test_access_token", resp.AccessToken)
	assert.Equal(t, "test_refresh_token", resp.RefreshToken)
	assert.Equal(t, int64(86400), resp.ExpiresIn)
	assert.Equal(t, int64(31536000), resp.RefreshTokenExpiresIn)
	assert.Equal(t, []string{"123456", "789012"}, resp.AdvertiserIDs)
	assert.Equal(t, "Bearer", resp.TokenType)
	assert.Equal(t, "user.info.basic", resp.Scope)
}

func TestGetAccessToken_Error(t *testing.T) {
	// Create test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"code":       40001,
			"message":    "Invalid auth_code",
			"request_id": "test_request_id",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create API client
	api := NewAPIWithConfig(server.URL, nil)

	// Test GetAccessToken with error
	req := &AccessTokenRequest{
		AppID:    "test_app_id",
		AuthCode: "invalid_code",
		Secret:   "test_secret",
	}

	_, err := api.GetAccessToken(context.Background(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid auth_code")
}

func TestGetAdvertisers(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		assert.Equal(t, http.MethodGet, r.Method)

		// Verify request path
		assert.Equal(t, "/open_api/v1.3/oauth2/advertiser/get/", r.URL.Path)

		// Verify query parameters
		assert.Equal(t, "test_app_id", r.URL.Query().Get("app_id"))
		assert.Equal(t, "test_secret", r.URL.Query().Get("secret"))

		// Verify Access-Token header
		assert.Equal(t, "test_access_token", r.Header.Get("Access-Token"))

		// Send response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{
						"advertiser_id":   "123456",
						"advertiser_name": "Test Advertiser 1",
					},
					{
						"advertiser_id":   "789012",
						"advertiser_name": "Test Advertiser 2",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create API client
	api := NewAPIWithConfig(server.URL, nil)

	// Test GetAdvertisers
	resp, err := api.GetAdvertisers(context.Background(), "test_app_id", "test_secret", "test_access_token")
	require.NoError(t, err)

	assert.Len(t, resp.List, 2)
	assert.Equal(t, "123456", resp.List[0].AdvertiserID)
	assert.Equal(t, "Test Advertiser 1", resp.List[0].AdvertiserName)
	assert.Equal(t, "789012", resp.List[1].AdvertiserID)
	assert.Equal(t, "Test Advertiser 2", resp.List[1].AdvertiserName)
}

func TestGetAdvertisers_Error(t *testing.T) {
	// Create test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"code":       40100,
			"message":    "Invalid access token",
			"request_id": "test_request_id",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create API client
	api := NewAPIWithConfig(server.URL, nil)

	// Test GetAdvertisers with error
	_, err := api.GetAdvertisers(context.Background(), "test_app_id", "test_secret", "invalid_token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid access token")
}

func TestNewAPI(t *testing.T) {
	api := NewAPI()
	assert.NotNil(t, api)
	assert.Equal(t, "https://business-api.tiktok.com", api.baseURL)
	assert.NotNil(t, api.httpClient)
}

func TestNewAPIWithConfig(t *testing.T) {
	// Test with custom config
	customClient := &http.Client{}
	api := NewAPIWithConfig("https://custom.example.com", customClient)
	assert.NotNil(t, api)
	assert.Equal(t, "https://custom.example.com", api.baseURL)
	assert.Equal(t, customClient, api.httpClient)

	// Test with nil client (should create default)
	api = NewAPIWithConfig("", nil)
	assert.NotNil(t, api)
	assert.Equal(t, "https://business-api.tiktok.com", api.baseURL)
	assert.NotNil(t, api.httpClient)
}
