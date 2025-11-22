package ad

import (
	"context"
	"fmt"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Ad API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Ad API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// AdInfo represents ad information
type AdInfo struct {
	AdID            string   `json:"ad_id"`
	AdName          string   `json:"ad_name"`
	AdgroupID       string   `json:"adgroup_id"`
	CampaignID      string   `json:"campaign_id"`
	AdvertiserID    string   `json:"advertiser_id"`
	ImageIDs        []string `json:"image_ids,omitempty"`
	VideoID         string   `json:"video_id,omitempty"`
	AdText          string   `json:"ad_text,omitempty"`
	CallToAction    string   `json:"call_to_action,omitempty"`
	OperationStatus string   `json:"operation_status"`
	PrimaryStatus   string   `json:"primary_status,omitempty"`
	SecondaryStatus string   `json:"secondary_status,omitempty"`
	CreateTime      string   `json:"create_time"`
	ModifyTime      string   `json:"modify_time"`
}

// GetAdResponse represents the response for getting ads
type GetAdResponse struct {
	List     []AdInfo        `json:"list"`
	PageInfo tiktok.PageInfo `json:"page_info"`
}

// GetAdRequest represents the request to get ads
type GetAdRequest struct {
	AdvertiserID string     `json:"advertiser_id"`
	Filtering    *Filtering `json:"filtering,omitempty"`
	Page         *int64     `json:"page,omitempty"`
	PageSize     *int64     `json:"page_size,omitempty"`
	Fields       []string   `json:"fields,omitempty"`
}

// Filtering represents filtering options for ads
type Filtering struct {
	AdIDs           []string `json:"ad_ids,omitempty"`
	AdgroupIDs      []string `json:"adgroup_ids,omitempty"`
	CampaignIDs     []string `json:"campaign_ids,omitempty"`
	PrimaryStatus   *string  `json:"primary_status,omitempty"`
	SecondaryStatus *string  `json:"secondary_status,omitempty"`
	ObjectiveType   *string  `json:"objective_type,omitempty"`
	CreateTimeMin   *string  `json:"create_time_min,omitempty"`
	CreateTimeMax   *string  `json:"create_time_max,omitempty"`
}

// GetAds gets ad information
// Reference: https://business-api.tiktok.com/portal/docs?id=1735735588640770
func (a *API) GetAds(ctx context.Context, req *GetAdRequest) (*GetAdResponse, error) {
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
	var resp GetAdResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/ad/get/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get ads: %w", err)
	}

	return &resp, nil
}

// AdCreative represents a creative for an ad
type AdCreative struct {
	AdName         string   `json:"ad_name"`
	AdText         string   `json:"ad_text"`
	AdFormat       string   `json:"ad_format"`
	VideoID        *string  `json:"video_id,omitempty"`
	ImageIDs       []string `json:"image_ids,omitempty"`
	CallToAction   *string  `json:"call_to_action,omitempty"`
	DisplayName    *string  `json:"display_name,omitempty"`
	LandingPageURL *string  `json:"landing_page_url,omitempty"`
	IdentityID     *string  `json:"identity_id,omitempty"`
	IdentityType   *string  `json:"identity_type,omitempty"`
}

// CreateAdRequest represents a simplified request to create an ad
// For full field list, see: https://business-api.tiktok.com/portal/docs?id=1737172488964097
type CreateAdRequest struct {
	AdvertiserID    string       `json:"advertiser_id"`
	AdGroupID       string       `json:"adgroup_id"`
	Creatives       []AdCreative `json:"creatives"`
	OperationStatus *string      `json:"operation_status,omitempty"`
	IdentityID      *string      `json:"identity_id,omitempty"`
	IdentityType    *string      `json:"identity_type,omitempty"`
}

// CreateAdResponse represents the response from creating an ad
type CreateAdResponse struct {
	AdID string `json:"ad_id"`
}

// CreateAd creates a new ad with simplified parameters
// Reference: https://business-api.tiktok.com/portal/docs?id=1737172488964097
func (a *API) CreateAd(ctx context.Context, req *CreateAdRequest) (*CreateAdResponse, error) {
	// Use generic DoPost helper
	var resp CreateAdResponse
	if err := tiktok.DoPost(ctx, a.client, "/open_api/v1.3/ad/create/", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create ad: %w", err)
	}

	return &resp, nil
}
