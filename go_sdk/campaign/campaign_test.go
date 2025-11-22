package campaign

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

func TestGetCampaigns_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/campaign/get/")
		assert.Equal(t, "123456789", r.URL.Query().Get("advertiser_id"))

		campaignData := GetCampaignResponse{
			List: []CampaignStatus{
				{
					CampaignID:      "campaign-001",
					CampaignName:    "Test Campaign 1",
					AdvertiserID:    "123456789",
					ObjectiveType:   "CONVERSIONS",
					Budget:          5000.00,
					BudgetMode:      "BUDGET_MODE_TOTAL",
					OperationStatus: "ENABLE",
					CreateTime:      "2024-01-01 10:00:00",
					ModifyTime:      "2024-01-02 15:30:00",
				},
				{
					CampaignID:      "campaign-002",
					CampaignName:    "Test Campaign 2",
					AdvertiserID:    "123456789",
					ObjectiveType:   "TRAFFIC",
					Budget:          1000.00,
					BudgetMode:      "BUDGET_MODE_DAY",
					OperationStatus: "ENABLE",
					CreateTime:      "2024-01-03 09:15:00",
					ModifyTime:      "2024-01-03 09:15:00",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 2,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(campaignData)
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

	result, err := api.GetCampaigns(context.Background(), &GetCampaignRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "campaign-001", result.List[0].CampaignID)
	assert.Equal(t, "Test Campaign 1", result.List[0].CampaignName)
	assert.Equal(t, "CONVERSIONS", result.List[0].ObjectiveType)
	assert.Equal(t, 5000.00, result.List[0].Budget)
	assert.Equal(t, "BUDGET_MODE_TOTAL", result.List[0].BudgetMode)
	assert.Equal(t, int64(1), result.PageInfo.Page)
	assert.Equal(t, int64(2), result.PageInfo.TotalNumber)
}

func TestGetCampaigns_WithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "25", r.URL.Query().Get("page_size"))

		campaignData := GetCampaignResponse{
			List: []CampaignStatus{},
			PageInfo: tiktok.PageInfo{
				Page:        2,
				PageSize:    25,
				TotalNumber: 75,
				TotalPage:   3,
			},
		}

		responseData, _ := json.Marshal(campaignData)
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
	pageSize := int64(25)
	result, err := api.GetCampaigns(context.Background(), &GetCampaignRequest{
		AdvertiserID: "123456789",
		Page:         &page,
		PageSize:     &pageSize,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.PageInfo.Page)
	assert.Equal(t, int64(25), result.PageInfo.PageSize)
}

func TestGetCampaigns_WithFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		assert.NotEmpty(t, filteringParam)

		var filtering Filtering
		json.Unmarshal([]byte(filteringParam), &filtering)
		assert.Contains(t, filtering.CampaignIDs, "campaign-123")
		assert.Contains(t, filtering.CampaignIDs, "campaign-456")

		campaignData := GetCampaignResponse{
			List: []CampaignStatus{
				{
					CampaignID:   "campaign-123",
					CampaignName: "Filtered Campaign 1",
				},
				{
					CampaignID:   "campaign-456",
					CampaignName: "Filtered Campaign 2",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 2,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(campaignData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetCampaigns(context.Background(), &GetCampaignRequest{
		AdvertiserID: "123456789",
		Filtering: &Filtering{
			CampaignIDs: []string{"campaign-123", "campaign-456"},
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "campaign-123", result.List[0].CampaignID)
}

func TestGetCampaigns_WithAllFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		assert.NotEmpty(t, filteringParam)

		var filtering Filtering
		json.Unmarshal([]byte(filteringParam), &filtering)
		assert.Contains(t, filtering.CampaignIDs, "campaign-001")
		assert.NotNil(t, filtering.CampaignName)
		assert.NotNil(t, filtering.ObjectiveType)
		assert.NotNil(t, filtering.PrimaryStatus)
		assert.NotNil(t, filtering.SecondaryStatus)
		assert.NotNil(t, filtering.CreateTimeMin)
		assert.NotNil(t, filtering.CreateTimeMax)

		campaignData := GetCampaignResponse{
			List:     []CampaignStatus{},
			PageInfo: tiktok.PageInfo{},
		}

		responseData, _ := json.Marshal(campaignData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	campaignName := "Test Campaign"
	objectiveType := "CONVERSIONS"
	primaryStatus := "ACTIVE"
	secondaryStatus := "CAMPAIGN_STATUS_DELIVERY_OK"
	createTimeMin := "2024-01-01 00:00:00"
	createTimeMax := "2024-12-31 23:59:59"

	result, err := api.GetCampaigns(context.Background(), &GetCampaignRequest{
		AdvertiserID: "123456789",
		Filtering: &Filtering{
			CampaignIDs:     []string{"campaign-001", "campaign-002"},
			CampaignName:    &campaignName,
			ObjectiveType:   &objectiveType,
			PrimaryStatus:   &primaryStatus,
			SecondaryStatus: &secondaryStatus,
			CreateTimeMin:   &createTimeMin,
			CreateTimeMax:   &createTimeMax,
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetCampaigns_ByObjectiveType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		var filtering Filtering
		json.Unmarshal([]byte(filteringParam), &filtering)
		assert.Equal(t, "TRAFFIC", *filtering.ObjectiveType)

		campaignData := GetCampaignResponse{
			List: []CampaignStatus{
				{
					CampaignID:    "campaign-traffic",
					CampaignName:  "Traffic Campaign",
					ObjectiveType: "TRAFFIC",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 1,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(campaignData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	objectiveType := "TRAFFIC"
	result, err := api.GetCampaigns(context.Background(), &GetCampaignRequest{
		AdvertiserID: "123456789",
		Filtering: &Filtering{
			ObjectiveType: &objectiveType,
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "TRAFFIC", result.List[0].ObjectiveType)
}

func TestGetCampaigns_ByStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		var filtering Filtering
		json.Unmarshal([]byte(filteringParam), &filtering)
		assert.Equal(t, "ENABLE", *filtering.PrimaryStatus)

		campaignData := GetCampaignResponse{
			List: []CampaignStatus{
				{
					CampaignID:      "campaign-active",
					OperationStatus: "ENABLE",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 1,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(campaignData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	primaryStatus := "ENABLE"
	result, err := api.GetCampaigns(context.Background(), &GetCampaignRequest{
		AdvertiserID: "123456789",
		Filtering: &Filtering{
			PrimaryStatus: &primaryStatus,
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "ENABLE", result.List[0].OperationStatus)
}

func TestGetCampaigns_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		campaignData := GetCampaignResponse{
			List: []CampaignStatus{},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 0,
				TotalPage:   0,
			},
		}

		responseData, _ := json.Marshal(campaignData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetCampaigns(context.Background(), &GetCampaignRequest{
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
