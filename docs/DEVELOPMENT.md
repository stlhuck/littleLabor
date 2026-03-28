# Development Guide

## Setting Up Development Environment

### Go
- Install Go 1.21 or later from https://golang.org/dl/

### Building
```bash
cd ~/go/src/littleLabor
go mod download
go build -o bin/littleLabor ./cmd/server
```

### Running
```bash
./bin/littleLabor
```

Visit `http://localhost:8080` to verify the API is running.

## Testing
```bash
go test ./...
```

## Frontend Development
(To be added - React/Vue web interface)

## Mobile Development
(To be added - iOS/Android apps)
