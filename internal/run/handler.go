package run

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Declan-Tokash/workflow-automation/internal/clone"
	"github.com/Declan-Tokash/workflow-automation/internal/github"
	"github.com/Declan-Tokash/workflow-automation/internal/runner"
	"github.com/Declan-Tokash/workflow-automation/internal/session"
)

type Handler struct {
	Sessions *session.Store
	Clone    *clone.Service
	Runner   *runner.DockerRunner
}

func NewHandler(
	sessions *session.Store,
	cloneService *clone.Service,
	dockerRunner *runner.DockerRunner,
) *Handler {
	return &Handler{
		Sessions: sessions,
		Clone:    cloneService,
		Runner:   dockerRunner,
	}
}

type RunRequest struct {
	Repository string `json:"repository"`
	Runtime    string `json:"runtime"`
	Command    string `json:"command"`
}

type RunResponse struct {
	Repository string `json:"repository"`
	Runtime    string `json:"runtime"`
	Output     string `json:"output"`
}

func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	// 1. Get session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userSession, ok := h.Sessions.Get(cookie.Value)
	if !ok {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	// 2. Parse request
	var req RunRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Validate request
	if req.Repository == "" ||
		req.Runtime == "" ||
		req.Command == "" {

		http.Error(
			w,
			"repository, runtime, and command are required",
			http.StatusBadRequest,
		)
		return
	}

	// Expected format:
	// owner/repository
	parts := strings.Split(req.Repository, "/")

	if len(parts) != 2 {
		http.Error(
			w,
			"repository must be in owner/repository format",
			http.StatusBadRequest,
		)
		return
	}

	owner := parts[0]
	repo := parts[1]

	// 4. Create GitHub client using user's token
	githubClient := github.NewClient(
		r.Context(),
		userSession.AccessToken,
	)

	// 5. Find repository
	repositories, err := githubClient.GetRepositories(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to fetch repositories",
			http.StatusInternalServerError,
		)
		return
	}

	var cloneURL string

	for _, repository := range repositories {
		if repository.FullName == owner+"/"+repo {
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

	// 6. Clone repository
	repoPath, err := h.Clone.Clone(
		r.Context(),
		cloneURL,
		userSession.AccessToken,
	)

	if err != nil {
		http.Error(
			w,
			"failed to clone repository",
			http.StatusInternalServerError,
		)
		return
	}

	// 7. Select Docker image
	image, err := runtimeImage(req.Runtime)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	// 8. Run command
	output, err := h.Runner.Run(
		r.Context(),
		repoPath,
		image,
		req.Command,
	)

	if err != nil {
		http.Error(
			w,
			output,
			http.StatusInternalServerError,
		)
		return
	}

	// 9. Return result
	response := RunResponse{
		Repository: req.Repository,
		Runtime:    req.Runtime,
		Output:     output,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}