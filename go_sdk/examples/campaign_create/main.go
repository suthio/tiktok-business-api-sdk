package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	// Create a new campaign
	fmt.Println("=== Creating New Campaign ===")
	campaignName := fmt.Sprintf("Test Campaign %s", time.Now().Format("20060102150405"))
	budgetMode := "BUDGET_MODE_INFINITE"

	createReq := &campaign.CreateCampaignRequest{
		AdvertiserID:  advertiserID,
		CampaignName:  campaignName,
		ObjectiveType: "TRAFFIC",
		BudgetMode:    &budgetMode,
	}

	createResp, err := campaignAPI.CreateCampaign(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create campaign: %v", err)
	}

	fmt.Printf("Successfully created campaign!\n")
	fmt.Printf("Campaign ID: %s\n", createResp.CampaignID)
	fmt.Printf("Campaign Name: %s\n", campaignName)

	// Verify creation by fetching the campaign
	fmt.Println("\n=== Verifying Campaign Creation ===")
	getReq := &campaign.GetCampaignRequest{
		AdvertiserID: advertiserID,
		Filtering: &campaign.Filtering{
			CampaignIDs: []string{createResp.CampaignID},
		},
	}

	getResp, err := campaignAPI.GetCampaigns(ctx, getReq)
	if err != nil {
		log.Fatalf("Failed to get campaign: %v", err)
	}

	if len(getResp.List) > 0 {
		c := getResp.List[0]
		fmt.Printf("Campaign ID: %s\n", c.CampaignID)
		fmt.Printf("Campaign Name: %s\n", c.CampaignName)
		fmt.Printf("Objective Type: %s\n", c.ObjectiveType)
		fmt.Printf("Budget Mode: %s\n", c.BudgetMode)
		fmt.Printf("Operation Status: %s\n", c.OperationStatus)
	}
}
