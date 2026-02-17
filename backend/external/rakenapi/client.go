package rakenapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type RakenClient interface {
	GetTimeCards(from, to string) (*TimeCardResponse, error)
	GetEquipmentLogs(from, to string) (*EquipmentLogResponse, error)
	GetProjects() (*ProjectResponse, error)
	GetEmployees() (*EmployeeResponse, error)
	GetClasses() (*ClassResponse, error)
	GetMaterialLogs(projectUuid, from, to string) (*MaterialLogResponse, error)
}

type Client struct {
	config     *Config
	httpClient *http.Client
	mu         sync.Mutex
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
	c := &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	return c, nil
}

func (c *Client) refreshAccessToken() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("refreshing access token")

	// Another goroutine may have refreshed already
	if time.Now().Before(c.config.ExpiresAt) {
		return nil
	}

	payload := url.Values{}
	payload.Set("grant_type", "refresh_token")
	payload.Set("client_id", c.config.ClientID)
	payload.Set("client_secret", c.config.ClientSecret)
	payload.Set("refresh_token", c.config.RefreshToken)

	req, err := http.NewRequest(
		"POST",
		c.config.RefreshURL,
		strings.NewReader(payload.Encode()),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("refresh failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("error decoding token response: %w", err)
	}

	c.config.AccessToken = tokenResp.AccessToken

	if tokenResp.RefreshToken != "" {
		c.config.RefreshToken = tokenResp.RefreshToken
	}

	c.config.ExpiresAt = time.Now().Add(
		time.Duration(tokenResp.ExpiresIn) * time.Second,
	)

	return nil
}

func (c *Client) doRequest(req *http.Request, respSchema interface{}) error {
	if time.Now().After(c.config.ExpiresAt) {
		if err := c.refreshAccessToken(); err != nil {
			return err
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	print("Requesting ", req.URL.String(), "\n")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if err := c.refreshAccessToken(); err != nil {
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
