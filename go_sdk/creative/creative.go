package creative

import (
	"context"
	"fmt"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Creative API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Creative API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// CreativeInfo represents creative information
type CreativeInfo struct {
	CreativeID       string   `json:"creative_id"`
	CreativeName     string   `json:"creative_name,omitempty"`
	AdID             string   `json:"ad_id,omitempty"`
	AdgroupID        string   `json:"adgroup_id,omitempty"`
	CampaignID       string   `json:"campaign_id,omitempty"`
	AdvertiserID     string   `json:"advertiser_id"`
	CreativeType     string   `json:"creative_type,omitempty"`
	ImageIDs         []string `json:"image_ids,omitempty"`
	VideoID          string   `json:"video_id,omitempty"`
	AdText           string   `json:"ad_text,omitempty"`
	AdFormat         string   `json:"ad_format,omitempty"`
	CallToAction     string   `json:"call_to_action,omitempty"`
	LandingPageURL   string   `json:"landing_page_url,omitempty"`
	DisplayName      string   `json:"display_name,omitempty"`
	IdentityID       string   `json:"identity_id,omitempty"`
	IdentityType     string   `json:"identity_type,omitempty"`
	CardID           string   `json:"card_id,omitempty"`
	OperationStatus  string   `json:"operation_status,omitempty"`
	CreateTime       string   `json:"create_time,omitempty"`
	ModifyTime       string   `json:"modify_time,omitempty"`
	VideoViewTrackingURL string `json:"video_view_tracking_url,omitempty"`
	ClickTrackingURL string   `json:"click_tracking_url,omitempty"`
	ImpressionTrackingURL string `json:"impression_tracking_url,omitempty"`
}

// GetCreativesResponse represents the response for getting creatives
type GetCreativesResponse struct {
	List     []CreativeInfo  `json:"list"`
	PageInfo tiktok.PageInfo `json:"page_info"`
}

// GetCreativesRequest represents the request to get creatives
type GetCreativesRequest struct {
	AdvertiserID string     `json:"advertiser_id"`
	Filtering    *Filtering `json:"filtering,omitempty"`
	Page         *int64     `json:"page,omitempty"`
	PageSize     *int64     `json:"page_size,omitempty"`
	Fields       []string   `json:"fields,omitempty"`
}

// Filtering represents filtering options for creatives
type Filtering struct {
	CreativeIDs     []string `json:"creative_ids,omitempty"`
	AdIDs           []string `json:"ad_ids,omitempty"`
	AdgroupIDs      []string `json:"adgroup_ids,omitempty"`
	CampaignIDs     []string `json:"campaign_ids,omitempty"`
	CreativeType    *string  `json:"creative_type,omitempty"`
	ObjectiveType   *string  `json:"objective_type,omitempty"`
	OperationStatus *string  `json:"operation_status,omitempty"`
	CreateTimeMin   *string  `json:"create_time_min,omitempty"`
	CreateTimeMax   *string  `json:"create_time_max,omitempty"`
}

// GetCreatives gets creative information
// Reference: https://business-api.tiktok.com/portal/docs?id=1740051721711618
func (a *API) GetCreatives(ctx context.Context, req *GetCreativesRequest) (*GetCreativesResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	// Add pagination using helper
	tiktok.AddPagination(params, &tiktok.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	// Add fields using helper
	if err := tiktok.AddStringSlice(params, "fields", req.Fields); err != nil {
		return nil, err
	}

	// Add filtering using helper
	if err := tiktok.AddJSONParam(params, "filtering", req.Filtering); err != nil {
		return nil, err
	}

	// Use generic DoGet helper
	var resp GetCreativesResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/creative/get/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get creatives: %w", err)
	}

	return &resp, nil
}

// GetAllCreatives retrieves all creatives by automatically handling pagination
// This is a convenience method that calls GetCreatives multiple times if needed
func (a *API) GetAllCreatives(ctx context.Context, req *GetCreativesRequest) ([]CreativeInfo, error) {
	var allCreatives []CreativeInfo
	page := int64(1)
	pageSize := int64(100) // Use maximum page size for efficiency

	// Override pagination parameters in request
	req.Page = &page
	req.PageSize = &pageSize

	for {
		resp, err := a.GetCreatives(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to get creatives page %d: %w", page, err)
		}

		allCreatives = append(allCreatives, resp.List...)

		// Check if we've fetched all pages
		if page >= resp.PageInfo.TotalPage {
			break
		}

		// Move to next page
		page++
		req.Page = &page
	}

	return allCreatives, nil
}
