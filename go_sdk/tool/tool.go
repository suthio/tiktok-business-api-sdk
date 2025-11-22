package tool

import (
	"context"
	"fmt"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Tool API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Tool API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// CarrierResponse represents the response for the carrier endpoint
type CarrierResponse struct {
	Carriers []Carrier `json:"carriers"`
}

// Carrier represents a carrier
type Carrier struct {
	CarrierID   string `json:"carrier_id"`
	CarrierName string `json:"carrier_name"`
}

// GetCarrier gets carriers for the specified advertiser
// Reference: https://business-api.tiktok.com/portal/docs?id=1737168013095938
func (a *API) GetCarrier(ctx context.Context, advertiserID string) (*CarrierResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", advertiserID)

	// Use generic DoGet helper
	var resp CarrierResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/tool/carrier/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get carriers: %w", err)
	}

	return &resp, nil
}

// LanguageResponse represents the response for the language endpoint
type LanguageResponse struct {
	Languages []Language `json:"languages"`
}

// Language represents a language
type Language struct {
	LanguageCode string `json:"language_code"`
	LanguageName string `json:"language_name"`
}

// GetLanguage gets supported languages
// Reference: https://business-api.tiktok.com/portal/docs?id=1737188554152962
func (a *API) GetLanguage(ctx context.Context, advertiserID string) (*LanguageResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", advertiserID)

	// Use generic DoGet helper
	var resp LanguageResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/tool/language/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	return &resp, nil
}

// ActionCategoryResponse represents the response for the action category endpoint
type ActionCategoryResponse struct {
	ActionCategories []ActionCategory `json:"action_categories"`
}

// ActionCategory represents an action category
type ActionCategory struct {
	ActionCategoryID   string `json:"action_category_id"`
	ActionCategoryName string `json:"action_category_name"`
}

// GetActionCategory gets action categories
// Reference: https://business-api.tiktok.com/portal/docs?id=1737166752522241
func (a *API) GetActionCategory(ctx context.Context, advertiserID string, specialIndustries []string) (*ActionCategoryResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", advertiserID)

	if len(specialIndustries) > 0 {
		for _, industry := range specialIndustries {
			params.Add("special_industries", industry)
		}
	}

	// Use generic DoGet helper
	var resp ActionCategoryResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/tool/action_category/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get action categories: %w", err)
	}

	return &resp, nil
}
