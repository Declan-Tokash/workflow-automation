package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type Client struct {
	OAuthClient *http.Client
}

type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	Private  bool   `json:"private"`
}

func NewClient(ctx context.Context, accessToken string) *Client {
	oauth2Token := &oauth2.Token{
		AccessToken: accessToken,
	}

	httpClient := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(oauth2Token),
	)

	return &Client{
		OAuthClient: httpClient,
	}
}

type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

func (c *Client) GetUser(ctx context.Context) (*User, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.OAuthClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"github API returned status %d",
			resp.StatusCode,
		)
	}

	var user User

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) GetRepositories(ctx context.Context) ([]Repository, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user/repos",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.OAuthClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"github API returned status %d",
			resp.StatusCode,
		)
	}

	var repos []Repository

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}