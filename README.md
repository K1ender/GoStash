# GoStash

<img src="./logo.png" />

GoStash is a lightweight in-memory key-value cache inspired by Redis and Memcached. It aims to provide a simple, fast, and easy-to-run caching server for development and small production use-cases. The project is written in Go and follows a minimal design that focuses on low-latency operations and simple configuration.

## Features

- **In-memory key-value storage** - Fast HashMap-based storage with thread-safe operations
- **Binary protocol support** - Custom binary protocol for efficient communication
- **Configurable server** - Support for both file-based and CLI configuration
- **Concurrent client handling** - Each client connection handled in a separate goroutine
- **Basic commands** - GET and SET operations with proper serialization
- **Small codebase** - Intended for learning, experimentation and lightweight caching

> Note: GoStash is an independent project and not a drop-in replacement for Redis or Memcached. Protocol and feature set are intentionally minimal.

## Architecture

GoStash follows a clean, modular architecture:

### Core Components

- **Server** (`internal/server/`) - TCP server with concurrent connection handling
- **Handler** (`internal/handler/`) - Command processing and protocol implementation
- **Store** (`internal/store/`) - Thread-safe in-memory storage backend
- **Config** (`internal/config/`) - Configuration management with file and CLI support

### Protocol

GoStash uses a custom binary protocol:

- **GET**: `GET\0<keyLen>\0<key>\r\n`
- **SET**: `SET\0<keyLen>\0<key>\0<valueLen>\0<value>\r\n`

## Project Structure

```
├── cmd/
│   ├── server/          # Server binary entrypoint
│   └── client/          # Example client implementation
├── internal/
│   ├── config/          # Configuration loading and CLI helpers
│   │   ├── config.go    # Core configuration logic
│   │   ├── cli.go       # Command-line argument parsing
│   │   └── file.go      # File-based configuration
│   ├── handler/         # Command handlers and protocol
│   │   ├── handler.go   # Main handler coordination
│   │   ├── get.go       # GET command implementation
│   │   ├── set.go       # SET command implementation
│   │   ├── commands.go  # Command definitions
│   │   └── responses.go # Response utilities
│   ├── server/          # TCP server implementation
│   └── store/           # Storage backends
│       ├── store.go     # Storage interface
│       └── hashmap.go   # HashMap implementation
├── .config.stash.example # Example configuration file
└── go.mod
```

## Quick Start

### Requirements

- Go 1.25.1+ installed

### Building and Running

**Build the server:**

```powershell
go build -o gostash.exe ./cmd/server
```

**Run with default settings:**

```powershell
# Using the built binary
./gostash.exe

# Or directly with go run
go run ./cmd/server
```

**Run with custom configuration:**

```powershell
# Using command line flags
./gostash.exe --host localhost --port 8080

# Using configuration file
./gostash.exe --config .config.stash.example
```

The server will start and listen on the configured address (default: `localhost:19201`).

### Testing the Server

You can test the server using the included client example:

```powershell
# In another terminal, run the test client
go run ./cmd/client
```

Or manually connect using telnet/netcat to test the binary protocol.

## Configuration

GoStash supports flexible configuration through both configuration files and command-line arguments:

### Configuration Options

- `host` - Server listen address (default: `localhost`)
- `port` - Server listen port (default: `19201`)

### Configuration Methods

**1. Command Line Flags:**

```powershell
./gostash.exe --host 0.0.0.0 --port 8080
```

**2. Configuration File:**

Create a configuration file (e.g., `.config.stash`):

```text
host=localhost
port=8080
```

Then run:

```powershell
./gostash.exe --config .config.stash
```

**3. Default Configuration:**

If no configuration is provided, the server uses defaults and can be customized via CLI flags.

### Example Configuration File

An example configuration file is provided as `.config.stash.example`:

```text
host=localhost
port=8080
```

## Protocol Documentation

GoStash implements a simple binary protocol for client-server communication:

### GET Command

**Format:** `GET\0<keyLen>\0<key>\r\n`

**Example:** To get the value for key "mykey":

```
GET\0005\0mykey\r\n
```

### SET Command

**Format:** `SET\0<keyLen>\0<key>\0<valueLen>\0<value>\r\n`

**Example:** To set key "mykey" to value "myvalue":

```
SET\0005\0mykey\0007\0myvalue\r\n
```

### Response Format

- **Success:** Returns the requested value followed by `\r\n`
- **Error:** Returns `ERR` status code

## Development

### Running Tests

```powershell
go test ./...
```

### Building for Production

```powershell
go build -ldflags="-s -w" -o gostash.exe ./cmd/server
```

### Client Example

The repository includes a working client example in `cmd/client/main.go` that demonstrates:

- Connecting to the server
- Sending SET command
- Sending GET command
- Reading responses

You can use this as a reference for implementing your own clients.

## Contributing

This project is designed for educational purposes and experimentation. Feel free to:

- Add new commands
- Implement additional storage backends
- Improve the protocol
- Add monitoring and metrics
- Enhance configuration options
