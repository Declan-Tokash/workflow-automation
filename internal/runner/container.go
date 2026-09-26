package runner

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type ContainerRunner struct{}

func NewContainerRunner() *ContainerRunner {
	return &ContainerRunner{}
}

func (r *ContainerRunner) Create(
	ctx context.Context,
	image string,
) (string, error) {

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"create",
		"-w",
		"/workspace",
		image,
		"sleep",
		"3600",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"failed to create container: %s",
			string(output),
		)
	}

	containerID := strings.TrimSpace(string(output))

	return containerID, nil
}

func (r *ContainerRunner) Start(
	ctx context.Context,
	containerID string,
) error {

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"start",
		containerID,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to start container: %s",
			string(output),
		)
	}

	return nil
}

func (r *ContainerRunner) Exec(
	ctx context.Context,
	containerID string,
	command string,
) (string, error) {

	args := []string{
		"exec",
		"-w",
		"/workspace",
		containerID,
		"/bin/sh",
		"-c",
		command,
	}

	cmd := exec.CommandContext(ctx, "docker", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"container command failed: %w\noutput: %s",
			err,
			string(output),
		)
	}

	return string(output), nil
}

func (r *ContainerRunner) Remove(
	ctx context.Context,
	containerID string,
) error {

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"rm",
		"-f",
		containerID,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to remove container: %s",
			string(output),
		)
	}

	return nil
}

func (r *ContainerRunner) Clone(
    ctx context.Context,
    containerID string,
    cloneURL string,
    token string,
) (string, error) {

    askpass := fmt.Sprintf(`#!/bin/sh

    case "$1" in
        Username*)
            echo "x-access-token"
            ;;
        Password*)
            echo '%s'
            ;;
    esac
    `, token)

    // Create the git-askpass script inside the container
    _, err := r.Exec(
        ctx,
        containerID,
        fmt.Sprintf(
            "printf '%%s' '%s' > /tmp/git-askpass && chmod 700 /tmp/git-askpass",
            askpass,
        ),
    )

    if err != nil {
        return "", fmt.Errorf(
            "failed to create git auth helper: %w",
            err,
        )
    }

    // Clone the repository using the git-askpass script
    output, err := r.Exec(
        ctx,
        containerID,
        fmt.Sprintf(
            "GIT_ASKPASS=/tmp/git-askpass "+
                "GIT_TERMINAL_PROMPT=0 "+
                "git clone '%s' /workspace/repo",
            cloneURL,
        ),
    )

    if err != nil {
        return output, fmt.Errorf(
            "failed to clone repository: %w",
            err,
        )
    }

    // Remove the git-askpass script after cloning
    _, _ = r.Exec(
        ctx,
        containerID,
        "rm -f /tmp/git-askpass",
    )

    return output, nil
}