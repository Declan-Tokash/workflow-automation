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
	command ...string,
) (string, error) {

	args := []string{
		"exec",
		containerID,
	}

	args = append(args, command...)

	cmd := exec.CommandContext(
		ctx,
		"docker",
		args...,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf(
			"container command failed: %w",
			err,
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