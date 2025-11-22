package audience

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Audience API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Audience API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// CustomAudienceInfo represents custom audience information
type CustomAudienceInfo struct {
	CustomAudienceID string `json:"custom_audience_id"`
	Name             string `json:"name"`
	AudienceType     string `json:"audience_type"`
	Size             int64  `json:"size,omitempty"`
	Status           string `json:"status"`
	ShareStatus      string `json:"share_status,omitempty"`
	CreateTime       string `json:"create_time"`
	ModifyTime       string `json:"modify_time"`
	AdvertiserID     string `json:"advertiser_id"`
	LookalikeType    string `json:"lookalike_type,omitempty"`
}

// CustomAudienceGetResponse represents the response for getting custom audiences
type CustomAudienceGetResponse struct {
	List []CustomAudienceInfo `json:"list"`
}

// CustomAudienceGetRequest represents the request to get custom audiences
type CustomAudienceGetRequest struct {
	AdvertiserID      string   `json:"advertiser_id"`
	CustomAudienceIDs []string `json:"custom_audience_ids"`
	HistorySize       *int64   `json:"history_size,omitempty"`
}

// GetCustomAudiences obtains the details of specified audiences
// Reference: https://business-api.tiktok.com/portal/docs?id=1739940507792385
func (a *API) GetCustomAudiences(ctx context.Context, req *CustomAudienceGetRequest) (*CustomAudienceGetResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	for _, id := range req.CustomAudienceIDs {
		params.Add("custom_audience_ids", id)
	}

	if req.HistorySize != nil {
		params.Set("history_size", strconv.FormatInt(*req.HistorySize, 10))
	}

	// Use generic DoGet helper
	var resp CustomAudienceGetResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/dmp/custom_audience/get/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get custom audiences: %w", err)
	}

	return &resp, nil
}

// CustomAudienceListResponse represents the response for listing custom audiences
type CustomAudienceListResponse struct {
	List     []CustomAudienceInfo `json:"list"`
	PageInfo tiktok.PageInfo      `json:"page_info"`
}

// CustomAudienceListRequest represents the request to list custom audiences
type CustomAudienceListRequest struct {
	AdvertiserID      string   `json:"advertiser_id"`
	CustomAudienceIDs []string `json:"custom_audience_ids,omitempty"`
	Page              *int64   `json:"page,omitempty"`
	PageSize          *int64   `json:"page_size,omitempty"`
}

// ListCustomAudiences gets all audiences
// Reference: https://business-api.tiktok.com/portal/docs?id=1739940506015746
func (a *API) ListCustomAudiences(ctx context.Context, req *CustomAudienceListRequest) (*CustomAudienceListResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	if len(req.CustomAudienceIDs) > 0 {
		for _, id := range req.CustomAudienceIDs {
			params.Add("custom_audience_ids", id)
		}
	}

	// Add pagination using helper
	tiktok.AddPagination(params, &tiktok.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	// Use generic DoGet helper
	var resp CustomAudienceListResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/dmp/custom_audience/list/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to list custom audiences: %w", err)
	}

	return &resp, nil
}
