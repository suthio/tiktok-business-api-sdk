package reporting

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	tiktok "github.com/suthio/tiktok-business-api-sdk/go_sdk"
)

// API represents the Reporting API client
type API struct {
	client *tiktok.Client
}

// NewAPI creates a new Reporting API client
func NewAPI(client *tiktok.Client) *API {
	return &API{
		client: client,
	}
}

// IntegratedGetRequest represents the request for getting integrated reports
type IntegratedGetRequest struct {
	ReportType              string      `json:"report_type"`
	AdvertiserID            *string     `json:"advertiser_id,omitempty"`
	AdvertiserIDs           []string    `json:"advertiser_ids,omitempty"`
	BcID                    *string     `json:"bc_id,omitempty"`
	ServiceType             *string     `json:"service_type,omitempty"`
	DataLevel               *string     `json:"data_level,omitempty"`
	Dimensions              []string    `json:"dimensions,omitempty"`
	Metrics                 []string    `json:"metrics,omitempty"`
	StartDate               *string     `json:"start_date,omitempty"`
	EndDate                 *string     `json:"end_date,omitempty"`
	QueryLifetime           *bool       `json:"query_lifetime,omitempty"`
	Page                    *int64      `json:"page,omitempty"`
	PageSize                *int64      `json:"page_size,omitempty"`
	OrderField              *string     `json:"order_field,omitempty"`
	OrderType               *string     `json:"order_type,omitempty"`
	EnableTotalMetrics      *bool       `json:"enable_total_metrics,omitempty"`
	MultiAdvReportInUTCTime *bool       `json:"multi_adv_report_in_utc_time,omitempty"`
	QueryMode               *string     `json:"query_mode,omitempty"`
	Filtering               interface{} `json:"filtering,omitempty"`
}

// IntegratedGetResponse represents the response from integrated report
type IntegratedGetResponse struct {
	List         []map[string]interface{} `json:"list"`
	PageInfo     tiktok.PageInfo          `json:"page_info"`
	TotalMetrics map[string]interface{}   `json:"total_metrics,omitempty"`
}

// GetIntegratedReport runs a synchronous report.
// Reference: https://business-api.tiktok.com/portal/docs?id=1740302848100353
func (a *API) GetIntegratedReport(ctx context.Context, req *IntegratedGetRequest) (*IntegratedGetResponse, error) {
	params := url.Values{}
	params.Set("report_type", req.ReportType)

	if req.AdvertiserID != nil {
		params.Set("advertiser_id", *req.AdvertiserID)
	}

	if len(req.AdvertiserIDs) > 0 {
		for _, id := range req.AdvertiserIDs {
			params.Add("advertiser_ids", id)
		}
	}

	if req.BcID != nil {
		params.Set("bc_id", *req.BcID)
	}

	if req.ServiceType != nil {
		params.Set("service_type", *req.ServiceType)
	}

	if req.DataLevel != nil {
		params.Set("data_level", *req.DataLevel)
	}

	// Add dimensions as JSON array
	if err := tiktok.AddStringSlice(params, "dimensions", req.Dimensions); err != nil {
		return nil, err
	}

	// Add metrics as JSON array
	if err := tiktok.AddStringSlice(params, "metrics", req.Metrics); err != nil {
		return nil, err
	}

	if req.StartDate != nil {
		params.Set("start_date", *req.StartDate)
	}

	if req.EndDate != nil {
		params.Set("end_date", *req.EndDate)
	}

	if req.QueryLifetime != nil {
		params.Set("query_lifetime", strconv.FormatBool(*req.QueryLifetime))
	}

	// Add pagination using helper
	tiktok.AddPagination(params, &tiktok.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if req.OrderField != nil {
		params.Set("order_field", *req.OrderField)
	}

	if req.OrderType != nil {
		params.Set("order_type", *req.OrderType)
	}

	if req.EnableTotalMetrics != nil {
		params.Set("enable_total_metrics", strconv.FormatBool(*req.EnableTotalMetrics))
	}

	if req.MultiAdvReportInUTCTime != nil {
		params.Set("multi_adv_report_in_utc_time", strconv.FormatBool(*req.MultiAdvReportInUTCTime))
	}

	if req.QueryMode != nil {
		params.Set("query_mode", *req.QueryMode)
	}

	// Add filtering using helper
	if err := tiktok.AddJSONParam(params, "filtering", req.Filtering); err != nil {
		return nil, err
	}

	// Use generic DoGet helper
	var resp IntegratedGetResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/report/integrated/get/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get integrated report: %w", err)
	}

	return &resp, nil
}

// TaskCheckResponse represents the response for task check
type TaskCheckResponse struct {
	TaskID      string `json:"task_id"`
	Status      string `json:"status"`
	DownloadURL string `json:"download_url,omitempty"`
	TotalCount  int64  `json:"total_count,omitempty"`
}

