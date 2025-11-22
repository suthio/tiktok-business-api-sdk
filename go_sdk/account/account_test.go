package account

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

func TestGetAdvertiserInfo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/advertiser/info/")

		advertiserData := AdvertiserInfoResponse{
			List: []AdvertiserInfo{
				{
					AdvertiserID:   "adv-123",
					AdvertiserName: "Test Advertiser",
					Currency:       "USD",
					Status:         "ACTIVE",
				},
			},
		}

		responseData, _ := json.Marshal(advertiserData)
		response := tiktok.Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("request-123"),
			Data:      json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAdvertiserInfo(context.Background(), []string{"adv-123"}, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "adv-123", result.List[0].AdvertiserID)
	assert.Equal(t, "Test Advertiser", result.List[0].AdvertiserName)
	assert.Equal(t, "USD", result.List[0].Currency)
}

func TestGetAdvertiserInfo_WithFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "fields")

		advertiserData := AdvertiserInfoResponse{
			List: []AdvertiserInfo{
				{
					AdvertiserID:   "adv-456",
					AdvertiserName: "Another Advertiser",
				},
			},
		}

		responseData, _ := json.Marshal(advertiserData)
		response := tiktok.Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("request-456"),
			Data:      json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	fields := []string{"advertiser_id", "advertiser_name"}
	result, err := api.GetAdvertiserInfo(context.Background(), []string{"adv-456"}, fields)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
}

// Helper functions
func ptrInt64(i int64) *int64 {
	return &i
}

func ptrString(s string) *string {
	return &s
}
