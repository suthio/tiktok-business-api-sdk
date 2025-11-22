package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/creative"
)

func main() {
	// Get credentials from environment variables
	accessToken := os.Getenv("TIKTOK_ACCESS_TOKEN")
	advertiserID := os.Getenv("TIKTOK_ADVERTISER_ID")

	if accessToken == "" || advertiserID == "" {
		log.Fatal("Please set TIKTOK_ACCESS_TOKEN and TIKTOK_ADVERTISER_ID environment variables")
	}

	// Create client
	client := tiktok.NewClient(accessToken)

	// Create Creative API instance
	creativeAPI := creative.NewAPI(client)

	ctx := context.Background()

	// Example 1: Get creatives with pagination
	fmt.Println("Example 1: Get creatives with pagination")
	page := int64(1)
	pageSize := int64(10)
	resp, err := creativeAPI.GetCreatives(ctx, &creative.GetCreativesRequest{
		AdvertiserID: advertiserID,
		Page:         &page,
		PageSize:     &pageSize,
	})
	if err != nil {
		log.Printf("Warning: Failed to get creatives: %v", err)
		fmt.Println("Note: This may be expected if the sandbox account has no creatives yet.")
		fmt.Println("Skipping remaining creative examples.")
		return
	}

	fmt.Printf("Total creatives: %d\n", resp.PageInfo.TotalNumber)
	for i, c := range resp.List {
		fmt.Printf("%d. Creative ID: %s, Name: %s, Type: %s\n",
			i+1, c.CreativeID, c.CreativeName, c.CreativeType)
	}
	fmt.Println()

	// Example 2: Get creatives with filtering
	fmt.Println("Example 2: Get creatives with filtering by video creatives")
	creativeType := "VIDEO"
	respFiltered, err := creativeAPI.GetCreatives(ctx, &creative.GetCreativesRequest{
		AdvertiserID: advertiserID,
		Filtering: &creative.Filtering{
			CreativeType: &creativeType,
		},
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil {
		log.Fatalf("Failed to get filtered creatives: %v", err)
	}

	fmt.Printf("Total video creatives: %d\n", respFiltered.PageInfo.TotalNumber)
	for i, c := range respFiltered.List {
		fmt.Printf("%d. Creative ID: %s, Video ID: %s, Ad Text: %s\n",
			i+1, c.CreativeID, c.VideoID, c.AdText)
	}
	fmt.Println()

	// Example 3: Get all creatives (auto-pagination)
	fmt.Println("Example 3: Get all creatives with auto-pagination")
	allCreatives, err := creativeAPI.GetAllCreatives(ctx, &creative.GetCreativesRequest{
		AdvertiserID: advertiserID,
	})
	if err != nil {
		log.Fatalf("Failed to get all creatives: %v", err)
	}

	fmt.Printf("Retrieved %d creatives in total\n", len(allCreatives))
	fmt.Println()

	// Example 4: Get creatives by specific ad IDs
	fmt.Println("Example 4: Get creatives by specific ad IDs")
	// Note: Replace with actual ad IDs from your account
	adIDs := []string{"your_ad_id_1", "your_ad_id_2"}
	respByAds, err := creativeAPI.GetCreatives(ctx, &creative.GetCreativesRequest{
		AdvertiserID: advertiserID,
		Filtering: &creative.Filtering{
			AdIDs: adIDs,
		},
	})
	if err != nil {
		log.Fatalf("Failed to get creatives by ad IDs: %v", err)
	}

	fmt.Printf("Found %d creatives for specified ads\n", len(respByAds.List))
	for _, c := range respByAds.List {
		fmt.Printf("- Creative ID: %s, Ad ID: %s, Video ID: %s\n",
			c.CreativeID, c.AdID, c.VideoID)
	}
}
