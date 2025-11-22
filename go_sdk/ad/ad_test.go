package ad

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

func TestGetAds_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/ad/get/")
		assert.Equal(t, "123456789", r.URL.Query().Get("advertiser_id"))

		adData := GetAdResponse{
			List: []AdInfo{
				{
					AdID:            "ad-001",
					AdName:          "Test Ad 1",
					AdgroupID:       "adgroup-001",
					CampaignID:      "campaign-001",
					AdvertiserID:    "123456789",
					AdText:          "Buy now!",
					CallToAction:    "SHOP_NOW",
					OperationStatus: "ENABLE",
					PrimaryStatus:   "ACTIVE",
					CreateTime:      "2024-01-01 10:00:00",
					ModifyTime:      "2024-01-02 15:30:00",
				},
				{
					AdID:            "ad-002",
					AdName:          "Test Ad 2",
					AdgroupID:       "adgroup-001",
					CampaignID:      "campaign-001",
					AdvertiserID:    "123456789",
					VideoID:         "video-123",
					AdText:          "Limited offer",
					CallToAction:    "LEARN_MORE",
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

		responseData, _ := json.Marshal(adData)
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

	result, err := api.GetAds(context.Background(), &GetAdRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "ad-001", result.List[0].AdID)
	assert.Equal(t, "Test Ad 1", result.List[0].AdName)
	assert.Equal(t, "Buy now!", result.List[0].AdText)
	assert.Equal(t, "SHOP_NOW", result.List[0].CallToAction)
	assert.Equal(t, int64(1), result.PageInfo.Page)
	assert.Equal(t, int64(2), result.PageInfo.TotalNumber)
}

func TestGetAds_WithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "50", r.URL.Query().Get("page_size"))

		adData := GetAdResponse{
			List: []AdInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        2,
				PageSize:    50,
				TotalNumber: 150,
				TotalPage:   3,
			},
		}

		responseData, _ := json.Marshal(adData)
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
	pageSize := int64(50)
	result, err := api.GetAds(context.Background(), &GetAdRequest{
		AdvertiserID: "123456789",
		Page:         &page,
		PageSize:     &pageSize,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.PageInfo.Page)
	assert.Equal(t, int64(50), result.PageInfo.PageSize)
}

func TestGetAds_WithFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fieldsParam := r.URL.Query().Get("fields")
		assert.NotEmpty(t, fieldsParam)
		assert.Contains(t, fieldsParam, "ad_id")
		assert.Contains(t, fieldsParam, "ad_name")

		adData := GetAdResponse{
			List: []AdInfo{
				{
					AdID:   "ad-001",
					AdName: "Test Ad",
				},
			},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 1,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(adData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAds(context.Background(), &GetAdRequest{
		AdvertiserID: "123456789",
		Fields:       []string{"ad_id", "ad_name"},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
}

func TestGetAds_WithFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		assert.NotEmpty(t, filteringParam)

		var filtering Filtering
		json.Unmarshal([]byte(filteringParam), &filtering)
		assert.Contains(t, filtering.AdgroupIDs, "adgroup-123")
		assert.Contains(t, filtering.CampaignIDs, "campaign-456")

		adData := GetAdResponse{
			List: []AdInfo{
				{
					AdID:       "ad-filtered",
					AdName:     "Filtered Ad",
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

		responseData, _ := json.Marshal(adData)
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
	result, err := api.GetAds(context.Background(), &GetAdRequest{
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
	assert.Equal(t, "ad-filtered", result.List[0].AdID)
	assert.Equal(t, "adgroup-123", result.List[0].AdgroupID)
}

func TestGetAds_WithAllFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		assert.NotEmpty(t, filteringParam)

		var filtering Filtering
		json.Unmarshal([]byte(filteringParam), &filtering)
		assert.Contains(t, filtering.AdIDs, "ad-001")
		assert.NotNil(t, filtering.PrimaryStatus)
		assert.NotNil(t, filtering.SecondaryStatus)
		assert.NotNil(t, filtering.ObjectiveType)
		assert.NotNil(t, filtering.CreateTimeMin)
		assert.NotNil(t, filtering.CreateTimeMax)

		adData := GetAdResponse{
			List:     []AdInfo{},
			PageInfo: tiktok.PageInfo{},
		}

		responseData, _ := json.Marshal(adData)
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
	secondaryStatus := "AD_STATUS_DELIVERY_OK"
	objectiveType := "CONVERSIONS"
	createTimeMin := "2024-01-01 00:00:00"
	createTimeMax := "2024-12-31 23:59:59"

	result, err := api.GetAds(context.Background(), &GetAdRequest{
		AdvertiserID: "123456789",
		Filtering: &Filtering{
			AdIDs:           []string{"ad-001", "ad-002"},
			AdgroupIDs:      []string{"adgroup-001"},
			CampaignIDs:     []string{"campaign-001"},
			PrimaryStatus:   &primaryStatus,
			SecondaryStatus: &secondaryStatus,
			ObjectiveType:   &objectiveType,
			CreateTimeMin:   &createTimeMin,
			CreateTimeMax:   &createTimeMax,
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAds_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		adData := GetAdResponse{
			List: []AdInfo{},
			PageInfo: tiktok.PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 0,
				TotalPage:   0,
			},
		}

		responseData, _ := json.Marshal(adData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAds(context.Background(), &GetAdRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.List)
	assert.Equal(t, int64(0), result.PageInfo.TotalNumber)
}

func TestGetAds_WithImageAd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		adData := GetAdResponse{
			List: []AdInfo{
				{
					AdID:            "ad-image-001",
					AdName:          "Image Ad",
					ImageIDs:        []string{"img-001", "img-002", "img-003"},
					AdText:          "Check out our products",
					CallToAction:    "SHOP_NOW",
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

		responseData, _ := json.Marshal(adData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAds(context.Background(), &GetAdRequest{
		AdvertiserID: "123456789",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 1)
	assert.Len(t, result.List[0].ImageIDs, 3)
	assert.Empty(t, result.List[0].VideoID)
}

// Helper functions
func ptrInt64(i int64) *int64 {
	return &i
}

func ptrString(s string) *string {
	return &s
}
