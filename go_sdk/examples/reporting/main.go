package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/reporting"
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

	// Create Reporting API client
	reportingAPI := reporting.NewAPI(client)

	ctx := context.Background()

	// Get integrated report
	fmt.Println("=== Getting Integrated Report ===")

	// Define the metrics you want to retrieve
	metrics := []string{
		"spend",
		"impressions",
		"clicks",
		"ctr",
		"cpc",
		"cpm",
		"conversion",
		"cost_per_conversion",
	}

	// Define the date range
	startDate := os.Getenv("TIKTOK_START_DATE") // Format: YYYY-MM-DD
	endDate := os.Getenv("TIKTOK_END_DATE")     // Format: YYYY-MM-DD

	if startDate == "" || endDate == "" {
		log.Println("TIKTOK_START_DATE and TIKTOK_END_DATE not set, using default values")
		startDate = "2024-01-01"
		endDate = "2024-01-31"
	}

	serviceType := "AUCTION"
	dataLevel := "AUCTION_CAMPAIGN"

	req := &reporting.IntegratedGetRequest{
		ReportType:   "BASIC",
		AdvertiserID: &advertiserID,
		ServiceType:  &serviceType,
		DataLevel:    &dataLevel,
		Dimensions:   []string{"campaign_id", "stat_time_day"},
		Metrics:      metrics,
		StartDate:    &startDate,
		EndDate:      &endDate,
	}

	resp, err := reportingAPI.GetIntegratedReport(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get report: %v", err)
	}

	fmt.Printf("Page Info - Total: %d, Page: %d, PageSize: %d\n",
		resp.PageInfo.TotalNumber, resp.PageInfo.Page, resp.PageInfo.PageSize)

	if len(resp.List) > 0 {
		fmt.Printf("\nFound %d report entries:\n", len(resp.List))
		for i, entry := range resp.List {
			if i >= 5 { // Show first 5 entries only
				break
			}
			fmt.Printf("\nEntry %d: %v\n", i+1, entry)
		}
	} else {
		fmt.Println("No report data available for the specified period")
	}

	// Example with filtering
	campaignID := os.Getenv("TIKTOK_CAMPAIGN_ID")
	if campaignID != "" {
		fmt.Println("\n=== Getting Campaign-specific Report ===")
		filtering := map[string]interface{}{
			"campaign_ids": []string{campaignID},
		}

		adDataLevel := "AUCTION_AD"
		reqFiltered := &reporting.IntegratedGetRequest{
			ReportType:   "BASIC",
			AdvertiserID: &advertiserID,
			ServiceType:  &serviceType,
			DataLevel:    &adDataLevel,
			Dimensions:   []string{"ad_id", "stat_time_day"},
			Metrics:      metrics,
			Filtering:    filtering,
			StartDate:    &startDate,
			EndDate:      &endDate,
		}

		respFiltered, err := reportingAPI.GetIntegratedReport(ctx, reqFiltered)
		if err != nil {
			log.Fatalf("Failed to get filtered report: %v", err)
		}

		fmt.Printf("Found %d ads in campaign %s\n", len(respFiltered.List), campaignID)
	}
}
