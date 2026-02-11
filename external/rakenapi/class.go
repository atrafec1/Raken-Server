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

func (c *Client) UpdateClassMap() error {
	classResp, err := c.GetClasses()
	if err != nil {
		return fmt.Errorf("error getting classifications: %v", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.classMap = make(map[string]Class)
	for _, class := range classResp.Collection {
		c.classMap[class.UUID] = class
	}
	return nil
}

func (c *Client) GetClassByUUID(uuid string) (Class, error) {
	class, exists := c.classMap[uuid]
	if !exists {
		if err := c.UpdateClassMap(); err != nil {
			return Class{}, fmt.Errorf("failed to refresh class map: %w", err)
		}
		class, exists = c.classMap[uuid]
		if !exists {
			return Class{}, fmt.Errorf("class with UUID %s not found after refresh", uuid)
		}
	}
	return class, nil
}

