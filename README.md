# Project go-audio-stream

This project is a microservices-based audio streaming application built with Go. It uses a **Go Workspace** structure to manage multiple services and shared packages in a single monorepo.

## Workspace Structure

- **`go.work`**: Workspace definition.
- **`services/`**: Microservices.
  - `catalog-service`: The main entry point for the API.
- **`pkg/`**: Shared libraries.
  - `database`: Database connection and helpers.
  - `models`: Shared data models.
  - `middlewares`: Shared HTTP middlewares.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go 1.18+
- Docker & Docker Compose

## Makefile Commands

Run build make command with tests
```bash
make all
```

Build the API Gateway
```bash
make build
```

Run the API Gateway
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application (requires `air`):
```bash
make watch
```

Run the test suite for all services and packages:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

Tidy all Go modules:
```bash
make tidy
```
