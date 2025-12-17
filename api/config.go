package api

import (
	"fmt"
	"os"
)

type Config struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	BaseURL      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		BaseURL:      os.Getenv("API_BASE_URL"),
	}

	// Validate required fields
	if cfg.ClientID == "" {
		return nil, fmt.Errorf("CLIENT_ID environment variable not set")
	}
	if cfg.ClientSecret == "" {
		return nil, fmt.Errorf("CLIENT_SECRET environment variable not set")
	}
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("API_BASE_URL environment variable not set")
	}

	return cfg, nil
}
