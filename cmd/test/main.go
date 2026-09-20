package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Declan-Tokash/workflow-automation/internal/runner"
)

func main() {
	ctx := context.Background()

	r := runner.NewContainerRunner()

	containerID, err := r.Create(
		ctx,
		"python:3.12",
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		fmt.Println("Removing container...")

		if err := r.Remove(ctx, containerID); err != nil {
			fmt.Println("Failed to remove:", err)
		}
	}()

	if err := r.Start(ctx, containerID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Container started")

	output, err := r.Exec(
		ctx,
		containerID,
		"python",
		"--version",
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Output:", output)
}