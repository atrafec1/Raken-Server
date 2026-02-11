package rakenapi

import (
	"fmt"
	"net/http"
)

type TimeCardResponse struct {
	Collection    []TimeCard `json:"collection"`
	TotalElements int        `json:"totalElements"`
}
type TimeCard struct {
	TimeEntries []TimeEntries `json:"timeEntries"`
	Worker      Worker        `json:"worker"`
	Project     Project       `json:"project"`
	Date        string        `json:"date"`
}
type TimeEntries struct {
	Hours          float64        `json:"hours"`
	PayType        PayType        `json:"payType"`
	Classification Classification `json:"classification"`
	CostCode       CostCode       `json:"costCode"`
}

type Worker struct {
	UUID string `json:"uuid"`
}

type PayType struct {
	Code string `json:"code"`
}

type Classification struct {
	Name string `json:"name"`
}

type CostCode struct {
	Code     string `json:"code"`
	Division string `json:"division"`
}

func (c *Client) GetTimeCards(fromDate, toDate string) (*TimeCardResponse, error) {
	limit := "1000"
	fullURL := c.config.BaseURL + "timeCards"

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("fromDate", fromDate)
	queryParams.Set("toDate", toDate)
	queryParams.Set("limit", limit)

	req.URL.RawQuery = queryParams.Encode()

	var timeCardResp TimeCardResponse
	// Pass the pointer so DoRequest can fill it
	err = c.doRequest(req, &timeCardResp)
	if err != nil {
		return nil, fmt.Errorf("error getting timecards: %w", err)
	}

	return &timeCardResp, nil
}

