# TikTok Business API SDK - Go Implementation

This document describes the Go SDK implementation for TikTok Business API.

## Overview

The Go SDK provides a simple and idiomatic way to interact with TikTok's Business API. It supports various advertising operations including campaign management, ad creation, audience targeting, and performance reporting.

## Architecture

### Package Structure

```
go_sdk/
├── client.go              # Base HTTP client
├── common.go              # Common types and responses
├── account/              # Account & Advertiser management
├── ad/                   # Ad operations
├── adgroup/              # Ad Group operations
├── audience/             # Audience & DMP operations
├── authentication/       # OAuth authentication
├── bc/                   # Business Center operations
├── campaign/             # Campaign operations
├── measurement/          # Pixel & Offline event tracking
├── reporting/            # Reporting & Smart Plus analytics
└── tool/                 # Utility APIs (carriers, languages, etc.)
```

### Core Components

#### Client (`client.go`)

The base HTTP client handles all API requests:

```go
type Client struct {
    baseURL     string
    httpClient  *http.Client
    accessToken string
}
```

- Base URL: `https://business-api.tiktok.com`
- Default timeout: 30 seconds
- Authentication: Via `Access-Token` header

#### Common Types (`common.go`)

Standard response wrapper used across all APIs:

```go
type Response struct {
    Code      *int64
    Message   *string
    RequestID *string
    Data      json.RawMessage
}
```

## API Reference

### 1. Campaign API (`campaign/`)

**Location:** `go_sdk/campaign/campaign.go`

**Methods:**
- `GetCampaigns(ctx, req)` - Retrieve campaign information with filtering

**Reference:** https://business-api.tiktok.com/portal/docs?id=1739315828649986

**Example:**
```go
campaignAPI := campaign.NewAPI(client)
resp, err := campaignAPI.GetCampaigns(ctx, &campaign.GetCampaignRequest{
    AdvertiserID: "123456789",
    Page: ptr(1),
    PageSize: ptr(10),
})
```

---

### 2. Ad API (`ad/`)

**Location:** `go_sdk/ad/ad.go`

**Methods:**
- `GetAds(ctx, req)` - Get regular ads and ACO ads data

**Reference:** https://business-api.tiktok.com/portal/docs?id=1735735588640770

**Example:**
```go
adAPI := ad.NewAPI(client)
resp, err := adAPI.GetAds(ctx, &ad.GetAdRequest{
    AdvertiserID: "123456789",
    Filtering: &ad.Filtering{
        AdgroupIDs: []string{"adgroup_id_1", "adgroup_id_2"},
    },
    Page: ptr(1),
    PageSize: ptr(20),
})
```

---

### 3. AdGroup API (`adgroup/`)

**Location:** `go_sdk/adgroup/adgroup.go`

**Methods:**
- `GetAdGroups(ctx, req)` - Obtain detailed information of ad groups

**Reference:** https://business-api.tiktok.com/portal/docs?id=1739314558673922

**Example:**
```go
adgroupAPI := adgroup.NewAPI(client)
resp, err := adgroupAPI.GetAdGroups(ctx, &adgroup.GetAdGroupRequest{
    AdvertiserID: "123456789",
    Filtering: &adgroup.Filtering{
        CampaignIDs: []string{"campaign_id_1"},
    },
})
```

---

### 4. Measurement API (`measurement/`)

**Location:** `go_sdk/measurement/measurement.go`

**Methods:**
- `ListPixels(ctx, req)` - Obtain a list of Pixel information
- `GetOfflineEventSets(ctx, req)` - Get Offline Event sets

**References:**
- Pixel List: https://business-api.tiktok.com/portal/docs?id=1740858697598978
- Offline Get: https://business-api.tiktok.com/portal/docs?id=1765596808589313

**Example:**
```go
measurementAPI := measurement.NewAPI(client)

// List pixels
pixelResp, err := measurementAPI.ListPixels(ctx, &measurement.PixelListRequest{
    AdvertiserID: "123456789",
    Page: ptr(1),
    PageSize: ptr(10),
})

// Get offline event sets
offlineResp, err := measurementAPI.GetOfflineEventSets(ctx, &measurement.OfflineGetRequest{
    AdvertiserID: "123456789",
})
```

---

### 5. Audience API (`audience/`)

**Location:** `go_sdk/audience/audience.go`

**Methods:**
- `GetCustomAudiences(ctx, req)` - Obtain details of specified audiences
- `ListCustomAudiences(ctx, req)` - Get all audiences

**References:**
- Get: https://business-api.tiktok.com/portal/docs?id=1739940507792385
- List: https://business-api.tiktok.com/portal/docs?id=1739940506015746

