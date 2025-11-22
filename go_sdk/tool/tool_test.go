package tool

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

func TestGetCarrier_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/tool/carrier/")
		assert.Equal(t, "test-advertiser-id", r.URL.Query().Get("advertiser_id"))

		carrierData := CarrierResponse{
			Carriers: []Carrier{
				{CarrierID: "carrier-1", CarrierName: "Carrier 1"},
				{CarrierID: "carrier-2", CarrierName: "Carrier 2"},
			},
		}

		responseData, _ := json.Marshal(carrierData)
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

	result, err := api.GetCarrier(context.Background(), "test-advertiser-id")
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Carriers, 2)
	assert.Equal(t, "carrier-1", result.Carriers[0].CarrierID)
	assert.Equal(t, "Carrier 1", result.Carriers[0].CarrierName)
}

func TestGetLanguage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/tool/language/")
		assert.Equal(t, "test-advertiser-id", r.URL.Query().Get("advertiser_id"))

		languageData := LanguageResponse{
			Languages: []Language{
				{LanguageCode: "en", LanguageName: "English"},
				{LanguageCode: "ja", LanguageName: "Japanese"},
			},
		}

		responseData, _ := json.Marshal(languageData)
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

	result, err := api.GetLanguage(context.Background(), "test-advertiser-id")
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Languages, 2)
	assert.Equal(t, "en", result.Languages[0].LanguageCode)
	assert.Equal(t, "English", result.Languages[0].LanguageName)
}

func TestGetActionCategory_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/tool/action_category/")
		assert.Equal(t, "test-advertiser-id", r.URL.Query().Get("advertiser_id"))

		actionCategoryData := ActionCategoryResponse{
			ActionCategories: []ActionCategory{
				{ActionCategoryID: "cat-1", ActionCategoryName: "Category 1"},
			},
		}

		responseData, _ := json.Marshal(actionCategoryData)
		response := tiktok.Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("request-789"),
			Data:      json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetActionCategory(context.Background(), "test-advertiser-id", []string{"HOUSING"})
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.ActionCategories, 1)
	assert.Equal(t, "cat-1", result.ActionCategories[0].ActionCategoryID)
}

// Helper functions
func ptrInt64(i int64) *int64 {
	return &i
}

func ptrString(s string) *string {
	return &s
}
