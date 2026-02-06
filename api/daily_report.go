package api

import (
	"fmt"
	"net/http"
)

type DailyReportResponse struct {
	Collection string `json:"collection"`
}
type DailyReport struct {
	Status      string      `json:"status"`
	ReportDate  string      `json:"reportDate"`
	ProjectUuid string      `json:"projectUuid"`
	ReportLinks ReportLinks `json:"reportLinks"`
}

type ReportLinks struct {
	Link string `json:"daily"`
}

func (c *Client) GetDailyReports(projectUuid string, fromDate string, toDate string) (DailyReportResponse, error) {
	limit := "1000"
	requestURL := c.config.BaseURL + fmt.Sprintf("dailyReports?projectUuid=%s&fromDate=%s&toDate=%s", projectUuid, fromDate, toDate)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return DailyReportResponse{}, fmt.Errorf("error making daily reports request: %v", err)
	}
	queryParams := req.URL.Query()
	queryParams.Set("limit", limit)
	req.URL.RawQuery = queryParams.Encode()
	var dailyReports DailyReportResponse
	if err := c.doRequest(req, &dailyReports); err != nil {
		return DailyReportResponse{}, fmt.Errorf("error retrieving daily reports: %v", err)
	}
	return dailyReports, nil
}
