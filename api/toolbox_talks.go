package api

import (
	"fmt"
	"net/http"
)


type ToolboxTalkResponse struct {
	Collection []ToolboxTalkEntry `json:"collection"`
}

type ToolboxTalkEntry struct {
	Project Project `json:"project"`
	Attendees []AttendeeRecord `json:"attendees"`
	Date string `json:"scheduleDate"`
}

type AttendeeRecord struct {
	Employee Employee `json:"member"`
}

func (c *Client) GetToolboxTalks() (*ToolboxTalkResponse, error) {
	requestURL := c.config.BaseURL + "toolboxTalks/past"
	limit := "1000"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making toolbox talk request: %v", err)
	}

	queryParams := req.URL.Query()
	queryParams.Set("limit", limit)
	req.URL.RawQuery = queryParams.Encode()

	var toolboxTalkResp ToolboxTalkResponse 

	if err := c.doRequest(req, &toolboxTalkResp); err != nil {
		return nil, fmt.Errorf("error retrieving toolbox talks: %v", err)
	}
	return &toolboxTalkResp, nil
}