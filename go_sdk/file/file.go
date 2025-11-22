package file

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the File API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new File API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// VideoInfo represents video information
type VideoInfo struct {
	VideoID              string      `json:"video_id"`
	FileName             string      `json:"file_name"`
	Format               string      `json:"format"`
	Width                int64       `json:"width"`
	Height               int64       `json:"height"`
	Duration             float64     `json:"duration"`
	Size                 int64       `json:"size"`
	MaterialID           string      `json:"material_id,omitempty"`
	PosterURL            string      `json:"poster_url,omitempty"`
	PreviewURL           string      `json:"preview_url,omitempty"`
	PreviewURLExpireTime interface{} `json:"preview_url_expire_time,omitempty"` // Can be int64 or string
	BitRate              int64       `json:"bit_rate,omitempty"`
	AllowDownload        bool        `json:"allow_download,omitempty"`
	AllowedPlacements    []string    `json:"allowed_placements,omitempty"`
	CreateTime           string      `json:"create_time"`
	ModifyTime           string      `json:"modify_time"`
}

// GetVideoInfoResponse represents the response for getting video info
type GetVideoInfoResponse struct {
	List []VideoInfo `json:"list"`
}

// GetVideoInfoRequest represents the request to get video info
type GetVideoInfoRequest struct {
	AdvertiserID string   `json:"advertiser_id"`
	VideoIDs     []string `json:"video_ids"`
}

// GetVideoInfo gets video information
// Reference: https://business-api.tiktok.com/portal/docs?id=1740050161973250
func (a *API) GetVideoInfo(ctx context.Context, req *GetVideoInfoRequest) (*GetVideoInfoResponse, error) {
	if len(req.VideoIDs) == 0 {
		return nil, fmt.Errorf("video_ids cannot be empty")
	}

	if len(req.VideoIDs) > 60 {
		return nil, fmt.Errorf("video_ids cannot exceed 60 items")
	}

	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	// Add video IDs using helper
	if err := tiktok.AddStringSlice(params, "video_ids", req.VideoIDs); err != nil {
		return nil, err
	}

	// Use generic DoGet helper
	var resp GetVideoInfoResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/file/video/ad/info/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	return &resp, nil
}

// VideoSearchFiltering represents filtering options for video search
type VideoSearchFiltering struct {
	VideoIDs      []string `json:"video_ids,omitempty"`
	Width         *int64   `json:"width,omitempty"`
	Height        *int64   `json:"height,omitempty"`
	Ratio         []string `json:"ratio,omitempty"`
	VideoTags     []string `json:"video_tags,omitempty"`
	CreateTimeMin *string  `json:"create_time_min,omitempty"`
	CreateTimeMax *string  `json:"create_time_max,omitempty"`
}

// SearchVideosRequest represents the request to search videos
type SearchVideosRequest struct {
	AdvertiserID string                `json:"advertiser_id"`
	Filtering    *VideoSearchFiltering `json:"filtering,omitempty"`
	Page         *int64                `json:"page,omitempty"`
	PageSize     *int64                `json:"page_size,omitempty"`
}

// SearchVideosResponse represents the response for searching videos
type SearchVideosResponse struct {
	List     []VideoInfo     `json:"list"`
	PageInfo tiktok.PageInfo `json:"page_info"`
}

// SearchVideos searches for video creatives in the Asset Library
// Reference: https://business-api.tiktok.com/portal/docs?id=1740050472410114
func (a *API) SearchVideos(ctx context.Context, req *SearchVideosRequest) (*SearchVideosResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	// Add pagination using helper
	tiktok.AddPagination(params, &tiktok.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	// Add filtering using helper
	if err := tiktok.AddJSONParam(params, "filtering", req.Filtering); err != nil {
		return nil, err
	}

	// Use generic DoGet helper
	var resp SearchVideosResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/file/video/ad/search/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to search videos: %w", err)
	}

	return &resp, nil
}

// DownloadVideoRequest represents the request to download a video
type DownloadVideoRequest struct {
	URL        string
	OutputPath string
	FileName   string
}

// DownloadVideo downloads a video from the given URL to the specified path
func (a *API) DownloadVideo(ctx context.Context, req *DownloadVideoRequest) error {
	if req.URL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if req.OutputPath == "" {
		req.OutputPath = "."
	}

	if req.FileName == "" {
		req.FileName = "video.mp4"
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", req.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download video: status code %d", resp.StatusCode)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(req.OutputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	outputFile := filepath.Join(req.OutputPath, req.FileName)
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() { _ = out.Close() }()

	// Copy content to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write video to file: %w", err)
	}

	return nil
}

// ImageInfo represents image information
type ImageInfo struct {
	ImageID           string   `json:"image_id"`
	FileName          string   `json:"file_name"`
	Format            string   `json:"format"`
	Width             int64    `json:"width"`
	Height            int64    `json:"height"`
	Size              int64    `json:"size"`
	MaterialID        string   `json:"material_id,omitempty"`
	ImageURL          string   `json:"image_url,omitempty"`
	Signature         string   `json:"signature,omitempty"`
	AllowedPlacements []string `json:"allowed_placements,omitempty"`
	CreateTime        string   `json:"create_time"`
	ModifyTime        string   `json:"modify_time"`
}

// GetImageInfoResponse represents the response for getting image info
type GetImageInfoResponse struct {
	List []ImageInfo `json:"list"`
}

// GetImageInfoRequest represents the request to get image info
type GetImageInfoRequest struct {
	AdvertiserID string   `json:"advertiser_id"`
	ImageIDs     []string `json:"image_ids"`
}

// GetImageInfo gets image information
// Reference: https://business-api.tiktok.com/portal/docs?id=1740051721711618
func (a *API) GetImageInfo(ctx context.Context, req *GetImageInfoRequest) (*GetImageInfoResponse, error) {
	if len(req.ImageIDs) == 0 {
		return nil, fmt.Errorf("image_ids cannot be empty")
	}

	if len(req.ImageIDs) > 100 {
		return nil, fmt.Errorf("image_ids cannot exceed 100 items")
	}

	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	// Add image IDs using helper
	if err := tiktok.AddStringSlice(params, "image_ids", req.ImageIDs); err != nil {
		return nil, err
	}

	// Use generic DoGet helper
	var resp GetImageInfoResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/file/image/ad/info/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get image info: %w", err)
	}

	return &resp, nil
}
