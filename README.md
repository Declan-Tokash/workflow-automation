# cloud-job-runner

An ephemeral containerized job execution platform built in Go.

This project is a learning-oriented backend for running isolated CI-style jobs: clone a GitHub repository, execute a command in a selected runtime, capture logs and exit code, then destroy the container.

## Current architecture

```text
HTTP API  →  Job Service  →  Executor (stub)  →  (Docker later)
```

## Prerequisites

- Go 1.24+
- Docker (for future executor work)

## Getting started

```bash
go run ./cmd/api
```

The server listens on port `8080` by default. Override with:

```bash
PORT=8080 LOG_LEVEL=info go run ./cmd/api
```

## API

### Health check

```bash
curl http://localhost:8080/health
```

### Create a job

```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "repository_url": "https://github.com/example/demo",
    "runtime": "go",
    "command": "go test ./..."
  }'
```

### Get a job

```bash
curl http://localhost:8080/api/v1/jobs/<job-id>
```

## Supported runtimes

| Runtime | Docker image           |
|---------|------------------------|
| python  | python:3.12            |
| java    | eclipse-temurin:21     |
| go      | golang:1.24            |
| node    | node:22                |

## Project layout

```text
cmd/api/              Application entrypoint
internal/api/         HTTP handlers and routing
internal/config/      Environment-based configuration
internal/jobs/        Job model and service layer
internal/executor/    Execution abstraction (Docker impl later)
internal/runtime/     Approved runtime → Docker image mapping
```

## Next step

Implement a real Docker-backed `Executor` that creates an ephemeral container and runs a command. See the architecture notes in the repository docs or ask the agent for guided help on that step.
