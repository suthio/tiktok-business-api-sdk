package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
	"github.com/suthio/tiktok-business-api-sdk/go_sdk/file"
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

	// Create File API client
	fileAPI := file.NewAPI(client)

	ctx := context.Background()

	// Example 1: Search for videos
	fmt.Println("=== Example 1: Search Videos ===")
	page := int64(1)
	pageSize := int64(10)
	searchReq := &file.SearchVideosRequest{
		AdvertiserID: advertiserID,
		Page:         &page,
		PageSize:     &pageSize,
	}

	searchResp, err := fileAPI.SearchVideos(ctx, searchReq)
	if err != nil {
		log.Fatalf("Failed to search videos: %v", err)
	}

	fmt.Printf("Found %d video(s)\n", len(searchResp.List))
	fmt.Printf("Total videos: %d\n", searchResp.PageInfo.TotalNumber)
	fmt.Printf("Total pages: %d\n", searchResp.PageInfo.TotalPage)

	// Example 2: Get video info for specific videos
	if len(searchResp.List) > 0 {
		fmt.Println("\n=== Example 2: Get Video Info ===")
		videoIDs := make([]string, 0)
		for i, v := range searchResp.List {
			if i < 5 { // Get info for first 5 videos
				videoIDs = append(videoIDs, v.VideoID)
			}
		}

		if len(videoIDs) > 0 {
			infoReq := &file.GetVideoInfoRequest{
				AdvertiserID: advertiserID,
				VideoIDs:     videoIDs,
			}

			infoResp, err := fileAPI.GetVideoInfo(ctx, infoReq)
			if err != nil {
				log.Fatalf("Failed to get video info: %v", err)
			}

			for _, v := range infoResp.List {
				fmt.Printf("\nVideo ID: %s\n", v.VideoID)
				fmt.Printf("File Name: %s\n", v.FileName)
				fmt.Printf("Format: %s\n", v.Format)
				fmt.Printf("Duration: %.2f seconds\n", v.Duration)
				fmt.Printf("Size: %d bytes\n", v.Size)
				fmt.Printf("Width: %d\n", v.Width)
				fmt.Printf("Height: %d\n", v.Height)
				if v.PreviewURL != "" {
					fmt.Printf("Preview URL: %s\n", v.PreviewURL)
				}
			}
		}
	}

	// Example 3: Download a video (if preview URL is available)
	videoID := os.Getenv("TIKTOK_VIDEO_ID")
	if videoID != "" {
		fmt.Println("\n=== Example 3: Download Video ===")
		infoReq := &file.GetVideoInfoRequest{
			AdvertiserID: advertiserID,
			VideoIDs:     []string{videoID},
		}

		infoResp, err := fileAPI.GetVideoInfo(ctx, infoReq)
		if err != nil {
			log.Fatalf("Failed to get video info: %v", err)
		}

		if len(infoResp.List) > 0 && infoResp.List[0].PreviewURL != "" {
			downloadReq := &file.DownloadVideoRequest{
				URL:        infoResp.List[0].PreviewURL,
				OutputPath: "./downloads",
				FileName:   fmt.Sprintf("%s.mp4", videoID),
			}

			err = fileAPI.DownloadVideo(ctx, downloadReq)
			if err != nil {
				log.Fatalf("Failed to download video: %v", err)
			}

			fmt.Printf("Video downloaded successfully to ./downloads/%s.mp4\n", videoID)
		} else {
			fmt.Println("Preview URL not available for this video")
		}
	}

	// Example 4: Get image info
	imageID := os.Getenv("TIKTOK_IMAGE_ID")
	if imageID != "" {
		fmt.Println("\n=== Example 4: Get Image Info ===")
		imageReq := &file.GetImageInfoRequest{
			AdvertiserID: advertiserID,
			ImageIDs:     []string{imageID},
		}

		imageResp, err := fileAPI.GetImageInfo(ctx, imageReq)
		if err != nil {
			log.Fatalf("Failed to get image info: %v", err)
		}

		for _, img := range imageResp.List {
			fmt.Printf("\nImage ID: %s\n", img.ImageID)
			fmt.Printf("File Name: %s\n", img.FileName)
			fmt.Printf("Format: %s\n", img.Format)
			fmt.Printf("Size: %d bytes\n", img.Size)
			fmt.Printf("Width: %d\n", img.Width)
			fmt.Printf("Height: %d\n", img.Height)
			if img.ImageURL != "" {
				fmt.Printf("Image URL: %s\n", img.ImageURL)
			}
		}
	}
}
