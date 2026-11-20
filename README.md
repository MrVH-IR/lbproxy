# lbproxy

A lightweight **Reverse Proxy and Load Balancer** written in Go.

`lbproxy` reads its configuration from a YAML file and forwards incoming HTTP requests, regardless of their paths or methods, to one of the configured backend servers.

## Features

- HTTP reverse proxy
- Load balancing across multiple backend servers
- `least_conn` and `round_robin` load balancing algorithms
- YAML-based configuration
- Daily request and error logs
- Automatic log rotation and compression (`tar.gz`)
- Configurable log retention

## Requirements

- Go (compatible with the project's `go.mod` version)
- Linux, macOS, or Windows

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd lbproxy
```

### 2. Run the application

Start the load balancer using the configuration file:

```bash
go run ./cmd/lbproxy -config config/config.yml
```

### 3. Run sample backend servers

To test load balancing, open separate terminal sessions and start the sample backend servers:

```bash
go run ./cmd/testbackend -port 9101 -name b1
```

```bash
go run ./cmd/testbackend -port 9102 -name b2
```

```bash
go run ./cmd/testbackend -port 9103 -name b3
```

The load balancer will distribute incoming requests across the configured backends.

## Configuration

Example configuration file:

```yaml
listen: ":8080"

algorithm: least_conn # least_conn or round_robin (default: least_conn)

backends:
  - "http://localhost:9101"
  - "http://localhost:9102"
  - "http://localhost:9103"

log_dir: logs
log_retention_days: 7
```

### Configuration Options

| Option | Description | Default |
|---|---|---|
| `listen` | Address and port the proxy listens on | `:8080` |
| `algorithm` | Load balancing algorithm | `least_conn` |
| `backends` | List of backend server URLs | Required |
| `log_dir` | Directory where log files are stored | `logs` |
| `log_retention_days` | Number of days to retain log archives | `7` |

### Load Balancing Algorithms

**`least_conn`**

Selects the backend with the lowest number of currently active requests.

**`round_robin`**

Distributes incoming requests across backends sequentially.

If only one backend is configured, all requests are forwarded to that backend.

## Building

Build the application binary:

```bash
mkdir -p dist

GOCACHE=/tmp/libproxy-go-cache go build -o dist/lbproxy ./cmd/lbproxy
```

Verify the generated binary:

```bash
file dist/lbproxy
ls -lh dist/lbproxy
```

## Running the Binary

After building, start the application directly:

```bash
./dist/lbproxy -config ./config/config.yml
```

You can also use a configuration file in the current directory:

```bash
./lbproxy -config ./config.yml
```

Example startup output:

```text
2026/11/21 00:15:56 proxy listening on :8080 with 1 backends (least_conn)
```

## Logging and Retention

Request and error logs are written to separate daily log files.

When a new day begins:

1. Previous daily log files are compressed into `.tar.gz` archives.
2. Request and error logs are archived separately.
3. Archives older than the configured `log_retention_days` are automatically deleted.

Each day produces up to two archives:

- Request logs
- Error logs

## Testing

Once the proxy and sample backends are running, send HTTP requests to:

```bash
curl http://localhost:8080/
```

You can also test different HTTP methods and paths:

```bash
curl http://localhost:8080/api/users
```

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"test"}'
```

All HTTP paths and methods are forwarded to the selected backend.

---

**lbproxy** — A small, configurable HTTP reverse proxy and load balancer built with Go.