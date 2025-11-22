package authentication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Authentication API client
type API struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPI creates a new Authentication API client
// Note: Authentication endpoints don't require an access token in the client
func NewAPI() *API {
	return &API{
		baseURL: "https://business-api.tiktok.com",
		httpClient: &http.Client{
			Timeout: 30 * 1000000000, // 30 seconds
		},
	}
}

// NewAPIWithConfig creates a new Authentication API client with custom configuration
func NewAPIWithConfig(baseURL string, httpClient *http.Client) *API {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * 1000000000, // 30 seconds
		}
	}
	if baseURL == "" {
		baseURL = "https://business-api.tiktok.com"
	}
	return &API{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// AccessTokenRequest represents the request to get access token
type AccessTokenRequest struct {
	AppID    string `json:"app_id"`
	AuthCode string `json:"auth_code"`
	Secret   string `json:"secret"`
}

// AccessTokenResponse represents the response for access token request
type AccessTokenResponse struct {
	AccessToken           string   `json:"access_token"`
	AdvertiserIDs         []string `json:"advertiser_ids,omitempty"`
	AdvertiserID          string   `json:"advertiser_id,omitempty"`
	RefreshToken          string   `json:"refresh_token"`
	ExpiresIn             int64    `json:"expires_in"`
	RefreshTokenExpiresIn int64    `json:"refresh_token_expires_in"`
	TokenType             string   `json:"token_type,omitempty"`
	Scope                 string   `json:"scope,omitempty"`
}

// GetAccessToken gets access_token and refresh_token by auth_code.
// The creator access token is valid for 24 hours and the refresh token is valid for one year.
// Within one year you will need to refresh the access token with the refresh token on a daily basis.
// After one year you will need to ask the creator to reauthorize.
// Reference: https://ads.tiktok.com/marketing_api/docs?id=1739965703387137
func (a *API) GetAccessToken(ctx context.Context, req *AccessTokenRequest) (*AccessTokenResponse, error) {
	// Build URL
	fullURL := a.baseURL + "/open_api/v1.3/oauth2/access_token/"

	// Marshal request body
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var apiResp tiktok.Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors
	if apiResp.Code != nil && *apiResp.Code != 0 {
		errResp := &tiktok.ErrorResponse{
			Code:    *apiResp.Code,
			Message: *apiResp.Message,
		}
		if apiResp.RequestID != nil {
			errResp.RequestID = *apiResp.RequestID
		}
		return nil, errResp
	}

	// Unmarshal data
	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(apiResp.Data, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access token response: %w", err)
	}

	return &tokenResp, nil
}

// AdvertiserInfo represents advertiser information in the OAuth response
type AdvertiserInfo struct {
	AdvertiserID   string `json:"advertiser_id"`
	AdvertiserName string `json:"advertiser_name"`
}

// GetAdvertisersResponse represents the response for getting authorized advertisers
type GetAdvertisersResponse struct {
	List []AdvertiserInfo `json:"list"`
}

// GetAdvertisers gets a list of advertisers that have granted you permission to manage their accounts.
// Reference: https://business-api.tiktok.com/portal/docs?id=1738455508553729
func (a *API) GetAdvertisers(ctx context.Context, appID, secret, accessToken string) (*GetAdvertisersResponse, error) {
	// Build URL
	fullURL := a.baseURL + "/open_api/v1.3/oauth2/advertiser/get/"

	// Add query parameters
	params := url.Values{}
	params.Set("app_id", appID)
	params.Set("secret", secret)
	fullURL += "?" + params.Encode()

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Access-Token header
	httpReq.Header.Set("Access-Token", accessToken)

	// Execute request
	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var apiResp tiktok.Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors
	if apiResp.Code != nil && *apiResp.Code != 0 {
		errResp := &tiktok.ErrorResponse{
			Code:    *apiResp.Code,
			Message: *apiResp.Message,
		}
		if apiResp.RequestID != nil {
			errResp.RequestID = *apiResp.RequestID
		}
		return nil, errResp
	}

	// Unmarshal data
	var advertisersResp GetAdvertisersResponse
	if err := json.Unmarshal(apiResp.Data, &advertisersResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal advertisers response: %w", err)
	}

	return &advertisersResp, nil
}
