package research

import (
	"context"
	"fmt"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Research Adlib API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Research Adlib API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// AdReportData represents ad report data from Research API
type AdReportData struct {
	AdID             string  `json:"ad_id,omitempty"`
	AdName           string  `json:"ad_name,omitempty"`
	AdvertiserID     string  `json:"advertiser_id,omitempty"`
	AdvertiserName   string  `json:"advertiser_name,omitempty"`
	CampaignID       string  `json:"campaign_id,omitempty"`
	CampaignName     string  `json:"campaign_name,omitempty"`
	AdgroupID        string  `json:"adgroup_id,omitempty"`
	AdgroupName      string  `json:"adgroup_name,omitempty"`
	Country          string  `json:"country,omitempty"`
	Region           string  `json:"region,omitempty"`
	Language         string  `json:"language,omitempty"`
	Platform         string  `json:"platform,omitempty"`
	ObjectiveType    string  `json:"objective_type,omitempty"`
	CallToAction     string  `json:"call_to_action,omitempty"`
	VideoID          string  `json:"video_id,omitempty"`
	VideoTitle       string  `json:"video_title,omitempty"`
	VideoDuration    float64 `json:"video_duration,omitempty"`
	ThumbnailURL     string  `json:"thumbnail_url,omitempty"`
	LandingPageURL   string  `json:"landing_page_url,omitempty"`
	DisplayName      string  `json:"display_name,omitempty"`
	ProfileImage     string  `json:"profile_image,omitempty"`
	AdText           string  `json:"ad_text,omitempty"`
	Impressions      int64   `json:"impressions,omitempty"`
	Clicks           int64   `json:"clicks,omitempty"`
	CTR              float64 `json:"ctr,omitempty"`
	Reach            int64   `json:"reach,omitempty"`
	Frequency        float64 `json:"frequency,omitempty"`
	Likes            int64   `json:"likes,omitempty"`
	Comments         int64   `json:"comments,omitempty"`
	Shares           int64   `json:"shares,omitempty"`
	VideoViews       int64   `json:"video_views,omitempty"`
	VideoViewRate    float64 `json:"video_view_rate,omitempty"`
	AverageVideoPlay float64 `json:"average_video_play,omitempty"`
	FirstShownDate   string  `json:"first_shown_date,omitempty"`
	LastShownDate    string  `json:"last_shown_date,omitempty"`
	StatTimePeriod   string  `json:"stat_time_period,omitempty"`
}

// GetAdReportResponse represents the response for getting ad report
type GetAdReportResponse struct {
	List     []AdReportData  `json:"list"`
	PageInfo tiktok.PageInfo `json:"page_info"`
}

// AdReportFiltering represents filtering options for ad report
type AdReportFiltering struct {
	CountryCodes      []string `json:"country_codes,omitempty"`
	RegionIDs         []string `json:"region_ids,omitempty"`
	Languages         []string `json:"languages,omitempty"`
	Platforms         []string `json:"platforms,omitempty"`
	ObjectiveTypes    []string `json:"objective_types,omitempty"`
	AdvertiserIDs     []string `json:"advertiser_ids,omitempty"`
	AdvertiserName    *string  `json:"advertiser_name,omitempty"`
	AdText            *string  `json:"ad_text,omitempty"`
	VideoTitle        *string  `json:"video_title,omitempty"`
	FirstShownDateMin *string  `json:"first_shown_date_min,omitempty"`
	FirstShownDateMax *string  `json:"first_shown_date_max,omitempty"`
	LastShownDateMin  *string  `json:"last_shown_date_min,omitempty"`
	LastShownDateMax  *string  `json:"last_shown_date_max,omitempty"`
}

// GetAdReportRequest represents the request to get ad report from Research API
type GetAdReportRequest struct {
	// Required: The search term to query ads
	SearchTerm string `json:"search_term"`

	// Optional: Country code (ISO 3166-1 alpha-2)
	CountryCode *string `json:"country_code,omitempty"`

	// Optional: Filtering criteria
	Filtering *AdReportFiltering `json:"filtering,omitempty"`

	// Optional: Fields to return in response
	Fields []string `json:"fields,omitempty"`

	// Optional: Pagination
	Page     *int64 `json:"page,omitempty"`
	PageSize *int64 `json:"page_size,omitempty"`

	// Optional: Sorting
	OrderBy    *string `json:"order_by,omitempty"`
	OrderField *string `json:"order_field,omitempty"`
}

// GetAdReport gets ad report from TikTok Research Adlib API
// This API provides access to TikTok's ad library for research purposes
// Reference: https://business-api.tiktok.com/portal/docs?id=1758579480845313
func (a *API) GetAdReport(ctx context.Context, req *GetAdReportRequest) (*GetAdReportResponse, error) {
	params := url.Values{}
	params.Set("search_term", req.SearchTerm)

	// Add optional country code
	if req.CountryCode != nil {
		params.Set("country_code", *req.CountryCode)
	}

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

	// Add sorting parameters
	if req.OrderBy != nil {
		params.Set("order_by", *req.OrderBy)
	}
	if req.OrderField != nil {
		params.Set("order_field", *req.OrderField)
	}

	// Use generic DoGet helper
	var resp GetAdReportResponse
	if err := tiktok.DoGet(ctx, a.client, "/v2/research/adlib/ad/report/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get ad report: %w", err)
	}

	return &resp, nil
}

// GetAllAdReports retrieves all ad reports by automatically handling pagination
// This is a convenience method that calls GetAdReport multiple times if needed
func (a *API) GetAllAdReports(ctx context.Context, req *GetAdReportRequest) ([]AdReportData, error) {
	var allReports []AdReportData
	page := int64(1)
	pageSize := int64(100) // Use maximum page size for efficiency

	// Override pagination parameters in request
	req.Page = &page
	req.PageSize = &pageSize

	for {
		resp, err := a.GetAdReport(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to get ad report page %d: %w", page, err)
		}

		allReports = append(allReports, resp.List...)

		// Check if we've fetched all pages
		if page >= resp.PageInfo.TotalPage {
			break
		}

		// Move to next page
		page++
		req.Page = &page
	}

	return allReports, nil
}
