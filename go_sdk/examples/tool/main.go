package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/tool"
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

	// Create Tool API client
	toolAPI := tool.NewAPI(client)

	ctx := context.Background()

	// Example 1: Get available languages
	fmt.Println("=== Example 1: Get Available Languages ===")
	langResp, err := toolAPI.GetLanguage(ctx, advertiserID)
	if err != nil {
		log.Fatalf("Failed to get languages: %v", err)
	}

	fmt.Printf("Found %d languages\n", len(langResp.Languages))
	for i, lang := range langResp.Languages {
		if i >= 10 { // Show first 10 only
			break
		}
		fmt.Printf("  - %s (Code: %s)\n", lang.LanguageName, lang.LanguageCode)
	}

	// Example 2: Get carriers
	fmt.Println("\n=== Example 2: Get Carriers ===")
	carrierResp, err := toolAPI.GetCarrier(ctx, advertiserID)
	if err != nil {
		log.Fatalf("Failed to get carriers: %v", err)
	}

	fmt.Printf("Found %d carriers\n", len(carrierResp.Carriers))
	for i, carrier := range carrierResp.Carriers {
		if i >= 10 { // Show first 10 only
			break
		}
		fmt.Printf("  - %s (ID: %s)\n", carrier.CarrierName, carrier.CarrierID)
	}

	// Example 3: Get action categories
	fmt.Println("\n=== Example 3: Get Action Categories ===")
	actionResp, err := toolAPI.GetActionCategory(ctx, advertiserID, nil)
	if err != nil {
		log.Fatalf("Failed to get action categories: %v", err)
	}

	fmt.Printf("Found %d action categories\n", len(actionResp.ActionCategories))
	for i, category := range actionResp.ActionCategories {
		if i >= 10 { // Show first 10 only
			break
		}
		fmt.Printf("  - %s (ID: %s)\n", category.ActionCategoryName, category.ActionCategoryID)
	}
}
