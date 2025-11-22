package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/ad"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/adgroup"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/campaign"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/file"
)

func ptr[T any](v T) *T {
	return &v
}

func main() {
	// Get access token and advertiser ID from environment variables
	accessToken := os.Getenv("TIKTOK_ACCESS_TOKEN")
	advertiserID := os.Getenv("TIKTOK_ADVERTISER_ID")
	if accessToken == "" || advertiserID == "" {
		log.Fatal("TIKTOK_ACCESS_TOKEN and TIKTOK_ADVERTISER_ID environment variables must be set")
	}

	// Create TikTok client
	client := tiktok.NewClient(accessToken)
	ctx := context.Background()

	// Step 1: Create a Campaign
	fmt.Println("=== Step 1: Creating Campaign ===")
	campaignAPI := campaign.NewAPI(client)
	campaignName := fmt.Sprintf("Full Flow Test %s", time.Now().Format("20060102150405"))
	budgetMode := "BUDGET_MODE_INFINITE"

	campaignReq := &campaign.CreateCampaignRequest{
		AdvertiserID:  advertiserID,
		CampaignName:  campaignName,
		ObjectiveType: "TRAFFIC",
		BudgetMode:    &budgetMode,
	}

	campaignResp, err := campaignAPI.CreateCampaign(ctx, campaignReq)
	if err != nil {
		log.Fatalf("Failed to create campaign: %v", err)
	}

	fmt.Printf("✓ Campaign created: %s (ID: %s)\n\n", campaignName, campaignResp.CampaignID)

	// Step 2: Get video for creative
	fmt.Println("=== Step 2: Finding Video for Creative ===")
	fileAPI := file.NewAPI(client)

	searchReq := &file.SearchVideosRequest{
		AdvertiserID: advertiserID,
		Page:         ptr(int64(1)),
		PageSize:     ptr(int64(10)),
	}

	videoResp, err := fileAPI.SearchVideos(ctx, searchReq)
	if err != nil {
		log.Fatalf("Failed to search videos: %v", err)
	}

	if len(videoResp.List) == 0 {
		log.Fatal("No videos found in account. Please upload a video first.")
	}

	videoID := videoResp.List[0].VideoID
	fmt.Printf("✓ Found video: %s\n\n", videoID)

	// Step 3: Create AdGroup
	fmt.Println("=== Step 3: Creating AdGroup ===")
	adgroupAPI := adgroup.NewAPI(client)
	adgroupName := fmt.Sprintf("AdGroup Test %s", time.Now().Format("150405"))
	promotionType := "WEBSITE"
	budget := 2000.0 // Minimum daily budget for JPY
	bidPrice := 10.0 // Minimum bid price for JPY
	scheduleType := "SCHEDULE_FROM_NOW"
	scheduleStartTime := time.Now().Format("2006-01-02 15:04:05")
	pacing := "PACING_MODE_SMOOTH"

	adgroupReq := &adgroup.CreateAdGroupRequest{
		AdvertiserID:      advertiserID,
		CampaignID:        campaignResp.CampaignID,
		AdGroupName:       adgroupName,
		PromotionType:     &promotionType,
		PlacementType:     "PLACEMENT_TYPE_AUTOMATIC",
		Placements:        []string{"PLACEMENT_TIKTOK"},
		LocationIDs:       []string{"6252001"}, // Japan
		BudgetMode:        "BUDGET_MODE_DAY",
		Budget:            &budget,
		ScheduleType:      &scheduleType,
		ScheduleStartTime: &scheduleStartTime,
		BillingEvent:      "CPC",
		OptimizationGoal:  "CLICK",
		BidPrice:          &bidPrice,
		Pacing:            &pacing,
	}

	adgroupResp, err := adgroupAPI.CreateAdGroup(ctx, adgroupReq)
	if err != nil {
		log.Fatalf("Failed to create adgroup: %v", err)
	}

	fmt.Printf("✓ AdGroup created: %s (ID: %s)\n\n", adgroupName, adgroupResp.AdGroupID)

	// Step 4: Create Ad
	fmt.Println("=== Step 4: Creating Ad ===")
	adAPI := ad.NewAPI(client)
	adName := fmt.Sprintf("Ad Test %s", time.Now().Format("150405"))
	displayName := "Test Advertiser"
	callToAction := "LEARN_MORE"
	landingPageURL := "https://www.example.com"
	identityType := "CUSTOMIZED_USER"
	identityID := advertiserID

	creative := ad.AdCreative{
		AdName:         adName,
		AdText:         "Check out this amazing product!",
		AdFormat:       "SINGLE_VIDEO",
		VideoID:        &videoID,
		DisplayName:    &displayName,
		CallToAction:   &callToAction,
		LandingPageURL: &landingPageURL,
		IdentityType:   &identityType,
		IdentityID:     &identityID,
	}

	adReq := &ad.CreateAdRequest{
		AdvertiserID: advertiserID,
		AdGroupID:    adgroupResp.AdGroupID,
		Creatives:    []ad.AdCreative{creative},
	}

	adCreateResp, err := adAPI.CreateAd(ctx, adReq)
	if err != nil {
		log.Fatalf("Failed to create ad: %v", err)
	}

	fmt.Printf("✓ Ad created: %s (ID: %s)\n\n", adName, adCreateResp.AdID)

	// Step 5: Verify everything was created
	fmt.Println("=== Step 5: Verification ===")

	// Verify Campaign
	getCampaignReq := &campaign.GetCampaignRequest{
		AdvertiserID: advertiserID,
		Filtering: &campaign.Filtering{
			CampaignIDs: []string{campaignResp.CampaignID},
		},
	}

	campaigns, err := campaignAPI.GetCampaigns(ctx, getCampaignReq)
	if err != nil {
		log.Printf("Warning: Failed to verify campaign: %v", err)
	} else if len(campaigns.List) > 0 {
		fmt.Printf("✓ Campaign verified: %s\n", campaigns.List[0].CampaignName)
	}

	// Verify AdGroup
	getAdgroupReq := &adgroup.GetAdGroupRequest{
		AdvertiserID: advertiserID,
		Filtering: &adgroup.Filtering{
			AdgroupIDs: []string{adgroupResp.AdGroupID},
		},
	}

	adgroups, err := adgroupAPI.GetAdGroups(ctx, getAdgroupReq)
	if err != nil {
		log.Printf("Warning: Failed to verify adgroup: %v", err)
	} else if len(adgroups.List) > 0 {
		fmt.Printf("✓ AdGroup verified: %s\n", adgroups.List[0].AdgroupName)
	}

	// Verify Ad
	getAdReq := &ad.GetAdRequest{
		AdvertiserID: advertiserID,
		Filtering: &ad.Filtering{
			AdIDs: []string{adCreateResp.AdID},
		},
	}

	ads, err := adAPI.GetAds(ctx, getAdReq)
	if err != nil {
		log.Printf("Warning: Failed to verify ad: %v", err)
	} else if len(ads.List) > 0 {
		fmt.Printf("✓ Ad verified: %s\n", ads.List[0].AdName)
	}

	// Summary
	fmt.Println("\n=== Summary ===")
	fmt.Printf("Campaign ID:  %s\n", campaignResp.CampaignID)
	fmt.Printf("AdGroup ID:   %s\n", adgroupResp.AdGroupID)
	fmt.Printf("Ad ID:        %s\n", adCreateResp.AdID)
	fmt.Println("\n✓ Full flow test completed successfully!")
}
