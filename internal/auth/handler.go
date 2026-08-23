package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/declantokash/workflow-automation-dt/internal/config"
)

type Handler struct {
	GitHubApp *GitHubApp
	Config    config.Config
}

func NewHandler(githubApp *GitHubApp, cfg config.Config) *Handler {
	return &Handler{
		GitHubApp: githubApp,
		Config:    cfg,
	}
}

func (h *Handler) GitHubLogin(w http.ResponseWriter, r *http.Request) {

	stateBytes := make([]byte, 32)

	if _, err := rand.Read(stateBytes); err != nil {
		http.Error(w, "failed to generate state", http.StatusInternalServerError)
		return
	}

	state := base64.URLEncoding.EncodeToString(stateBytes)

	oauthConfig := &oauth2.Config{
		ClientID:     h.Config.GitHubClientID,
		ClientSecret: h.Config.GitHubClientSecret,
		RedirectURL:  h.Config.GitHubRedirectURL,

		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},

		Scopes: []string{
			"repo",
		},
	}

	url := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	oauthConfig := &oauth2.Config{
		ClientID:     h.Config.GitHubClientID,
		ClientSecret: h.Config.GitHubClientSecret,
		RedirectURL:  h.Config.GitHubRedirectURL,

		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Successfully authenticated with GitHub!"))

	_ = token
}