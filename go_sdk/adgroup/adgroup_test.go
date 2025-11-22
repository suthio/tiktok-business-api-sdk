package adgroup

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

func TestGetAdGroups_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/adgroup/get/")
		assert.Equal(t, "123456789", r.URL.Query().Get("advertiser_id"))

		adgroupData := GetAdGroupResponse{
			List: []AdGroupInfo{
				{
					AdgroupID:        "adgroup-001",
					AdgroupName:      "Test AdGroup 1",
					CampaignID:       "campaign-001",
					AdvertiserID:     "123456789",
					ObjectiveType:    "CONVERSIONS",
					Budget:           100.50,
					BudgetMode:       "BUDGET_MODE_DAY",
					BillingEvent:     "CPC",
					OptimizationGoal: "CLICK",
					Placements:       []string{"PLACEMENT_TIKTOK"},
					Locations:        []string{"US", "CA"},
					Age:              []string{"AGE_25_34", "AGE_35_44"},
					Gender:           "GENDER_UNLIMITED",
					Languages:        []string{"en"},
					OperationStatus:  "ENABLE",
					PrimaryStatus:    "ACTIVE",
					CreateTime:       "2024-01-01 10:00:00",
					ModifyTime:       "2024-01-02 15:30:00",
				},
				{
					AdgroupID:       "adgroup-002",
					AdgroupName:     "Test AdGroup 2",
					CampaignID:      "campaign-001",
					AdvertiserID:    "123456789",
					ObjectiveType:   "TRAFFIC",
					Budget:          250.00,
					BudgetMode:      "BUDGET_MODE_TOTAL",
					OperationStatus: "ENABLE",
					PrimaryStatus:   "ACTIVE",
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

		responseData, _ := json.Marshal(adgroupData)
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

	result, err := api.GetAdGroups(context.Background(), &GetAdGroupRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "adgroup-001", result.List[0].AdgroupID)
	assert.Equal(t, "Test AdGroup 1", result.List[0].AdgroupName)
	assert.Equal(t, "CONVERSIONS", result.List[0].ObjectiveType)
	assert.Equal(t, 100.50, result.List[0].Budget)
	assert.Contains(t, result.List[0].Placements, "PLACEMENT_TIKTOK")
	assert.Contains(t, result.List[0].Locations, "US")
	assert.Equal(t, int64(1), result.PageInfo.Page)
	assert.Equal(t, int64(2), result.PageInfo.TotalNumber)
}

func TestGetAdGroups_WithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "3", r.URL.Query().Get("page"))
		assert.Equal(t, "25", r.URL.Query().Get("page_size"))

		adgroupData := GetAdGroupResponse{
			List: []AdGroupInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        3,
				PageSize:    25,
				TotalNumber: 75,
				TotalPage:   3,
			},
		}

		responseData, _ := json.Marshal(adgroupData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	page := int64(3)
	pageSize := int64(25)
	result, err := api.GetAdGroups(context.Background(), &GetAdGroupRequest{
		AdvertiserID: "123456789",
		Page:         &page,
		PageSize:     &pageSize,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(3), result.PageInfo.Page)
	assert.Equal(t, int64(25), result.PageInfo.PageSize)
}

func TestGetAdGroups_WithFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fieldsParam := r.URL.Query().Get("fields")
		assert.NotEmpty(t, fieldsParam)
		assert.Contains(t, fieldsParam, "adgroup_id")
		assert.Contains(t, fieldsParam, "adgroup_name")
		assert.Contains(t, fieldsParam, "budget")

		adgroupData := GetAdGroupResponse{
			List: []AdGroupInfo{
				{
					AdgroupID:   "adgroup-001",
					AdgroupName: "Test AdGroup",
					Budget:      500.00,
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 1,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(adgroupData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAdGroups(context.Background(), &GetAdGroupRequest{
		AdvertiserID: "123456789",
		Fields:       []string{"adgroup_id", "adgroup_name", "budget"},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "adgroup-001", result.List[0].AdgroupID)
}

func TestGetAdGroups_WithFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		assert.NotEmpty(t, filteringParam)

		var filtering Filtering
		json.Unmarshal([]byte(filteringParam), &filtering)
		assert.Contains(t, filtering.AdgroupIDs, "adgroup-123")
		assert.Contains(t, filtering.CampaignIDs, "campaign-456")

		adgroupData := GetAdGroupResponse{
			List: []AdGroupInfo{
				{
					AdgroupID:  "adgroup-123",
					CampaignID: "campaign-456",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 1,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(adgroupData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	primaryStatus := "ACTIVE"
	result, err := api.GetAdGroups(context.Background(), &GetAdGroupRequest{
		AdvertiserID: "123456789",
		Filtering: &Filtering{
			AdgroupIDs:    []string{"adgroup-123"},
			CampaignIDs:   []string{"campaign-456"},
			PrimaryStatus: &primaryStatus,
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "adgroup-123", result.List[0].AdgroupID)
}

func TestGetAdGroups_WithAllFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		assert.NotEmpty(t, filteringParam)

		var filtering Filtering
		json.Unmarshal([]byte(filteringParam), &filtering)
		assert.Contains(t, filtering.AdgroupIDs, "adgroup-001")
		assert.NotNil(t, filtering.PrimaryStatus)
		assert.NotNil(t, filtering.SecondaryStatus)
		assert.NotNil(t, filtering.ObjectiveType)
		assert.NotNil(t, filtering.BillingEvent)
		assert.NotNil(t, filtering.CreateTimeMin)
		assert.NotNil(t, filtering.CreateTimeMax)

		adgroupData := GetAdGroupResponse{
			List:     []AdGroupInfo{},
			PageInfo: tiktok.PageInfo{},
		}

		responseData, _ := json.Marshal(adgroupData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	primaryStatus := "ACTIVE"
	secondaryStatus := "ADGROUP_STATUS_DELIVERY_OK"
	objectiveType := "CONVERSIONS"
	billingEvent := "CPC"
	createTimeMin := "2024-01-01 00:00:00"
	createTimeMax := "2024-12-31 23:59:59"

	result, err := api.GetAdGroups(context.Background(), &GetAdGroupRequest{
		AdvertiserID: "123456789",
		Filtering: &Filtering{
			AdgroupIDs:      []string{"adgroup-001", "adgroup-002"},
			CampaignIDs:     []string{"campaign-001"},
			PrimaryStatus:   &primaryStatus,
			SecondaryStatus: &secondaryStatus,
			ObjectiveType:   &objectiveType,
			BillingEvent:    &billingEvent,
			CreateTimeMin:   &createTimeMin,
			CreateTimeMax:   &createTimeMax,
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAdGroups_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adgroupData := GetAdGroupResponse{
			List: []AdGroupInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 0,
				TotalPage:   0,
			},
		}

		responseData, _ := json.Marshal(adgroupData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAdGroups(context.Background(), &GetAdGroupRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.List)
	assert.Equal(t, int64(0), result.PageInfo.TotalNumber)
}

func TestGetAdGroups_WithSchedule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adgroupData := GetAdGroupResponse{
			List: []AdGroupInfo{
				{
					AdgroupID:         "adgroup-scheduled",
					AdgroupName:       "Scheduled AdGroup",
					ScheduleStartTime: "2024-01-01 00:00:00",
					ScheduleEndTime:   "2024-12-31 23:59:59",
					OperationStatus:   "ENABLE",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 1,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(adgroupData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAdGroups(context.Background(), &GetAdGroupRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "2024-01-01 00:00:00", result.List[0].ScheduleStartTime)
	assert.Equal(t, "2024-12-31 23:59:59", result.List[0].ScheduleEndTime)
}

// Helper functions
func ptrInt64(i int64) *int64 {
	return &i
}

func ptrString(s string) *string {
	return &s
}
