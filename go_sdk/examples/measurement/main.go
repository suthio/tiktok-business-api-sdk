package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/measurement"
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

	// Create Measurement API client
	measurementAPI := measurement.NewAPI(client)

	ctx := context.Background()

	// List pixels
	fmt.Println("=== Listing Pixels ===")
	page := int64(1)
	pageSize := int64(20)
	req := &measurement.PixelListRequest{
		AdvertiserID: advertiserID,
		Page:         &page,
		PageSize:     &pageSize,
	}

	resp, err := measurementAPI.ListPixels(ctx, req)
	if err != nil {
		log.Fatalf("Failed to list pixels: %v", err)
	}

	fmt.Printf("Found %d pixel(s)\n", len(resp.List))
	for _, pixel := range resp.List {
		fmt.Printf("\nPixel ID: %s\n", pixel.PixelID)
		fmt.Printf("Pixel Name: %s\n", pixel.PixelName)
		fmt.Printf("Pixel Code: %s\n", pixel.PixelCode)
		fmt.Printf("Status: %s\n", pixel.PixelStatus)
		fmt.Printf("Create Time: %s\n", pixel.CreateTime)
	}
}
