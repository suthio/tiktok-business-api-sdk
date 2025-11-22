package bc

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

func TestGetAccountTransactions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/bc/account/transaction/get/")
		assert.Equal(t, "bc-123", r.URL.Query().Get("bc_id"))

		transactionData := AccountTransactionResponse{
			List: []TransactionInfo{
				{
					TransactionID:   "trans-001",
					TransactionTime: "2024-01-01 10:00:00",
					TransactionType: "DEPOSIT",
					Amount:          1000.50,
					Currency:        "USD",
					AdvertiserID:    "adv-001",
					AdvertiserName:  "Test Advertiser",
					Description:     "Initial deposit",
					Balance:         1000.50,
				},
				{
					TransactionID:   "trans-002",
					TransactionTime: "2024-01-02 15:30:00",
					TransactionType: "CHARGE",
					Amount:          -50.25,
					Currency:        "USD",
					AdvertiserID:    "adv-001",
					AdvertiserName:  "Test Advertiser",
					Description:     "Ad spend",
					Balance:         950.25,
				},
			},
			PageInfo: PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 2,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(transactionData)
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

	bcID := "bc-123"
	result, err := api.GetAccountTransactions(context.Background(), &AccountTransactionRequest{
		BcID: &bcID,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "trans-001", result.List[0].TransactionID)
	assert.Equal(t, "DEPOSIT", result.List[0].TransactionType)
	assert.Equal(t, 1000.50, result.List[0].Amount)
	assert.Equal(t, "USD", result.List[0].Currency)
	assert.Equal(t, int64(1), result.PageInfo.Page)
	assert.Equal(t, int64(2), result.PageInfo.TotalNumber)
}

func TestGetAccountTransactions_WithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "50", r.URL.Query().Get("page_size"))

		transactionData := AccountTransactionResponse{
			List: []TransactionInfo{},
			PageInfo: PageInfo{
				Page:        2,
				PageSize:    50,
				TotalNumber: 100,
				TotalPage:   2,
			},
		}

		responseData, _ := json.Marshal(transactionData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	bcID := "bc-456"
	page := int64(2)
	pageSize := int64(50)
	result, err := api.GetAccountTransactions(context.Background(), &AccountTransactionRequest{
		BcID:     &bcID,
		Page:     &page,
		PageSize: &pageSize,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.PageInfo.Page)
	assert.Equal(t, int64(50), result.PageInfo.PageSize)
}

func TestGetAccountTransactions_WithChildBC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "bc-parent", r.URL.Query().Get("bc_id"))
		assert.Equal(t, "bc-child", r.URL.Query().Get("child_bc_id"))

		transactionData := AccountTransactionResponse{
			List:     []TransactionInfo{},
			PageInfo: PageInfo{},
		}

		responseData, _ := json.Marshal(transactionData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	bcID := "bc-parent"
	childBcID := "bc-child"
	result, err := api.GetAccountTransactions(context.Background(), &AccountTransactionRequest{
		BcID:      &bcID,
		ChildBcID: &childBcID,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAccountTransactions_WithTransactionLevel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "BC", r.URL.Query().Get("transaction_level"))

		transactionData := AccountTransactionResponse{
			List:     []TransactionInfo{},
			PageInfo: PageInfo{},
		}

		responseData, _ := json.Marshal(transactionData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	bcID := "bc-789"
	transactionLevel := "BC"
	result, err := api.GetAccountTransactions(context.Background(), &AccountTransactionRequest{
		BcID:             &bcID,
		TransactionLevel: &transactionLevel,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAccountTransactions_WithFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filteringParam := r.URL.Query().Get("filtering")
		assert.NotEmpty(t, filteringParam)

		transactionData := AccountTransactionResponse{
			List:     []TransactionInfo{},
			PageInfo: PageInfo{},
		}

		responseData, _ := json.Marshal(transactionData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	bcID := "bc-filter"
	filtering := map[string]interface{}{
		"transaction_type": "DEPOSIT",
		"start_time":       "2024-01-01",
		"end_time":         "2024-12-31",
	}
	result, err := api.GetAccountTransactions(context.Background(), &AccountTransactionRequest{
		BcID:      &bcID,
		Filtering: filtering,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetAssets_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/open_api/v1.3/bc/asset/get/")
		assert.Equal(t, "bc-asset-123", r.URL.Query().Get("bc_id"))
		assert.Equal(t, "ADVERTISER", r.URL.Query().Get("asset_type"))

		assetData := AssetGetResponse{
			List: []AssetInfo{
				{
					AssetID:   "asset-001",
					AssetType: "ADVERTISER",
					AssetName: "Test Advertiser 1",
				},
				{
					AssetID:   "asset-002",
					AssetType: "ADVERTISER",
					AssetName: "Test Advertiser 2",
				},
			},
			PageInfo: PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 2,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(assetData)
		response := tiktok.Response{
			Code:      ptrInt64(0),
			Message:   ptrString("Success"),
			RequestID: ptrString("req-asset-123"),
			Data:      json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAssets(context.Background(), &AssetGetRequest{
		BcID:      "bc-asset-123",
		AssetType: "ADVERTISER",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
	assert.Equal(t, "asset-001", result.List[0].AssetID)
	assert.Equal(t, "ADVERTISER", result.List[0].AssetType)
	assert.Equal(t, "Test Advertiser 1", result.List[0].AssetName)
}

func TestGetAssets_WithAssetIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assetIDsParam := r.URL.Query().Get("asset_ids")
		assert.NotEmpty(t, assetIDsParam)

		var assetIDs []string
		json.Unmarshal([]byte(assetIDsParam), &assetIDs)
		assert.Contains(t, assetIDs, "asset-specific-1")
		assert.Contains(t, assetIDs, "asset-specific-2")

		assetData := AssetGetResponse{
			List: []AssetInfo{
				{
					AssetID:   "asset-specific-1",
					AssetType: "PIXEL",
					AssetName: "Pixel 1",
				},
				{
					AssetID:   "asset-specific-2",
					AssetType: "PIXEL",
					AssetName: "Pixel 2",
				},
			},
			PageInfo: PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 2,
				TotalPage:   1,
			},
		}

		responseData, _ := json.Marshal(assetData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAssets(context.Background(), &AssetGetRequest{
		BcID:      "bc-123",
		AssetType: "PIXEL",
		AssetIDs:  []string{"asset-specific-1", "asset-specific-2"},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, 2)
}

func TestGetAssets_WithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "3", r.URL.Query().Get("page"))
		assert.Equal(t, "20", r.URL.Query().Get("page_size"))

		assetData := AssetGetResponse{
			List: []AssetInfo{},
			PageInfo: PageInfo{
				Page:        3,
				PageSize:    20,
				TotalNumber: 60,
				TotalPage:   3,
			},
		}

		responseData, _ := json.Marshal(assetData)
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
	pageSize := int64(20)
	result, err := api.GetAssets(context.Background(), &AssetGetRequest{
		BcID:      "bc-page-test",
		AssetType: "ADVERTISER",
		Page:      &page,
		PageSize:  &pageSize,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(3), result.PageInfo.Page)
	assert.Equal(t, int64(20), result.PageInfo.PageSize)
}

func TestGetAssets_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assetData := AssetGetResponse{
			List: []AssetInfo{},
			PageInfo: PageInfo{
				Page:        1,
				PageSize:    10,
				TotalNumber: 0,
				TotalPage:   0,
			},
		}

		responseData, _ := json.Marshal(assetData)
		response := tiktok.Response{
			Code: ptrInt64(0),
			Data: json.RawMessage(responseData),
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := tiktok.NewClientWithConfig("test-token", server.URL, nil)
	api := NewAPI(client)

	result, err := api.GetAssets(context.Background(), &AssetGetRequest{
		BcID:      "bc-empty",
		AssetType: "CATALOG",
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
