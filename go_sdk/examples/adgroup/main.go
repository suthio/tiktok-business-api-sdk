package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/adgroup"
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

	// Create Adgroup API client
	adgroupAPI := adgroup.NewAPI(client)

	ctx := context.Background()

	// Example 1: Get all adgroups
	fmt.Println("=== Example 1: Get All Adgroups ===")
	page := int64(1)
	pageSize := int64(10)
	req := &adgroup.GetAdGroupRequest{
		AdvertiserID: advertiserID,
		Page:         &page,
		PageSize:     &pageSize,
	}

	resp, err := adgroupAPI.GetAdGroups(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get adgroups: %v", err)
	}

	fmt.Printf("Found %d adgroup(s)\n", len(resp.List))
	for _, ag := range resp.List {
		fmt.Printf("\nAdgroup ID: %s\n", ag.AdgroupID)
		fmt.Printf("Adgroup Name: %s\n", ag.AdgroupName)
		fmt.Printf("Campaign ID: %s\n", ag.CampaignID)
		fmt.Printf("Operation Status: %s\n", ag.OperationStatus)
		if ag.Budget > 0 {
			fmt.Printf("Budget: %.2f\n", ag.Budget)
		}
	}

	// Example 2: Get specific adgroups by campaign ID
	campaignID := os.Getenv("TIKTOK_CAMPAIGN_ID")
	if campaignID != "" {
		fmt.Println("\n=== Example 2: Get Adgroups by Campaign ===")
		reqCampaign := &adgroup.GetAdGroupRequest{
			AdvertiserID: advertiserID,
			Filtering: &adgroup.Filtering{
				CampaignIDs: []string{campaignID},
			},
		}

		respCampaign, err := adgroupAPI.GetAdGroups(ctx, reqCampaign)
		if err != nil {
			log.Fatalf("Failed to get adgroups by campaign: %v", err)
		}

		fmt.Printf("Found %d adgroup(s) in campaign %s\n", len(respCampaign.List), campaignID)
	}
}
