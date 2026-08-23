package auth

import (
	"os"

	// "github.com/golang-jwt/jwt/v5"

	"github.com/declantokash/workflow-automation-dt/internal/config"
)

type GitHubApp struct {
	AppID      string
	PrivateKey []byte
}

func NewGitHubApp(cfg config.Config) (*GitHubApp, error) {
	privateKey, err := os.ReadFile(cfg.GitHubPrivateKeyPath)
	if err != nil {
		return nil, err
	}

	return &GitHubApp{
		AppID:      cfg.GitHubAppID,
		PrivateKey: privateKey,
	}, nil
}