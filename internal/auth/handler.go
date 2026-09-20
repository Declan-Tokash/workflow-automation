package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/Declan-Tokash/workflow-automation/internal/config"
	"github.com/Declan-Tokash/workflow-automation/internal/github"
	"github.com/Declan-Tokash/workflow-automation/internal/session"
)

type Handler struct {
	GitHubApp *GitHubApp
	Config    config.Config
	Sessions  *session.Store
}

func NewHandler(
	githubApp *GitHubApp,
	cfg config.Config,
	sessionStore *session.Store,
) *Handler {
	return &Handler{
		GitHubApp: githubApp,
		Config:    cfg,
		Sessions:  sessionStore,
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

	githubClient := github.NewClient(
		r.Context(),
		token.AccessToken,
	)

	user, err := githubClient.GetUser(r.Context())
	if err != nil {
		http.Error(w, "failed to get GitHub user", http.StatusInternalServerError)
		return
	}

	sessionID := h.Sessions.Create(session.Session{
		UserID:      user.ID,
		GitHubLogin: user.Login,
		AccessToken: token.AccessToken,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(
		w,
		r,
		"/api/me",
		http.StatusSeeOther,
	)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	userSession, ok := h.Sessions.Get(cookie.Value)
	if !ok {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(
		w,
		"Logged in as %s (GitHub ID: %d)",
		userSession.GitHubLogin,
		userSession.UserID,
	)
}