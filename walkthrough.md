# Walkthrough: Go Workspace Migration

We have successfully migrated the project to a Go Workspace structure. This allows for better management of microservices and shared packages.

## New Structure

- **`go.work`**: The workspace definition file.
- **`pkg/`**: Shared libraries.
  - `database`: Database connection and helpers.
  - `models`: Shared data models.
  - `middlewares`: Shared HTTP middlewares.
- **`services/`**: Independent microservices.
  - `api-gateway`: The main API service (formerly the root application).

## How to Run

To run the API Gateway:

```bash
go run ./services/api-gateway/cmd/main.go
```

## How to Build

To build the API Gateway:

```bash
go build -v ./services/api-gateway/...
```

## Adding New Services

1. Create a new directory in `services/` (e.g., `services/auth-service`).
2. Initialize a new module:
   ```bash
   cd services/auth-service
   go mod init go-audio-stream/services/auth-service
   ```
3. Add it to the workspace:
   ```bash
   go work use ./services/auth-service
   ```
