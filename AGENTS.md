# Agent Instructions for go-info-share

This is a Go-based key-value store with WebSocket support and a CLI client.

## Build/Test/Lint Commands

### Build
```bash
# Build the server
go build -o server main.go

# Build the CLI
go build -o cli cli.go

# Download dependencies
go mod download

# Update dependencies
go mod tidy
```

### Test
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run a single test
go test -run TestFunctionName ./...

# Run tests in a specific package
go test ./package_name
```

### Lint/Format
```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run static analysis (if installed)
staticcheck ./...

# Run linter (if installed)
golint ./...
```

### Run
```bash
# Run the server
go run main.go
# or
./server

# Run the CLI
go run cli.go <key> <value>
# or
./cli <key> <value>

# CLI with custom URL
go run cli.go --url http://localhost:8080 <key> <value>
```

## Code Style Guidelines

### Formatting
- Use `go fmt` for automatic formatting
- Use tabs for indentation (Go standard)
- Maximum line length: ~100 characters when practical
- Group related imports together with empty lines between groups

### Imports
Standard Go import grouping:
```go
import (
    "standard/library"
    "packages"
    
    "github.com/third/party"
    
    "github.com/matst80/go-info-share/internal"
)
```
- Standard library imports first
- Third-party imports second
- Internal/project imports last
- Use `goimports` to manage imports automatically

### Naming Conventions
- **Packages**: lowercase, single word (e.g., `main`, `handlers`)
- **Types/Structs**: PascalCase (e.g., `KVStore`, `WebSocketHandler`)
- **Functions/Methods**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase (e.g., `connMu`, `upgrader`)
- **Constants**: PascalCase or ALL_CAPS for exported constants
- **Interfaces**: PascalCase with -er suffix when appropriate (e.g., `Reader`, `Writer`)
- **Acronyms**: Keep consistent case (e.g., `KVStore`, not `KvStore`; `HTTPHandler` not `HttpHandler`)

### Types and Structs
```go
type KVStore struct {
    data   map[string]string
    mu     sync.RWMutex
    conns  []*websocket.Conn
    connMu sync.Mutex
}
```
- Place mutex close to the data it protects
- Use meaningful field names
- Document exported types with comments starting with the type name

### Error Handling
- Always check errors and handle them explicitly
- Use `fmt.Errorf()` for wrapping errors with context
- Use `errors.New()` for simple error messages
- Log errors using `log.Println()` or `log.Printf()`
- Return HTTP errors using `http.Error()` with appropriate status codes
- Use `defer` for cleanup (e.g., `defer resp.Body.Close()`)

Example:
```go
if err != nil {
    log.Println("error:", err)
    http.Error(w, "message", http.StatusInternalServerError)
    return
}
```

### HTTP Handlers
- Wrap handlers to inject dependencies (e.g., `func handler(kv *KVStore) http.HandlerFunc`)
- Set CORS headers at the beginning of handlers
- Handle OPTIONS requests for CORS preflight
- Return appropriate HTTP status codes (200, 400, 404, 405)
- Use `http.Error()` for error responses
- Set `Content-Type` headers when returning JSON

### Concurrency
- Use `sync.RWMutex` for read-heavy workloads
- Use `sync.Mutex` for general synchronization
- Keep critical sections small
- Always unlock with `defer` when possible, or use explicit unlock immediately after operation

Example:
```go
func (k *KVStore) Get(key string) (string, bool) {
    k.mu.RLock()
    v, ok := k.data[key]
    k.mu.RUnlock()
    return v, ok
}
```

### WebSocket Handling
- Use `gorilla/websocket` for WebSocket connections
- Check origin in `Upgrader.CheckOrigin`
- Handle connection cleanup on error/disconnect
- Use mutexes when managing connection lists

### Comments
- All exported types, functions, and methods should have doc comments
- Comments start with the name of the exported identifier
- Use inline comments sparingly; prefer self-documenting code

Example:
```go
// KVStore manages key-value pairs with concurrent access and WebSocket broadcasting.
type KVStore struct {
    // ...
}
```

### Project Structure
- `main.go`: Server entry point with HTTP handlers and WebSocket support
- `cli.go`: CLI client entry point
- `go.mod`: Module definition
- `Dockerfile`: Multi-stage Docker build
- `.github/workflows/`: GitHub Actions for releases
- `*.yaml`: Kubernetes deployment manifests

## Docker Commands
```bash
# Build Docker image
docker build -t go-info-share .

# Run container
docker run -p 8080:8080 go-info-share
```

## Testing
When adding tests:
- Use `*_test.go` naming convention
- Use table-driven tests for multiple test cases
- Name test functions with `Test` prefix (e.g., `TestKVStore_Set`)
- Use `testing.T` for unit tests, `testing.B` for benchmarks
- Place test files in the same package as the code being tested

## Notes
- The project currently has two main entry points: `main.go` (server) and `cli.go` (CLI)
- Server runs on port 8080 by default
- CLI defaults to `http://localhost:8080` or uses `INFO_SERVER_URL` env var
