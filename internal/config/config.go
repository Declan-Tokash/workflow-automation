package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GitHubAppID          string
	GitHubClientID       string
	GitHubClientSecret   string
	GitHubPrivateKeyPath string
	GitHubRedirectURL    string
}

func Load() Config {
	godotenv.Load()

	return Config{
		GitHubAppID:          os.Getenv("GITHUB_APP_ID"),
		GitHubClientID:       os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
		GitHubPrivateKeyPath: os.Getenv("GITHUB_PRIVATE_KEY_PATH"),
		GitHubRedirectURL:    os.Getenv("GITHUB_REDIRECT_URL"),
	}
}