**Example:**
```go
audienceAPI := audience.NewAPI(client)

// Get specific audiences
resp, err := audienceAPI.GetCustomAudiences(ctx, &audience.CustomAudienceGetRequest{
    AdvertiserID: "123456789",
    CustomAudienceIDs: []string{"audience_id_1", "audience_id_2"},
})

// List all audiences
listResp, err := audienceAPI.ListCustomAudiences(ctx, &audience.CustomAudienceListRequest{
    AdvertiserID: "123456789",
    Page: ptr(1),
    PageSize: ptr(10),
})
```

---

### 6. Tool API (`tool/`)

**Location:** `go_sdk/tool/tool.go`

**Methods:**
- `GetCarrier(ctx, advertiserID)` - Get carriers for targeting
- `GetLanguage(ctx, advertiserID)` - Get supported languages
- `GetActionCategory(ctx, advertiserID, specialIndustries)` - Get action categories

**References:**
- Carrier: https://business-api.tiktok.com/portal/docs?id=1737168013095938
- Language: https://business-api.tiktok.com/portal/docs?id=1737188554152962
- Action Category: https://business-api.tiktok.com/portal/docs?id=1737166752522241

**Example:**
```go
toolAPI := tool.NewAPI(client)

// Get carriers
carriers, err := toolAPI.GetCarrier(ctx, "123456789")

// Get languages
languages, err := toolAPI.GetLanguage(ctx, "123456789")

// Get action categories
categories, err := toolAPI.GetActionCategory(ctx, "123456789", []string{})
```

---

### 7. Account API (`account/`)

**Location:** `go_sdk/account/account.go`

**Methods:**
- `GetAdvertiserInfo(ctx, advertiserIDs, fields)` - Get advertiser information including balance

**Reference:** https://business-api.tiktok.com/portal/docs?id=1739593083610113

**Example:**
```go
accountAPI := account.NewAPI(client)
resp, err := accountAPI.GetAdvertiserInfo(ctx,
    []string{"123456789", "987654321"},
    []string{"advertiser_id", "advertiser_name", "balance"},
)
```

---

### 8. BC (Business Center) API (`bc/`)

**Location:** `go_sdk/bc/bc.go`

**Methods:**
- `GetAccountTransactions(ctx, req)` - Get transaction records of a BC or ad accounts
- `GetAssets(ctx, req)` - Get assets in a Business Center

**References:**
- Transaction: https://business-api.tiktok.com/portal/docs?id=1792849810925569
- Assets: https://business-api.tiktok.com/portal/docs?id=1739593603696641

**Example:**
```go
bcAPI := bc.NewAPI(client)

// Get transactions
transResp, err := bcAPI.GetAccountTransactions(ctx, &bc.AccountTransactionRequest{
    BcID: ptr("bc_id_123"),
    Page: ptr(1),
    PageSize: ptr(20),
})

// Get assets
assetResp, err := bcAPI.GetAssets(ctx, &bc.AssetGetRequest{
    BcID: "bc_id_123",
    AssetType: "ADVERTISER",
})
```

---

### 9. Reporting API (`reporting/`)

**Location:** `go_sdk/reporting/reporting.go`

**Methods:**
- `GetIntegratedReport(ctx, req)` - Run synchronous reports
- `CheckReportTask(ctx, taskID, advertiserID)` - Check async report task status
- `GetMaterialReportBreakdown(ctx, req)` - Get Smart Plus material report breakdown
- `GetMaterialReportOverview(ctx, req)` - Get Smart Plus material report overview

**References:**
- Integrated: https://business-api.tiktok.com/portal/docs?id=1740302848100353
- Task Check: https://business-api.tiktok.com/portal/docs?id=1740302781443073
- Material Breakdown: https://business-api.tiktok.com/portal/docs?id=1765936670832641
- Material Overview: https://business-api.tiktok.com/portal/docs?id=1765936643763201

**Example:**
```go
reportingAPI := reporting.NewAPI(client)

// Get integrated report
report, err := reportingAPI.GetIntegratedReport(ctx, &reporting.IntegratedGetRequest{
    ReportType: "BASIC",
    AdvertiserID: ptr("123456789"),
    DataLevel: ptr("AUCTION_CAMPAIGN"),
    Dimensions: []string{"campaign_id", "stat_time_day"},
    Metrics: []string{"spend", "impressions", "clicks"},
    StartDate: ptr("2024-01-01"),
    EndDate: ptr("2024-01-31"),
})

// Get Smart Plus material report
materialReport, err := reportingAPI.GetMaterialReportBreakdown(ctx, &reporting.MaterialReportBreakdownRequest{
    AdvertiserID: "123456789",
    Dimensions: []string{"material_id", "stat_time_day"},
    StartDate: "2024-01-01",
    EndDate: "2024-01-31",
})
```

