package measurement

import (
	"context"
	"fmt"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Measurement API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Measurement API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// PixelInfo represents pixel information
type PixelInfo struct {
	PixelID        string `json:"pixel_id"`
	PixelName      string `json:"pixel_name"`
	PixelCode      string `json:"pixel_code"`
	AdvertiserID   string `json:"advertiser_id"`
	PixelStatus    string `json:"pixel_status"`
	CreateTime     string `json:"create_time"`
	LastUpdateTime string `json:"last_update_time"`
}

// PixelListResponse represents the response for listing pixels
type PixelListResponse struct {
	List     []PixelInfo     `json:"list"`
	PageInfo tiktok.PageInfo `json:"page_info"`
}

// PixelListRequest represents the request to list pixels
type PixelListRequest struct {
	AdvertiserID string      `json:"advertiser_id"`
	PixelID      *string     `json:"pixel_id,omitempty"`
	Code         *string     `json:"code,omitempty"`
	Name         *string     `json:"name,omitempty"`
	OrderBy      *string     `json:"order_by,omitempty"`
	Filtering    interface{} `json:"filtering,omitempty"`
	Page         *int64      `json:"page,omitempty"`
	PageSize     *int64      `json:"page_size,omitempty"`
}

// ListPixels obtains a list of Pixel information
// Reference: https://business-api.tiktok.com/portal/docs?id=1740858697598978
func (a *API) ListPixels(ctx context.Context, req *PixelListRequest) (*PixelListResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	if req.PixelID != nil {
		params.Set("pixel_id", *req.PixelID)
	}

	if req.Code != nil {
		params.Set("code", *req.Code)
	}

	if req.Name != nil {
		params.Set("name", *req.Name)
	}

	if req.OrderBy != nil {
		params.Set("order_by", *req.OrderBy)
	}

	// Add filtering using helper
	if err := tiktok.AddJSONParam(params, "filtering", req.Filtering); err != nil {
		return nil, err
	}

	// Add pagination using helper
	tiktok.AddPagination(params, &tiktok.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	// Use generic DoGet helper
	var resp PixelListResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/pixel/list/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to list pixels: %w", err)
	}

	return &resp, nil
}

// OfflineEventSetInfo represents offline event set information
type OfflineEventSetInfo struct {
	EventSetID   string `json:"event_set_id"`
	Name         string `json:"name"`
	AdvertiserID string `json:"advertiser_id"`
	Status       string `json:"status"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}

// OfflineGetResponse represents the response for getting offline event sets
type OfflineGetResponse struct {
	List     []OfflineEventSetInfo `json:"list"`
	PageInfo tiktok.PageInfo       `json:"page_info,omitempty"`
}

// OfflineGetRequest represents the request to get offline event sets
type OfflineGetRequest struct {
	AdvertiserID string   `json:"advertiser_id,omitempty"`
	EventSetIDs  []string `json:"event_set_ids,omitempty"`
	Name         *string  `json:"name,omitempty"`
}

// GetOfflineEventSets gets Offline Event sets
// Reference: https://business-api.tiktok.com/portal/docs?id=1765596808589313
func (a *API) GetOfflineEventSets(ctx context.Context, req *OfflineGetRequest) (*OfflineGetResponse, error) {
	params := url.Values{}

	if req.AdvertiserID != "" {
		params.Set("advertiser_id", req.AdvertiserID)
	}

	// Add event_set_ids using helper
	if err := tiktok.AddStringSlice(params, "event_set_ids", req.EventSetIDs); err != nil {
		return nil, err
	}

	if req.Name != nil {
		params.Set("name", *req.Name)
	}

	// Use generic DoGet helper
	var resp OfflineGetResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/offline/get/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get offline event sets: %w", err)
	}

	return &resp, nil
}
