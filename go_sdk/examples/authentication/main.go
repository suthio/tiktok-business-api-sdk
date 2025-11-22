package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/suthio/tiktok-business-api-sdk/go_sdk/authentication"
)

func main() {
	// Create Authentication API client (no access token needed for auth)
	authAPI := authentication.NewAPI()

	// Get OAuth2 Access Token
	ctx := context.Background()
	req := &authentication.AccessTokenRequest{
		AppID:    os.Getenv("TIKTOK_APP_ID"),
		Secret:   os.Getenv("TIKTOK_SECRET"),
		AuthCode: os.Getenv("TIKTOK_AUTH_CODE"),
	}

	if req.AppID == "" || req.Secret == "" || req.AuthCode == "" {
		log.Fatal("TIKTOK_APP_ID, TIKTOK_SECRET, and TIKTOK_AUTH_CODE environment variables must be set")
	}

	resp, err := authAPI.GetAccessToken(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get OAuth2 access token: %v", err)
	}

	fmt.Printf("Access Token: %s\n", resp.AccessToken)
	fmt.Printf("Advertiser IDs: %v\n", resp.AdvertiserIDs)
	fmt.Printf("Expires In: %d seconds (24 hours)\n", resp.ExpiresIn)
	fmt.Printf("Refresh Token: %s\n", resp.RefreshToken)
	fmt.Printf("Refresh Token Expires In: %d seconds (1 year)\n", resp.RefreshTokenExpiresIn)

	// Example: Refresh the access token
	// This should be done when the access token expires (after 24 hours)
	fmt.Println("\n--- Refreshing Access Token ---")

	refreshReq := &authentication.RefreshTokenRequest{
		AppID:        req.AppID,
		Secret:       req.Secret,
		RefreshToken: resp.RefreshToken,
	}

	refreshResp, err := authAPI.RefreshToken(ctx, refreshReq)
	if err != nil {
		log.Fatalf("Failed to refresh access token: %v", err)
	}

	fmt.Printf("New Access Token: %s\n", refreshResp.AccessToken)
	fmt.Printf("New Refresh Token: %s\n", refreshResp.RefreshToken)
	fmt.Printf("Expires In: %d seconds (24 hours)\n", refreshResp.ExpiresIn)
	fmt.Println("\nNote: The refresh token is valid for 1 year. You should refresh daily.")
	fmt.Println("After 1 year, you need to ask the user to reauthorize.")

	// Example: Get Advertiser List
	fmt.Println("\n--- Getting Advertiser List ---")

	advertisers, err := authAPI.GetAdvertisers(ctx, req.AppID, req.Secret, refreshResp.AccessToken)
	if err != nil {
		log.Fatalf("Failed to get advertisers: %v", err)
	}

	fmt.Printf("Found %d advertisers:\n", len(advertisers.List))
	for i, adv := range advertisers.List {
		fmt.Printf("%d. ID: %s, Name: %s\n", i+1, adv.AdvertiserID, adv.AdvertiserName)
	}
}