---

### 10. Authentication API (`authentication/`)

**Location:** `go_sdk/authentication/authentication.go`

**Methods:**
- `GetAccessToken(ctx, req)` - Get OAuth access token
- `GetAdvertisers(ctx, appID, secret, accessToken)` - Get authorized advertiser accounts

**Reference:** OAuth documentation

**Example:**
```go
authAPI := authentication.NewAPI(client)

// Get access token
tokenResp, err := authAPI.GetAccessToken(ctx, &authentication.AccessTokenRequest{
    AppID: "your_app_id",
    Secret: "your_secret",
    AuthCode: "auth_code_from_oauth",
})

// Get advertiser list
advertisers, err := authAPI.GetAdvertisers(ctx, "app_id", "secret", "access_token")
```

---

## Usage Patterns

### Initialization

```go
package main

import (
    "context"
    tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
    "github.com/suthio/tiktok-business-api-sdk/go_sdk/campaign"
)

func main() {
    // Create client
    client := tiktok.NewClient("your_access_token")

    // Create API instance
    campaignAPI := campaign.NewAPI(client)

    // Make API call
    ctx := context.Background()
    resp, err := campaignAPI.GetCampaigns(ctx, &campaign.GetCampaignRequest{
        AdvertiserID: "123456789",
    })
    if err != nil {
        // Handle error
    }

    // Use response
    for _, c := range resp.List {
        println(c.CampaignName)
    }
}
```

### Error Handling

All API methods return an error as the second return value:

```go
resp, err := api.GetCampaigns(ctx, req)
if err != nil {
    // API error or network error
    log.Printf("Error: %v", err)
    return
}

// Use resp
```

### Pagination

Most list endpoints support pagination:

```go
page := int64(1)
pageSize := int64(100)

resp, err := api.GetAds(ctx, &ad.GetAdRequest{
    AdvertiserID: "123456789",
    Page: &page,
    PageSize: &pageSize,
})

// Access pagination info
fmt.Printf("Total: %d, Pages: %d\n", resp.PageInfo.TotalNumber, resp.PageInfo.TotalPage)
```

### Filtering

Many endpoints support filtering:

```go
resp, err := api.GetCampaigns(ctx, &campaign.GetCampaignRequest{
    AdvertiserID: "123456789",
    Filtering: campaign.Filtering{
        CampaignIDs: []string{"campaign_1", "campaign_2"},
        PrimaryStatus: ptr("STATUS_ENABLE"),
    },
})
```

## Helper Function

For optional pointer parameters:

```go
func ptr[T any](v T) *T {
    return &v
}

// Usage
Page: ptr(1),
PageSize: ptr(10),
```

## Design Patterns

1. **Factory Pattern**: Each API package has `NewAPI(client)` constructor
2. **Context Support**: All methods accept `context.Context` for cancellation
3. **Error Wrapping**: Errors are wrapped with descriptive messages using `fmt.Errorf`
4. **JSON Marshaling**: Flexible data handling with `json.RawMessage`
5. **Pagination**: Consistent `PageInfo` struct across all list endpoints

## Smart Plus Features

Smart Plus campaigns are fully supported through the Reporting API:

- Material performance breakdown by creative assets
- Overview reports with lifetime or date-range queries
- Filtering and sorting capabilities
- All standard metrics available

## Requirements

- Go 1.21 or higher
- Valid TikTok Business API access token
- Active TikTok Ads account

## Module

```
module github.com/suthio/tiktok-business-api-sdk/go_sdk

go 1.21
```

## Implementation Status

### ✅ Completed (Read/Get Operations)

- Campaign API - Get campaigns with Smart Plus support
- Ad API - Get ads data
- AdGroup API - Get ad groups
- Measurement API - List pixels, get offline event sets
- Audience API - Get and list custom audiences
- Tool API - Get carriers, languages, action categories
- Account API - Get advertiser info and balance
- BC API - Get transactions and assets
- Reporting API - Integrated reports and Smart Plus analytics
- Authentication API - OAuth flow

### 📝 Future Enhancements

The SDK currently focuses on Read operations. Write operations (Create, Update, Delete) can be added following the same patterns.

## Contributing

When adding new APIs:

1. Create a new package under `go_sdk/`
2. Follow the existing patterns:
   - `NewAPI(client)` constructor
   - Context-aware methods
   - Proper error wrapping
   - Clear struct definitions
3. Add reference links to TikTok Business API docs
4. Update this document

## Support

- TikTok Business API Documentation: https://business-api.tiktok.com/portal/docs
- GitHub Issues: https://github.com/suthio/tiktok-business-api-sdk/issues

---

Generated with Claude Code
