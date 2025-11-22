package creative

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

func TestGetCreatives(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/open_api/v1.3/creative/get/" {
			t.Errorf("Expected path '/open_api/v1.3/creative/get/', got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify advertiser_id parameter
		advertiserID := r.URL.Query().Get("advertiser_id")
		if advertiserID != "123456789" {
			t.Errorf("Expected advertiser_id '123456789', got %s", advertiserID)
		}

		// Mock response
		response := map[string]interface{}{
			"code":       0,
			"message":    "OK",
			"request_id": "test_request_id",
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{
						"creative_id":   "creative_001",
						"creative_name": "Test Creative 1",
						"advertiser_id": "123456789",
						"ad_id":         "ad_001",
						"video_id":      "video_001",
						"ad_text":       "Test ad text",
						"creative_type": "VIDEO",
					},
					{
						"creative_id":   "creative_002",
						"creative_name": "Test Creative 2",
						"advertiser_id": "123456789",
						"ad_id":         "ad_002",
						"video_id":      "video_002",
						"ad_text":       "Test ad text 2",
						"creative_type": "VIDEO",
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
	req := &GetCreativesRequest{
		AdvertiserID: "123456789",
		Page:         &page,
		PageSize:     &pageSize,
	}

	resp, err := api.GetCreatives(context.Background(), req)
	if err != nil {
		t.Fatalf("GetCreatives failed: %v", err)
	}

	// Verify response
	if len(resp.List) != 2 {
		t.Errorf("Expected 2 creatives, got %d", len(resp.List))
	}

	if resp.List[0].CreativeID != "creative_001" {
		t.Errorf("Expected creative_id 'creative_001', got %s", resp.List[0].CreativeID)
	}

	if resp.PageInfo.TotalNumber != 2 {
		t.Errorf("Expected total_number 2, got %d", resp.PageInfo.TotalNumber)
	}
}

func TestGetAllCreatives(t *testing.T) {
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
							"creative_id":   "creative_001",
							"creative_name": "Test Creative 1",
							"advertiser_id": "123456789",
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
							"creative_id":   "creative_002",
							"creative_name": "Test Creative 2",
							"advertiser_id": "123456789",
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
	req := &GetCreativesRequest{
		AdvertiserID: "123456789",
	}

	creatives, err := api.GetAllCreatives(context.Background(), req)
	if err != nil {
		t.Fatalf("GetAllCreatives failed: %v", err)
	}

	// Verify response
	if len(creatives) != 2 {
		t.Errorf("Expected 2 creatives, got %d", len(creatives))
	}

	if callCount != 2 {
		t.Errorf("Expected 2 API calls, got %d", callCount)
	}

	if creatives[0].CreativeID != "creative_001" {
		t.Errorf("Expected first creative_id 'creative_001', got %s", creatives[0].CreativeID)
	}

	if creatives[1].CreativeID != "creative_002" {
		t.Errorf("Expected second creative_id 'creative_002', got %s", creatives[1].CreativeID)
	}
}
