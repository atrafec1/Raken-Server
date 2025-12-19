package api

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
	resp, err := c.doRequest(req, *response) 
	if err != nil {
		return nil, fmt.Errorf("error retrieving Class: %v", err)
	}

	return resp, nil
}

func (c *Client) UpdateClassMap() error {
	classResp, err := c.GetClasses()
	if err != nil {
		return fmt.Errof("error getting classifications: %v", err)
	}
	c.mu.Lock()
	defer.c.mu.Unlock()
	c.classMap = make(map[string]Class)
	for _, class := range classResp.Collection {
		c.classMap[class.UUID] = class
	}
	return nil
}

func (c *Client) GetClassByUUID(uuid string) (Class, err) {
	class, ok := c.classMap[uuid]
	if !ok {
		return nil, fmt.Errorf("Did not find a class for uuid: %v", uuid)
	}
	return class, nil
}

