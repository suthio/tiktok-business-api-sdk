package research

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

func TestGetAdReport(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/v2/research/adlib/ad/report/" {
			t.Errorf("Expected path '/v2/research/adlib/ad/report/', got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify search_term parameter
		searchTerm := r.URL.Query().Get("search_term")
		if searchTerm != "shoes" {
			t.Errorf("Expected search_term 'shoes', got %s", searchTerm)
		}

		// Verify country_code parameter
		countryCode := r.URL.Query().Get("country_code")
		if countryCode != "US" {
			t.Errorf("Expected country_code 'US', got %s", countryCode)
		}

		// Mock response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{
						"ad_id":           "ad_001",
						"ad_name":         "Test Ad 1",
						"advertiser_id":   "123456789",
						"advertiser_name": "Test Advertiser",
						"campaign_id":     "campaign_001",
						"campaign_name":   "Test Campaign",
						"ad_text":         "Buy our shoes!",
						"country":         "United States",
						"impressions":     10000,
						"clicks":          500,
						"ctr":             0.05,
						"video_id":        "video_001",
						"video_title":     "Shoes Ad Video",
					},
					{
						"ad_id":           "ad_002",
						"ad_name":         "Test Ad 2",
						"advertiser_id":   "987654321",
						"advertiser_name": "Another Advertiser",
						"ad_text":         "Best shoes ever!",
						"country":         "United States",
						"impressions":     5000,
						"clicks":          250,
						"ctr":             0.05,
					},
				},
				"page_info": map[string]interface{}{
					"page":         1,
					"page_size":    10,
					"total_number": 2,
					"total_page":   1,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := tiktok.NewClientWithConfig("test_access_token", server.URL, nil)
	api := NewAPI(client)

	// Make request
	page := int64(1)
	pageSize := int64(10)
	countryCode := "US"
	req := &GetAdReportRequest{
		SearchTerm:  "shoes",
		CountryCode: &countryCode,
		Page:        &page,
		PageSize:    &pageSize,
	}

	resp, err := api.GetAdReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetAdReport failed: %v", err)
	}

	// Verify response
	if len(resp.List) != 2 {
		t.Errorf("Expected 2 ad reports, got %d", len(resp.List))
	}

	if resp.List[0].AdID != "ad_001" {
		t.Errorf("Expected ad_id 'ad_001', got %s", resp.List[0].AdID)
	}

	if resp.List[0].Impressions != 10000 {
		t.Errorf("Expected impressions 10000, got %d", resp.List[0].Impressions)
	}

	if resp.PageInfo.TotalNumber != 2 {
		t.Errorf("Expected total_number 2, got %d", resp.PageInfo.TotalNumber)
	}
}

func TestGetAllAdReports(t *testing.T) {
	callCount := 0

	// Mock server that simulates pagination
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		page := r.URL.Query().Get("page")

		var response map[string]interface{}

		if page == "1" {
			response = map[string]interface{}{
				"code":       0,
				"message":    "OK",
				"request_id": "test_request_id",
				"data": map[string]interface{}{
					"list": []map[string]interface{}{
						{
							"ad_id":         "ad_001",
							"ad_name":       "Test Ad 1",
							"advertiser_id": "123456789",
							"impressions":   10000,
						},
					},
					"page_info": map[string]interface{}{
						"page":         1,
						"page_size":    100,
						"total_number": 2,
						"total_page":   2,
					},
				},
			}
		} else if page == "2" {
			response = map[string]interface{}{
				"code":       0,
				"message":    "OK",
				"request_id": "test_request_id",
				"data": map[string]interface{}{
					"list": []map[string]interface{}{
						{
							"ad_id":         "ad_002",
							"ad_name":       "Test Ad 2",
							"advertiser_id": "987654321",
							"impressions":   5000,
						},
					},
					"page_info": map[string]interface{}{
						"page":         2,
						"page_size":    100,
						"total_number": 2,
						"total_page":   2,
					},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := tiktok.NewClientWithConfig("test_access_token", server.URL, nil)
	api := NewAPI(client)

	// Make request
	req := &GetAdReportRequest{
		SearchTerm: "test",
	}

	reports, err := api.GetAllAdReports(context.Background(), req)
	if err != nil {
		t.Fatalf("GetAllAdReports failed: %v", err)
	}

	// Verify response
	if len(reports) != 2 {
		t.Errorf("Expected 2 ad reports, got %d", len(reports))
	}

	if callCount != 2 {
		t.Errorf("Expected 2 API calls, got %d", callCount)
	}

	if reports[0].AdID != "ad_001" {
		t.Errorf("Expected first ad_id 'ad_001', got %s", reports[0].AdID)
	}

	if reports[1].AdID != "ad_002" {
		t.Errorf("Expected second ad_id 'ad_002', got %s", reports[1].AdID)
	}
}

func TestGetAdReportWithFiltering(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify filtering parameter
		filteringParam := r.URL.Query().Get("filtering")
		if filteringParam == "" {
			t.Error("Expected filtering parameter to be set")
		}

		// Mock response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{
						"ad_id":         "ad_001",
						"ad_name":       "Filtered Ad",
						"advertiser_id": "specific_advertiser",
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
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := tiktok.NewClientWithConfig("test_access_token", server.URL, nil)
	api := NewAPI(client)

	// Make request with filtering
	req := &GetAdReportRequest{
		SearchTerm: "shoes",
		Filtering: &AdReportFiltering{
			AdvertiserIDs: []string{"specific_advertiser"},
			CountryCodes:  []string{"US"},
		},
	}

	resp, err := api.GetAdReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetAdReport with filtering failed: %v", err)
	}

	// Verify response
	if len(resp.List) != 1 {
		t.Errorf("Expected 1 ad report, got %d", len(resp.List))
	}

	if resp.List[0].AdvertiserID != "specific_advertiser" {
		t.Errorf("Expected advertiser_id 'specific_advertiser', got %s", resp.List[0].AdvertiserID)
	}
}
