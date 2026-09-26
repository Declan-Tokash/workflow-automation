package repos

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/Declan-Tokash/workflow-automation/internal/session"
	"github.com/Declan-Tokash/workflow-automation/internal/github"
)

type Handler struct {
	Sessions *session.Store
}

func NewHandler(sessions *session.Store) *Handler {
	return &Handler{
		Sessions: sessions,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
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

	githubClient := github.NewClient(
		r.Context(),
		userSession.AccessToken,
	)

	repos, err := githubClient.GetRepositories(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get repositories",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(repos)
}