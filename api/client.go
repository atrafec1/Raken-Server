package api

import (
	// "bytes"
	// "encoding/json"
	"daily_check_in/config"
	"fmt"

	//"io"
	"net/http"
	"time"
)

// Client handles API communication with OAuth2 bearer token authentication
type Client struct {
	config     *config.Config
	httpClient *http.Client
	headers map[string]string

}

// NewClient creates and initializes a new API client
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	headers := {
		"Authorization": f"Bearer cfg.AccessToken"
	}
	client {
		config: cfg,
		httpClient: httpClient,
		headers
	}

	return client, nil


func get(client *Client, request string) (*http.Response, error) {
	url := client.config.BaseURL + request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Errorf("Error making GET request: %v", err)
	}

	req.Header := 


}
