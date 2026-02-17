package rakenapi

import (
	"fmt"
	"net/http"
)

type MaterialLogResponse struct {
	Collection []MaterialLog `json:"collection"`
}

type MaterialLog struct {
	Date     string   `json:"date"`
	Material Material `json:"material"`
	Quantity float64  `json:"quantity"`
}

type Material struct {
	UUID string       `json:"uuid"`
	Name string       `json:"name"`
	Unit MaterialUnit `json:"materialUnit"`
}

type MaterialUnit struct {
	Name string `json:"name"`
}

func (c *Client) GetMaterialLogs(projectUuid, fromDate, toDate string) (*MaterialLogResponse, error) {
	requestURL := c.config.BaseURL + "materialLogs"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("bad request creation: %w", err)
	}
	queryParams := req.URL.Query()
	queryParams.Set("fromDate", fromDate)
	queryParams.Set("toDate", toDate)
	queryParams.Set("projectUuid", projectUuid)
	req.URL.RawQuery = queryParams.Encode()

	var resp MaterialLogResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to get material logs: %w", err)
	}
	return &resp, nil
}
