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
	fmt.Printf("Expires In: %d seconds\n", resp.ExpiresIn)
	fmt.Printf("Refresh Token: %s\n", resp.RefreshToken)
}
