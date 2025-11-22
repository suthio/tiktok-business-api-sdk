package campaign

import (
	"context"
	"fmt"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Campaign API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Campaign API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// CampaignStatus represents the status of a campaign
type CampaignStatus struct {
	CampaignID      string  `json:"campaign_id"`
	CampaignName    string  `json:"campaign_name"`
	AdvertiserID    string  `json:"advertiser_id"`
	ObjectiveType   string  `json:"objective_type"`
	Budget          float64 `json:"budget"`
	BudgetMode      string  `json:"budget_mode"`
	OperationStatus string  `json:"operation_status"`
	CreateTime      string  `json:"create_time"`
	ModifyTime      string  `json:"modify_time"`
}

// GetCampaignResponse represents the response for getting campaigns
type GetCampaignResponse struct {
	List     []CampaignStatus `json:"list"`
	PageInfo tiktok.PageInfo  `json:"page_info"`
}

// GetCampaignRequest represents the request to get campaigns
type GetCampaignRequest struct {
	AdvertiserID string     `json:"advertiser_id"`
	Filtering    *Filtering `json:"filtering,omitempty"`
	Page         *int64     `json:"page,omitempty"`
	PageSize     *int64     `json:"page_size,omitempty"`
}

// Filtering represents filtering options
type Filtering struct {
	CampaignIDs     []string `json:"campaign_ids,omitempty"`
	CampaignName    *string  `json:"campaign_name,omitempty"`
	ObjectiveType   *string  `json:"objective_type,omitempty"`
	PrimaryStatus   *string  `json:"primary_status,omitempty"`
	SecondaryStatus *string  `json:"secondary_status,omitempty"`
	CreateTimeMin   *string  `json:"create_time_min,omitempty"`
	CreateTimeMax   *string  `json:"create_time_max,omitempty"`
}

// GetCampaigns gets campaign information
// Reference: https://business-api.tiktok.com/portal/docs?id=1739315828649986
func (a *API) GetCampaigns(ctx context.Context, req *GetCampaignRequest) (*GetCampaignResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	// Add pagination using helper
	tiktok.AddPagination(params, &tiktok.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	// Add filtering using helper
	if err := tiktok.AddJSONParam(params, "filtering", req.Filtering); err != nil {
		return nil, err
	}

	// Use generic DoGet helper
	var resp GetCampaignResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/campaign/get/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}

	return &resp, nil
}

// CreateCampaignRequest represents the request to create a campaign
type CreateCampaignRequest struct {
	AdvertiserID      string   `json:"advertiser_id"`
	CampaignName      string   `json:"campaign_name"`
	ObjectiveType     string   `json:"objective_type"`
	AppID             *string  `json:"app_id,omitempty"`
	AppPromotionType  *string  `json:"app_promotion_type,omitempty"`
	Budget            *float64 `json:"budget,omitempty"`
	BudgetMode        *string  `json:"budget_mode,omitempty"`
	BudgetOptimizeOn  *bool    `json:"budget_optimize_on,omitempty"`
	CampaignType      *string  `json:"campaign_type,omitempty"`
	OperationStatus   *string  `json:"operation_status,omitempty"`
	OptimizationGoal  *string  `json:"optimization_goal,omitempty"`
	RfCampaignType    *string  `json:"rf_campaign_type,omitempty"`
	SpecialIndustries []string `json:"special_industries,omitempty"`
}

// CreateCampaignResponse represents the response from creating a campaign
type CreateCampaignResponse struct {
	CampaignID string `json:"campaign_id"`
}

// CreateCampaign creates a new campaign
// Reference: https://business-api.tiktok.com/portal/docs?id=1739318962329602
func (a *API) CreateCampaign(ctx context.Context, req *CreateCampaignRequest) (*CreateCampaignResponse, error) {
	// Use generic DoPost helper
	var resp CreateCampaignResponse
	if err := tiktok.DoPost(ctx, a.client, "/open_api/v1.3/campaign/create/", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return &resp, nil
}
