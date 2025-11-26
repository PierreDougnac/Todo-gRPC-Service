# Todo gRPC Service

A **high-performance Todo service** implemented in **Go** using **gRPC**.  
It provides a fully functional server and CLI client for creating, listing, retrieving, and deleting todos.

---

## Features

- gRPC server and client in Go
- In-memory storage for todos
- Fully tested with unit tests
- Easy build and test automation via `build.ps1`

---

## Quick Start

### Prerequisites

- Go 1.25+
- Protocol Buffers (`protoc`)
- `protoc-gen-go` and `protoc-gen-go-grpc`

### Build & Run

```powershell
# Generate gRPC code, build server/client, and run tests
.\build.ps1

# Start the server
go run ./server/main.go

# Use the client
go run ./client/main.go add "Buy milk"
go run ./client/main.go list
go run ./client/main.go get 1
go run ./client/main.go delete 1

Project Structure
grpc-todo/
├── client/       # CLI client and tests
├── server/       # gRPC server and tests
├── proto/        # Protocol Buffers definitions
├── build.ps1     # Build/test automation
├── go.mod
└── go.sum

Testing
# Run unit tests separately
go test ./server/...
go test ./client/...

# Or rely on the build script
.\build.ps1