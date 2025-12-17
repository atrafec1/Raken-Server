package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	config     *Config
	httpClient *http.Client
	headers    map[string]string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + cfg.AccessToken,
	}

	client := &Client{
		config:     cfg,
		httpClient: httpClient,
		headers:    headers,
	}

	return client, nil
}

func (client *Client) RefreshAccessToken() error {
	payload := url.Values{}
	payload.Set("grant_type", "refresh_token")
	payload.Set("client_id", client.config.ClientID)
	payload.Set("client_secret", client.config.ClientSecret)
	payload.Set("refresh_token", client.config.RefreshToken)
	fmt.Println("Refreshing access token...")
	req, err := http.NewRequest("POST", client.config.RefreshURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error retrieving refresh token: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return fmt.Errorf("error converting body to token response: %v", err)
	}
	fmt.Println("Access token refreshed successfully.")
	client.config.AccessToken = tokenResp.AccessToken
	client.config.RefreshToken = tokenResp.RefreshToken
	client.headers["Authorization"] = "Bearer " + tokenResp.AccessToken
	client.config.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}


func (c *Client) DoRequest(req *http.Request, respSchema interface{}) error {

    // Refresh if expired
    if time.Now().After(c.config.ExpiresAt) {
        if err := c.RefreshAccessToken(); err != nil {
            return err
        }
    }

    // Always set auth header from source of truth
    req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }

    // Retry once on 403
    if resp.StatusCode == http.StatusForbidden {
        resp.Body.Close()

        if err := c.RefreshAccessToken(); err != nil {
            return err
        }

        req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
        resp, err = c.httpClient.Do(req)
        if err != nil {
            return err
        }
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("request failed with status %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    return json.Unmarshal(body, respSchema)
}


type TimeCardResponse struct {
	Collection []TimeCard `json:"collection"`
	TotalElements int		`json:"totalElements"`
} 
type TimeCard struct {
	TimeEntries []TimeEntries `json:"timeEntries"`
	Worker      Worker        `json:"worker"`
} 
type TimeEntries struct {
	Hours float64 `json:"hours"`
	Date  string  `json:"date"`
	PayType PayType `json:"payType"`
	Classification Classification `json:"classification"`
	CostCode CostCode `json:"costCode"`
	Project Project `json:"project"`
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
	Code string `json:"code"`
	Division string `json:"division"`
}

type Project struct {
	UUID string `json:"uuid"`
}

func (c *Client) GetTimecards(fromDate, toDate string) (*TimeCardResponse, error) {
    // Use := for new variables
	limit := 1000
    endpoint := fmt.Sprintf("timeCards?fromDate=%s&toDate=%s&limit=%v", fromDate, toDate, limit)
    fullURL := c.config.BaseURL + endpoint
    
    // Use http.NewRequest (your snippet had httpClient.NewRequest)
    req, err := http.NewRequest("GET", fullURL, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    var timeCardResp TimeCardResponse
    // Pass the pointer so DoRequest can fill it
    err = c.DoRequest(req, &timeCardResp) 
    if err != nil {
        return nil, fmt.Errorf("error getting timecards: %w", err)
    }

    return &timeCardResp, nil
}


