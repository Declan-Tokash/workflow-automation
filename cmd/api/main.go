package main

import (
    "log"
    "net/http"

    "github.com/Declan-Tokash/workflow-automation/internal/auth"
    "github.com/Declan-Tokash/workflow-automation/internal/config"
	"github.com/Declan-Tokash/workflow-automation/internal/session"
	"github.com/Declan-Tokash/workflow-automation/internal/repos"
	"github.com/Declan-Tokash/workflow-automation/internal/clone"
)

func main() {
	cfg := config.Load()

	githubApp, err := auth.NewGitHubApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sessionStore := session.NewStore()

	authHandler := auth.NewHandler(
		githubApp,
		cfg,
		sessionStore,
	)

	repoHandler := repos.NewHandler(sessionStore)

	cloneService := clone.NewService()

	cloneHandler := clone.NewHandler(
		sessionStore,
		cloneService,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/github", authHandler.GitHubLogin)
	mux.HandleFunc("/auth/github/callback", authHandler.GitHubCallback)
	mux.HandleFunc("/api/me", authHandler.Me)
	mux.HandleFunc("/api/repos", repoHandler.List)
	mux.HandleFunc("/api/repos/", cloneHandler.Clone)

	log.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
