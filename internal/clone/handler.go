package clone

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Declan-Tokash/workflow-automation/internal/github"
	"github.com/Declan-Tokash/workflow-automation/internal/session"
)

type Handler struct {
	Sessions *session.Store
	Service  *Service
}

func NewHandler(
	sessions *session.Store,
	service *Service,
) *Handler {
	return &Handler{
		Sessions: sessions,
		Service:  service,
	}
}

type CloneResponse struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
}

func (h *Handler) Clone(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get session
	userSession, ok := h.Sessions.Get(cookie.Value)
	if !ok {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	// Get owner/repo from URL
	parts := strings.Split(
		strings.Trim(r.URL.Path, "/"),
		"/",
	)

	// Expected:
	// /api/repos/{owner}/{repo}/clone
	if len(parts) != 5 ||
		parts[0] != "api" ||
		parts[1] != "repos" ||
		parts[4] != "clone" {

		http.Error(w, "invalid repository path", http.StatusBadRequest)
		return
	}

	owner := parts[2]
	repo := parts[3]

	fullName := owner + "/" + repo

	// Create GitHub client using the token
	githubClient := github.NewClient(
		r.Context(),
		userSession.AccessToken,
	)

	// Get repositories belonging to authenticated user
	repositories, err := githubClient.GetRepositories(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to fetch repositories",
			http.StatusInternalServerError,
		)
		return
	}

	// Find requested repository
	var cloneURL string

	for _, repository := range repositories {
		if repository.FullName == fullName {
			cloneURL = repository.CloneURL
			break
		}
	}

	if cloneURL == "" {
		http.Error(
			w,
			"repository not found",
			http.StatusNotFound,
		)
		return
	}

	// Clone repository
	path, err := h.Service.Clone(
		r.Context(),
		cloneURL,
		userSession.AccessToken,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := CloneResponse{
		Repository: fullName,
		Path:       path,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}