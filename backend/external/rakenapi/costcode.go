package rakenapi

import (
	"fmt"
	"net/http"
)

type CostCodeResponse struct {
	Collection []CostCode `json:"collection"`
}

type CostCode struct {
	UUID        string `json:"uuid"`
	Division    string `json:"division"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (c *Client) GetCostCodes() (*CostCodeResponse, error) {
	limit := "1000"
	requestURL := c.config.BaseURL + "costCodes"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cost code request: %w", err)
	}
	queryParams := req.URL.Query()
	queryParams.Set("limit", limit)
	req.URL.RawQuery = queryParams.Encode()
	var response CostCodeResponse
	if err = c.doRequest(req, &response); err != nil {
		return nil, fmt.Errorf("error retrieving cost codes: %w", err)
	}
	return &response, nil
}
