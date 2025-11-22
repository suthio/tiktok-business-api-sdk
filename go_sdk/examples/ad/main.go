package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/ad"
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

	// Create Ad API client
	adAPI := ad.NewAPI(client)

	ctx := context.Background()

	// Example 1: Get all ads
	fmt.Println("=== Example 1: Get All Ads ===")
	page := int64(1)
	pageSize := int64(10)
	req := &ad.GetAdRequest{
		AdvertiserID: advertiserID,
		Page:         &page,
		PageSize:     &pageSize,
	}

	resp, err := adAPI.GetAds(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get ads: %v", err)
	}

	fmt.Printf("Found %d ad(s)\n", len(resp.List))
	for _, adInfo := range resp.List {
		fmt.Printf("\nAd ID: %s\n", adInfo.AdID)
		fmt.Printf("Ad Name: %s\n", adInfo.AdName)
		fmt.Printf("Campaign ID: %s\n", adInfo.CampaignID)
		fmt.Printf("Adgroup ID: %s\n", adInfo.AdgroupID)
		fmt.Printf("Operation Status: %s\n", adInfo.OperationStatus)
		if adInfo.VideoID != "" {
			fmt.Printf("Video ID: %s\n", adInfo.VideoID)
		}
	}

	// Example 2: Get specific ads by ID
	if len(resp.List) > 0 {
		fmt.Println("\n=== Example 2: Get Specific Ad by ID ===")
		specificAdID := resp.List[0].AdID
		reqSpecific := &ad.GetAdRequest{
			AdvertiserID: advertiserID,
			Filtering: &ad.Filtering{
				AdIDs: []string{specificAdID},
			},
		}

		respSpecific, err := adAPI.GetAds(ctx, reqSpecific)
		if err != nil {
			log.Fatalf("Failed to get specific ad: %v", err)
		}

		if len(respSpecific.List) > 0 {
			adInfo := respSpecific.List[0]
			fmt.Printf("Ad ID: %s\n", adInfo.AdID)
			fmt.Printf("Ad Name: %s\n", adInfo.AdName)
			fmt.Printf("Ad Text: %s\n", adInfo.AdText)
			fmt.Printf("Call to Action: %s\n", adInfo.CallToAction)
		}
	}

	// Example 3: Filter ads by campaign
	campaignID := os.Getenv("TIKTOK_CAMPAIGN_ID")
	if campaignID != "" {
		fmt.Println("\n=== Example 3: Get Ads by Campaign ===")
		reqCampaign := &ad.GetAdRequest{
			AdvertiserID: advertiserID,
			Filtering: &ad.Filtering{
				CampaignIDs: []string{campaignID},
			},
		}

		respCampaign, err := adAPI.GetAds(ctx, reqCampaign)
		if err != nil {
			log.Fatalf("Failed to get ads by campaign: %v", err)
		}

		fmt.Printf("Found %d ad(s) in campaign %s\n", len(respCampaign.List), campaignID)
	}
}
