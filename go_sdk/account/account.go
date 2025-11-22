package account

import (
	"context"
	"fmt"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Account API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Account API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// AdvertiserInfo represents advertiser information
type AdvertiserInfo struct {
	AdvertiserID          string  `json:"advertiser_id"`
	AdvertiserName        string  `json:"name"`
	Address               string  `json:"address"`
	Brand                 string  `json:"brand"`
	Company               string  `json:"company"`
	ContactPerson         string  `json:"contact_person"`
	Country               string  `json:"country"`
	Currency              string  `json:"currency"`
	Description           string  `json:"description"`
	Email                 string  `json:"email"`
	Industry              string  `json:"industry"`
	Language              string  `json:"language"`
	LicenseNo             string  `json:"license_no"`
	PromotionArea         string  `json:"promotion_area"`
	PromotionCenterCity   string  `json:"promotion_center_city"`
	ReasonForAdvertising  string  `json:"reason_for_advertising"`
	Telephone             string  `json:"telephone"`
	Timezone              string  `json:"timezone"`
	DisplayTimezone       string  `json:"display_timezone"`
	AdvertiserAccountType string  `json:"advertiser_account_type"`
	BalanceMode           string  `json:"balance_mode"`
	CreateTime            int64   `json:"create_time"`
	Status                string  `json:"status"`
	Balance               float64 `json:"balance"`
}

// Balance represents balance information
type Balance struct {
	BalanceType string  `json:"balance_type"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
}

// AdvertiserInfoResponse represents the response for advertiser info
type AdvertiserInfoResponse struct {
	List []AdvertiserInfo `json:"list"`
}

// GetAdvertiserInfo gets advertiser information
// Reference: https://business-api.tiktok.com/portal/docs?id=1739593083610113
func (a *API) GetAdvertiserInfo(ctx context.Context, advertiserIDs []string, fields []string) (*AdvertiserInfoResponse, error) {
	params := url.Values{}

	// Add advertiser_ids using helper
	if err := tiktok.AddStringSlice(params, "advertiser_ids", advertiserIDs); err != nil {
		return nil, err
	}

	// Add fields using helper
	if err := tiktok.AddStringSlice(params, "fields", fields); err != nil {
		return nil, err
	}

	// Use generic DoGet helper
	var resp AdvertiserInfoResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/advertiser/info/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get advertiser info: %w", err)
	}

	return &resp, nil
}