// CheckReportTask gets the status of an async report task.
// Reference: https://business-api.tiktok.com/portal/docs?id=1740302781443073
func (a *API) CheckReportTask(ctx context.Context, taskID, advertiserID string) (*TaskCheckResponse, error) {
	params := url.Values{}
	params.Set("task_id", taskID)
	params.Set("advertiser_id", advertiserID)

	// Use generic DoGet helper
	var resp TaskCheckResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/report/task/check/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to check report task: %w", err)
	}

	return &resp, nil
}

// MaterialReportBreakdownRequest represents the request for Smart Plus material report breakdown
type MaterialReportBreakdownRequest struct {
	AdvertiserID string      `json:"advertiser_id"`
	Dimensions   []string    `json:"dimensions"`
	StartDate    string      `json:"start_date"`
	EndDate      string      `json:"end_date"`
	Metrics      []string    `json:"metrics,omitempty"`
	Filtering    interface{} `json:"filtering,omitempty"`
	SortField    *string     `json:"sort_field,omitempty"`
	SortType     *string     `json:"sort_type,omitempty"`
	Page         *int64      `json:"page,omitempty"`
	PageSize     *int64      `json:"page_size,omitempty"`
}

// MaterialReportResponse represents the response for material reports
type MaterialReportResponse struct {
	List     []map[string]interface{} `json:"list"`
	PageInfo tiktok.PageInfo          `json:"page_info"`
}

// GetMaterialReportBreakdown gets breakdown of Smart Plus material reports.
// Reference: https://business-api.tiktok.com/portal/docs?id=1765936670832641
func (a *API) GetMaterialReportBreakdown(ctx context.Context, req *MaterialReportBreakdownRequest) (*MaterialReportResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)
	params.Set("start_date", req.StartDate)
	params.Set("end_date", req.EndDate)

	// Add dimensions as JSON array
	if err := tiktok.AddStringSlice(params, "dimensions", req.Dimensions); err != nil {
		return nil, err
	}

	// Add metrics as JSON array
	if err := tiktok.AddStringSlice(params, "metrics", req.Metrics); err != nil {
		return nil, err
	}

	// Add filtering using helper
	if err := tiktok.AddJSONParam(params, "filtering", req.Filtering); err != nil {
		return nil, err
	}

	if req.SortField != nil {
		params.Set("sort_field", *req.SortField)
	}

	if req.SortType != nil {
		params.Set("sort_type", *req.SortType)
	}

	// Add pagination using helper
	tiktok.AddPagination(params, &tiktok.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	// Use generic DoGet helper
	var resp MaterialReportResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/smart_plus/material_report/breakdown/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get material report breakdown: %w", err)
	}

	return &resp, nil
}

// MaterialReportOverviewRequest represents the request for Smart Plus material report overview
type MaterialReportOverviewRequest struct {
	AdvertiserID  string      `json:"advertiser_id"`
	Dimensions    []string    `json:"dimensions"`
	Metrics       []string    `json:"metrics,omitempty"`
	StartDate     *string     `json:"start_date,omitempty"`
	EndDate       *string     `json:"end_date,omitempty"`
	QueryLifetime *bool       `json:"query_lifetime,omitempty"`
	Filtering     interface{} `json:"filtering,omitempty"`
	SortField     *string     `json:"sort_field,omitempty"`
	SortType      *string     `json:"sort_type,omitempty"`
	Page          *int64      `json:"page,omitempty"`
	PageSize      *int64      `json:"page_size,omitempty"`
}

// GetMaterialReportOverview gets overview of Smart Plus material reports.
// Reference: https://business-api.tiktok.com/portal/docs?id=1765936643763201
func (a *API) GetMaterialReportOverview(ctx context.Context, req *MaterialReportOverviewRequest) (*MaterialReportResponse, error) {
	params := url.Values{}
	params.Set("advertiser_id", req.AdvertiserID)

	// Add dimensions as JSON array
	if err := tiktok.AddStringSlice(params, "dimensions", req.Dimensions); err != nil {
		return nil, err
	}

	// Add metrics as JSON array
	if err := tiktok.AddStringSlice(params, "metrics", req.Metrics); err != nil {
		return nil, err
	}

	if req.StartDate != nil {
		params.Set("start_date", *req.StartDate)
	}

	if req.EndDate != nil {
		params.Set("end_date", *req.EndDate)
	}

	if req.QueryLifetime != nil {
		params.Set("query_lifetime", strconv.FormatBool(*req.QueryLifetime))
	}

	// Add filtering using helper
	if err := tiktok.AddJSONParam(params, "filtering", req.Filtering); err != nil {
		return nil, err
	}

	if req.SortField != nil {
		params.Set("sort_field", *req.SortField)
	}

	if req.SortType != nil {
		params.Set("sort_type", *req.SortType)
	}

	// Add pagination using helper
	tiktok.AddPagination(params, &tiktok.PaginationParams{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	// Use generic DoGet helper
	var resp MaterialReportResponse
	if err := tiktok.DoGet(ctx, a.client, "/open_api/v1.3/smart_plus/material_report/overview/", params, &resp); err != nil {
		return nil, fmt.Errorf("failed to get material report overview: %w", err)
	}

	return &resp, nil
}
