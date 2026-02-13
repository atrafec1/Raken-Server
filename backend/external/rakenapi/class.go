package rakenapi

import (
	"fmt"
	"net/http"
)

type ClassResponse struct {
	Collection []Class `json:"collection"`
}

type Class struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (c *Client) GetClasses() (*ClassResponse, error) {
	limit := "1000"
	requestURL := c.config.BaseURL + "classifications"

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making class request: %v", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("limit", limit)
	req.URL.RawQuery = queryParams.Encode()
	var response ClassResponse
	err = c.doRequest(req, &response)
	if err != nil {
		return nil, fmt.Errorf("error retrieving classifications: %v", err)
	}

	return &response, nil
}
