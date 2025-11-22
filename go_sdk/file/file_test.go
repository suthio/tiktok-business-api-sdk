package file

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

func TestNewAPI(t *testing.T) {
	client := &tiktok.Client{}
	api := NewAPI(client)
	assert.NotNil(t, api)
	assert.Equal(t, client, api.client)
}

func TestGetVideoInfo_EmptyVideoIDs(t *testing.T) {
	client := &tiktok.Client{}
	api := NewAPI(client)

	req := &GetVideoInfoRequest{
		AdvertiserID: "123456",
		VideoIDs:     []string{},
	}

	_, err := api.GetVideoInfo(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "video_ids cannot be empty")
}

func TestGetVideoInfo_TooManyVideoIDs(t *testing.T) {
	client := &tiktok.Client{}
	api := NewAPI(client)

	videoIDs := make([]string, 61)
	for i := 0; i < 61; i++ {
		videoIDs[i] = "video_id"
	}

	req := &GetVideoInfoRequest{
		AdvertiserID: "123456",
		VideoIDs:     videoIDs,
	}

	_, err := api.GetVideoInfo(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "video_ids cannot exceed 60 items")
}

func TestGetImageInfo_EmptyImageIDs(t *testing.T) {
	client := &tiktok.Client{}
	api := NewAPI(client)

	req := &GetImageInfoRequest{
		AdvertiserID: "123456",
		ImageIDs:     []string{},
	}

	_, err := api.GetImageInfo(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "image_ids cannot be empty")
}

func TestGetImageInfo_TooManyImageIDs(t *testing.T) {
	client := &tiktok.Client{}
	api := NewAPI(client)

	imageIDs := make([]string, 101)
	for i := 0; i < 101; i++ {
		imageIDs[i] = "image_id"
	}

	req := &GetImageInfoRequest{
		AdvertiserID: "123456",
		ImageIDs:     imageIDs,
	}

	_, err := api.GetImageInfo(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "image_ids cannot exceed 100 items")
}

func TestDownloadVideo_EmptyURL(t *testing.T) {
	client := &tiktok.Client{}
	api := NewAPI(client)

	req := &DownloadVideoRequest{
		URL: "",
	}

	err := api.DownloadVideo(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "URL cannot be empty")
}
