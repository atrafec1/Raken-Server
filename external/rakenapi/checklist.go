package rakenapi

import (
	"fmt"
	"net/http"
)

type ChecklistResponse struct {
	Collection []ChecklistEntries
}

type ChecklistEntries struct {
	ProjectUUID string `json:"projectUuid"`
}

func (c *Client) GetCompletedChecklists() (*ChecklistResponse, error) {
	limit := "1000"
	requestURL := c.config.BaseURL + "checklists"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making checklist request: %v", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("limit", limit)
	queryParams.Set("statuses", "COMPLETED")
	req.URL.RawQuery = queryParams.Encode()
	var checklistResponse ChecklistResponse
	if err := c.doRequest(req, &checklistResponse); err != nil {
		return nil, fmt.Errorf("error retrieving checklist response: %v", err)
	}
	return &checklistResponse, nil
}
