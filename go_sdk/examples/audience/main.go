package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/audience"
)

func main() {
	// Get access token and advertiser ID from environment variables
	accessToken := os.Getenv("TIKTOK_ACCESS_TOKEN")
	advertiserID := os.Getenv("TIKTOK_ADVERTISER_ID")
	if accessToken == "" || advertiserID == "" {
		log.Fatal("TIKTOK_ACCESS_TOKEN and TIKTOK_ADVERTISER_ID environment variables must be set")
	}

	// Create TikTok client
	client := tiktok.NewClient(accessToken)

	// Create Audience API client
	audienceAPI := audience.NewAPI(client)

	ctx := context.Background()

	// List custom audiences
	fmt.Println("=== Listing Custom Audiences ===")
	page := int64(1)
	pageSize := int64(20)
	req := &audience.CustomAudienceListRequest{
		AdvertiserID: advertiserID,
		Page:         &page,
		PageSize:     &pageSize,
	}

	resp, err := audienceAPI.ListCustomAudiences(ctx, req)
	if err != nil {
		log.Fatalf("Failed to list custom audiences: %v", err)
	}

	fmt.Printf("Found %d custom audience(s)\n", len(resp.List))
	for _, aud := range resp.List {
		fmt.Printf("\nAudience ID: %s\n", aud.CustomAudienceID)
		fmt.Printf("Audience Name: %s\n", aud.Name)
		fmt.Printf("Audience Type: %s\n", aud.AudienceType)
		fmt.Printf("Size: %d\n", aud.Size)
		fmt.Printf("Status: %s\n", aud.Status)
		fmt.Printf("Create Time: %s\n", aud.CreateTime)
	}
}
