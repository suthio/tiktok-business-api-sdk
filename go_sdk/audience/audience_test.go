package audience

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

func TestGetCustomAudiences_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/dmp/custom_audience/get/")
		assert.Equal(t, "123456789", r.URL.Query().Get("advertiser_id"))

		customAudienceIDs := r.URL.Query()["custom_audience_ids"]
		assert.Contains(t, customAudienceIDs, "audience-001")
		assert.Contains(t, customAudienceIDs, "audience-002")

		audienceData := CustomAudienceGetResponse{
			List: []CustomAudienceInfo{
				{
					CustomAudienceID: "audience-001",
					Name:             "Test Audience 1",
					AudienceType:     "PIXEL",
					Size:             5000,
					Status:           "ACTIVE",
					ShareStatus:      "SHARED",
					CreateTime:       "2024-01-01 10:00:00",
					ModifyTime:       "2024-01-02 15:30:00",
					AdvertiserID:     "123456789",
				},
				{
					CustomAudienceID: "audience-002",
					Name:             "Test Audience 2",
					AudienceType:     "LOOKALIKE",
					Size:             10000,
					Status:           "ACTIVE",
					CreateTime:       "2024-01-03 09:15:00",
					ModifyTime:       "2024-01-03 09:15:00",
					AdvertiserID:     "123456789",
					LookalikeType:    "SMART",
				},
			},
		}

		responseData, _ := json.Marshal(audienceData)
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

	result, err := api.GetCustomAudiences(context.Background(), &CustomAudienceGetRequest{
		AdvertiserID:      "123456789",
		CustomAudienceIDs: []string{"audience-001", "audience-002"},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "audience-001", result.List[0].CustomAudienceID)
	assert.Equal(t, "Test Audience 1", result.List[0].Name)
	assert.Equal(t, "PIXEL", result.List[0].AudienceType)
	assert.Equal(t, int64(5000), result.List[0].Size)
	assert.Equal(t, "LOOKALIKE", result.List[1].AudienceType)
	assert.Equal(t, "SMART", result.List[1].LookalikeType)
}

func TestGetCustomAudiences_WithHistorySize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "30", r.URL.Query().Get("history_size"))

		audienceData := CustomAudienceGetResponse{
			List: []CustomAudienceInfo{
				{
					CustomAudienceID: "audience-001",
					Name:             "Test Audience",
					Status:           "ACTIVE",
				},
			},
		}

		responseData, _ := json.Marshal(audienceData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	historySize := int64(30)
	result, err := api.GetCustomAudiences(context.Background(), &CustomAudienceGetRequest{
		AdvertiserID:      "123456789",
		CustomAudienceIDs: []string{"audience-001"},
		HistorySize:       &historySize,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
}

func TestGetCustomAudiences_SingleAudience(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customAudienceIDs := r.URL.Query()["custom_audience_ids"]
		assert.Len(t, customAudienceIDs, 1)
		assert.Equal(t, "audience-single", customAudienceIDs[0])

		audienceData := CustomAudienceGetResponse{
			List: []CustomAudienceInfo{
				{
					CustomAudienceID: "audience-single",
					Name:             "Single Audience",
					AudienceType:     "FILE",
					Status:           "ACTIVE",
				},
			},
		}

		responseData, _ := json.Marshal(audienceData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetCustomAudiences(context.Background(), &CustomAudienceGetRequest{
		AdvertiserID:      "123456789",
		CustomAudienceIDs: []string{"audience-single"},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "audience-single", result.List[0].CustomAudienceID)
}

func TestListCustomAudiences_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/dmp/custom_audience/list/")
		assert.Equal(t, "123456789", r.URL.Query().Get("advertiser_id"))

		audienceData := CustomAudienceListResponse{
			List: []CustomAudienceInfo{
				{
					CustomAudienceID: "audience-001",
					Name:             "Audience 1",
					AudienceType:     "PIXEL",
					Status:           "ACTIVE",
				},
				{
					CustomAudienceID: "audience-002",
					Name:             "Audience 2",
					AudienceType:     "FILE",
					Status:           "ACTIVE",
				},
				{
					CustomAudienceID: "audience-003",
					Name:             "Audience 3",
					AudienceType:     "LOOKALIKE",
					Status:           "ACTIVE",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 3,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(audienceData)
		response := tiktok.Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("req-list-123"),
			Data:      json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.ListCustomAudiences(context.Background(), &CustomAudienceListRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 3)
	assert.Equal(t, "audience-001", result.List[0].CustomAudienceID)
	assert.Equal(t, "Audience 1", result.List[0].Name)
	assert.Equal(t, int64(1), result.PageInfo.Page)
	assert.Equal(t, int64(3), result.PageInfo.TotalNumber)
}

func TestListCustomAudiences_WithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "20", r.URL.Query().Get("page_size"))

		audienceData := CustomAudienceListResponse{
			List: []CustomAudienceInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        2,
				PageSize:    20,
				TotalNumber: 50,
				TotalPage:   3,
			},
		}

		responseData, _ := json.Marshal(audienceData)
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
	result, err := api.ListCustomAudiences(context.Background(), &CustomAudienceListRequest{
		AdvertiserID: "123456789",
		Page:         &page,
		PageSize:     &pageSize,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.PageInfo.Page)
	assert.Equal(t, int64(20), result.PageInfo.PageSize)
	assert.Equal(t, int64(50), result.PageInfo.TotalNumber)
}

func TestListCustomAudiences_WithAudienceIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customAudienceIDs := r.URL.Query()["custom_audience_ids"]
		assert.Contains(t, customAudienceIDs, "audience-a")
		assert.Contains(t, customAudienceIDs, "audience-b")

		audienceData := CustomAudienceListResponse{
			List: []CustomAudienceInfo{
				{
					CustomAudienceID: "audience-a",
					Name:             "Audience A",
				},
				{
					CustomAudienceID: "audience-b",
					Name:             "Audience B",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 2,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(audienceData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.ListCustomAudiences(context.Background(), &CustomAudienceListRequest{
		AdvertiserID:      "123456789",
		CustomAudienceIDs: []string{"audience-a", "audience-b"},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
}

func TestListCustomAudiences_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		audienceData := CustomAudienceListResponse{
			List: []CustomAudienceInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 0,
				TotalPage:   0,
			},
		}

		responseData, _ := json.Marshal(audienceData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.ListCustomAudiences(context.Background(), &CustomAudienceListRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.List)
	assert.Equal(t, int64(0), result.PageInfo.TotalNumber)
}

// Helper functions
func ptrInt64(i int64) *int64 {
	return &i
}

func ptrString(s string) *string {
	return &s
}
