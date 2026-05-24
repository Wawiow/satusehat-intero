package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AuthURL      string
	BaseURL      string
	ConsentURL   string
	ClientID     string
	ClientSecret string
	OrgID        string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("No .env file found: " + err.Error())
	}

	authURL := os.Getenv("AUTH_URL")
	baseURL := os.Getenv("BASE_URL")
	consentURL := os.Getenv("CONSENT_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	orgID := os.Getenv("ORG_ID")

	return &Config{
		AuthURL:      authURL,
		BaseURL:      baseURL,
		ConsentURL:   consentURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		OrgID:        orgID,
	}
}
