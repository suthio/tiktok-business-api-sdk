package measurement

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

func TestNewAPI(t *testing.T) {
	client := tiktok.NewClient("test-token")
	api := NewAPI(client)

	assert.NotNil(t, api)
	assert.NotNil(t, api.client)
}

func TestListPixels_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/pixel/list/")
		assert.Equal(t, "123456789", r.URL.Query().Get("advertiser_id"))

		pixelData := PixelListResponse{
			List: []PixelInfo{
				{
					PixelID:        "pixel-001",
					PixelName:      "Test Pixel 1",
					PixelCode:      "ABC123",
					AdvertiserID:   "123456789",
					PixelStatus:    "ACTIVE",
					CreateTime:     "2024-01-01 10:00:00",
					LastUpdateTime: "2024-01-02 15:30:00",
				},
				{
					PixelID:        "pixel-002",
					PixelName:      "Test Pixel 2",
					PixelCode:      "DEF456",
					AdvertiserID:   "123456789",
					PixelStatus:    "ACTIVE",
					CreateTime:     "2024-01-03 09:15:00",
					LastUpdateTime: "2024-01-03 09:15:00",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 2,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(pixelData)
		response := tiktok.Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("req-123"),
			Data:      json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.ListPixels(context.Background(), &PixelListRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "pixel-001", result.List[0].PixelID)
	assert.Equal(t, "Test Pixel 1", result.List[0].PixelName)
	assert.Equal(t, "ABC123", result.List[0].PixelCode)
	assert.Equal(t, "ACTIVE", result.List[0].PixelStatus)
	assert.Equal(t, int64(1), result.PageInfo.Page)
	assert.Equal(t, int64(2), result.PageInfo.TotalNumber)
}

func TestListPixels_WithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "20", r.URL.Query().Get("page_size"))

		pixelData := PixelListResponse{
			List: []PixelInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        2,
				PageSize:    20,
				TotalNumber: 50,
				TotalPage:   3,
			},
		}

		responseData, _ := json.Marshal(pixelData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	page := int64(2)
	pageSize := int64(20)
	result, err := api.ListPixels(context.Background(), &PixelListRequest{
		AdvertiserID: "123456789",
		Page:         &page,
		PageSize:     &pageSize,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.PageInfo.Page)
	assert.Equal(t, int64(20), result.PageInfo.PageSize)
}

func TestListPixels_WithPixelID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "pixel-specific", r.URL.Query().Get("pixel_id"))

		pixelData := PixelListResponse{
			List: []PixelInfo{
				{
					PixelID:     "pixel-specific",
					PixelName:   "Specific Pixel",
					PixelStatus: "ACTIVE",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 1,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(pixelData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	pixelID := "pixel-specific"
	result, err := api.ListPixels(context.Background(), &PixelListRequest{
		AdvertiserID: "123456789",
		PixelID:      &pixelID,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "pixel-specific", result.List[0].PixelID)
}

func TestListPixels_WithCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "XYZ789", r.URL.Query().Get("code"))

		pixelData := PixelListResponse{
			List: []PixelInfo{
				{
					PixelID:   "pixel-code",
					PixelCode: "XYZ789",
				},
			},
			PageInfo: tiktok.PageInfo{},
		}

		responseData, _ := json.Marshal(pixelData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	code := "XYZ789"
	result, err := api.ListPixels(context.Background(), &PixelListRequest{
		AdvertiserID: "123456789",
		Code:         &code,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "XYZ789", result.List[0].PixelCode)
}

func TestListPixels_WithName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "My Pixel", r.URL.Query().Get("name"))

		pixelData := PixelListResponse{
			List: []PixelInfo{
				{
					PixelID:   "pixel-name",
					PixelName: "My Pixel",
				},
			},
			PageInfo: tiktok.PageInfo{},
		}

		responseData, _ := json.Marshal(pixelData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	name := "My Pixel"
	result, err := api.ListPixels(context.Background(), &PixelListRequest{
		AdvertiserID: "123456789",
		Name:         &name,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "My Pixel", result.List[0].PixelName)
}

func TestListPixels_WithOrderBy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "CREATE_TIME", r.URL.Query().Get("order_by"))

		pixelData := PixelListResponse{
			List:     []PixelInfo{},
			PageInfo: tiktok.PageInfo{},
		}

		responseData, _ := json.Marshal(pixelData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	orderBy := "CREATE_TIME"
	result, err := api.ListPixels(context.Background(), &PixelListRequest{
		AdvertiserID: "123456789",
		OrderBy:      &orderBy,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListPixels_WithFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		assert.NotEmpty(t, filteringParam)

		pixelData := PixelListResponse{
			List:     []PixelInfo{},
			PageInfo: tiktok.PageInfo{},
		}

		responseData, _ := json.Marshal(pixelData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	filtering := map[string]interface{}{
		"status": "ACTIVE",
	}
	result, err := api.ListPixels(context.Background(), &PixelListRequest{
		AdvertiserID: "123456789",
		Filtering:    filtering,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListPixels_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pixelData := PixelListResponse{
			List: []PixelInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 0,
				TotalPage:   0,
			},
		}

		responseData, _ := json.Marshal(pixelData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.ListPixels(context.Background(), &PixelListRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.List)
	assert.Equal(t, int64(0), result.PageInfo.TotalNumber)
}

func TestGetOfflineEventSets_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/offline/get/")
		assert.Equal(t, "123456789", r.URL.Query().Get("advertiser_id"))

		offlineData := OfflineGetResponse{
			List: []OfflineEventSetInfo{
				{
					EventSetID:   "event-set-001",
					Name:         "Test Event Set 1",
					AdvertiserID: "123456789",
					Status:       "ACTIVE",
					CreateTime:   "2024-01-01 10:00:00",
					UpdateTime:   "2024-01-02 15:30:00",
				},
				{
					EventSetID:   "event-set-002",
					Name:         "Test Event Set 2",
					AdvertiserID: "123456789",
					Status:       "ACTIVE",
					CreateTime:   "2024-01-03 09:15:00",
					UpdateTime:   "2024-01-03 09:15:00",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 2,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(offlineData)
		response := tiktok.Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("req-offline-123"),
			Data:      json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetOfflineEventSets(context.Background(), &OfflineGetRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "event-set-001", result.List[0].EventSetID)
	assert.Equal(t, "Test Event Set 1", result.List[0].Name)
	assert.Equal(t, "ACTIVE", result.List[0].Status)
}

func TestGetOfflineEventSets_WithEventSetIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventSetIDsParam := r.URL.Query().Get("event_set_ids")
		assert.NotEmpty(t, eventSetIDsParam)

		var eventSetIDs []string
		json.Unmarshal([]byte(eventSetIDsParam), &eventSetIDs)
		assert.Contains(t, eventSetIDs, "event-set-a")
		assert.Contains(t, eventSetIDs, "event-set-b")

		offlineData := OfflineGetResponse{
			List: []OfflineEventSetInfo{
				{
					EventSetID: "event-set-a",
					Name:       "Event Set A",
				},
				{
					EventSetID: "event-set-b",
					Name:       "Event Set B",
				},
			},
		}

		responseData, _ := json.Marshal(offlineData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetOfflineEventSets(context.Background(), &OfflineGetRequest{
		AdvertiserID: "123456789",
		EventSetIDs:  []string{"event-set-a", "event-set-b"},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "event-set-a", result.List[0].EventSetID)
}

func TestGetOfflineEventSets_WithName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "My Event Set", r.URL.Query().Get("name"))

		offlineData := OfflineGetResponse{
			List: []OfflineEventSetInfo{
				{
					EventSetID: "event-set-name",
					Name:       "My Event Set",
				},
			},
		}

		responseData, _ := json.Marshal(offlineData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	name := "My Event Set"
	result, err := api.GetOfflineEventSets(context.Background(), &OfflineGetRequest{
		AdvertiserID: "123456789",
		Name:         &name,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "My Event Set", result.List[0].Name)
}

func TestGetOfflineEventSets_EmptyAdvertiserID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.Query().Get("advertiser_id"))

		offlineData := OfflineGetResponse{
			List: []OfflineEventSetInfo{},
		}

		responseData, _ := json.Marshal(offlineData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetOfflineEventSets(context.Background(), &OfflineGetRequest{})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetOfflineEventSets_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offlineData := OfflineGetResponse{
			List: []OfflineEventSetInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 0,
				TotalPage:   0,
			},
		}

		responseData, _ := json.Marshal(offlineData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetOfflineEventSets(context.Background(), &OfflineGetRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.List)
}

// Helper functions
func ptrInt64(i int64) *int64 {
	return &i
}

func ptrString(s string) *string {
	return &s
}
