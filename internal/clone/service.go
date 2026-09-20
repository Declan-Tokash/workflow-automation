package clone

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Clone(
	ctx context.Context,
	cloneURL string,
	token string,
) (string, error) {

	// Create temporary workspace
	dir, err := os.MkdirTemp("", "sandbox-repo-*")
	if err != nil {
		return "", fmt.Errorf("create workspace: %w", err)
	}

	// Create a temporary askpass script.
	// Git executes this when it needs credentials.
	askpassPath := filepath.Join(dir, "git-askpass.sh")

	askpass := fmt.Sprintf(`#!/bin/sh echo "%s"`, token)

	if err := os.WriteFile(askpassPath, []byte(askpass), 0700); err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("create git auth helper: %w", err)
	}

	// Clone into a separate directory so the askpass
	// script isn't inside the repository.
	repoDir := filepath.Join(dir, "repo")

	cmd := exec.CommandContext(
		ctx,
		"git",
		"clone",
		cloneURL,
		repoDir,
	)

	cmd.Env = append(
		os.Environ(),
		"GIT_ASKPASS="+askpassPath,
		"GIT_TERMINAL_PROMPT=0",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(dir)

		return "", fmt.Errorf(
			"git clone failed: %s",
			string(output),
		)
	}

	// Remove the authentication helper after cloning.
	os.Remove(askpassPath)

	return repoDir, nil
}