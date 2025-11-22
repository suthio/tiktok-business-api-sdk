package adgroup

import (
	"context"
	"fmt"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the AdGroup API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new AdGroup API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// AdGroupInfo represents ad group information
type AdGroupInfo struct {
	AdgroupID         string   `json:"adgroup_id"`
	AdgroupName       string   `json:"adgroup_name"`
	CampaignID        string   `json:"campaign_id"`
	AdvertiserID      string   `json:"advertiser_id"`
	ObjectiveType     string   `json:"objective_type,omitempty"`
	Budget            float64  `json:"budget,omitempty"`
	BudgetMode        string   `json:"budget_mode,omitempty"`
	BillingEvent      string   `json:"billing_event,omitempty"`
	OptimizationGoal  string   `json:"optimization_goal,omitempty"`
	Placements        []string `json:"placements,omitempty"`
	Locations         []string `json:"locations,omitempty"`
	Age               []string `json:"age,omitempty"`
	Gender            string   `json:"gender,omitempty"`
	Languages         []string `json:"languages,omitempty"`
	OperationStatus   string   `json:"operation_status"`
	PrimaryStatus     string   `json:"primary_status,omitempty"`
	SecondaryStatus   string   `json:"secondary_status,omitempty"`
	CreateTime        string   `json:"create_time"`
	ModifyTime        string   `json:"modify_time"`
	ScheduleStartTime string   `json:"schedule_start_time,omitempty"`
	ScheduleEndTime   string   `json:"schedule_end_time,omitempty"`
}

// GetAdGroupResponse represents the response for getting ad groups
type GetAdGroupResponse struct {
	List     []AdGroupInfo   `json:"list"`
	PageInfo tiktok.PageInfo `json:"page_info"`
}

// GetAdGroupRequest represents the request to get ad groups
type GetAdGroupRequest struct {
	AdvertiserID string     `json:"advertiser_id"`
	Filtering    *Filtering `json:"filtering,omitempty"`
	Page         *int64     `json:"page,omitempty"`
	PageSize     *int64     `json:"page_size,omitempty"`
	Fields       []string   `json:"fields,omitempty"`
}

// Filtering represents filtering options for ad groups
type Filtering struct {
	AdgroupIDs      []string `json:"adgroup_ids,omitempty"`
	CampaignIDs     []string `json:"campaign_ids,omitempty"`
	PrimaryStatus   *string  `json:"primary_status,omitempty"`
	SecondaryStatus *string  `json:"secondary_status,omitempty"`
	ObjectiveType   *string  `json:"objective_type,omitempty"`
	BillingEvent    *string  `json:"billing_event,omitempty"`
	CreateTimeMin   *string  `json:"create_time_min,omitempty"`
	CreateTimeMax   *string  `json:"create_time_max,omitempty"`
}

// GetAdGroups gets ad group information
// Reference: https://business-api.tiktok.com/portal/docs?id=1739314558673922
func (a *API) GetAdGroups(ctx context.Context, req *GetAdGroupRequest) (*GetAdGroupResponse, error) {
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
	var resp GetAdGroupResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/adgroup/get/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get ad groups: %w", err)
	}

	return &resp, nil
}

// CreateAdGroupRequest represents a simplified request to create an ad group
// For full field list, see: https://business-api.tiktok.com/portal/docs?id=1739499616346114
type CreateAdGroupRequest struct {
	AdvertiserID      string   `json:"advertiser_id"`
	CampaignID        string   `json:"campaign_id"`
	AdGroupName       string   `json:"adgroup_name"`
	PromotionType     *string  `json:"promotion_type,omitempty"`
	PlacementType     string   `json:"placement_type"`
	Placements        []string `json:"placements"`
	LocationIDs       []string `json:"location_ids"`
	Languages         []string `json:"languages,omitempty"`
	Gender            *string  `json:"gender,omitempty"`
	AgeGroups         []string `json:"age_groups,omitempty"`
	BudgetMode        string   `json:"budget_mode"`
	Budget            *float64 `json:"budget,omitempty"`
	ScheduleType      *string  `json:"schedule_type,omitempty"`
	ScheduleStartTime *string  `json:"schedule_start_time,omitempty"`
	ScheduleEndTime   *string  `json:"schedule_end_time,omitempty"`
	BillingEvent      string   `json:"billing_event"`
	BidPrice          *float64 `json:"bid_price,omitempty"`
	OptimizationGoal  string   `json:"optimization_goal"`
	Pacing            *string  `json:"pacing,omitempty"`
	PixelID           *string  `json:"pixel_id,omitempty"`
	OperationStatus   *string  `json:"operation_status,omitempty"`
}

// CreateAdGroupResponse represents the response from creating an ad group
type CreateAdGroupResponse struct {
	AdGroupID string `json:"adgroup_id"`
}

// CreateAdGroup creates a new ad group with simplified parameters
// Reference: https://business-api.tiktok.com/portal/docs?id=1739499616346114
func (a *API) CreateAdGroup(ctx context.Context, req *CreateAdGroupRequest) (*CreateAdGroupResponse, error) {
	// Use generic DoPost helper
	var resp CreateAdGroupResponse
	if err := tiktok.DoPost(ctx, a.client, "/open_api/v1.3/adgroup/create/", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create ad group: %w", err)
	}

	return &resp, nil
}
