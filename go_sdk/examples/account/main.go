package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/account"
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

	// Create Account API client
	accountAPI := account.NewAPI(client)

	// Get advertiser info
	ctx := context.Background()
	advertiserIDs := []string{advertiserID}
	fields := []string{} // Empty means return all fields

	resp, err := accountAPI.GetAdvertiserInfo(ctx, advertiserIDs, fields)
	if err != nil {
		log.Fatalf("Failed to get advertiser info: %v", err)
	}

	fmt.Printf("Found %d advertiser(s)\n", len(resp.List))
	for _, adv := range resp.List {
		fmt.Printf("\nAdvertiser ID: %s\n", adv.AdvertiserID)
		fmt.Printf("Name: %s\n", adv.AdvertiserName)
		fmt.Printf("Company: %s\n", adv.Company)
		fmt.Printf("Status: %s\n", adv.Status)
		fmt.Printf("Currency: %s\n", adv.Currency)
		fmt.Printf("Balance: %.2f\n", adv.Balance)
		fmt.Printf("Timezone: %s\n", adv.Timezone)
		fmt.Printf("Display Timezone: %s\n", adv.DisplayTimezone)
	}
}
