package reporting

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

func TestGetIntegratedReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/open_api/v1.3/report/integrated/get/", r.URL.Path)

		// Verify query parameters
		assert.Equal(t, "BASIC", r.URL.Query().Get("report_type"))
		assert.Equal(t, "123456", r.URL.Query().Get("advertiser_id"))
		assert.Equal(t, "AUCTION_CAMPAIGN", r.URL.Query().Get("data_level"))
		assert.Equal(t, "2024-01-01", r.URL.Query().Get("start_date"))
		assert.Equal(t, "2024-01-31", r.URL.Query().Get("end_date"))

		// Send response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{
						"advertiser_id": "123456",
						"campaign_id":   "987654",
						"spend":         100.50,
						"impressions":   10000,
					},
				},
				"page_info": map[string]interface{}{
					"page":         1,
					"page_size":    10,
					"total_number": 1,
					"total_page":   1,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test_token", server.URL, nil)
	api := NewAPI(client)

	advertiserID := "123456"
	dataLevel := "AUCTION_CAMPAIGN"
	startDate := "2024-01-01"
	endDate := "2024-01-31"

	req := &IntegratedGetRequest{
		ReportType:   "BASIC",
		AdvertiserID: &advertiserID,
		DataLevel:    &dataLevel,
		StartDate:    &startDate,
		EndDate:      &endDate,
	}

	resp, err := api.GetIntegratedReport(context.Background(), req)
	require.NoError(t, err)

	assert.Len(t, resp.List, 1)
	assert.Equal(t, "123456", resp.List[0]["advertiser_id"])
	assert.Equal(t, float64(100.50), resp.List[0]["spend"])
	assert.Equal(t, int64(1), resp.PageInfo.Page)
	assert.Equal(t, int64(10), resp.PageInfo.PageSize)
}

func TestCheckReportTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/open_api/v1.3/report/task/check/", r.URL.Path)

		// Verify query parameters
		assert.Equal(t, "task123", r.URL.Query().Get("task_id"))
		assert.Equal(t, "123456", r.URL.Query().Get("advertiser_id"))

		// Send response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"task_id":      "task123",
				"status":       "COMPLETED",
				"download_url": "https://example.com/download",
				"total_count":  100,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test_token", server.URL, nil)
	api := NewAPI(client)

	resp, err := api.CheckReportTask(context.Background(), "task123", "123456")
	require.NoError(t, err)

	assert.Equal(t, "task123", resp.TaskID)
	assert.Equal(t, "COMPLETED", resp.Status)
	assert.Equal(t, "https://example.com/download", resp.DownloadURL)
	assert.Equal(t, int64(100), resp.TotalCount)
}

func TestGetMaterialReportBreakdown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/open_api/v1.3/smart_plus/material_report/breakdown/", r.URL.Path)

		// Verify query parameters
		assert.Equal(t, "123456", r.URL.Query().Get("advertiser_id"))
		assert.Equal(t, "2024-01-01", r.URL.Query().Get("start_date"))
		assert.Equal(t, "2024-01-31", r.URL.Query().Get("end_date"))

		// Send response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{
						"material_id": "mat123",
						"spend":       50.25,
						"clicks":      500,
					},
				},
				"page_info": map[string]interface{}{
					"page":         1,
					"page_size":    10,
					"total_number": 1,
					"total_page":   1,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test_token", server.URL, nil)
	api := NewAPI(client)

	req := &MaterialReportBreakdownRequest{
		AdvertiserID: "123456",
		Dimensions:   []string{"material_id"},
		StartDate:    "2024-01-01",
		EndDate:      "2024-01-31",
	}

	resp, err := api.GetMaterialReportBreakdown(context.Background(), req)
	require.NoError(t, err)

	assert.Len(t, resp.List, 1)
	assert.Equal(t, "mat123", resp.List[0]["material_id"])
	assert.Equal(t, float64(50.25), resp.List[0]["spend"])
	assert.Equal(t, int64(1), resp.PageInfo.Page)
}

func TestGetMaterialReportOverview(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/open_api/v1.3/smart_plus/material_report/overview/", r.URL.Path)

		// Verify query parameters
		assert.Equal(t, "123456", r.URL.Query().Get("advertiser_id"))

		// Send response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{
						"material_type": "VIDEO",
						"total_spend":   1000.00,
						"total_clicks":  10000,
					},
				},
				"page_info": map[string]interface{}{
					"page":         1,
					"page_size":    10,
					"total_number": 1,
					"total_page":   1,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test_token", server.URL, nil)
	api := NewAPI(client)

	queryLifetime := true
	req := &MaterialReportOverviewRequest{
		AdvertiserID:  "123456",
		Dimensions:    []string{"material_type"},
		QueryLifetime: &queryLifetime,
	}

	resp, err := api.GetMaterialReportOverview(context.Background(), req)
	require.NoError(t, err)

	assert.Len(t, resp.List, 1)
	assert.Equal(t, "VIDEO", resp.List[0]["material_type"])
	assert.Equal(t, float64(1000.00), resp.List[0]["total_spend"])
	assert.Equal(t, int64(1), resp.PageInfo.Page)
}
