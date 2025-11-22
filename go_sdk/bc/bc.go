package bc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the BC (Business Center) API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new BC API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// TransactionInfo represents transaction information
type TransactionInfo struct {
	TransactionID   string  `json:"transaction_id"`
	TransactionTime string  `json:"transaction_time"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	AdvertiserID    string  `json:"advertiser_id,omitempty"`
	AdvertiserName  string  `json:"advertiser_name,omitempty"`
	Description     string  `json:"description,omitempty"`
	Balance         float64 `json:"balance,omitempty"`
}

// AccountTransactionResponse represents the response for account transactions
type AccountTransactionResponse struct {
	List     []TransactionInfo `json:"list"`
	PageInfo PageInfo          `json:"page_info"`
}

// PageInfo represents pagination information
type PageInfo struct {
	Page        int64 `json:"page"`
	PageSize    int64 `json:"page_size"`
	TotalNumber int64 `json:"total_number"`
	TotalPage   int64 `json:"total_page"`
}

// AccountTransactionRequest represents the request to get account transactions
type AccountTransactionRequest struct {
	BcID             *string     `json:"bc_id,omitempty"`
	ChildBcID        *string     `json:"child_bc_id,omitempty"`
	TransactionLevel *string     `json:"transaction_level,omitempty"`
	Filtering        interface{} `json:"filtering,omitempty"`
	Page             *int64      `json:"page,omitempty"`
	PageSize         *int64      `json:"page_size,omitempty"`
}

// GetAccountTransactions gets the transaction records of a BC or ad accounts
// Reference: https://business-api.tiktok.com/portal/docs?id=1792849810925569
func (a *API) GetAccountTransactions(ctx context.Context, req *AccountTransactionRequest) (*AccountTransactionResponse, error) {
	params := url.Values{}

	if req.BcID != nil {
		params.Set("bc_id", *req.BcID)
	}

	if req.ChildBcID != nil {
		params.Set("child_bc_id", *req.ChildBcID)
	}

	if req.TransactionLevel != nil {
		params.Set("transaction_level", *req.TransactionLevel)
	}

	if req.Filtering != nil {
		filteringJSON, err := json.Marshal(req.Filtering)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal filtering: %w", err)
		}
		params.Set("filtering", string(filteringJSON))
	}

	if req.Page != nil {
		params.Set("page", strconv.FormatInt(*req.Page, 10))
	}

	if req.PageSize != nil {
		params.Set("page_size", strconv.FormatInt(*req.PageSize, 10))
	}

	resp, err := a.client.Get(ctx, "/open_api/v1.3/bc/account/transaction/get/", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get account transactions: %w", err)
	}

	var transResp AccountTransactionResponse
	if err := json.Unmarshal(resp.Data, &transResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account transaction response: %w", err)
	}

	return &transResp, nil
}

// AssetInfo represents asset information
type AssetInfo struct {
	AssetID   string `json:"asset_id"`
	AssetType string `json:"asset_type"`
	AssetName string `json:"asset_name,omitempty"`
}

// AssetGetResponse represents the response for getting assets
type AssetGetResponse struct {
	List     []AssetInfo `json:"list"`
	PageInfo PageInfo    `json:"page_info,omitempty"`
}

// AssetGetRequest represents the request to get assets
type AssetGetRequest struct {
	BcID      string   `json:"bc_id"`
	AssetType string   `json:"asset_type"`
	AssetIDs  []string `json:"asset_ids,omitempty"`
	Page      *int64   `json:"page,omitempty"`
	PageSize  *int64   `json:"page_size,omitempty"`
}

// GetAssets gets assets in a Business Center
// Reference: https://business-api.tiktok.com/portal/docs?id=1739593603696641
func (a *API) GetAssets(ctx context.Context, req *AssetGetRequest) (*AssetGetResponse, error) {
	params := url.Values{}
	params.Set("bc_id", req.BcID)
	params.Set("asset_type", req.AssetType)

	if len(req.AssetIDs) > 0 {
		assetIDsJSON, err := json.Marshal(req.AssetIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal asset_ids: %w", err)
		}
		params.Set("asset_ids", string(assetIDsJSON))
	}

	if req.Page != nil {
		params.Set("page", strconv.FormatInt(*req.Page, 10))
	}

	if req.PageSize != nil {
		params.Set("page_size", strconv.FormatInt(*req.PageSize, 10))
	}

	resp, err := a.client.Get(ctx, "/open_api/v1.3/bc/asset/get/", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}

	var assetResp AssetGetResponse
	if err := json.Unmarshal(resp.Data, &assetResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset get response: %w", err)
	}

	return &assetResp, nil
}
