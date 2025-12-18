package api

import (
	"fmt"
	"net/http"
)

type ProjectResponse struct {
	Collection []Project `json:"collection"`
}

type Project struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Number string `json:"number"`
}

func (c *Client) GetProjects() (*ProjectResponse, error) {
	requestURL := c.config.BaseURL + "projects"
	limit := "1000"
	req, err := http.NewRequest("GET",requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making project request: %v", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("limit",limit)
	req.URL.RawQuery = queryParams.Encode()	
	var projects ProjectResponse
	if err := c.DoRequest(req, &projects); err != nil {
		return nil, fmt.Errorf("error retrieving projects %v", err)
	}
	return projects, nil
}