package main

import (
    "log"
    "net/http"

    "github.com/Declan-Tokash/workflow-automation/internal/auth"
    "github.com/Declan-Tokash/workflow-automation/internal/config"
)

func main() {
	cfg := config.Load()

	githubApp, err := auth.NewGitHubApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authHandler := auth.NewHandler(githubApp, cfg)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/github", authHandler.GitHubLogin)
	mux.HandleFunc("/auth/github/callback", authHandler.GitHubCallback)

	log.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
