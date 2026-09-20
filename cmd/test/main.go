package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Declan-Tokash/workflow-automation/internal/runner"
)

func main() {
	ctx := context.Background()

	r := runner.NewContainerRunner()

	containerID, err := r.Create(
		ctx,
		"workflow-python:latest",
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

	_, err = r.Clone(
		ctx,
		containerID,
		"https://github.com/Declan-Tokash/workflow-automation.git",
		os.Getenv("ACCESS_TOKEN"),
	)
	
	if err != nil {
		log.Fatal(err)
	}

	output, err := r.Exec(
		ctx,
		containerID,
		"pwd",
	)
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(output)

	output, err = r.Exec(
		ctx,
		containerID,
		"sh",
		"-c",
		"cd repo && ls",
	)
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(output)

	// output, err := r.Exec(
	// 	ctx,
	// 	containerID,
	// 	"python",
	// 	"--version",
	// )
	// output, err := r.Exec(
	// 	ctx,
	// 	containerID,
	// 	"git",
	// 	"--version",
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Output:", output)
}