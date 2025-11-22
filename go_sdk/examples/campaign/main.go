package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/campaign"
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

	// Create Campaign API client
	campaignAPI := campaign.NewAPI(client)

	ctx := context.Background()

	// Example 1: Get all campaigns
	fmt.Println("=== Example 1: Get All Campaigns ===")
	page := int64(1)
	pageSize := int64(10)
	req := &campaign.GetCampaignRequest{
		AdvertiserID: advertiserID,
		Page:         &page,
		PageSize:     &pageSize,
	}

	resp, err := campaignAPI.GetCampaigns(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get campaigns: %v", err)
	}

	fmt.Printf("Found %d campaign(s)\n", len(resp.List))
	for _, c := range resp.List {
		fmt.Printf("\nCampaign ID: %s\n", c.CampaignID)
		fmt.Printf("Campaign Name: %s\n", c.CampaignName)
		fmt.Printf("Objective Type: %s\n", c.ObjectiveType)
		fmt.Printf("Operation Status: %s\n", c.OperationStatus)
		fmt.Printf("Budget: %.2f\n", c.Budget)
		fmt.Printf("Budget Mode: %s\n", c.BudgetMode)
	}

	// Example 2: Get specific campaign by ID
	campaignID := os.Getenv("TIKTOK_CAMPAIGN_ID")
	if campaignID != "" {
		fmt.Println("\n=== Example 2: Get Specific Campaign ===")
		reqSpecific := &campaign.GetCampaignRequest{
			AdvertiserID: advertiserID,
			Filtering: &campaign.Filtering{
				CampaignIDs: []string{campaignID},
			},
		}

		respSpecific, err := campaignAPI.GetCampaigns(ctx, reqSpecific)
		if err != nil {
			log.Fatalf("Failed to get specific campaign: %v", err)
		}

		if len(respSpecific.List) > 0 {
			c := respSpecific.List[0]
			fmt.Printf("Campaign ID: %s\n", c.CampaignID)
			fmt.Printf("Campaign Name: %s\n", c.CampaignName)
			fmt.Printf("Objective Type: %s\n", c.ObjectiveType)
			fmt.Printf("Budget: %.2f %s\n", c.Budget, "USD")
		}
	}
}
