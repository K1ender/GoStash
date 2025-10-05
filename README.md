# GoStash

<img src="./logo.png" />

GoStash is a lightweight in-memory key-value cache inspired by Redis and Memcached. It aims to provide a simple, fast, and easy-to-run caching server for development and small production use-cases. The project is written in Go and follows a minimal design that focuses on low-latency operations and simple configuration.

## Features

- **In-memory key-value storage** - Fast HashMap-based storage with thread-safe operations
- **Binary protocol support** - Custom binary protocol for efficient communication
- **Multiple commands** - GET, SET, INCR, DECR, DEL operations with proper serialization
- **Configurable server** - Support for both file-based and CLI configuration
- **Concurrent client handling** - Each client connection handled in a separate goroutine
- **High performance** - Sub-microsecond operation latency for core commands
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

GoStash uses a custom binary protocol for all operations:

- **GET**: `GET\0<keyLen>\0<key>\r\n`
- **SET**: `SET\0<keyLen>\0<key>\0<valueLen>\0<value>\r\n`
- **INCR**: `INC\0<keyLen>\0<key>\r\n`
- **DECR**: `DEC\0<keyLen>\0<key>\r\n`
- **DEL**: `DEL\0<keyLen>\0<key>\r\n`

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
│   │   ├── incr.go      # INCR command implementation
│   │   ├── decr.go      # DECR command implementation
│   │   ├── del.go       # DEL command implementation
│   │   ├── commands.go  # Command definitions
│   │   └── responses.go # Response utilities
│   ├── server/          # TCP server implementation
│   └── store/           # Storage backends
│       ├── store.go     # Storage interface
│       ├── hashmap.go   # HashMap implementation
│       └── sharded.go   # Sharded implementation
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

GoStash implements a simple binary protocol for client-server communication. All commands use null bytes (`\0`) as delimiters and end with `\r\n`.

### Supported Commands

#### GET Command

**Format:** `GET\0<keyLen>\0<key>\r\n`
**Example:** To get the value for key "mykey":

```
GET\0005\0mykey\r\n
```

#### SET Command

**Format:** `SET\0<keyLen>\0<key>\0<valueLen>\0<value>\r\n`
**Example:** To set key "mykey" to value "myvalue":

```
SET\0005\0mykey\0007\0myvalue\r\n
```

#### INCR Command

**Format:** `INC\0<keyLen>\0<key>\r\n`
**Example:** To increment key "counter":

```
INC\0007\0counter\r\n
```

#### DECR Command

**Format:** `DEC\0<keyLen>\0<key>\r\n`
**Example:** To decrement key "counter":

```
DEC\0007\0counter\r\n
```

#### DEL Command

**Format:** `DEL\0<keyLen>\0<key>\r\n`
**Example:** To delete key "mykey":

```
DEL\0005\0mykey\r\n
```

### Response Format

- **Success:** Returns the requested value followed by `\r\n`
- **Error:** Returns `ERR` status code

## Benchmarks

Performance benchmarks (12th Gen Intel(R) Core(TM) i5-12400F, Go 1.25.1):

### Direct Handler Performance

```
BenchmarkGetHandler-12      13851630     89.23 ns/op    70 B/op    4 allocs/op
BenchmarkSetHandler-12      10332882    116.1 ns/op     89 B/op    5 allocs/op
BenchmarkIncrHandler-12      9525033    125.7 ns/op     72 B/op    4 allocs/op
BenchmarkDecrHandler-12      9213865    126.8 ns/op     72 B/op    5 allocs/op
BenchmarkDelHandler-12       8784843    135.5 ns/op     70 B/op    4 allocs/op
```

### Socket Handler Performance

```
BenchmarkSocketGetHandler-12         840018	   1255 ns/op	368 B/op	 7 allocs/op
BenchmarkSocketSetHandler-12         826018	   1268 ns/op	368 B/op	 7 allocs/op
BenchmarkSocketIncrHandler-12        929475	   1238 ns/op	240 B/op	 6 allocs/op
BenchmarkSocketDecrHandler-12        965814	   1233 ns/op	240 B/op	 6 allocs/op
BenchmarkSocketDelHandler-12         434486	   2651 ns/op	480 B/op	12 allocs/op
BenchmarkSocketRandomKeyInserts-12   666988	   1597 ns/op	341 B/op	11 allocs/op
```

### Performance Notes

- Direct handler calls achieve sub-microsecond latency (~90-140ns)
- Socket operations include network overhead but still maintain excellent performance
- All operations are thread-safe with minimal allocation overhead

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
