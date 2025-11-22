package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/research"
)

func main() {
	// Get credentials from environment variables
	accessToken := os.Getenv("TIKTOK_ACCESS_TOKEN")

	if accessToken == "" {
		log.Fatal("Please set TIKTOK_ACCESS_TOKEN environment variable")
	}

	// Create client
	client := tiktok.NewClient(accessToken)

	// Create Research API instance
	researchAPI := research.NewAPI(client)

	ctx := context.Background()

	// Example 1: Basic search for ads
	fmt.Println("Example 1: Search for ads related to 'shoes'")
	page := int64(1)
	pageSize := int64(10)
	countryCode := "US"

	resp, err := researchAPI.GetAdReport(ctx, &research.GetAdReportRequest{
		SearchTerm:  "shoes",
		CountryCode: &countryCode,
		Page:        &page,
		PageSize:    &pageSize,
	})
	if err != nil {
		log.Printf("Warning: Failed to get ad report: %v", err)
		fmt.Println("Note: Research Adlib API may not be available in sandbox environment.")
		fmt.Println("This API typically requires production credentials and special access.")
		fmt.Println("Skipping remaining research examples.")
		return
	}

	fmt.Printf("Total ads found: %d\n", resp.PageInfo.TotalNumber)
	for i, ad := range resp.List {
		fmt.Printf("%d. Ad ID: %s\n", i+1, ad.AdID)
		fmt.Printf("   Advertiser: %s (%s)\n", ad.AdvertiserName, ad.AdvertiserID)
		fmt.Printf("   Ad Text: %s\n", ad.AdText)
		fmt.Printf("   Impressions: %d, Clicks: %d, CTR: %.2f%%\n",
			ad.Impressions, ad.Clicks, ad.CTR*100)
		if ad.VideoID != "" {
			fmt.Printf("   Video ID: %s, Title: %s\n", ad.VideoID, ad.VideoTitle)
		}
		fmt.Println()
	}

	// Example 2: Search with filtering
	fmt.Println("Example 2: Search with filtering by advertiser")
	respFiltered, err := researchAPI.GetAdReport(ctx, &research.GetAdReportRequest{
		SearchTerm:  "fashion",
		CountryCode: &countryCode,
		Filtering: &research.AdReportFiltering{
			CountryCodes: []string{"US", "GB"},
			Platforms:    []string{"TikTok"},
		},
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil {
		log.Fatalf("Failed to get filtered ad report: %v", err)
	}

	fmt.Printf("Total filtered ads found: %d\n", respFiltered.PageInfo.TotalNumber)
	for i, ad := range respFiltered.List {
		fmt.Printf("%d. %s - %s (%s)\n",
			i+1, ad.AdName, ad.AdvertiserName, ad.Country)
	}
	fmt.Println()

	// Example 3: Search with date filtering
	fmt.Println("Example 3: Search ads shown in a specific date range")
	firstShownMin := "2024-01-01"
	firstShownMax := "2024-12-31"

	respByDate, err := researchAPI.GetAdReport(ctx, &research.GetAdReportRequest{
		SearchTerm:  "technology",
		CountryCode: &countryCode,
		Filtering: &research.AdReportFiltering{
			FirstShownDateMin: &firstShownMin,
			FirstShownDateMax: &firstShownMax,
		},
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil {
		log.Fatalf("Failed to get ad report by date: %v", err)
	}

	fmt.Printf("Total ads shown in 2024: %d\n", respByDate.PageInfo.TotalNumber)
	for i, ad := range respByDate.List {
		fmt.Printf("%d. %s - First shown: %s, Last shown: %s\n",
			i+1, ad.AdName, ad.FirstShownDate, ad.LastShownDate)
	}
	fmt.Println()

	// Example 4: Get all ads with auto-pagination
	fmt.Println("Example 4: Get all ads for a search term (auto-pagination)")
	allAds, err := researchAPI.GetAllAdReports(ctx, &research.GetAdReportRequest{
		SearchTerm:  "gaming",
		CountryCode: &countryCode,
	})
	if err != nil {
		log.Fatalf("Failed to get all ad reports: %v", err)
	}

	fmt.Printf("Retrieved %d ads in total\n", len(allAds))
	fmt.Println()

	// Example 5: Search with sorting
	fmt.Println("Example 5: Search with sorting by impressions")
	orderBy := "DESC"
	orderField := "impressions"

	respSorted, err := researchAPI.GetAdReport(ctx, &research.GetAdReportRequest{
		SearchTerm:  "travel",
		CountryCode: &countryCode,
		OrderBy:     &orderBy,
		OrderField:  &orderField,
		Page:        &page,
		PageSize:    &pageSize,
	})
	if err != nil {
		log.Fatalf("Failed to get sorted ad report: %v", err)
	}

	fmt.Printf("Top ads by impressions:\n")
	for i, ad := range respSorted.List {
		fmt.Printf("%d. %s - Impressions: %d\n",
			i+1, ad.AdName, ad.Impressions)
	}
	fmt.Println()

	// Example 6: Get detailed ad information
	fmt.Println("Example 6: Get detailed ad information with specific fields")
	fields := []string{
		"ad_id", "ad_name", "advertiser_name", "ad_text",
		"impressions", "clicks", "ctr", "video_id", "video_title",
		"landing_page_url", "first_shown_date", "last_shown_date",
	}

	respDetailed, err := researchAPI.GetAdReport(ctx, &research.GetAdReportRequest{
		SearchTerm:  "fitness",
		CountryCode: &countryCode,
		Fields:      fields,
		Page:        &page,
		PageSize:    &pageSize,
	})
	if err != nil {
		log.Fatalf("Failed to get detailed ad report: %v", err)
	}

	fmt.Printf("Detailed ad information:\n")
	for i, ad := range respDetailed.List {
		fmt.Printf("\n%d. Ad ID: %s\n", i+1, ad.AdID)
		fmt.Printf("   Name: %s\n", ad.AdName)
		fmt.Printf("   Advertiser: %s\n", ad.AdvertiserName)
		fmt.Printf("   Text: %s\n", ad.AdText)
		fmt.Printf("   Performance: %d impressions, %d clicks, %.2f%% CTR\n",
			ad.Impressions, ad.Clicks, ad.CTR*100)
		if ad.VideoID != "" {
			fmt.Printf("   Video: %s (%s)\n", ad.VideoTitle, ad.VideoID)
		}
		fmt.Printf("   Landing Page: %s\n", ad.LandingPageURL)
		fmt.Printf("   Active Period: %s to %s\n", ad.FirstShownDate, ad.LastShownDate)
	}
}
