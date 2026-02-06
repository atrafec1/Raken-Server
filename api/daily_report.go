package api

import (
	"fmt"
	"net/http"
)

type DailyReportResponse struct {
	Collection []DailyReport `json:"collection"`
}
type DailyReport struct {
	Status      string      `json:"status"`
	ReportDate  string      `json:"reportDate"`
	ProjectUuid string      `json:"projectUuid"`
	ReportLinks ReportLinks `json:"reportLinks"`
	SignedBy    Creator     `json:"signedBy"`
}

type ReportLinks struct {
	Link string `json:"daily"`
}

type Creator struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

func (c *Client) GetDailyReports(projectUuid string, fromDate string, toDate string) (*DailyReportResponse, error) {
	limit := "1000"
	requestURL := c.config.BaseURL + fmt.Sprintf("dailyReports?projectUuid=%s&fromDate=%s&toDate=%s", projectUuid, fromDate, toDate)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making daily reports request: %v", err)
	}
	queryParams := req.URL.Query()
	queryParams.Set("limit", limit)
	req.URL.RawQuery = queryParams.Encode()
	var dailyReports DailyReportResponse
	if err := c.doRequest(req, &dailyReports); err != nil {
		return nil, fmt.Errorf("error retrieving daily reports: %v", err)
	}
	return &dailyReports, nil
}
