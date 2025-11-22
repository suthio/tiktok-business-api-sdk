# TikTok Business API Go SDK Examples

This directory contains examples demonstrating how to use the TikTok Business API Go SDK.

## Prerequisites

Before running the examples, you need to:

1. Install Go 1.25 or later
2. Set up your TikTok Business API credentials
3. Set the required environment variables

## Environment Variables

Most examples require the following environment variables:

```bash
export TIKTOK_ACCESS_TOKEN="your_access_token"
export TIKTOK_ADVERTISER_ID="your_advertiser_id"
```

Some examples may require additional environment variables. Check each example's documentation below.

## Examples

### Authentication

Demonstrates OAuth2 authentication flow.

**Required environment variables:**
- `TIKTOK_APP_ID`
- `TIKTOK_SECRET`
- `TIKTOK_AUTH_CODE`

**Run:**
```bash
cd examples/authentication
go run main.go
```

### Account

Retrieves advertiser account information.

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Run:**
```bash
cd examples/account
go run main.go
```

### Campaign

Manages campaign information and retrieval.

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Optional:**
- `TIKTOK_CAMPAIGN_ID` - for filtering specific campaigns

**Run:**
```bash
cd examples/campaign
go run main.go
```

### Adgroup

Retrieves adgroup information with filtering options.

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Optional:**
- `TIKTOK_CAMPAIGN_ID` - for filtering adgroups by campaign

**Run:**
```bash
cd examples/adgroup
go run main.go
```

### Ad

Demonstrates various ad retrieval operations including filtering by specific ad IDs.

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Optional:**
- `TIKTOK_CAMPAIGN_ID` - for filtering ads by campaign

**Run:**
```bash
cd examples/ad
go run main.go
```

### File (Video & Image)

Shows how to:
- Search for videos
- Get video information
- Download videos
- Get image information

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Optional:**
- `TIKTOK_VIDEO_ID` - for video download example
- `TIKTOK_IMAGE_ID` - for image info example

**Run:**
```bash
cd examples/file
go run main.go
```

**Features:**
- Video search with pagination
- Retrieve video metadata (duration, size, dimensions)
- Download videos using preview URLs
- Get image information

### Reporting

Retrieves reporting data with various metrics and dimensions.

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Optional:**
- `TIKTOK_START_DATE` - start date in YYYY-MM-DD format (default: 2024-01-01)
- `TIKTOK_END_DATE` - end date in YYYY-MM-DD format (default: 2024-01-31)
- `TIKTOK_CAMPAIGN_ID` - for campaign-specific reports

**Run:**
```bash
cd examples/reporting
go run main.go
```

**Metrics included:**
- Spend
- Impressions
- Clicks
- CTR (Click-Through Rate)
- CPC (Cost Per Click)
- CPM (Cost Per Mille)
- Conversion
- Cost Per Conversion

### Audience

Retrieves custom audience information.

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Run:**
```bash
cd examples/audience
go run main.go
```

### Measurement

Gets app list for measurement and tracking.

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Run:**
```bash
cd examples/measurement
go run main.go
```

### Tool

Demonstrates targeting and tool APIs:
- Get targeting categories (interests, behaviors)
- Get available languages

**Required environment variables:**
- `TIKTOK_ACCESS_TOKEN`
- `TIKTOK_ADVERTISER_ID`

**Run:**
```bash
cd examples/tool
go run main.go
```

## Building Examples

To build all examples:

```bash
cd go_sdk
for dir in examples/*/; do
  (cd "$dir" && go build -o example)
done
```

## Common Patterns

### Error Handling

All examples use proper error handling:

```go
resp, err := api.GetSomething(ctx, req)
if err != nil {
    log.Fatalf("Failed to get something: %v", err)
}
```

### Pagination

Many APIs support pagination:

```go
page := int64(1)
pageSize := int64(10)
req := &SomeRequest{
    AdvertiserID: advertiserID,
    Page:         &page,
    PageSize:     &pageSize,
}
```

### Filtering

Most GET APIs support filtering:

```go
req := &GetRequest{
    AdvertiserID: advertiserID,
    Filtering: &Filtering{
        IDs: []string{"id1", "id2"},
    },
}
```

## Notes

- These examples are for demonstration purposes only
- Do not commit your credentials to version control
- Use environment variables or a configuration file to manage credentials
- Rate limits apply to all TikTok Business API endpoints
- Some features may require specific permissions in your TikTok Business account

## Support

For more information, refer to:
- [TikTok Business API Documentation](https://business-api.tiktok.com/portal/docs)
- [Go SDK Repository](https://github.com/suthio/tiktok-business-api-sdk)
